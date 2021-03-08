package cluster

import (
	"strings"

	"github.com/liyiysng/scatter/cluster/registry"
	"google.golang.org/grpc"
)

// Options grpc server 选项
type Options struct {
	id string

	grpcOpts []grpc.ServerOption

	// 向registry注册时,过滤无需注册的服务,如健康检测服务等
	registryFillter func(srvName string) bool
	// 检测过滤,过滤无需关注的服务,或者节点
	// 默认跳过本节点的服务
	watchFillter func(srvName string, n *registry.Node) bool

	// 节点元数据
	nodeMeta map[string]string

	// endpoint metas
	endpointMetas map[string] /*service name*/ map[string] /*method name*/ map[string]string
}

var defaultOptions = Options{
	registryFillter: func(srvName string) bool { return !strings.Contains(srvName, "grpc.health") },              // 跳过健康服务注册
	watchFillter:    func(srvName string, n *registry.Node) bool { return !strings.Contains(srvName, "consul") }, // 过滤consul服务
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
