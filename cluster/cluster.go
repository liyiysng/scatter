package cluster

import (
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
	name      string
	opts      *Options
	healthSrv *health.Server

	cachedSrv *registry.Cache
}

// Serve 开始服务
func (s *GrpcServer) Serve(lis net.Listener) error {
	// 开启健康检测
	healthgrpc.RegisterHealthServer(s, s.healthSrv)
	defer s.healthSrv.Shutdown()

	// 开始注册服务

	return s.Serve(lis)
}

// NewGrpcServer 创建grpc服务器
func NewGrpcServer(name string, reg registry.Registry, o ...IOption) *GrpcServer {

	opts := defaultOptions

	for _, v := range o {
		v.apply(&opts)
	}

	return &GrpcServer{
		Server:    grpc.NewServer(opts.grpcOpts...),
		name:      name,
		opts:      &opts,
		healthSrv: health.NewServer(),
	}
}
