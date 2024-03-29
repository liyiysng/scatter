package handle

import (
	"context"
	"errors"
	"reflect"
	"unicode"
	"unicode/utf8"

	"github.com/liyiysng/scatter/encoding"
)

var (
	// ErrRequstTypeError 服务请求参数类型错误
	ErrRequstTypeError = errors.New("handle: service request argument type error")
	// ErrResponseTypeError 服务回复参数类型错误
	ErrResponseTypeError = errors.New("handle: service response argument type error")
	// ErrSessionTypeError session 类型错误
	ErrSessionTypeError = errors.New("handle: service session argument type error")
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
	ArgsTypeValidator func(srvName string, methodName string, argsType []reflect.Type) error
	Call              func(session interface{}, srvName string, methodName string, callee func(argValues ...interface{}) error) error
}

type CallHookType func(ctx context.Context, session interface{}, srv interface{}, srvName string, methodName string, req interface{}, callee func(req interface{}) (res interface{}, err error)) error
type NotifyHookType func(ctx context.Context, session interface{}, srv interface{}, srvName string, methodName string, req interface{}, callee func(req interface{}) (err error)) error

var DefaultCallHook = func(ctx context.Context, session interface{}, srv interface{}, srvName string, methodName string, req interface{}, callee func(req interface{}) (res interface{}, err error)) error {
	_, err := callee(req)
	return err
}

var DefaultNotifyHook = func(ctx context.Context, session interface{}, srv interface{}, srvName string, methodName string, req interface{}, callee func(req interface{}) (err error)) error {
	err := callee(req)
	return err
}

// Option handle 选项设置
type Option struct {
	Codec encoding.Codec

	OptArgs *OptionalArgs

	HookCall   CallHookType
	HookNofify NotifyHookType
}

type serviceHandler struct {
	*Option
	serviceMap map[string]*service
}

// NewServiceHandle 创建服务处理
func NewServiceHandle(opt *Option) IHandler {
	if opt.Codec == nil {
		panic("invalid handle option")
	}

	if opt.HookCall == nil {
		opt.HookCall = DefaultCallHook
	}

	if opt.HookNofify == nil {
		opt.HookNofify = DefaultNotifyHook
	}

	return &serviceHandler{
		Option:     opt,
		serviceMap: make(map[string]*service),
	}
}

func (s *serviceHandler) AllServiceName() []string {
	ret := []string{}
	for k := range s.serviceMap {
		ret = append(ret, k)
	}
	return ret
}

func (s *serviceHandler) AllRecv() []interface{} {
	ret := []interface{}{}
	for _, v := range s.serviceMap {
		ret = append(ret, v.recv.Interface())
	}
	return ret
}

func (s *serviceHandler) Register(recv interface{}, opt ...RegisterOpt) error {
	return s.register(recv, "", false, opt...)
}

func (s *serviceHandler) RegisterName(name string, recv interface{}, opt ...RegisterOpt) error {
	return s.register(recv, name, true, opt...)
}

func (s *serviceHandler) Call(ctx context.Context, session interface{}, serviceName string, methodName string, req []byte) (res []byte, err error) {
	if svci, ok := s.serviceMap[serviceName]; ok {
		mtype, ok := svci.method[methodName]
		if !ok {
			err = NewCriticalErrorf("[serviceHandler.Call] can't find %s.%s", serviceName, methodName)
			return
		}
		if mtype.sgType != callType {
			err = NewCriticalErrorf("[serviceHandler.Call] %s.%s is a notify method", serviceName, methodName)
			return
		}

		argReq := reflect.New(mtype.reqType.Elem())
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

		callee := func(req interface{}) (callRes interface{}, callErr error) {
			// args
			varArgs := []reflect.Value{svci.recv, argCtx, argSession, reflect.ValueOf(req)}
			if mtype.hasOptArgs {

				optCallee := func(args ...interface{}) error {
					for i := 0; i < len(args); i++ {
						varArgs = append(varArgs, reflect.ValueOf(args[i]))
					}
					retValues = function.Call(varArgs)
					return errGetter()
				}

				callErr = s.OptArgs.Call(session, serviceName, methodName, optCallee)

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

		err = s.HookCall(ctx, session, svci, serviceName, methodName, argReq.Interface(), callee)

		if err != nil {
			return
		}

		callerRes := retValues[0].Interface()
		if callerRes == nil {
			err = NewCriticalErrorf("[serviceHandler.Call] %s.%s nil response", serviceName, methodName)
			return
		}

		res, err = s.Codec.Marshal(callerRes)
		return

	}

	err = NewCriticalErrorf("[serviceHandler.Call] can't find %s.%s", serviceName, methodName)
	return
}

func (s *serviceHandler) Notify(ctx context.Context, session interface{}, serviceName string, methodName string, req []byte) (err error) {
	if svci, ok := s.serviceMap[serviceName]; ok {
		mtype, ok := svci.method[methodName]
		if !ok {
			err = NewCriticalErrorf("[serviceHandler.Notify] can't find %s.%s", serviceName, methodName)
			return
		}
		if mtype.sgType != notifyType {
			err = NewCriticalErrorf("[serviceHandler.Notify] %s.%s is a call method", serviceName, methodName)
			return
		}

		argReq := reflect.New(mtype.reqType.Elem())
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

		callee := func(req interface{}) (callErr error) {
			// args
			varArgs := []reflect.Value{svci.recv, argCtx, argSession, reflect.ValueOf(req)}
			if mtype.hasOptArgs {

				optCallee := func(args ...interface{}) error {
					for i := 0; i < len(args); i++ {
						varArgs = append(varArgs, reflect.ValueOf(args[i]))
					}
					retValues = function.Call(varArgs)
					return errGetter()
				}

				callErr = s.OptArgs.Call(session, serviceName, methodName, optCallee)

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

		err = s.HookNofify(ctx, session, svci, serviceName, methodName, argReq.Interface(), callee)

		if err != nil {
			err = NewCriticalErrorf("[serviceHandler.Call] %s.%s error %v", serviceName, methodName, err)
			return
		}
		return

	}

	err = NewCriticalErrorf("[serviceHandler.Call] can't find %s.%s", serviceName, methodName)
	return
}

func (s *serviceHandler) register(recv interface{}, name string, useName bool, opt ...RegisterOpt) error {

	opts := defaultRegisterOpt

	for _, v := range opt {
		v(&opts)
	}

	if !opts.allowCall && !opts.allowNotify {
		return errors.New("invalid register option: all method forbiden")
	}

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
	srv.method = s.suitableMethods(srv.typ, sname, true, &opts)

	if len(srv.method) == 0 {
		str := ""

		// To help the user, see if a pointer receiver would work.
		method := s.suitableMethods(reflect.PtrTo(srv.typ), sname, false, &opts)
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

func (s *serviceHandler) suitableMethods(typ reflect.Type, srvName string, reportErr bool, opts *regiserOpt) map[string]*methodType {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name

		if opts.IgnoreMethod(mname) {
			continue
		}

		mt := &methodType{method: method}

		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		// Method needs more than 4 in arg num (include recv itself)
		if mtype.NumIn() < 4 {
			if reportErr {
				myLog.Errorf("method %s.%s has wrong number of ins: %d", srvName, mname, mtype.NumIn())
			}
			continue
		}

		// First arg need be a context.Context
		if argCtx := mtype.In(1); argCtx != typeOfContext {
			if reportErr {
				myLog.Errorf("method %s.%s first argument type not a context.Context: %v", srvName, mname, argCtx)
			}
			continue
		}

		// second arg need be a session
		argSession := mtype.In(2)
		if serr := opts.SessionTypeValidator(argSession); serr != nil {
			if reportErr {
				myLog.Errorf("method %s.%s first argument type not a session but a %v , %v", srvName, mname, argSession, serr)
			}
			continue
		}

		// forth req arg need to be a pointer an valid
		argReq := mtype.In(3)
		if argReq.Kind() != reflect.Ptr {
			if reportErr {
				myLog.Errorf("method %s.%s argument type not a pointer: %v", srvName, mname, argReq)
			}
			continue
		}
		if !isExportedOrBuiltinType(argReq) {
			if reportErr {
				myLog.Errorf("method %s.%s argument type not exported: %v", srvName, mname, argReq)
			}
			continue
		}
		err := opts.ReqTypeValidator(argReq)
		if err != nil {
			if reportErr {
				myLog.Errorf("method %s.%s argument type %v invalid %v ", srvName, mname, argReq, err)
			}
			continue
		}

		// 请求参数类型
		mt.reqType = argReq

		// 有可选参数
		if mtype.NumIn() > 4 {
			if s.OptArgs == nil {
				if reportErr {
					myLog.Errorf("method %s.%s optional args not support", srvName, mname)
				}
				continue
			}

			optArgsInTypes := []reflect.Type{}

			for i := 4; i < mtype.NumIn(); i++ {
				optArgsInTypes = append(optArgsInTypes, mtype.In(i))
			}

			err = s.OptArgs.ArgsTypeValidator(srvName, mname, optArgsInTypes)
			if err != nil {
				myLog.Errorf("method %s.%s optional args type invalid %v", srvName, mname, err)
				continue
			}
			mt.hasOptArgs = true
		}

		// Method needs one or two out.
		if mtype.NumOut() != 1 && mtype.NumOut() != 2 {
			if reportErr {
				myLog.Errorf("method %s.%s has wrong number of outs: %d", srvName, mname, mtype.NumOut())
			}
			continue
		}

		if !opts.allowCall && mtype.NumOut() != 1 { // just notfiy
			if reportErr {
				myLog.Errorf("method %s.%s has wrong number of outs: %d, call function not allow", srvName, mname, mtype.NumOut())
			}
			continue
		}

		if !opts.allowNotify && mtype.NumOut() != 2 { // just call
			if reportErr {
				myLog.Errorf("method %s.%s has wrong number of outs: %d , notify function not allow", srvName, mname, mtype.NumOut())
			}
			continue
		}

		// an notify method
		if mtype.NumOut() == 1 {
			// The return type of the method must be error.
			returnType := mtype.Out(0)
			if returnType != typeOfError {
				if reportErr {
					myLog.Errorf("method %s.%s returns %v not an error type", srvName, mname, returnType)
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
					myLog.Errorf("method %s.%s return type not a pointer: %v", srvName, mname, returnResType)
				}
				continue
			}
			if !isExportedOrBuiltinType(returnResType) {
				if reportErr {
					myLog.Errorf("method %s.%s return type not exported: %v", srvName, mname, returnResType)
				}
				continue
			}
			err = opts.ResTypeValidator(returnResType)
			if err != nil {
				if reportErr {
					myLog.Errorf("method %s.%s return type invalid :%v err: %v", srvName, mname, returnResType, err)
				}
				continue
			}
			// The second return type of the method must be error.
			if returnErrType := mtype.Out(1); returnErrType != typeOfError {
				if reportErr {
					myLog.Errorf("method %s.%s returns: %v not an error type", srvName, mname, returnErrType)
				}
				continue
			}
			mt.sgType = callType
			methods[mname] = mt
		}
	}
	return methods
}
