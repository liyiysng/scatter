package selector

import (
	"fmt"
	"strings"

	"github.com/liyiysng/scatter/cluster/registry"
	"github.com/liyiysng/scatter/cluster/registry/publisher"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

func init() {
	// 注册解析
	resolver.Register(&discoverResolverBuilder{})
}

// GetServiceID 获取serviceID
func GetServiceID(nodeID string, srvName string) string {
	return fmt.Sprintf("%s/%s", nodeID, srvName)
}

// GetServiceNameAndNodeID 根据 service ID 获取service name 和 nodeID
func GetServiceNameAndNodeID(srvID string) (nodeID string, srvName string, err error) {
	srvAndID := strings.Split(srvID, "/")
	if len(srvAndID) != 2 {
		return "", "", fmt.Errorf("invalid format service id %s", srvID)
	}
	return srvAndID[0], srvAndID[1], nil
}

type discoverResolverBuilder struct {
}

// target.Scheme	:	scatter
// target.Authority	:	nodeID
// target.Endpoint	:	serviceName
func (d *discoverResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	myLog.Info("[discoverResolverBuilder.Build] build target ", target)

	serviceName := target.Endpoint

	ret := &discoverResolver{
		cc: cc,
		sub: publisher.GetPublisher().Subscribe(func(srvName string, node *registry.Node) bool {
			return srvName == serviceName // 所有该服务的节点
		}),
		srvName: serviceName,
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
}

func (r *discoverResolver) updateCC() {
	nodes := publisher.GetPublisher().FindAllNodes(func(srv *registry.Service, node *registry.Node) bool {
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
			Attributes: attributes.New("serviceID", n.SrvNodeID, "srvName", r.srvName, "nodeID", nid, "meta", n.Metadata),
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
