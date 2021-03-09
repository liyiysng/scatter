package cluster

import (
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

type grpcService struct {
	srvInfo  *registry.Service
	selector selector.Selector
}

// GrpcServer grpc服务器
type GrpcServer struct {
	*grpc.Server
	opts      *Options
	healthSrv *health.Server

	reg registry.Registry

	// remote grpc services
	remoteSrvs map[string] /*name*/ *grpcService
	// remote grpc clients
	clients map[string] /*name*/ map[string] /*node id*/ grpc.ClientConnInterface
	mu      sync.RWMutex
}

// Serve 开始服务
func (s *GrpcServer) Serve(lis net.Listener, cfg *config.Config) error {
	// 开启健康检测
	healthgrpc.RegisterHealthServer(s, s.healthSrv)
	s.healthSrv.SetServingStatus("service_health", healthgrpc.HealthCheckResponse_SERVING)
	defer s.healthSrv.Shutdown()

	srvCfg := make(map[string]*config.Service)
	err := cfg.UnmarshalKey("scatter.service", &srvCfg)
	if err != nil {
		return err
	}

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
					ID:       s.opts.id,
					Address:  lis.Addr().String(),
					Metadata: s.opts.nodeMeta,
				},
			},
			Metadata: map[string]string{},
		}

		// 设置服务元数据 和 选择策略
		selectPolic := selector.DefaultPolicy
		srvName := strings.ToLower(k)
		if svC, ok := srvCfg[srvName]; ok {
			// copy
			for ksc, vsc := range svC.Meta {
				srv.Metadata[ksc] = vsc
			}
			// 设置选择策略
			if svC.SelectPolicy != "" {
				selectPolic = svC.SelectPolicy
			}
		} else {
			s.opts.logerr.Errorf("service config %s not found", srvName)
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

		err := s.reg.Register(srv)
		if err != nil {
			return err
		}
	}

	// 检测服务
	watcher, err := s.reg.Watch()
	if err != nil {
		return err
	}
	defer watcher.Stop()

	go s.watchSrvs(watcher)

	return s.Server.Serve(lis)
}

func (s *GrpcServer) watchSrvs(watcher registry.Watcher) {
	for {
		res, err := watcher.Next()
		if err == registry.ErrWatcherStopped {
			return
		}

		if len(res.Service.Nodes) == 0 { // 若无节点数据,无需关心
			continue
		}

		// filltter
		nodes := make([]*registry.Node, 0, len(res.Service.Nodes)) // copy remove
		for _, n := range res.Service.Nodes {
			// 跳过本节点服务
			if n.ID == s.opts.id {
				continue
			}
			if !s.opts.watchFillter(res.Service.Name, n) {
				continue
			}
			nodes = append(nodes, n)
		}

		if len(nodes) == 0 {
			continue
		}

		switch res.Action {
		case "create":
			//1.{"Action":"create","Service":{"name":"SrvStrings","version":"","endpoints":null,"nodes":null}}
			//2.{"Action":"create","Service":{"name":"SrvStrings","version":"0.0.1","endpoints":[{"name":"ToLower"},{"name":"ToUpper"},{"name":"Split"}],"nodes":null}}
			//3.{"Action":"update","Service":{"name":"SrvStrings","version":"0.0.1","endpoints":[{"name":"ToLower"},{"name":"ToUpper"},{"name":"Split"}],"nodes":[{"id":"11010","address":"127.0.0.1:1155"}]}}
			fallthrough
		case "update": // create or update
			{
				s.onUpdate(res.Service, nodes)
			}
		case "delete":
			{
				s.onDelete(res.Service, nodes)
			}
		default:
			{
				s.opts.logerr.Errorf("invalid watch action %s", res.Action)
			}
		}
	}
}

func (s *GrpcServer) onUpdate(newSrv *registry.Service, nodes []*registry.Node) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 比对现有服务
	if srv, ok := s.remoteSrvs[newSrv.Name]; ok { // 服务存在
		for _, newNode := range nodes { // 是否有新节点添加
			found := false
			for _, oldNode := range srv.srvInfo.Nodes {
				if oldNode.ID == newNode.ID { // 已存在改节点 , 做更新操作 , 替换元数据
					oldNode.Metadata = newNode.Metadata
					found = true
					break
				}
			}
			if !found { // 有新节点加入
				s.addNode(newSrv, newNode)
			}
		}

	} else { //创建服务,并且添加节点
		s.addSrv(newSrv)
		for _, n := range nodes {
			s.addNode(newSrv, n)
		}
	}
}

func (s *GrpcServer) onDelete(newSrv *registry.Service, nodes []*registry.Node) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, n := range nodes {
		s.delNode(newSrv, n)
	}
}

func (s *GrpcServer) addSrv(srv *registry.Service) {

	gs := &grpcService{
		srvInfo: &registry.Service{
			Name:      srv.Name,
			Version:   srv.Version,
			Endpoints: srv.Endpoints,
			Metadata:  srv.Metadata,
		},
	}

	s.remoteSrvs[srv.Name] = gs
	buf, _ := json.Marshal(gs.srvInfo)
	s.opts.logerr.Infof("find servcie %s", string(buf))
}

func (s *GrpcServer) addNode(srv *registry.Service, node *registry.Node) {
	if rs, ok := s.remoteSrvs[srv.Name]; ok {
		if len(srv.Endpoints) > 0 {
			rs.srvInfo.Endpoints = srv.Endpoints
		}
		if len(srv.Metadata) > 0 {
			rs.srvInfo.Metadata = srv.Metadata
		}
		s.opts.logerr.Infof("find node servcice %s : %v", rs.srvInfo.Name, node)
		rs.srvInfo.Nodes = append(rs.srvInfo.Nodes, node)
	} else {
		s.opts.logerr.Errorf("[GrpcServer.addNode] servers %s not found", srv.Name)
	}
}

func (s *GrpcServer) delNode(srv *registry.Service, node *registry.Node) {
	s.opts.logerr.Infof("delete servcie %s node %v", srv.Name, node)
}

// NewGrpcServer 创建grpc服务器
func NewGrpcServer(id string, reg registry.Registry, o ...IOption) *GrpcServer {

	opts := defaultOptions

	for _, v := range o {
		v.apply(&opts)
	}

	opts.id = id

	if opts.logerr == nil {
		opts.logerr = logger.NewPrefixLogger(myLog, fmt.Sprintf("node %s:", id))
	}

	return &GrpcServer{
		Server:     grpc.NewServer(opts.grpcOpts...),
		opts:       &opts,
		reg:        reg,
		healthSrv:  health.NewServer(),
		remoteSrvs: make(map[string]*grpcService),
		clients:    make(map[string] /*name*/ map[string] /*node id*/ grpc.ClientConnInterface),
	}
}
