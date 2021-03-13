package policy

import (
	"context"
	"strings"

	"github.com/liyiysng/scatter/util/hash"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
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

func newConsistentHashBuilder() balancer.Builder {
	return base.NewBalancerBuilder(_consistentHashName, &consistentHashBuilder{}, base.Config{HealthCheck: false})
}

type consistentHashBuilder struct {
}

func (b *consistentHashBuilder) Build(info base.PickerBuildInfo) balancer.Picker {

	myLog.Info("[consistentHashBuilder.Build]", info)

	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	consistent := hash.New()
	subConns := map[string]balancer.SubConn{}

	for k, v := range info.ReadySCs {
		myLog.Info(v.Address.Addr)
		subConns[v.Address.Addr] = k
		consistent.Add(v.Address.Addr)
	}

	return &consistentHashPicker{
		subConns:   subConns,
		consistent: consistent,
	}
}

type consistentHashPicker struct {
	subConns   map[string]balancer.SubConn
	srvName    string
	consistent *hash.Consistent
}

func (p *consistentHashPicker) Pick(info balancer.PickInfo) (res balancer.PickResult, err error) {

	myLog.Infof("[consistentHashPicker.Pick] %v", info.FullMethodName)

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
