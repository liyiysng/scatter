package cluster

import (
	"strings"

	"github.com/liyiysng/scatter/cluster/registry"
	"github.com/liyiysng/scatter/config"
	"github.com/liyiysng/scatter/logger"
	"google.golang.org/grpc"
)

// Options grpc server 选项
type Options struct {
	id string

	grpcOpts []grpc.ServerOption

	// 向registry注册时,过滤无需注册的服务,如健康检测服务等
	registryFillter func(srvName string) bool

	// 节点元数据
	nodeMeta map[string]string

	// endpoint metas
	endpointMetas map[string] /*service name*/ map[string] /*method name*/ map[string]string

	// 注册器
	reg registry.Registry

	// 日志
	logerr logger.Logger

	// 配置
	cfg *config.Config
}

var defaultOptions = Options{
	registryFillter: func(srvName string) bool { return !strings.Contains(srvName, "grpc.health") }, // 跳过健康服务注册
}

// IOption 服务器选项
type IOption interface {
	apply(*Options)
}

type funcServerOption struct {
	f func(*Options)
}

func (fdo *funcServerOption) apply(do *Options) {
	fdo.f(do)
}

func newFuncServerOption(f func(*Options)) IOption {
	return &funcServerOption{
		f: f,
	}
}

// OptWithRegistry 设置注册器
func OptWithRegistry(reg registry.Registry) IOption {
	return newFuncServerOption(func(o *Options) {
		o.reg = reg
	})
}

// OptWithConfig 配置
func OptWithConfig(cfg *config.Config) IOption {
	return newFuncServerOption(func(o *Options) {
		o.cfg = cfg
	})
}

// OptWithLogger 配置
func OptWithLogger(l logger.Logger) IOption {
	return newFuncServerOption(func(o *Options) {
		o.logerr = l
	})
}
