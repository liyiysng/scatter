package subsrv

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/liyiysng/scatter/cluster/session"
	"github.com/liyiysng/scatter/cluster/subsrvpb"
	"github.com/liyiysng/scatter/encoding"
	"github.com/liyiysng/scatter/handle"
	"github.com/liyiysng/scatter/logger"

	//proto编码
	_ "github.com/liyiysng/scatter/encoding/proto"
)

var (
	myLog = logger.Component("subsrv")
)

const (
	SubSrvGrpcName = "scatter.service.SubService"
)

// Meta key,用于从Meta中获取相关联的值
const (
	// MetaKeySubSrv 获取所能处理的子服务(handle包中,非grpc服务)
	MetaKeySubSrv = "_MetaKeySubSrv"
)

type Session interface {
	session.ITransferSession
}

func GetSubSrvFromMeta(meta map[string]string) (srvs []string, err error) {
	if str, ok := meta[MetaKeySubSrv]; ok {
		err = json.Unmarshal([]byte(str), &srvs)
	}
	return
}

func SetSubSrvToMeta(srvs []string, meta map[string]string) error {
	buf, err := json.Marshal(srvs)
	if err != nil {
		return err
	}
	meta[MetaKeySubSrv] = string(buf)
	return nil
}

var typeProtoMessage = reflect.TypeOf((*proto.Message)(nil)).Elem()

type SubServiceImp struct {
	subsrvpb.UnimplementedSubServiceServer
	// 子服务处理
	SubSrvHandle handle.IHandler
}

var argValidator = func(t reflect.Type) error {
	if t.Implements(typeProtoMessage) {
		return nil
	}
	return handle.ErrRequstTypeError
}

var transforSessionValidator = func(t reflect.Type) error {
	stype := reflect.TypeOf((*session.ITransferSession)(nil)).Elem()
	myLog.Infof("---------valid session %v", t)
	if t != stype && !t.Implements(stype) {
		return handle.ErrSessionTypeError
	}
	return nil
}

var pubSessionValidator = func(t reflect.Type) error {
	stype := reflect.TypeOf((*session.ISession)(nil)).Elem()
	if t != stype && !t.Implements(stype) {
		return handle.ErrSessionTypeError
	}
	return nil
}

func NewSubServiceImp(callHook handle.CallHookType, notifyHook handle.NotifyHookType) *SubServiceImp {
	srv := &SubServiceImp{}

	srv.SubSrvHandle = handle.NewServiceHandle(&handle.Option{
		Codec: encoding.GetCodec("proto"),
	})

	return srv
}

// RegisterTansfered 注册转发接受
// 非协程安全
func (ss *SubServiceImp) RegisterTansfered(recv interface{}) error {
	return ss.SubSrvHandle.Register(recv,
		handle.OptWithReqTypeValidator(argValidator),
		handle.OptWithResTypeValidator(argValidator),
		handle.OptWithSessionTypeValidator(transforSessionValidator))
}

// RegisterTansferedName 注册转发接受
// 非协程安全
func (ss *SubServiceImp) RegisterTansferedName(name string, recv interface{}) error {
	return ss.SubSrvHandle.RegisterName(name, recv,
		handle.OptWithReqTypeValidator(argValidator),
		handle.OptWithResTypeValidator(argValidator),
		handle.OptWithSessionTypeValidator(transforSessionValidator))
}

// Subscribe 订阅
// 非协程安全
func (ss *SubServiceImp) Subscribe(topic string, recv interface{}) error {
	return ss.SubSrvHandle.RegisterName(topic, recv,
		handle.OptWithReqTypeValidator(argValidator),
		handle.OptWithResTypeValidator(argValidator),
		handle.OptWithSessionTypeValidator(pubSessionValidator),
		handle.OptWithForbidenCall(),
	)
}

func (ss *SubServiceImp) Call(ctx context.Context, req *subsrvpb.CallReq) (res *subsrvpb.CallRes, err error) {

	s, err := session.GetSession(req.Sinfo)
	if err != nil {
		return nil, err
	}

	resPayload, cerr := ss.SubSrvHandle.Call(ctx, s, req.ServiceName, req.MethodName, req.Payload)

	if cerr != nil {
		errInfo := &subsrvpb.ErrorInfo{
			Err: cerr.Error(),
		}
		if _, ok := cerr.(handle.ICustomError); ok {
			errInfo.ErrType = subsrvpb.ErrorInfo_ErrorTypeCustom
		} else if _, ok := cerr.(handle.ICriticalError); ok {
			errInfo.ErrType = subsrvpb.ErrorInfo_ErrorTypeCritical
		} else {
			errInfo.ErrType = subsrvpb.ErrorInfo_ErrorTypeCommon
		}
		res = &subsrvpb.CallRes{
			ErrInfo: errInfo,
		}
		return
	}

	res = &subsrvpb.CallRes{
		Payload: resPayload,
	}

	return
}

func (ss *SubServiceImp) Notify(ctx context.Context, req *subsrvpb.NotifyReq) (res *subsrvpb.NotifyRes, err error) {

	s, err := session.GetSession(req.Sinfo)
	if err != nil {
		return nil, err
	}

	cerr := ss.SubSrvHandle.Notify(ctx, s, req.ServiceName, req.MethodName, req.Payload)

	if cerr != nil {
		errInfo := &subsrvpb.ErrorInfo{
			Err: cerr.Error(),
		}
		if _, ok := cerr.(handle.ICustomError); ok {
			errInfo.ErrType = subsrvpb.ErrorInfo_ErrorTypeCustom
		} else if _, ok := cerr.(handle.ICriticalError); ok {
			errInfo.ErrType = subsrvpb.ErrorInfo_ErrorTypeCritical
		} else {
			errInfo.ErrType = subsrvpb.ErrorInfo_ErrorTypeCommon
		}
		res = &subsrvpb.NotifyRes{
			ErrInfo: errInfo,
		}
		return
	}
	res = &subsrvpb.NotifyRes{}
	return
}

func (ss *SubServiceImp) Pub(ctx context.Context, req *subsrvpb.PubReq) (res *subsrvpb.PubRes, err error) {
	s, err := session.GetSession(req.Sinfo)
	if err != nil {
		return nil, err
	}

	cerr := ss.SubSrvHandle.Notify(ctx, s, req.Topic, req.Cmd, req.Payload)

	if cerr != nil {
		errInfo := &subsrvpb.ErrorInfo{
			Err: cerr.Error(),
		}
		if _, ok := cerr.(handle.ICustomError); ok {
			errInfo.ErrType = subsrvpb.ErrorInfo_ErrorTypeCustom
		} else if _, ok := cerr.(handle.ICriticalError); ok {
			errInfo.ErrType = subsrvpb.ErrorInfo_ErrorTypeCritical
		} else {
			errInfo.ErrType = subsrvpb.ErrorInfo_ErrorTypeCommon
		}
		res = &subsrvpb.PubRes{
			ErrInfo: errInfo,
		}
		return
	}
	res = &subsrvpb.PubRes{}
	return
}
