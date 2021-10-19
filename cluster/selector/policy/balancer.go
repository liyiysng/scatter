package policy

import (
	"context"
	"errors"

	"github.com/liyiysng/scatter/logger"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"
)

var (
	myLog = logger.Component("policy")
)

var (
	//ErrorSessionFuncNotFount session未绑定
	ErrorSessionFuncNotFount = errors.New("session bind function not found")
	//ErrorContextIDNotBind ID未绑定
	ErrorContextIDNotBind = errors.New("id value does not in context")
	//ErrorContextNodeIDNotBind node id 未绑定
	ErrorContextNodeIDNotBind = errors.New("node id value does not in context")
	//ErrorContextBackendSessionNotBind backend session not found in context
	ErrorContextBackendSessionNotBind = errors.New("backend session does not in context")
	//ErrorContextHashAffinityCtxValueNotFound 未绑定值
	ErrorContextHashAffinityCtxValueNotFound = errors.New("hash affinity value not in context")
	//ErrorServerUnvaliable 服务器不可用
	ErrorServerUnvaliable = errors.New("server unabliable")
	//ErrorServiceFormatError 服务名称错误
	ErrorServiceFormatError = errors.New("service format error")
)

func init() {
	// 注册 session_affinity
	balancer.Register(newSessionAffinityBuilder())
	balancer.Register(newConsistentHashBuilder())
	balancer.Register(newP2CBuilder())
	balancer.Register(newPubBuilder())
	balancer.Register(newBackendSessionBuilder())
	balancer.Register(newRoundRobinBuilder())
}

// ErrorAcceptable checks if given error is acceptable.
func ErrorAcceptable(err error) bool {
	switch status.Code(err) {
	case codes.DeadlineExceeded, codes.Internal, codes.Unavailable, codes.DataLoss:
		return false
	default:
		return true
	}
}

type _nodeIDKeyType string

const _nodeIDKey _nodeIDKeyType = "_nodeIDKey"

// 指定nid ctx 需求
func WithNodeID(ctx context.Context, nid string) context.Context {
	return context.WithValue(ctx, _nodeIDKey, nid)
}

func getNodeID(ctx context.Context) (nodeID string, ok bool) {
	nodeID, ok = ctx.Value(_nodeIDKey).(string)
	return
}

// 当所有连接状态都为READY , 聚合连接状态变为READY
type allReadyConnectivityStateEvaluator struct {
	numWantToConn uint64
	numReady      uint64 // Number of addrConns in ready state.
	numConnecting uint64 // Number of addrConns in connecting state.
}

func (cse *allReadyConnectivityStateEvaluator) SetReadyConn(count uint64) {
	cse.numWantToConn = count
}

func (cse *allReadyConnectivityStateEvaluator) RecordTransition(oldState, newState connectivity.State) connectivity.State {
	// Update counters.
	for idx, state := range []connectivity.State{oldState, newState} {
		updateVal := 2*uint64(idx) - 1 // -1 for oldState and +1 for new.
		switch state {
		case connectivity.Ready:
			cse.numReady += updateVal
		case connectivity.Connecting:
			cse.numConnecting += updateVal
		}
	}

	// Evaluate.
	if cse.numReady == cse.numWantToConn {
		return connectivity.Ready
	}
	if cse.numConnecting > 0 {
		return connectivity.Connecting
	}
	return connectivity.TransientFailure
}
