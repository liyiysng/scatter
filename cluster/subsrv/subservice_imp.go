package subsrv

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/liyiysng/scatter/cluster/subsrvpb"
	"github.com/liyiysng/scatter/encoding"
	"github.com/liyiysng/scatter/handle"
)

type DummySession int64

var dummySession DummySession = 0

const (
	SubSrvGrpcName = "scatter.service.SubService"
)

// Meta key,用于从Meta中获取相关联的值
const (
	// MetaKeySubSrv 获取所能处理的子服务(handle包中,非grpc服务)
	MetaKeySubSrv = "_MetaKeySubSrv"
)

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

func NewSubServiceImp(codec encoding.Codec, callHook handle.CallHookType, notifyHook handle.NotifyHookType) *SubServiceImp {
	srv := &SubServiceImp{}

	srv.SubSrvHandle = handle.NewServiceHandle(&handle.Option{
		Codec: codec,
		ReqTypeValidator: func(reqType reflect.Type) error {
			if reqType.Implements(typeProtoMessage) {
				return nil
			}
			return handle.ErrRequstTypeError
		},
		ResTypeValidator: func(reqType reflect.Type) error {
			if reqType.Implements(typeProtoMessage) {
				return nil
			}
			return handle.ErrResponseTypeError
		},
		SessionType: reflect.TypeOf((*DummySession)(nil)),
		HookCall:    callHook,
		HookNofify:  notifyHook,
	})

	return srv
}

func (ss *SubServiceImp) Call(ctx context.Context, req *subsrvpb.CallReq) (res *subsrvpb.CallRes, err error) {

	resPayload, cerr := ss.SubSrvHandle.Call(ctx, &dummySession, req.ServiceName, req.MethodName, req.Payload)

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

	cerr := ss.SubSrvHandle.Notify(ctx, &dummySession, req.ServiceName, req.MethodName, req.Payload)

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
