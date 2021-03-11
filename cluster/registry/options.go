package registry

import (
	"context"
	"crypto/tls"
	"time"
)

// Options 选项
type Options struct {
	Addrs     []string
	Timeout   time.Duration
	Secure    bool
	TLSConfig *tls.Config
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// RegisterOptions 注册选项
type RegisterOptions struct {
	TTL     time.Duration
	GrpcTTL time.Duration
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// WatchOptions 监视选项
type WatchOptions struct {
	// Specify a service to watch
	// If blank, the watch is for all services
	Service string
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// DeregisterOptions 取消注册
type DeregisterOptions struct {
	Context context.Context
}

// GetOptions 获取服务
type GetOptions struct {
	Context context.Context
}

// ListOptions 列举
type ListOptions struct {
	Context context.Context
}

// Addrs is the registry addresses to use
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// Timeout 超时选项
func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

// Secure communication with the registry
func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

// TLSConfig Specify TLS Config
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

// RegisterTTL ttl选项
func RegisterTTL(t time.Duration) RegisterOption {
	return func(o *RegisterOptions) {
		o.TTL = t
	}
}

// RegisterGrpcTTL grpc ttl
func RegisterGrpcTTL(t time.Duration) RegisterOption {
	return func(o *RegisterOptions) {
		o.GrpcTTL = t
	}
}

// RegisterContext context 选项
func RegisterContext(ctx context.Context) RegisterOption {
	return func(o *RegisterOptions) {
		o.Context = ctx
	}
}

// WatchService Watch a service
func WatchService(name string) WatchOption {
	return func(o *WatchOptions) {
		o.Service = name
	}
}

// WatchContext context 选项
func WatchContext(ctx context.Context) WatchOption {
	return func(o *WatchOptions) {
		o.Context = ctx
	}
}

// DeregisterContext context 选项
func DeregisterContext(ctx context.Context) DeregisterOption {
	return func(o *DeregisterOptions) {
		o.Context = ctx
	}
}

// GetContext context 选项
func GetContext(ctx context.Context) GetOption {
	return func(o *GetOptions) {
		o.Context = ctx
	}
}

// ListContext context 选项
func ListContext(ctx context.Context) ListOption {
	return func(o *ListOptions) {
		o.Context = ctx
	}
}
