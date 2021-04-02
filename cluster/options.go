package cluster

import (
	"strings"

	"github.com/liyiysng/scatter/cluster/registry"
	"github.com/liyiysng/scatter/config"

	"github.com/liyiysng/scatter/handle"
	"github.com/liyiysng/scatter/logger"
	"google.golang.org/grpc"
)

const (
	NoneSubService = "none"
)

// Options grpc server 选项
type Options struct {
	id string

	// grpc 配置
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

	// sub service hook
	callHook   handle.CallHookType
	notifyHook handle.NotifyHookType
}

var defaultOptions = Options{
	registryFillter: func(srvName string) bool { return !strings.Contains(srvName, "grpc.health") }, // 跳过健康服务注册
	nodeMeta:        map[string]string{},
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

// OptWithSubSrvCallHook 配置
func OptWithSubSrvCallHook(f handle.CallHookType) IOption {
	return newFuncServerOption(func(o *Options) {
		o.callHook = f
	})
}

// OptWithSubSrvNotifyHook 配置
func OptWithSubSrvNotifyHook(f handle.NotifyHookType) IOption {
	return newFuncServerOption(func(o *Options) {
		o.notifyHook = f
	})
}

// OptWithGrpcOption 配置
func OptWithGrpcOption(gopt ...grpc.ServerOption) IOption {
	return newFuncServerOption(func(o *Options) {
		o.grpcOpts = append(o.grpcOpts, gopt...)
	})
}

// GetClientOption 获取客户端链接选项
type GetClientOption struct {
	subSrv string
	policy string
}

var defaultGetClientOption = GetClientOption{
	subSrv: NoneSubService,
}

// IGetClientOption 获取客户端链接选项
type IGetClientOption interface {
	apply(*GetClientOption)
}

type funcGetClientOption struct {
	f func(*GetClientOption)
}

func (fdo *funcGetClientOption) apply(do *GetClientOption) {
	fdo.f(do)
}

func newFuncGetClientOptio(f func(*GetClientOption)) IGetClientOption {
	return &funcGetClientOption{
		f: f,
	}
}

// GetClientOptWithSubService 子服务
func GetClientOptWithSubService(srv string) IGetClientOption {
	return newFuncGetClientOptio(func(gco *GetClientOption) {
		gco.subSrv = srv
	})
}

// GetClientOptWithPolicy 负载均衡策略
func GetClientOptWithPolicy(p string) IGetClientOption {
	return newFuncGetClientOptio(func(gco *GetClientOption) {
		gco.policy = p
	})
}

// grpc client opt
type GrpcClientOpt struct {
	// 日志
	logerr logger.Logger

	// 配置
	cfg *config.Config
}

// IGrpcClientOpt 客户端选项
type IGrpcClientOption interface {
	apply(*GrpcClientOpt)
}

type funcGrpcClientOption struct {
	f func(*GrpcClientOpt)
}

func (fdo *funcGrpcClientOption) apply(do *GrpcClientOpt) {
	fdo.f(do)
}

func newFuncGrpcClientOption(f func(*GrpcClientOpt)) IGrpcClientOption {
	return &funcGrpcClientOption{
		f: f,
	}
}

// OptGrpcClientWithLogger logger 选项
func OptGrpcClientWithLogger(logerr logger.Logger) IGrpcClientOption {
	return newFuncGrpcClientOption(func(gco *GrpcClientOpt) {
		gco.logerr = logerr
	})
}

// OptGrpcClientWithCfg 配置 选项
func OptGrpcClientWithCfg(cfg *config.Config) IGrpcClientOption {
	return newFuncGrpcClientOption(func(gco *GrpcClientOpt) {
		gco.cfg = cfg
	})
}
