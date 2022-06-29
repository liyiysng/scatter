package node

import (
	"context"
	"errors"
	"reflect"

	"github.com/liyiysng/scatter/cluster/sessionpb"
	"github.com/liyiysng/scatter/cluster/subsrvpb"
	"github.com/liyiysng/scatter/handle"
	nsession "github.com/liyiysng/scatter/node/session"
)

type ISubSrvClientBuilder interface {
	GetSubServiceClient(ctx context.Context, session interface{}, srvName string) (subsrvpb.SubServiceClient, context.Context, error)
}

// 若内部handle不支持该服务,则转发该服务
type srvHanleProxy struct {
	innerHandle         handle.IHandler
	supportSrvs         map[string]struct{}
	subsrvClientBuilder ISubSrvClientBuilder
	srvValidator        func(srvName string) bool
}

func newSrvHandleProxy(innerHandle handle.IHandler, subsrvClientBuilder ISubSrvClientBuilder, srvValidator func(srvName string) bool) *srvHanleProxy {
	return &srvHanleProxy{
		innerHandle:         innerHandle,
		supportSrvs:         map[string]struct{}{},
		subsrvClientBuilder: subsrvClientBuilder,
		srvValidator:        srvValidator,
	}
}

//! register rpc service
func (h *srvHanleProxy) Register(recv interface{}, opt ...handle.RegisterOpt) error {
	err := h.innerHandle.Register(recv, opt...)
	if err != nil {
		return nil
	}
	sname := reflect.Indirect(reflect.ValueOf(recv)).Type().Name()
	h.supportSrvs[sname] = struct{}{}
	return nil
}

//! register named service
func (h *srvHanleProxy) RegisterName(name string, recv interface{}, opt ...handle.RegisterOpt) error {
	err := h.innerHandle.RegisterName(name, recv, opt...)
	if err != nil {
		return err
	}
	h.supportSrvs[name] = struct{}{}
	return nil
}

func (h *srvHanleProxy) Call(ctx context.Context, session interface{}, serviceName string, methodName string, req []byte) (res []byte, err error) {

	if !h.srvValidator(serviceName) {
		return nil, ErrSubSrvNotSupport
	}

	if _, ok := h.supportSrvs[serviceName]; ok {
		return h.innerHandle.Call(ctx, session, serviceName, methodName, req)
	}
	client, ctx, err := h.subsrvClientBuilder.GetSubServiceClient(ctx, session, serviceName)
	if err != nil {
		return nil, err
	}

	sinfo := &sessionpb.SessionInfo{
		SType: sessionpb.SessionType_Transfer,
		TransferInfo: &sessionpb.TransferInfo{
			UID: session.(nsession.Session).GetUID(),
		},
	}

	cres, err := client.Call(ctx, &subsrvpb.CallReq{
		Sinfo:       sinfo,
		ServiceName: serviceName,
		MethodName:  methodName,
		Payload:     req,
	})
	if err != nil {
		return nil, err
	}
	if cres.ErrInfo != nil {
		if cres.ErrInfo.ErrType == subsrvpb.ErrorInfo_ErrorTypeCustom {
			return nil, handle.NewCustomErrorWithCode(cres.ErrInfo.Code, cres.ErrInfo.Err)
		} else if cres.ErrInfo.ErrType == subsrvpb.ErrorInfo_ErrorTypeCritical {
			return nil, handle.NewCriticalError(cres.ErrInfo.Err)
		}
		return nil, errors.New(cres.ErrInfo.Err)
	}
	return cres.Payload, nil
}

func (h *srvHanleProxy) Notify(ctx context.Context, session interface{}, serviceName string, methodName string, req []byte) error {

	if !h.srvValidator(serviceName) {
		return ErrSubSrvNotSupport
	}

	if _, ok := h.supportSrvs[serviceName]; ok {
		return h.innerHandle.Notify(ctx, session, serviceName, methodName, req)
	}
	client, ctx, err := h.subsrvClientBuilder.GetSubServiceClient(ctx, session, serviceName)
	if err != nil {
		return err
	}

	sinfo := &sessionpb.SessionInfo{
		SType: sessionpb.SessionType_Transfer,
		TransferInfo: &sessionpb.TransferInfo{
			UID: session.(nsession.Session).GetUID(),
		},
	}

	cres, err := client.Notify(ctx, &subsrvpb.NotifyReq{
		Sinfo:       sinfo,
		ServiceName: serviceName,
		MethodName:  methodName,
		Payload:     req,
	})
	if err != nil {
		return err
	}
	if cres.ErrInfo != nil {
		if cres.ErrInfo.ErrType == subsrvpb.ErrorInfo_ErrorTypeCustom {
			return handle.NewCustomErrorWithCode(cres.ErrInfo.Code, cres.ErrInfo.Err)
		} else if cres.ErrInfo.ErrType == subsrvpb.ErrorInfo_ErrorTypeCritical {
			return handle.NewCriticalError(cres.ErrInfo.Err)
		}
		return errors.New(cres.ErrInfo.Err)
	}
	return nil
}

//! 获取所有服务名
// 支持所有服务
func (h *srvHanleProxy) AllServiceName() []string {
	return h.innerHandle.AllServiceName()
}

func (h *srvHanleProxy) AllRecv() []interface{} {
	return h.innerHandle.AllRecv()
}
