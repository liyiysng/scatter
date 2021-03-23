package handle

import (
	"context"

	"github.com/liyiysng/scatter/logger"
)

var (
	myLog = logger.Component("handle")
)

// IHandler 前端消息处理
type IHandler interface {
	//! register rpc service
	Register(recv interface{}) error
	//! register named service
	RegisterName(name string, recv interface{}) error
	//!
	Call(ctx context.Context, session interface{}, serviceName string, methodName string, req []byte) (res []byte, err error)
	//!
	Notify(ctx context.Context, session interface{}, serviceName string, methodName string, req []byte) error
	//! 获取所有服务名
	AllServiceName() []string
}
