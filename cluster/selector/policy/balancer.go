package policy

import (
	"errors"

	"github.com/liyiysng/scatter/logger"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/codes"
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
