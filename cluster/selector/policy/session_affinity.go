package policy

import (
	"context"
	"math/rand"
	"strings"
	"sync/atomic"

	"github.com/liyiysng/scatter/logger"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
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
	return base.NewBalancerBuilder(_sessionAffinityName, &sessionAffinityBuilder{}, base.Config{HealthCheck: false})
}

type sessionAffinityBuilder struct {
}

func (b *sessionAffinityBuilder) Build(info base.PickerBuildInfo) balancer.Picker {

	if myLog.V(logger.VIMPORTENT) {
		myLog.Info("[sessionAffinityBuilder.Build]", info)
	}

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
	srvName  string
}

func (p *sessionAffinityPicker) Pick(info balancer.PickInfo) (res balancer.PickResult, err error) {

	if myLog.V(logger.VTRACE) {
		myLog.Infof("[sessionAffinityPicker.Pick] %v", info.FullMethodName)
	}

	// 服务名可以从builder获取,但balancer.SubConn未提供相关数据
	// base.SubConnInfo 在该版本中只提供了地址[IP]信息,未/不能 提供相关attribute和meta数据(meta和attribute不同 会导致重复链接)
	// TODO grpc.balancer修复后 , 修改相关实现
	if p.srvName == "" {
		v := strings.Split(info.FullMethodName, "/")
		if len(v) != 3 {
			return balancer.PickResult{}, ErrorServiceFormatError
		}
		p.srvName = v[1]
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
