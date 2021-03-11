package selector

import (
	"strings"

	"github.com/liyiysng/scatter/cluster/registry"
	"github.com/liyiysng/scatter/cluster/registry/publisher"
	"google.golang.org/grpc/resolver"
)

func init() {
	// 注册解析
	resolver.Register(&discoverResolverBuilder{})
}

type discoverResolverBuilder struct {
}

// target.Scheme	:	scatter
// target.Authority	:	nodeID
// target.Endpoint	:	serviceName
func (d *discoverResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	myLog.Info("[discoverResolverBuilder.Build] build target ", target)

	nodeID := target.Authority
	serviceName := target.Endpoint

	ret := &discoverResolver{
		cc: cc,
		sub: publisher.GetPublisher().Subscribe(func(srvName string, node *registry.Node) bool {
			srvAndID := strings.Split(node.SrvNodeID, "/")
			if len(srvAndID) != 2 {
				myLog.Errorf("srvNodeID errror %s", node.SrvNodeID)
				return srvName == serviceName
			}
			return srvAndID[0] != nodeID && srvName == serviceName
		}),
		nodeID: nodeID,
	}

	ret.run()
	ret.updateCC()

	return ret, nil
}

func (d *discoverResolverBuilder) Scheme() string {
	return "scatter"
}

type discoverResolver struct {
	cc     resolver.ClientConn
	sub    chan interface{}
	nodeID string
}

func (r *discoverResolver) updateCC() {
	nodes := publisher.GetPublisher().FindAllNodes(func(srv *registry.Service, node *registry.Node) bool {
		srvAndID := strings.Split(node.SrvNodeID, "/")
		if len(srvAndID) != 2 {
			myLog.Errorf("srvNodeID errror %s", node.SrvNodeID)
			return true
		}
		return srvAndID[0] != r.nodeID && srvAndID[1] == srv.Name // 除去本节点之外的所有节点
	})

	addrs := []resolver.Address{}

	for _, n := range nodes {
		addrs = append(addrs, resolver.Address{
			Addr: n.Address,
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
