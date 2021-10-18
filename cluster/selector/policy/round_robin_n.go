package policy

import (
	"sync"
	"fmt"
	"math/rand"

	"github.com/liyiysng/scatter/cluster/selector"
	"github.com/liyiysng/scatter/cluster/selector/policy/common"
	"github.com/liyiysng/scatter/logger"
	"google.golang.org/grpc/balancer"
)

const (
	_round_robin = "round_robin_n"
)

func newRoundRobinBuilder() balancer.Builder {
	return common.NewBalancerBuilder(_round_robin, &roundRobinBuilder{}, common.Config{HealthCheck: false})
}

type roundRobinBuilder struct {
}

func (b *roundRobinBuilder) Build(info common.PickerBuildInfo) balancer.Picker {
	if myLog.V(logger.VIMPORTENT) {
		myLog.Infof("[roundRobinBuilder.Build] %v ", info)
	}
	subConns := map[string]balancer.SubConn{}
	arrSubConns := []balancer.SubConn{}
	for k, v := range info.ReadySCs {
		if v.Address.Attributes == nil {
			myLog.Errorf("[roundRobinBuilder.Build] attributes not fount")
			continue
		}
		nodeID := v.Address.Attributes.Value(selector.AttrKeyNodeID)
		if nodeID == nil {
			myLog.Errorf("[roundRobinBuilder.Build] attributes nodeID not fount")
			continue
		}
		strNodeID := nodeID.(string)
		subConns[strNodeID] = k
		arrSubConns = append(arrSubConns,k)
	}

	if len(arrSubConns) == 0{
		return common.NewErrPicker(fmt.Errorf("no avaliable connection for %s",info.Target.Endpoint))
	}

	return &roundRobinPicker{
		subConns: subConns,
		arrSubConns:arrSubConns,
		next:int(rand.Int31n(int32(len(arrSubConns)))),
	}
}

type roundRobinPicker struct {
	arrSubConns []balancer.SubConn
	subConns map[string] /*node id*/ balancer.SubConn
	mu   sync.Mutex
	next int
}

func (p *roundRobinPicker) Pick(info balancer.PickInfo) (res balancer.PickResult, err error) {

	nodeID, ok := getNodeID(info.Ctx)
	if ok {
		if myLog.V(logger.VDEBUG) {
			myLog.Infof("[roundRobinPicker.Pick] %s select specified node %s", info.FullMethodName ,nodeID)
		}
		if subConn, ok := p.subConns[nodeID]; ok {
			return balancer.PickResult{SubConn: subConn}, nil
		}
	}else{
		if myLog.V(logger.VDEBUG) {
			myLog.Infof("[roundRobinPicker.Pick] %s select next node",info.FullMethodName)
		}
		p.mu.Lock()
		sc := p.arrSubConns[p.next]
		p.next = (p.next + 1) % len(p.subConns)
		p.mu.Unlock()
		return balancer.PickResult{SubConn: sc}, nil
	}

	return balancer.PickResult{}, ErrorServerUnvaliable
}
