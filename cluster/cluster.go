package cluster

import (
	"encoding/json"
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
	clients map[string] /*name*/ map[string] /*version*/ map[string] /*node id*/ grpc.ClientConnInterface
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
			myLog.Errorf("service config %s not found", srvName)
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

		buf, _ := json.Marshal(res)
		myLog.Info(string(buf))

		switch res.Action {
		case "create":
			//1.{"Action":"create","Service":{"name":"SrvStrings","version":"","endpoints":null,"nodes":null}}
			//2.{"Action":"create","Service":{"name":"SrvStrings","version":"0.0.1","endpoints":[{"name":"ToLower"},{"name":"ToUpper"},{"name":"Split"}],"nodes":null}}
			//3.{"Action":"update","Service":{"name":"SrvStrings","version":"0.0.1","endpoints":[{"name":"ToLower"},{"name":"ToUpper"},{"name":"Split"}],"nodes":[{"id":"11010","address":"127.0.0.1:1155"}]}}
			fallthrough
		case "update": // create or update
			{
				if len(res.Service.Nodes) == 0 { // 若无节点数据,无需关心
					return
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

				// 比对现有服务
				s.mu.Lock()
				defer s.mu.Unlock()

				// if srv, ok := s.remoteSrvs[res.Service.Name]; ok { // 服务存在

				// } else { //创建服务,并且添加节点
				// 	s.remoteSrvs[res.Service.Name] = &grpcService{}
				// }
			}
		case "delete":
			{
			}
		default:
			{
				myLog.Errorf("invalid watch action %s", res.Action)
				return
			}
		}

	}
}

func (s *GrpcServer) addSrv(srv *registry.Service) {

}

// NewGrpcServer 创建grpc服务器
func NewGrpcServer(id string, reg registry.Registry, o ...IOption) *GrpcServer {

	opts := defaultOptions

	for _, v := range o {
		v.apply(&opts)
	}

	opts.id = id

	return &GrpcServer{
		Server:     grpc.NewServer(opts.grpcOpts...),
		opts:       &opts,
		reg:        reg,
		healthSrv:  health.NewServer(),
		remoteSrvs: make(map[string]*grpcService),
		clients:    make(map[string] /*name*/ map[string] /*version*/ map[string] /*node id*/ grpc.ClientConnInterface),
	}
}
