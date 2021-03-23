package selector

import (
	"github.com/liyiysng/scatter/cluster/registry"
	"github.com/liyiysng/scatter/cluster/registry/publisher"
	"github.com/liyiysng/scatter/cluster/subsrv"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

func init() {
	// 注册解析
	resolver.Register(&discoverResolverBuilder{})
}

type discoverResolverBuilder struct {
}

// target.Scheme	:	scatter
// target.Authority	:	sub service
// target.Endpoint	:	serviceName
func (d *discoverResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	myLog.Info("[discoverResolverBuilder.Build] build target ", target)

	serviceName := target.Endpoint
	subSrv := target.Authority

	ret := &discoverResolver{
		cc: cc,
		sub: publisher.GetPublisher().Subscribe(func(srvName string, node *registry.Node) bool { // 需要订阅的节点
			if serviceName == subsrv.SubSrvGrpcName { // 子服务
				hsrvs, err := subsrv.GetSubSrvFromMeta(node.Metadata)
				if err != nil {
					myLog.Errorf("node meta error %v", err)
					return false
				}
				found := false
				for _, v := range hsrvs {
					if v == subSrv {
						found = true
						break
					}
				}
				return found
			}
			return srvName == serviceName // 所有该服务的节点
		}),
		srvName: serviceName,
		subSrv:  subSrv,
	}

	ret.run()
	ret.updateCC()

	return ret, nil
}

func (d *discoverResolverBuilder) Scheme() string {
	return "scatter"
}

type discoverResolver struct {
	cc      resolver.ClientConn
	sub     chan interface{}
	srvName string
	subSrv  string
}

func (r *discoverResolver) updateCC() {
	nodes := publisher.GetPublisher().FindAllNodes(func(srv *registry.Service, node *registry.Node) bool {
		if r.srvName == subsrv.SubSrvGrpcName { // 子服务
			hsrvs, err := subsrv.GetSubSrvFromMeta(node.Metadata)
			if err != nil {
				myLog.Errorf("node meta error %v", err)
				return false
			}
			found := false
			for _, v := range hsrvs {
				if v == r.subSrv {
					found = true
					break
				}
			}
			return found
		}
		return r.srvName == srv.Name // 所有该服务的节点
	})

	addrs := []resolver.Address{}

	for _, n := range nodes {
		nid, _, err := GetServiceNameAndNodeID(n.SrvNodeID)
		if err != nil {
			myLog.Error(err)
			continue
		}
		addrs = append(addrs, resolver.Address{
			Addr:       n.Address,
			Attributes: attributes.New(AttrKeyServiceID, n.SrvNodeID, AttrKeyServiceName, r.srvName, AttrKeyNodeID, nid, AttrKeyMeta, n.Metadata),
		})
	}

	if len(addrs) > 0 {
		myLog.Infof(" update address %v", addrs)
		r.cc.UpdateState(resolver.State{
			Addresses: addrs,
		})
	}
}

func (r *discoverResolver) run() {
	go func() {
		for {
			if _, ok := <-r.sub; ok {
				r.updateCC()
			} else {
				return
			}
		}
	}()
}

func (r *discoverResolver) Close() {
	publisher.GetPublisher().Evict(r.sub)
}

func (r *discoverResolver) ResolveNow(options resolver.ResolveNowOptions) { // nothing changed & noting we can do
	r.updateCC()
}
