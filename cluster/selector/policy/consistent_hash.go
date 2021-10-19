package policy

import (
	"context"

	"github.com/liyiysng/scatter/cluster/selector"
	"github.com/liyiysng/scatter/cluster/selector/policy/common"
	"github.com/liyiysng/scatter/logger"
	"github.com/liyiysng/scatter/util/hash"
	"google.golang.org/grpc/balancer"
)

type _consistentHashKeyType string

const (
	ConsistentHashName                            = "consistent_hash"
	_consistentHashBindKey _consistentHashKeyType = "_consistentHashBindKey"
)

// WithConsistentHashID ctx 需求
func WithConsistentHashID(ctx context.Context, ID string) context.Context {
	ctx = context.WithValue(ctx, _consistentHashBindKey, ID)
	return ctx
}

func newConsistentHashBuilder() balancer.Builder {
	return common.NewBalancerBuilderWithConnectivityStateEvaluator(ConsistentHashName, &consistentHashBuilder{}, common.Config{HealthCheck: false}, func() common.IConnectivityStateEvaluator {
		return &allReadyConnectivityStateEvaluator{}
	})
}

// 当前服务配置信息
var currentConfigInfo IServiceConfigInfo = nil

// 服务节点配置信息
// 区别于从 register 获取的实时节点信息, 该信息为静态配置信息
type IServiceConfigInfo interface {
	// 根据 target 获取节点id
	GetServiceNode(srvName string) []string
}

// 设置 consistent 所需的服务配置信息
func SetServiceConfigInfo(sinfo IServiceConfigInfo) {
	currentConfigInfo = sinfo
}

type consistentHashBuilder struct {
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

		if currentConfigInfo == nil {
			consistent.Add(strNodeID)
		}
	}

	if currentConfigInfo != nil {
		nodeIDs := currentConfigInfo.GetServiceNode(info.Target.Endpoint)
		if len(nodeIDs) == 0 {
			panic("[consistentHashBuilder.Build] invalid service config info")
		}
		for _, nodeID := range nodeIDs {
			consistent.Add(nodeID)
		}
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
				myLog.Infof("[consistentHashPicker.Pick] id:%v connStr:%s %v", id, connStr, info.FullMethodName)
			}
			return balancer.PickResult{SubConn: subConn}, nil
		}
		return balancer.PickResult{}, ErrorServerUnvaliable
	}

	// session 中未提供bind的context
	return balancer.PickResult{}, ErrorContextIDNotBind
}
