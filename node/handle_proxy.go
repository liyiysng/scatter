package node

import (
	"context"
	"errors"
	"reflect"
	"sync"

	"github.com/liyiysng/scatter/cluster"
	"github.com/liyiysng/scatter/cluster/subsrvpb"
	"github.com/liyiysng/scatter/handle"
	nsession "github.com/liyiysng/scatter/node/session"
)

// 若内部handle不支持该服务,则转发该服务
type srvHanleProxy struct {
	innerHandle handle.IHandler
	supportSrvs map[string]struct{}

	clientBuild  cluster.IGrpcSubSrvClient
	cmu          sync.Mutex
	clients      map[string] /*srv name*/ subsrvpb.SubServiceClient
	srvValidator func(srvName string) bool
}

func newSrvHandleProxy(innerHandle handle.IHandler, clientBuild cluster.IGrpcSubSrvClient, srvValidator func(srvName string) bool) *srvHanleProxy {
	return &srvHanleProxy{
		innerHandle:  innerHandle,
		supportSrvs:  map[string]struct{}{},
		clientBuild:  clientBuild,
		clients:      make(map[string]subsrvpb.SubServiceClient),
		srvValidator: srvValidator,
	}
}

//! register rpc service
func (h *srvHanleProxy) Register(recv interface{}) error {
	err := h.innerHandle.Register(recv)
	if err != nil {
		return nil
	}
	sname := reflect.Indirect(reflect.ValueOf(recv)).Type().Name()
	h.supportSrvs[sname] = struct{}{}
	return nil
}

//! register named service
func (h *srvHanleProxy) RegisterName(name string, recv interface{}) error {
	err := h.innerHandle.RegisterName(name, recv)
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
	client, err := h.getClient(serviceName)
	if err != nil {
		return nil, err
	}
	cres, err := client.Call(ctx, &subsrvpb.CallReq{
		UID:         session.(nsession.Session).GetUID().(int64),
		ServiceName: serviceName,
		MethodName:  methodName,
		Payload:     req,
	})
	if err != nil {
		return nil, err
	}
	if cres.ErrInfo != nil {
		if cres.ErrInfo.ErrType == subsrvpb.ErrorInfo_ErrorTypeCustom {
			return nil, handle.NewCustomError(cres.ErrInfo.Err)
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
	client, err := h.getClient(serviceName)
	if err != nil {
		return err
	}
	cres, err := client.Notify(ctx, &subsrvpb.NotifyReq{
		UID:         session.(nsession.Session).GetUID().(int64),
		ServiceName: serviceName,
		MethodName:  methodName,
		Payload:     req,
	})
	if err != nil {
		return err
	}
	if cres.ErrInfo != nil {
		if cres.ErrInfo.ErrType == subsrvpb.ErrorInfo_ErrorTypeCustom {
			return handle.NewCustomError(cres.ErrInfo.Err)
		} else if cres.ErrInfo.ErrType == subsrvpb.ErrorInfo_ErrorTypeCritical {
			return handle.NewCriticalError(cres.ErrInfo.Err)
		}
		return errors.New(cres.ErrInfo.Err)
	}
	return nil
}

func (h *srvHanleProxy) getClient(serviceName string) (c subsrvpb.SubServiceClient, err error) {
	h.cmu.Lock()
	defer h.cmu.Unlock()
	if c, ok := h.clients[serviceName]; ok {
		return c, nil
	}
	// dail
	c, err = h.clientBuild.GetSubSrvClient(serviceName)
	if err != nil {
		return nil, err
	}
	h.clients[serviceName] = c
	return c, nil
}

//! 获取所有服务名
// 支持所有服务
func (h *srvHanleProxy) AllServiceName() []string {
	return []string{}
}
