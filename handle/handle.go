package handle

import (
	"context"
	"reflect"

	"github.com/liyiysng/scatter/logger"
)

var (
	myLog = logger.Component("handle")
)

var ignoreMethods = map[string]struct{}{}

// AddIgnorMethod 添加忽略函数
// 协程不安全
// 一般在init函数中使用
func AddIgnorMethod(name string) {
	ignoreMethods[name] = struct{}{}
}

// DefaultIgnoreMethod 缺省过滤函数
var DefaultIgnoreMethod = func(methodName string) bool {
	_, ok := ignoreMethods[methodName]
	return ok
}

type regiserOpt struct {
	ReqTypeValidator     func(t reflect.Type) error   // 請求類型驗證
	ResTypeValidator     func(t reflect.Type) error   // 回复類型驗證
	SessionTypeValidator func(t reflect.Type) error   // session类型验证
	IgnoreMethod         func(methodName string) bool // 忽略的函数
	allowCall            bool
	allowNotify          bool
}

var defaultRegisterOpt = regiserOpt{
	ReqTypeValidator:     func(t reflect.Type) error { return nil },
	ResTypeValidator:     func(t reflect.Type) error { return nil },
	SessionTypeValidator: func(t reflect.Type) error { return nil },
	IgnoreMethod:         DefaultIgnoreMethod,
	allowCall:            true,
	allowNotify:          true,
}

type RegisterOpt func(o *regiserOpt)

// OptWithReqTypeValidator 请求类型验证
func OptWithReqTypeValidator(v func(t reflect.Type) error) RegisterOpt {
	return func(o *regiserOpt) {
		o.ReqTypeValidator = v
	}
}

// OptWithResTypeValidator 回复类型验证
func OptWithResTypeValidator(v func(t reflect.Type) error) RegisterOpt {
	return func(o *regiserOpt) {
		o.ResTypeValidator = v
	}
}

// OptWithSessionTypeValidator session类型验证
func OptWithSessionTypeValidator(v func(t reflect.Type) error) RegisterOpt {
	return func(o *regiserOpt) {
		o.SessionTypeValidator = v
	}
}

// OptWithIgnorMethod 忽略的函数
func OptWithIgnorMethod(f func(methodName string) bool) RegisterOpt {
	return func(o *regiserOpt) {
		o.IgnoreMethod = f
	}
}

// OptWithForbidenCall 禁止Call签名函数
func OptWithForbidenCall() RegisterOpt {
	return func(o *regiserOpt) {
		o.allowCall = false
	}
}

// OptWithForbidenNotify 禁止Notify签名函数
func OptWithForbidenNotify() RegisterOpt {
	return func(o *regiserOpt) {
		o.allowNotify = false
	}
}

// IHandler 前端消息处理
type IHandler interface {
	// Register register rpc service
	Register(recv interface{}, opt ...RegisterOpt) error
	// RegisterName register named service
	RegisterName(name string, recv interface{}, opt ...RegisterOpt) error
	// Call 调用所注册的函数
	Call(ctx context.Context, session interface{}, serviceName string, methodName string, req []byte) (res []byte, err error)
	// Notify 调用所注册的函数
	Notify(ctx context.Context, session interface{}, serviceName string, methodName string, req []byte) error
	// AllServiceName 获取所有服务名
	AllServiceName() []string
	// AllRecv 获取注册对象
	AllRecv() []interface{}
}
