package cluster

import "google.golang.org/grpc"

// Options grpc server 选项
type Options struct {
	grpcOpts []grpc.ServerOption

	regAddrs []string
}

var defaultOptions = Options{
	regAddrs: []string{"127.0.0.1:8500"},
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

// OptRegistryAddrs 注册服务地址
func OptRegistryAddrs(addrs ...string) IOption {
	return newFuncServerOption(func(o *Options) {
		o.regAddrs = addrs
	})
}
