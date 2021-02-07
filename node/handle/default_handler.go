package handle

import (
	"context"
	"errors"
	"reflect"
	"unicode"
	"unicode/utf8"

	"github.com/liyiysng/scatter/encoding"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
var typeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()

// Is this an exported - upper case - name?
func isExported(name string) bool {
	c, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(c)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

// 方法類型
type methodSgType int

const (
	callType   methodSgType = 1
	notifyType methodSgType = 2
)

type methodType struct {
	sgType     methodSgType
	method     reflect.Method
	reqType    reflect.Type // need to new
	hasOptArgs bool
}

type service struct {
	name   string                 // name of service
	recv   reflect.Value          // receiver of methods for the service
	typ    reflect.Type           // type of the receiver
	method map[string]*methodType // registered methods
}

// OptionalArgs 可选参数
type OptionalArgs struct {
	ArgsTypeValidator func(argsType []reflect.Type) error
	Call              func(session interface{}, srvName string, mName string, caller func(argValues ...interface{}) error) error
}

// Option handle 选项设置
type Option struct {
	Codec            encoding.Codec
	ReqTypeValidator func(reqType reflect.Type) error // 請求類型驗證
	ResTypeValidator func(reqType reflect.Type) error // 回复類型驗證
	SessionType      reflect.Type

	OptArgs *OptionalArgs

	HookCall   func(ctx context.Context, session interface{}, srv interface{}, srvName string, methodName string, req interface{}, caller func(req interface{}) (res interface{}, err error)) error
	HookNofify func(ctx context.Context, session interface{}, srv interface{}, srvName string, methodName string, req interface{}, caller func(req interface{}) (err error)) error
}

type serviceHandler struct {
	*Option
	serviceMap map[string]*service
}

// NewServiceHandle 创建服务处理
func NewServiceHandle(opt *Option) IHandler {
	if opt.Codec == nil || opt.HookCall == nil ||
		opt.HookNofify == nil || opt.ReqTypeValidator == nil ||
		opt.ResTypeValidator == nil || opt.SessionType == nil {
		panic("invalid handle option")
	}

	return &serviceHandler{
		Option: opt,
	}
}

func (s *serviceHandler) Register(recv interface{}) error {
	return s.register(recv, "", false)
}

func (s *serviceHandler) RegisterName(name string, recv interface{}) error {
	return s.register(recv, name, true)
}

func (s *serviceHandler) Call(ctx context.Context, session interface{}, serviceName string, methodName string, req []byte) (res []byte, err error) {
	if svci, ok := s.serviceMap[serviceName]; ok {
		mtype, ok := svci.method[methodName]
		if !ok {
			err = NewCriticalErrorf("[serviceHandler.Call] can't find service %s method %s", serviceName, methodName)
			return
		}
		if mtype.sgType != callType {
			err = NewCriticalErrorf("[serviceHandler.Call] service %s method %s is a notify method", serviceName, methodName)
			return
		}

		argReq := reflect.New(mtype.reqType)
		err = s.Codec.Unmarshal(req, argReq.Interface())
		if err != nil {
			err = NewCustomErrorWithError(err)
			return
		}
		argCtx := reflect.ValueOf(ctx)
		argSession := reflect.ValueOf(session)

		function := mtype.method.Func
		var retValues []reflect.Value

		errGetter := func() error {
			// 第二个返回值为错误对象
			retErrInterface := retValues[1].Interface()
			if retErrInterface != nil {
				return retErrInterface.(error)
			}
			return nil
		}

		caller := func(req interface{}) (callRes interface{}, callErr error) {
			// args
			varArgs := []reflect.Value{svci.recv, argCtx, argSession, reflect.ValueOf(req)}
			if mtype.hasOptArgs {

				optCaller := func(args ...interface{}) error {
					for i := 0; i < len(args); i++ {
						varArgs = append(varArgs, reflect.ValueOf(args[i]))
					}
					retValues = function.Call(varArgs)
					return errGetter()
				}

				callErr = s.OptArgs.Call(session, serviceName, methodName, optCaller)

				if callErr != nil {
					return
				}

			} else {
				// call func
				retValues = function.Call(varArgs)
			}

			// 检查错误
			callErr = errGetter()
			if callErr == nil {
				callRes = retValues[0].Interface()
				// 无错误,但是res为空
				if callRes == nil {
					callErr = errors.New("nil response")
					return
				}
			}
			return
		}

		err = s.HookCall(ctx, session, svci, serviceName, methodName, argReq.Interface(), caller)

		if err != nil {
			err = NewCriticalErrorf("[serviceHandler.Call] service %s method %s error %v", serviceName, methodName, err)
			return
		}

		callerRes := retValues[1].Interface()
		if callerRes == nil {
			err = NewCriticalErrorf("[serviceHandler.Call] service %s method %s nil response", serviceName, methodName)
			return
		}

		res, err = s.Codec.Marshal(callerRes)
		return

	}

	err = NewCriticalErrorf("[serviceHandler.Call] can't find service %s method %s", serviceName, methodName)
	return
}

func (s *serviceHandler) Notify(ctx context.Context, session interface{}, serviceName string, methodName string, req []byte) (err error) {
	if svci, ok := s.serviceMap[serviceName]; ok {
		mtype, ok := svci.method[methodName]
		if !ok {
			err = NewCriticalErrorf("[serviceHandler.Notify] can't find service %s method %s", serviceName, methodName)
			return
		}
		if mtype.sgType != notifyType {
			err = NewCriticalErrorf("[serviceHandler.Notify] service %s method %s is a call method", serviceName, methodName)
			return
		}

		argReq := reflect.New(mtype.reqType)
		err = s.Codec.Unmarshal(req, argReq.Interface())
		if err != nil {
			err = NewCustomErrorWithError(err)
			return
		}
		argCtx := reflect.ValueOf(ctx)
		argSession := reflect.ValueOf(session)

		function := mtype.method.Func
		var retValues []reflect.Value

		errGetter := func() error {
			// 第一个返回值为错误对象
			retErrInterface := retValues[0].Interface()
			if retErrInterface != nil {
				return retErrInterface.(error)
			}
			return nil
		}

		caller := func(req interface{}) (callErr error) {
			// args
			varArgs := []reflect.Value{svci.recv, argCtx, argSession, reflect.ValueOf(req)}
			if mtype.hasOptArgs {

				optCaller := func(args ...interface{}) error {
					for i := 0; i < len(args); i++ {
						varArgs = append(varArgs, reflect.ValueOf(args[i]))
					}
					retValues = function.Call(varArgs)
					return errGetter()
				}

				callErr = s.OptArgs.Call(session, serviceName, methodName, optCaller)

				if callErr != nil {
					return
				}

			} else {
				// call func
				retValues = function.Call(varArgs)
			}

			// 检查错误
			callErr = errGetter()
			return
		}

		err = s.HookNofify(ctx, session, svci, serviceName, methodName, argReq.Interface(), caller)

		if err != nil {
			err = NewCriticalErrorf("[serviceHandler.Call] service %s method %s error %v", serviceName, methodName, err)
			return
		}

		callerRes := retValues[1].Interface()
		if callerRes == nil {
			err = NewCriticalErrorf("[serviceHandler.Call] service %s method %s nil response", serviceName, methodName)
			return
		}
		return

	}

	err = NewCriticalErrorf("[serviceHandler.Call] can't find service %s method %s", serviceName, methodName)
	return
}

func (s *serviceHandler) register(recv interface{}, name string, useName bool) error {
	srv := new(service)
	srv.typ = reflect.TypeOf(recv)
	srv.recv = reflect.ValueOf(recv)
	sname := reflect.Indirect(srv.recv).Type().Name()
	if useName {
		sname = name
	}
	if sname == "" {
		return errors.New("serviceHandler.Register: no service name for type " + srv.typ.String())
	}
	if !isExported(sname) && !useName {
		return errors.New("serviceCall.Register: type " + sname + " is not exported")
	}
	srv.name = sname

	// Install the methods
	srv.method = s.suitableMethods(srv.typ, true)

	if len(srv.method) == 0 {
		str := ""

		// To help the user, see if a pointer receiver would work.
		method := s.suitableMethods(reflect.PtrTo(srv.typ), false)
		if len(method) != 0 {
			str = "serviceHandler.Register: type " + sname + " has no exported methods of suitable type (hint: pass a pointer to value of that type)"
		} else {
			str = "serviceHandler.Register: type " + sname + " has no exported methods of suitable type"
		}
		return errors.New(str)
	}

	if _, dup := s.serviceMap[sname]; dup {
		return errors.New("serviceCall: service already defined: " + sname)
	}
	s.serviceMap[sname] = srv

	return nil
}

func (s *serviceHandler) suitableMethods(typ reflect.Type, reportErr bool) map[string]*methodType {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name

		mt := &methodType{method: method}

		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		// Method needs more than 4 in arg num (include recv itself)
		if mtype.NumIn() < 4 {
			if reportErr {
				myLog.Error("method", mname, "has wrong number of ins:", mtype.NumIn())
			}
			continue
		}

		// First arg need be a context.Context
		if argCtx := mtype.In(1); argCtx != typeOfContext {
			if reportErr {
				myLog.Error("method", mname, " first argument type not a context.Context:", argCtx)
			}
			continue
		}

		// second arg need be a session
		if argSession := mtype.In(2); argSession != s.SessionType {
			if reportErr {
				myLog.Error("method", mname, " first argument type not a session:", argSession)
			}
			continue
		}

		// forth req arg need to be a pointer an valid
		argReq := mtype.In(3)
		if argReq.Kind() != reflect.Ptr {
			if reportErr {
				myLog.Error("method", mname, "argument type not a pointer:", argReq)
			}
			continue
		}
		if !isExportedOrBuiltinType(argReq) {
			if reportErr {
				myLog.Error(mname, "argument type not exported:", argReq)
			}
			continue
		}
		err := s.ReqTypeValidator(argReq)
		if err != nil {
			if reportErr {
				myLog.Error("method", mname, "argument type invalid  :", argReq, " ", err)
			}
			continue
		}

		// 请求参数类型
		mt.reqType = argReq

		// 有可选参数
		if mtype.NumIn() > 4 {
			if s.OptArgs == nil {
				if reportErr {
					myLog.Errorf("method %s optional args not support", mname)
				}
				continue
			}

			optArgsInTypes := []reflect.Type{}

			for i := 4; i < mtype.NumIn(); i++ {
				optArgsInTypes = append(optArgsInTypes, mtype.In(i))
			}

			err = s.OptArgs.ArgsTypeValidator(optArgsInTypes)
			if err != nil {
				myLog.Errorf("method %s optional args type invalid %v", mname, err)
				continue
			}
			mt.hasOptArgs = true
		}

		// Method needs one or two out.
		if mtype.NumOut() != 1 && mtype.NumOut() != 2 {
			if reportErr {
				myLog.Error("method ", mname, " has wrong number of outs:", mtype.NumOut())
			}
			continue
		}

		// an notify method
		if mtype.NumOut() == 1 {
			// The return type of the method must be error.
			returnType := mtype.Out(0)
			if returnType != typeOfError {
				if reportErr {
					myLog.Error("method", mname, "returns", returnType.String(), "not error")
				}
				continue
			}
			mt.sgType = notifyType
			methods[mname] = mt
		} else { // call method
			//The first return type of the method must be ptr
			returnResType := mtype.Out(0)
			if returnResType.Kind() != reflect.Ptr {
				if reportErr {
					myLog.Error("method", mname, "return type not a pointer:", returnResType)
				}
				continue
			}
			if !isExportedOrBuiltinType(returnResType) {
				if reportErr {
					myLog.Error("method", mname, "return type not exported:", returnResType)
				}
				continue
			}
			err = s.ResTypeValidator(returnResType)
			if err != nil {
				if reportErr {
					myLog.Error("method", mname, "return type invalid :", returnResType, " ", err)
				}
				continue
			}
			// The second return type of the method must be error.
			if returnErrType := mtype.Out(1); returnErrType != typeOfError {
				if reportErr {
					myLog.Error("method", mname, "returns", returnErrType.String(), "not error")
				}
				continue
			}
			mt.sgType = callType
			methods[mname] = mt
		}
	}
	return methods
}
