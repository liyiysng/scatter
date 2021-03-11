package policy

import (
	"context"
	"errors"
	"math/rand"
	"sync/atomic"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

var (
	// ErrorSessionFuncNotFount session未绑定
	ErrorSessionFuncNotFount = errors.New("session bind function not found")
	// ErrorServerUnvaliable 服务器不可用
	ErrorServerUnvaliable = errors.New("server unabliable")
)

type _sessionAffinityKeyType string

const (
	_sessionAffinityName                            = "session_affinity"
	_sessionAffinityBindKey _sessionAffinityKeyType = "_sessionAffinityBindKey"
	_sessionAffinityGetKey  _sessionAffinityKeyType = "_sessionAffinityGetKey"
)

// ConnBindFunc 绑定链接
type ConnBindFunc func(conn interface{})

// ConnGetFunc 链接获取
type ConnGetFunc func() (conn interface{})

// WithSessionAffinity ctx 需求
func WithSessionAffinity(ctx context.Context, bindFunc ConnBindFunc, getFunc ConnGetFunc) context.Context {
	ctx = context.WithValue(ctx, _sessionAffinityBindKey, bindFunc)
	ctx = context.WithValue(ctx, _sessionAffinityGetKey, getFunc)
	return ctx
}

func newSessionAffinityBuilder() balancer.Builder {
	return base.NewBalancerBuilder(_sessionAffinityName, &sessionAffinityBuilder{}, base.Config{HealthCheck: false})
}

type sessionAffinityBuilder struct {
}

func (b *sessionAffinityBuilder) Build(info base.PickerBuildInfo) balancer.Picker {

	myLog.Info("[sessionAffinityBuilder.Build]", info)

	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	var arr []balancer.SubConn
	for sc := range info.ReadySCs {
		arr = append(arr, sc)
	}

	return &sessionAffinityPicker{
		subConns: info.ReadySCs,
		arrConns: arr,
		next:     rand.Int63n(int64(len(info.ReadySCs))),
	}
}

type sessionAffinityPicker struct {
	subConns map[balancer.SubConn]base.SubConnInfo
	arrConns []balancer.SubConn
	next     int64
}

func (p *sessionAffinityPicker) Pick(info balancer.PickInfo) (res balancer.PickResult, err error) {

	myLog.Infof("[sessionAffinityPicker.Pick] %v", info.FullMethodName)

	getTargetConn, ok := info.Ctx.Value(_sessionAffinityGetKey).(ConnGetFunc)

	var targetConn balancer.SubConn

	if ok { // 已经存在
		conn := getTargetConn()

		if conn == nil { // 第一次调用 , 选择一个节点并绑定
			myLog.Infof("%s fist pick , bind conn", info.FullMethodName)
			// 选择一个
			index := atomic.AddInt64(&p.next, 1) % int64(len(p.arrConns))
			targetConn = p.arrConns[index]
			// 绑定
			if bindFunc, bindOK := info.Ctx.Value(_sessionAffinityBindKey).(ConnBindFunc); bindOK {
				bindFunc(targetConn)
			} else {
				return balancer.PickResult{}, ErrorSessionFuncNotFount
			}
		} else { // 链接已经存在
			myLog.Infof("%s pick conn in ctx", info.FullMethodName)
			// 验证状态是否正确
			targetConn = conn.(balancer.SubConn)
			if _, stateOK := p.subConns[targetConn]; !stateOK {
				return balancer.PickResult{}, ErrorServerUnvaliable
			}
		}

	} else {
		return balancer.PickResult{}, ErrorSessionFuncNotFount
	}

	return balancer.PickResult{SubConn: targetConn}, nil
}
