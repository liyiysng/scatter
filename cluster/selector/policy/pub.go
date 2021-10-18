package policy

import (
	"github.com/liyiysng/scatter/cluster/selector"
	"github.com/liyiysng/scatter/cluster/selector/policy/common"
	"github.com/liyiysng/scatter/logger"
	"google.golang.org/grpc/balancer"
)

const (
	_pubName                  = "pub"
)

func newPubBuilder() balancer.Builder {
	return common.NewBalancerBuilder(_pubName, &pubBuilder{}, common.Config{HealthCheck: false})
}

type pubBuilder struct {
}

func (b *pubBuilder) Build(info common.PickerBuildInfo) balancer.Picker {
	if myLog.V(logger.VIMPORTENT) {
		myLog.Infof("[pubBuilder.Build] %v ", info)
	}
	subConns := map[string]balancer.SubConn{}

	for k, v := range info.ReadySCs {
		if v.Address.Attributes == nil {
			myLog.Errorf("[pubBuilder.Build] attributes not fount")
			continue
		}
		nodeID := v.Address.Attributes.Value(selector.AttrKeyNodeID)
		if nodeID == nil {
			myLog.Errorf("[pubBuilder.Build] attributes nodeID not fount")
			continue
		}
		strNodeID := nodeID.(string)
		subConns[strNodeID] = k
	}

	return &pubPicker{
		subConns: subConns,
	}
}

type pubPicker struct {
	subConns map[string] /*node id*/ balancer.SubConn
}

func (p *pubPicker) Pick(info balancer.PickInfo) (res balancer.PickResult, err error) {

	nodeID, ok := info.Ctx.Value(_nodeIDKey).(string)
	if !ok {
		return balancer.PickResult{}, ErrorContextNodeIDNotBind
	}

	if myLog.V(logger.VDEBUG) {
		myLog.Infof("[pubPicker.Pick] select node %s",nodeID)
	}

	if subConn, ok := p.subConns[nodeID]; ok {
		return balancer.PickResult{SubConn: subConn}, nil
	}

	return balancer.PickResult{}, ErrorServerUnvaliable
}
