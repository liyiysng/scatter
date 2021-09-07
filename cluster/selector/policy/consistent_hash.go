package policy

import (
	"context"

	"github.com/liyiysng/scatter/cluster/selector"
	"github.com/liyiysng/scatter/cluster/selector/policy/common"
	"github.com/liyiysng/scatter/logger"
	"github.com/liyiysng/scatter/util/hash"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/resolver"
)

type _consistentHashKeyType string

const (
	_consistentHashName                           = "consistent_hash"
	_consistentHashBindKey _consistentHashKeyType = "_consistentHashBindKey"
)

// WithConsistentHashID ctx 需求
func WithConsistentHashID(ctx context.Context, ID string) context.Context {
	ctx = context.WithValue(ctx, _consistentHashBindKey, ID)
	return ctx
}

var defaultConsistentHashBuilder = &consistentHashBuilder{
	srvConsistent: map[resolver.Target]*hash.Consistent{},
}

func newConsistentHashBuilder() balancer.Builder {

	evaluator := &consistentHashConnectivityStateEvaluator{}

	return common.NewBalancerBuilderWithConnectivityStateEvaluator(_consistentHashName, &consistentHashBuilder{}, common.Config{HealthCheck: false}, evaluator)
}

// 当所有连接状态都为READY , 聚合连接状态变为READY
type consistentHashConnectivityStateEvaluator struct {
	numWantToConn uint64
	numReady      uint64 // Number of addrConns in ready state.
	numConnecting uint64 // Number of addrConns in connecting state.
}

func (cse *consistentHashConnectivityStateEvaluator) SetReadyConn(count uint64) {
	cse.numWantToConn = count
}

func (cse *consistentHashConnectivityStateEvaluator) RecordTransition(oldState, newState connectivity.State) connectivity.State {
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

type consistentHashBuilder struct {
	srvConsistent map[resolver.Target] /*service name*/ *hash.Consistent
}

func (b *consistentHashBuilder) Build(info common.PickerBuildInfo) balancer.Picker {

	if myLog.V(logger.VIMPORTENT) {
		myLog.Infof("[consistentHashBuilder.Build] %v ", info)
	}

	if myLog.V(logger.VDEBUG) {
		myLog.Info("[consistentHashBuilder.Build]--------------------------------------")
		for _, v := range info.ReadySCs {
			myLog.Infof("[consistentHashBuilder.Build] service_id[%s] srvName[%s] nodeID[%s] %p",
				v.Address.Attributes.Value(selector.AttrKeyServiceID),
				v.Address.Attributes.Value(selector.AttrKeyServiceName),
				v.Address.Attributes.Value(selector.AttrKeyNodeID),
				b)
		}
		myLog.Info("[consistentHashBuilder.Build]--------------------------------------")
	}

	consistent := hash.New()
	subConns := map[string]balancer.SubConn{}

	for k, v := range info.ReadySCs {
		if v.Address.Attributes == nil {
			myLog.Errorf("[consistentHashBuilder.Build] attributes not fount")
			continue
		}
		nodeID := v.Address.Attributes.Value(selector.AttrKeyNodeID)
		if nodeID == nil {
			myLog.Errorf("[consistentHashBuilder.Build] attributes nodeID not fount")
			continue
		}
		strNodeID := nodeID.(string)
		subConns[strNodeID] = k
		consistent.Add(strNodeID)
	}

	return &consistentHashPicker{
		srvName:    info.Target.Endpoint,
		consistent: consistent,
		subConns:   subConns,
	}
}

type consistentHashPicker struct {
	subConns   map[string] /*node id*/ balancer.SubConn
	srvName    string
	consistent *hash.Consistent
}

func (p *consistentHashPicker) Pick(info balancer.PickInfo) (res balancer.PickResult, err error) {

	//if myLog.V(logger.VTRACE) {
	//	myLog.Infof("[consistentHashPicker.Pick] %v", info.FullMethodName)
	//}

	id := info.Ctx.Value(_consistentHashBindKey)

	if id != nil {
		connStr, err := p.consistent.Get(id.(string))
		if err != nil {
			return balancer.PickResult{}, err
		}
		if subConn, ok := p.subConns[connStr]; ok {
			if myLog.V(logger.VTRACE) {
				myLog.Infof("[consistentHashPicker.Pick] id:%v connStr:%s %v",id, connStr, info.FullMethodName)
			}
			return balancer.PickResult{SubConn: subConn}, nil
		}
		return balancer.PickResult{}, ErrorServerUnvaliable
	}

	// session 中未提供bind的context
	return balancer.PickResult{}, ErrorContextIDNotBind
}
