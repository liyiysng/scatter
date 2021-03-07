package cluster

import (
	"encoding/json"
	"net"

	"github.com/liyiysng/scatter/cluster/registry"
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

// GrpcServer grpc服务器
type GrpcServer struct {
	*grpc.Server
	opts      *Options
	healthSrv *health.Server

	reg registry.Registry
}

// Serve 开始服务
func (s *GrpcServer) Serve(lis net.Listener) error {
	// 开启健康检测
	healthgrpc.RegisterHealthServer(s, s.healthSrv)
	s.healthSrv.SetServingStatus("service_health", healthgrpc.HealthCheckResponse_SERVING)
	defer s.healthSrv.Shutdown()

	// 检测服务
	watcher, err := s.reg.Watch()
	if err != nil {
		return err
	}
	defer watcher.Stop()

	go s.watchSrvs(watcher)

	// 开始注册服务
	infos := s.GetServiceInfo()
	for k, v := range infos {

		if !s.opts.registryFillter(k) {
			continue
		}

		srv := &registry.Service{
			Name:    k,
			Version: "0.0.1",
			Nodes: []*registry.Node{
				{
					ID:      s.opts.id,
					Address: lis.Addr().String(),
				},
			},
		}
		// 所有方法
		for _, m := range v.Methods {
			srv.Endpoints = append(srv.Endpoints, &registry.Endpoint{
				Name: m.Name,
			})
		}

		s.reg.Register(srv)
	}

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
	}
}

// NewGrpcServer 创建grpc服务器
func NewGrpcServer(id string, reg registry.Registry, o ...IOption) *GrpcServer {

	opts := defaultOptions

	for _, v := range o {
		v.apply(&opts)
	}

	opts.id = id

	return &GrpcServer{
		Server:    grpc.NewServer(opts.grpcOpts...),
		opts:      &opts,
		reg:       reg,
		healthSrv: health.NewServer(),
	}
}
