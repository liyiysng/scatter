package policy

import (
	"context"
	"fmt"

	"github.com/liyiysng/scatter/cluster/selector"
	"github.com/liyiysng/scatter/cluster/selector/policy/common"
	"github.com/liyiysng/scatter/logger"
	"github.com/liyiysng/scatter/util/hash"
	"google.golang.org/grpc/balancer"
)

type _backedSessionKeyType string

const (
	BackedSessionName                       = "backend_session"
	_backedSessionKey _backedSessionKeyType = "_backedSessionKey"
)

// WithBackendSessionID ctx 需求
func WithBackendSessionID(ctx context.Context, ID string) context.Context {
	ctx = context.WithValue(ctx, _backedSessionKey, ID)
	return ctx
}

type IBackendSessionMgr interface {
	GetNID(srvName, id string) (nid string, exists bool)
}

var _bsMgr IBackendSessionMgr = nil

func SetupBackendSessionMgr(smgr IBackendSessionMgr) {
	_bsMgr = smgr
}

func getBackendSessionMgr() IBackendSessionMgr {
	return _bsMgr
}

func newBackendSessionBuilder() balancer.Builder {
	return common.NewBalancerBuilderWithConnectivityStateEvaluator(BackedSessionName, &backendSessionBuilder{}, common.Config{HealthCheck: false}, func() common.IConnectivityStateEvaluator {
		return &allReadyConnectivityStateEvaluator{}
	})
}

type backendSessionBuilder struct {
}

func (b *backendSessionBuilder) Build(info common.PickerBuildInfo) balancer.Picker {

	if myLog.V(logger.VIMPORTENT) {
		myLog.Info("[backendSessionBuilder.Build]", info)
	}

	if len(info.ReadySCs) == 0 {
		return common.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	consistent := hash.New()
	subConns := map[string]balancer.SubConn{}

	for k, v := range info.ReadySCs {
		if v.Address.Attributes == nil {
			myLog.Errorf("[backendSessionBuilder.Build] attributes not fount")
			continue
		}
		nodeID := v.Address.Attributes.Value(selector.AttrKeyNodeID)
		if nodeID == nil {
			myLog.Errorf("[backendSessionBuilder.Build] attributes nodeID not fount")
			continue
		}
		strNodeID := nodeID.(string)
		subConns[strNodeID] = k

		if currentConfigInfo == nil {
			consistent.Add(strNodeID)
		}
	}

	return &backendSessionPicker{
		subConns:   subConns,
		consistent: consistent,
		srvName:    info.Target.Endpoint,
	}
}

type backendSessionPicker struct {
	subConns   map[string] /*node id*/ balancer.SubConn
	consistent *hash.Consistent
	srvName    string
}

func (p *backendSessionPicker) Pick(info balancer.PickInfo) (res balancer.PickResult, err error) {

	if myLog.V(logger.VTRACE) {
		myLog.Infof("[backendSessionPicker.Pick] %v", info.FullMethodName)
	}

	nodeID, ok := getNodeID(info.Ctx)
	if ok { // specified node
		if myLog.V(logger.VDEBUG) {
			myLog.Infof("[backendSessionPicker.Pick] %s select specified node %s", info.FullMethodName, nodeID)
		}
		if subConn, ok := p.subConns[nodeID]; ok {
			return balancer.PickResult{SubConn: subConn}, nil
		} else {
			return balancer.PickResult{}, fmt.Errorf("[backendSessionPicker.Pick] %s select specified node %s not found", info.FullMethodName, nodeID)
		}
	} else { //use backend session
		v := info.Ctx.Value(_backedSessionKey)
		if v != nil {
			if id, ok := v.(string); ok {
				smgr := getBackendSessionMgr()
				if smgr == nil {
					return balancer.PickResult{}, fmt.Errorf("[backendSessionPicker.Pick] session mgr not set , call SetupBackendSessionMgr first")
				}
				if nid, exists := smgr.GetNID(p.srvName, id); exists {
					//确保当前nid可用
					if conn, cok := p.subConns[nid]; cok {
						return balancer.PickResult{SubConn: conn}, nil
					} else {
						return balancer.PickResult{}, fmt.Errorf("[backendSessionPicker.Pick] node %s unavailable now please retry later", nid)
					}
				} else {
					return balancer.PickResult{}, fmt.Errorf("[backendSessionPicker.Pick] service %s id %s backend session not found", p.srvName, id)
				}
			}
		}
	}

	// 未提供bind的context
	return balancer.PickResult{}, ErrorServerUnvaliable
}
