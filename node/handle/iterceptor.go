package handle

import (
	"context"

	"github.com/liyiysng/scatter/node/message"
)

// ServiceInfo 包含一元调用信息
type ServiceInfo struct {
	Service    interface{}
	FullMethod string
}

// RequestHandler req message 处理
type RequestHandler func(ctx context.Context, req interface{}) (res interface{}, err error)

// SerrvicePushInterceptor message 推送 处理拦截
type SerrvicePushInterceptor func(ctx context.Context, msg message.Message) (needPush bool, err error)

// NotifyHandler notify message 处理
type NotifyHandler func(ctx context.Context, req interface{}) (err error)

// SerrviceRequestInterceptor message 处理拦截
type SerrviceRequestInterceptor func(ctx context.Context, req interface{}, info *ServiceInfo, handler RequestHandler) (res interface{}, err error)

// SerrviceNotifyInterceptor message 处理拦截
type SerrviceNotifyInterceptor func(ctx context.Context, req interface{}, info *ServiceInfo, handler NotifyHandler) (err error)
