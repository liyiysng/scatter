package policy

import (
	"context"

	"github.com/liyiysng/scatter/cluster/selector/policy/common"
	"github.com/liyiysng/scatter/logger"
	"github.com/liyiysng/scatter/util/hash"
	"google.golang.org/grpc/balancer"
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
	srvConsistent: map[string]*hash.Consistent{},
}

func newConsistentHashBuilder() balancer.Builder {
	return common.NewBalancerBuilder(_consistentHashName, &consistentHashBuilder{}, common.Config{HealthCheck: false})
}

type consistentHashBuilder struct {
	srvConsistent map[string] /*service name*/ *hash.Consistent
}

func (b *consistentHashBuilder) Build(info common.PickerBuildInfo) balancer.Picker {

	if myLog.V(logger.VIMPORTENT) {
		myLog.Infof("[consistentHashBuilder.Build] %v ", info)
	}

	if myLog.V(logger.VDEBUG) {
		myLog.Info("[consistentHashBuilder.Build]--------------------------------------")
		for _, v := range info.ReadySCs {
			myLog.Infof("[consistentHashBuilder.Build] service_id[%s] srvName[%s] nodeID[%s] %p",
				v.Address.Attributes.Value("serviceID"),
				v.Address.Attributes.Value("srvName"),
				v.Address.Attributes.Value("nodeID"),
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
		nodeID := v.Address.Attributes.Value("nodeID")
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
	subConns   map[string]balancer.SubConn
	srvName    string
	consistent *hash.Consistent
}

func (p *consistentHashPicker) Pick(info balancer.PickInfo) (res balancer.PickResult, err error) {

	if myLog.V(logger.VTRACE) {
		myLog.Infof("[consistentHashPicker.Pick] %v", info.FullMethodName)
	}

	id := info.Ctx.Value(_consistentHashBindKey)

	if id != nil {
		connStr, err := p.consistent.Get(id.(string))
		if err != nil {
			return balancer.PickResult{}, err
		}
		if subConn, ok := p.subConns[connStr]; ok {
			return balancer.PickResult{SubConn: subConn}, nil
		}
		return balancer.PickResult{}, ErrorServerUnvaliable
	}

	// session 中未提供bind的context
	return balancer.PickResult{}, ErrorContextIDNotBind
}
