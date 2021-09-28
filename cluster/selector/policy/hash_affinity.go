package policy

import (
	"context"
	"fmt"
	"sync"

	"github.com/liyiysng/scatter/cluster/selector"
	"github.com/liyiysng/scatter/cluster/selector/policy/common"
	"github.com/liyiysng/scatter/logger"
	"github.com/liyiysng/scatter/util/hash"
	"google.golang.org/grpc/balancer"
)

type _hashAffinityKeyType string

const (
	HashAffinityName                      = "hash_affinity"
	_hashAffinityKey _hashAffinityKeyType = "_hashAffinityKey"
)

//IHashAffinityCtx 对hash_affinity支持
//balancer调用,用户或balancer实现
type IHashAffinityCtxValue interface {
	//获取UID
	GetUID() string
	//绑定该节点
	BindToNode(srvName, nid string, unbind func())
}

//解绑
//用户调用
type IHashAffinityCtxValueUnbind interface {
	Unbind()
}

type HashAffinityCtxValue struct {
	uid    string
	unbind func()
}

func NewHashAffinityCtxValue(uid string) (IHashAffinityCtxValue, IHashAffinityCtxValueUnbind) {
	ret := &HashAffinityCtxValue{
		uid: uid,
	}
	return ret, ret
}

func (h *HashAffinityCtxValue) GetUID() string {
	return h.uid
}
func (h *HashAffinityCtxValue) BindToNode(srvName, nid string, unbind func()) {
	h.unbind = unbind
}
func (h *HashAffinityCtxValue) Unbind() {
	if h.unbind != nil {
		h.Unbind()
	} else {
		myLog.Error("[HashAffinityCtxValue.Unbind] not bind")
	}
}

// WithHashAffinityCtx 需求
func WithHashAffinityCtx(ctx context.Context, v IHashAffinityCtxValue) context.Context {
	ctx = context.WithValue(ctx, _hashAffinityKey, v)
	return ctx
}

type hashAffinityValues struct {
	m      sync.Mutex
	values map[string] /*srv name*/ map[string] /*uid*/ string /*node id*/
}

//若没有该信息则设置保存,若有该信息则返回原来的信息
func (h *hashAffinityValues) SetNx(srvName, uid, nid string) (string, bool) {
	h.m.Lock()
	defer h.m.Unlock()
	if uids, uok := h.values[srvName]; uok {
		if nodeID, ok := uids[uid]; ok {
			return nodeID, true
		} else {
			uids[uid] = nid
		}
	} else {
		h.values[srvName] = map[string]string{uid: nid}
	}
	return nid, false
}

func (h *hashAffinityValues) Remove(srvName string, uid string) {
	h.m.Lock()
	defer h.m.Unlock()
	if uids, uok := h.values[srvName]; uok {
		delete(uids, uid)
	}
}

func (h *hashAffinityValues) SetValue(srvName string, uid string, nid string) {
	h.m.Lock()
	defer h.m.Unlock()
	if uids, uok := h.values[srvName]; uok {
		uids[uid] = nid
	} else {
		h.values[srvName] = map[string]string{uid: nid}
	}
}

var _hvalues *hashAffinityValues
var _hvaluesOnce sync.Once

func getHashAffinityValues() *hashAffinityValues {
	_hvaluesOnce.Do(func() {
		_hvalues = &hashAffinityValues{
			values: map[string]map[string]string{},
		}
	})
	return _hvalues
}

func newHashAffinityBuilder() balancer.Builder {
	return common.NewBalancerBuilder(HashAffinityName, &hashAffinityBuilder{}, common.Config{HealthCheck: false})
}

type hashAffinityBuilder struct {
}

func (b *hashAffinityBuilder) Build(info common.PickerBuildInfo) balancer.Picker {

	if myLog.V(logger.VIMPORTENT) {
		myLog.Info("[hashAffinityBuilder.Build]", info)
	}

	if len(info.ReadySCs) == 0 {
		return common.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	consistent := hash.New()
	subConns := map[string]balancer.SubConn{}

	for k, v := range info.ReadySCs {
		if v.Address.Attributes == nil {
			myLog.Errorf("[hashAffinityBuilder.Build] attributes not fount")
			continue
		}
		nodeID := v.Address.Attributes.Value(selector.AttrKeyNodeID)
		if nodeID == nil {
			myLog.Errorf("[hashAffinityBuilder.Build] attributes nodeID not fount")
			continue
		}
		strNodeID := nodeID.(string)
		subConns[strNodeID] = k

		if currentConfigInfo == nil {
			consistent.Add(strNodeID)
		}
	}

	return &hashAffinityPicker{
		subConns:   subConns,
		consistent: consistent,
		srvName:    info.Target.Endpoint,
	}
}

type hashAffinityPicker struct {
	subConns   map[string] /*node id*/ balancer.SubConn
	consistent *hash.Consistent
	srvName    string
}

func (p *hashAffinityPicker) Pick(info balancer.PickInfo) (res balancer.PickResult, err error) {

	if myLog.V(logger.VTRACE) {
		myLog.Infof("[hashAffinityPicker.Pick] %v", info.FullMethodName)
	}

	v := info.Ctx.Value(_hashAffinityKey)
	if v != nil {
		if hvalue, ok := v.(IHashAffinityCtxValue); ok {
			uid := hvalue.GetUID()
			cnid, err := p.consistent.Get(uid)
			if err != nil {
				return balancer.PickResult{}, err
			}
			// 从当前hvalues中获取已经绑定的nid , 若没有则绑定
			nid, exists := getHashAffinityValues().SetNx(p.srvName, uid, cnid)
			if !exists {
				//绑定
				hvalue.BindToNode(p.srvName, nid, func() {
					getHashAffinityValues().Remove(p.srvName, uid)
				})
			}
			//确保当前nid可用
			if conn, cok := p.subConns[nid]; cok {
				return balancer.PickResult{SubConn: conn}, nil
			} else {
				return balancer.PickResult{}, fmt.Errorf("[hashAffinityPicker.Pick] node %s unavailable now please retry later", nid)
			}
		}
	}

	// 未提供bind的context
	return balancer.PickResult{}, ErrorContextHashAffinityCtxValueNotFound
}
