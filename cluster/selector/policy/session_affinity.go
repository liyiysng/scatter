package policy

import (
	"context"
	"math/rand"
	"sync/atomic"

	"github.com/liyiysng/scatter/cluster/selector/policy/common"
	"github.com/liyiysng/scatter/logger"
	"google.golang.org/grpc/balancer"
)

type _sessionAffinityKeyType string

const (
	_sessionAffinityName                            = "session_affinity"
	_sessionAffinityBindKey _sessionAffinityKeyType = "_sessionAffinityBindKey"
	_sessionAffinityGetKey  _sessionAffinityKeyType = "_sessionAffinityGetKey"
)

type sessionAffinityBindData struct {
	srvSessions map[string] /*srvName*/ balancer.SubConn
}

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
	return common.NewBalancerBuilder(_sessionAffinityName, &sessionAffinityBuilder{}, common.Config{HealthCheck: false})
}

type sessionAffinityBuilder struct {
}

func (b *sessionAffinityBuilder) Build(info common.PickerBuildInfo) balancer.Picker {

	if myLog.V(logger.VIMPORTENT) {
		myLog.Info("[sessionAffinityBuilder.Build]", info)
	}

	if len(info.ReadySCs) == 0 {
		return common.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	var arr []balancer.SubConn
	for sc := range info.ReadySCs {
		arr = append(arr, sc)
	}

	return &sessionAffinityPicker{
		subConns: info.ReadySCs,
		arrConns: arr,
		srvName:  info.Target.Endpoint,
		next:     rand.Int63n(int64(len(info.ReadySCs))),
	}
}

type sessionAffinityPicker struct {
	subConns map[balancer.SubConn]common.SubConnInfo
	arrConns []balancer.SubConn
	next     int64
	srvName  string
}

func (p *sessionAffinityPicker) Pick(info balancer.PickInfo) (res balancer.PickResult, err error) {

	if myLog.V(logger.VTRACE) {
		myLog.Infof("[sessionAffinityPicker.Pick] %v", info.FullMethodName)
	}

	getFunc, ok := info.Ctx.Value(_sessionAffinityGetKey).(ConnGetFunc)

	if ok {
		bindObj := getFunc()

		var bindData *sessionAffinityBindData

		if bindObj == nil { // 绑定数据
			bindData = &sessionAffinityBindData{
				srvSessions: make(map[string]balancer.SubConn),
			}
			// 绑定
			if bindFunc, bindOK := info.Ctx.Value(_sessionAffinityBindKey).(ConnBindFunc); bindOK {
				bindFunc(bindData)
			} else {
				return balancer.PickResult{}, ErrorSessionFuncNotFount
			}
		} else {
			bindData = bindObj.(*sessionAffinityBindData)
		}

		subConn, ok := bindData.srvSessions[p.srvName]
		if !ok { // 选择一个节点
			index := atomic.AddInt64(&p.next, 1) % int64(len(p.arrConns))
			subConn = p.arrConns[index]
			bindData.srvSessions[p.srvName] = subConn
			return balancer.PickResult{SubConn: subConn}, nil
		}
		// 验证状态是否正确
		if _, stateOK := p.subConns[subConn]; !stateOK {
			return balancer.PickResult{}, ErrorServerUnvaliable
		}

		if myLog.V(logger.VTRACE) {
			myLog.Infof("pick %s ", info.FullMethodName)
		}

		return balancer.PickResult{SubConn: subConn}, nil

	}

	// session 中未提供bind的context
	return balancer.PickResult{}, ErrorSessionFuncNotFount
}
