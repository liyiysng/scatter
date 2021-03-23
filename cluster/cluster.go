package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/liyiysng/scatter/cluster/registry"
	"github.com/liyiysng/scatter/cluster/selector"
	"github.com/liyiysng/scatter/cluster/subsrv"
	"github.com/liyiysng/scatter/cluster/subsrvpb"
	"github.com/liyiysng/scatter/config"

	// consul服务注册
	_ "github.com/liyiysng/scatter/cluster/registry/consul"
	"github.com/liyiysng/scatter/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	myLog = logger.Component("cluster")
)

const (
	_selectPolicyKey = "select_policy"
)

// IGrpcClient node 的客户端
// 用于服务调用
type IGrpcClient interface {
	// GetGrpcClient 根据服务名获得客户端
	GetClient(srvName string, opts ...IGetClientOption) (c grpc.ClientConnInterface, err error)
}

// IGrpcSubSrvClient 用于服务调用
type IGrpcSubSrvClient interface {
	// GetSubSrvClient 子服务调用
	GetSubSrvClient(name string, opts ...IGetClientOption) (c subsrvpb.SubServiceClient, err error)
}

type GrpcClient struct {
	opts     GrpcClientOpt
	clientMu sync.Mutex
	clients  map[string] /*srvice name*/ *grpc.ClientConn
}

func NewGrpcClient(o ...IGrpcClientOption) *GrpcClient {

	opts := GrpcClientOpt{}
	for _, v := range o {
		v.apply(&opts)
	}

	if opts.logerr == nil {
		opts.logerr = myLog
	}

	if opts.cfg == nil {
		opts.cfg = config.GetConfig()
	}

	return &GrpcClient{
		opts:    opts,
		clients: make(map[string]*grpc.ClientConn),
	}
}

// GetClient 获取grpc客户端链接
func (s *GrpcClient) GetClient(srvName string, opts ...IGetClientOption) (c grpc.ClientConnInterface, err error) {
	s.clientMu.Lock()
	defer s.clientMu.Unlock()

	if gc, ok := s.clients[srvName]; ok {
		return gc, nil
	}

	conn, err := s.dail(srvName, opts...)
	if err != nil {
		return nil, err
	}

	s.clients[srvName] = conn

	return conn, nil
}

func (s *GrpcClient) GetSubSrvClient(name string, opts ...IGetClientOption) (c subsrvpb.SubServiceClient, err error) {

	opt := defaultGetClientOption
	for _, v := range opts {
		v.apply(&opt)
	}

	policy := selector.DefaultPolicy

	if opt.policy != "" { // 优先从选项中获取
		policy = opt.policy
	} else { // 从配置中获取
		var subSrvCfg *config.SubService
		err = s.opts.cfg.UnmarshalKey("scatter.subservice."+strings.ToLower(name), &subSrvCfg)
		if err != nil {
			return
		}
		if subSrvCfg != nil && subSrvCfg.SelectPolicy != "" {
			policy = subSrvCfg.SelectPolicy
		}
	}

	conn, err := s.dail(subsrv.SubSrvGrpcName, GetClientOptWithSubService(name), GetClientOptWithPolicy(policy))
	if err != nil {
		return nil, err
	}
	return subsrvpb.NewSubServiceClient(conn), nil
}

func (s *GrpcClient) dail(srvName string, opts ...IGetClientOption) (c *grpc.ClientConn, err error) {
	opt := defaultGetClientOption
	for _, v := range opts {
		v.apply(&opt)
	}

	// 创建链接
	ctx := context.Background()
	dialOpts := []grpc.DialOption{}

	dialOpts = append(dialOpts, grpc.WithInsecure())

	// 超时设置
	timeout := s.opts.cfg.GetDuration("scatter.gnode.dial.timeout")
	if timeout > 0 {
		dialOpts = append(dialOpts, grpc.WithBlock())
		// In the non-blocking case, the ctx does not act against the connection. It
		// only controls the setup steps.
		c, cancel := context.WithTimeout(ctx, timeout)
		ctx = c
		defer cancel()
	}

	// 优先从选项中获取策略
	policy := selector.DefaultPolicy
	if opt.policy != "" {
		policy = opt.policy
	} else { // 从服务配置中获取策略
		var srvCfg *config.Service
		err = s.opts.cfg.UnmarshalKey(strings.ToLower(srvName), &srvCfg)
		if err != nil {
			return nil, err
		}

		if srvCfg != nil {
			if srvCfg.SelectPolicy != "" {
				policy = srvCfg.SelectPolicy
			}
		}
	}

	myLog.Infof("srvName[%s] use policy [%s]", strings.ToLower(srvName), policy)
	strCfg := fmt.Sprintf(
		`
		{
			"loadBalancingConfig":[ { "%s": {} } ]
		}`,
		policy)
	dialOpts = append(dialOpts, grpc.WithDefaultServiceConfig(strCfg))

	conn, err := grpc.DialContext(
		ctx,
		fmt.Sprintf("scatter://%s/%s", opt.subSrv, srvName),
		dialOpts...,
	)
	if err != nil {
		return
	}

	return conn, nil
}

func (s *GrpcClient) Close() {
	s.clientMu.Lock()
	defer s.clientMu.Unlock()

	for _, v := range s.clients {
		err := v.Close()
		if err != nil {
			s.opts.logerr.Warningf("[GrpcClient.Close] close client %v", err)
		}
	}
	s.clients = make(map[string]*grpc.ClientConn)
}

// GrpcNode grpc服务器
type GrpcNode struct {
	*grpc.Server
	opts      *Options
	healthSrv *health.Server

	// 子服务处理
	subSrv *subsrv.SubServiceImp
}

func (s *GrpcNode) RegisterSubService(recv interface{}) error {
	return s.subSrv.SubSrvHandle.Register(recv)
}

func (s *GrpcNode) RegisterSubServiceName(name string, recv interface{}) error {
	return s.subSrv.SubSrvHandle.RegisterName(name, recv)
}

// Serve 开始服务
func (s *GrpcNode) Serve(lis net.Listener) error {

	// 开启健康检测
	healthgrpc.RegisterHealthServer(s, s.healthSrv)
	s.healthSrv.SetServingStatus("service_health", healthgrpc.HealthCheckResponse_SERVING)
	defer s.healthSrv.Shutdown()

	// 注册子服务
	subsrvpb.RegisterSubServiceServer(s, s.subSrv)

	// 开始注册服务
	srvNeedRegister, err := s.getSrvNeedRegister(lis.Addr().String())
	if err != nil {
		return err
	}

	// 注册服务
	for _, srv := range srvNeedRegister {
		var err error
		if s.opts.reg != nil {
			err = s.opts.reg.Register(srv, registry.RegisterGrpcTTL(s.opts.cfg.GetDuration("scatter.register.grpc_check_interval")))
		} else {
			err = registry.Register(srv, registry.RegisterGrpcTTL(s.opts.cfg.GetDuration("scatter.register.grpc_check_interval")))
		}
		if err != nil {
			return err
		}
		buf, _ := json.Marshal(srv)
		s.opts.logerr.Infof("register service success : %s ", string(buf))
	}

	defer func() {
		// 取消注册
		for _, srv := range srvNeedRegister {
			var err error = nil
			if s.opts.reg != nil {
				err = s.opts.reg.Deregister(srv)
			} else {
				err = registry.Deregister(srv)
			}
			if err != nil {
				s.opts.logerr.Errorf("[GrpcNode.Serve] dereigster error %v", err)
			}
		}
	}()

	return s.Server.Serve(lis)
}

func (s *GrpcNode) getSrvNeedRegister(addr string) ([]*registry.Service, error) {
	infos := s.GetServiceInfo()
	srvNeedRegister := make([]*registry.Service, 0, len(infos))
	for k, v := range infos {
		if !s.opts.registryFillter(k) {
			continue
		}

		// copy meta
		meta := map[string]string{}
		for k, v := range s.opts.nodeMeta {
			meta[k] = v
		}

		// 初始信息
		srv := &registry.Service{
			Name:    k,
			Version: "0.0.1",
			Nodes: []*registry.Node{
				{
					SrvNodeID: selector.GetServiceID(s.opts.id, k),
					Address:   addr,
					Metadata:  meta,
				},
			},
			Metadata: map[string]string{},
		}

		if k == subsrv.SubSrvGrpcName { // 包含子服务
			subsrv.SetSubSrvToMeta(s.subSrv.SubSrvHandle.AllServiceName(), srv.Nodes[0].Metadata)
		}

		// 设置服务元数据
		srvName := strings.ToLower(k)
		var srvCfg *config.Service
		err := s.opts.cfg.UnmarshalKey(srvName, &srvCfg)
		if err != nil {
			return nil, err
		}
		if srvCfg != nil {
			// copy
			for ksc, vsc := range srvCfg.Meta {
				srv.Metadata[ksc] = vsc
			}
		}

		// endpoints 元数据
		if s.opts.endpointMetas != nil {
			// 所有包含meta的方法
			if optsSrv, ok := s.opts.endpointMetas[srv.Name]; ok {
				for _, m := range v.Methods {
					if optsMethodMeta, mok := optsSrv[m.Name]; mok {
						srv.Endpoints = append(srv.Endpoints, &registry.Endpoint{
							Name:     m.Name,
							Metadata: optsMethodMeta,
						})
					}
				}
			}
		}

		srvNeedRegister = append(srvNeedRegister, srv)
	}
	return srvNeedRegister, nil
}

// NewGrpcNode 创建grpc节点
func NewGrpcNode(id string, o ...IOption) *GrpcNode {

	opts := defaultOptions

	for _, v := range o {
		v.apply(&opts)
	}

	opts.id = id

	if opts.cfg == nil { // use global config
		opts.cfg = config.GetConfig()
	}

	if opts.logerr == nil {
		opts.logerr = logger.NewPrefixLogger(myLog, fmt.Sprintf("node %s:", id))
	}

	n := &GrpcNode{
		Server:    grpc.NewServer(opts.grpcOpts...),
		opts:      &opts,
		healthSrv: health.NewServer(),
		subSrv:    subsrv.NewSubServiceImp(opts.getCodec(), opts.callHook, opts.notifyHook),
	}

	return n
}
