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

func getServiceNodeID(nodeID, srvName string) string {
	return fmt.Sprintf("%s/%s", nodeID, srvName)
}

// GrpcNode grpc服务器
type GrpcNode struct {
	*grpc.Server
	opts      *Options
	healthSrv *health.Server

	clientMu sync.Mutex
	clients  map[string] /*srvice name*/ grpc.ClientConnInterface
}

// Serve 开始服务
func (s *GrpcNode) Serve(lis net.Listener) error {
	// 开启健康检测
	healthgrpc.RegisterHealthServer(s, s.healthSrv)
	s.healthSrv.SetServingStatus("service_health", healthgrpc.HealthCheckResponse_SERVING)
	defer s.healthSrv.Shutdown()

	// 开始注册服务
	infos := s.GetServiceInfo()
	for k, v := range infos {

		if !s.opts.registryFillter(k) {
			continue
		}

		// 初始信息
		srv := &registry.Service{
			Name:    k,
			Version: "0.0.1",
			Nodes: []*registry.Node{
				{
					SrvNodeID: getServiceNodeID(s.opts.id, k),
					Address:   lis.Addr().String(),
					Metadata:  s.opts.nodeMeta,
				},
			},
			Metadata: map[string]string{},
		}

		// 设置服务元数据 和 选择策略
		srvName := strings.ToLower(k)
		var srvCfg *config.Service
		err := s.opts.cfg.UnmarshalKey(srvName, &srvCfg)
		if err != nil {
			return err
		}
		selectPolic := selector.DefaultPolicy
		if srvCfg != nil {
			// copy
			for ksc, vsc := range srvCfg.Meta {
				srv.Metadata[ksc] = vsc
			}
			// 设置选择策略
			if srvCfg.SelectPolicy != "" {
				selectPolic = srvCfg.SelectPolicy
			}
		}
		srv.Metadata[_selectPolicyKey] = selectPolic

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

	return s.Server.Serve(lis)
}

// GetClient 获取grpc客户端链接
func (s *GrpcNode) GetClient(srvName string) (c grpc.ClientConnInterface, err error) {
	s.clientMu.Lock()
	defer s.clientMu.Unlock()

	if gc, ok := s.clients[srvName]; ok {
		return gc, nil
	}

	// 创建链接并保存
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

	// 服务配置
	var srvCfg *config.Service
	err = s.opts.cfg.UnmarshalKey(strings.ToLower(srvName), &srvCfg)
	if err != nil {
		return nil, err
	}
	policy := selector.DefaultPolicy
	if srvCfg != nil {
		if srvCfg.SelectPolicy != "" {
			policy = srvCfg.SelectPolicy
		}
	}
	myLog.Infof("cfg %v , srvName %s policy %s", srvCfg, strings.ToLower(srvName), policy)
	strCfg := fmt.Sprintf(
		`
		{
			"loadBalancingConfig":[ { "%s": {} } ]
		}`,
		policy)
	dialOpts = append(dialOpts, grpc.WithDefaultServiceConfig(strCfg))

	c, err = grpc.DialContext(
		ctx,
		fmt.Sprintf("scatter://%s/%s", s.opts.id, srvName),
		dialOpts...,
	)
	if err != nil {
		return
	}

	s.clients[srvName] = c

	return
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

	return &GrpcNode{
		Server:    grpc.NewServer(opts.grpcOpts...),
		opts:      &opts,
		healthSrv: health.NewServer(),
		clients:   make(map[string]grpc.ClientConnInterface),
	}
}
