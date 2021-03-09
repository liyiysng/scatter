package cluster

import (
	"context"
	"math/rand"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/liyiysng/scatter/cluster/cluster_testing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

var (
	srvAddr = []string{"127.0.0.1:5544", "127.0.0.1:5533"}
)

//////////////////////////////////////////////////////resover////////////////////////////////////////////////////////////////////////////
type resolverBuilderTest struct {
}

func (d *resolverBuilderTest) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	myLog.Info("build target ", target)
	cc.UpdateState(resolver.State{
		Addresses: []resolver.Address{{Addr: srvAddr[0]}, {Addr: srvAddr[1]}},
	})
	return &resolverTest{
		cc: cc,
	}, nil
}

func (d *resolverBuilderTest) Scheme() string {
	return "scatter"
}

type resolverTest struct {
	cc resolver.ClientConn
}

func (r *resolverTest) Close() {
	myLog.Info("[resolverTest.Close]")
}

func (r *resolverTest) ResolveNow(options resolver.ResolveNowOptions) {
	myLog.Info("[resolverTest.ResolveNow]")
}

//////////////////////////////////////////////////////resover////////////////////////////////////////////////////////////////////////////

// newBuilder creates a new roundrobin balancer builder.
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder("scatter", &pickerBuilderTest{}, base.Config{HealthCheck: false})
}

type pickerBuilderTest struct {
}

func (b *pickerBuilderTest) Build(info base.PickerBuildInfo) balancer.Picker {

	myLog.Info("[pickerBuilderTest.Build]", info)

	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	var scs []balancer.SubConn
	for sc := range info.ReadySCs {
		scs = append(scs, sc)
	}
	return &pickerTest{
		subConns: scs,
		// Start at a random index, as the same RR balancer rebuilds a new
		// picker when SubConn states change, and we don't want to apply excess
		// load to the first server in the list.
		next: rand.Intn(len(scs)),
	}
}

type pickerTest struct {
	subConns []balancer.SubConn

	mu   sync.Mutex
	next int
}

func (p *pickerTest) Pick(info balancer.PickInfo) (res balancer.PickResult, err error) {

	myLog.Infof("[pickerTest.Pick] %v", info)

	p.mu.Lock()
	sc := p.subConns[p.next]
	p.next = (p.next + 1) % len(p.subConns)
	p.mu.Unlock()
	return balancer.PickResult{SubConn: sc}, nil
}

////////////////////////////////////////////////////service/////////////////////////////////////////////////////////////////////////////
type srvStringsImp struct {
	cluster_testing.UnimplementedSrvStringsServer
	id string
}

func (s *srvStringsImp) ToLower(ctx context.Context, req *cluster_testing.String) (*cluster_testing.String, error) {
	return &cluster_testing.String{
		Str: strings.ToLower(req.Str),
	}, nil
}
func (s *srvStringsImp) ToUpper(ctx context.Context, req *cluster_testing.String) (*cluster_testing.String, error) {
	return &cluster_testing.String{
		Str: strings.ToUpper(req.Str),
	}, nil
}
func (s *srvStringsImp) Split(ctx context.Context, req *cluster_testing.String) (*cluster_testing.StringS, error) {
	myLog.Infof("[srvStringsImp.Split] call id %s", s.id)
	return &cluster_testing.StringS{
		Strs: strings.Split(req.Str, " "),
	}, nil
}

type srvIntsImp struct {
	cluster_testing.UnimplementedSrvIntsServer
}

func (s *srvIntsImp) Sum(ctx context.Context, req *cluster_testing.Ints) (*cluster_testing.Int, error) {
	sum := int64(0)
	for _, v := range req.I {
		sum += v
	}
	return &cluster_testing.Int{
		I: sum,
	}, nil
}

func (s *srvIntsImp) Multi(ctx context.Context, req *cluster_testing.Ints) (*cluster_testing.Int, error) {
	multi := int64(0)
	for _, v := range req.I {
		multi *= v
	}
	return &cluster_testing.Int{
		I: multi,
	}, nil
}

func TestGrpcBlancer(t *testing.T) {

	// 注册解析
	resolver.Register(&resolverBuilderTest{})
	// 注册负载均衡
	balancer.Register(newBuilder())

	s1 := grpc.NewServer()
	s2 := grpc.NewServer()

	lis1, err := net.Listen("tcp", srvAddr[0])
	if err != nil {
		t.Fatal(err)
	}

	lis2, err := net.Listen("tcp", srvAddr[1])
	if err != nil {
		t.Fatal(err)
	}

	cluster_testing.RegisterSrvStringsServer(s1, &srvStringsImp{id: "1"})
	cluster_testing.RegisterSrvStringsServer(s2, &srvStringsImp{id: "2"})
	cluster_testing.RegisterSrvIntsServer(s1, &srvIntsImp{})
	cluster_testing.RegisterSrvIntsServer(s2, &srvIntsImp{})

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer s1.Stop()
		err := s1.Serve(lis1)
		if err != nil {
			myLog.Error(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer s2.Stop()
		err := s2.Serve(lis2)
		if err != nil {
			myLog.Error(err)
		}
	}()

	time.Sleep(time.Second * 2)

	client, err := grpc.Dial("scatter://auth/math", grpc.WithInsecure(), grpc.WithDefaultServiceConfig(
		`{
			"loadBalancingConfig":[ { "scatter": {} } ]
		}
		`,
	))
	if err != nil {
		t.Fatal(err)
	}

	stringsClient := cluster_testing.NewSrvStringsClient(client)

	go func() {
		time.Sleep(time.Second * 10)
		lis1.Close()
	}()

	for i := 0; i < 50; i++ {
		time.Sleep(time.Second)
		_, err := stringsClient.Split(context.Background(), &cluster_testing.String{
			Str: "hello world",
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	//time.Sleep(time.Second * 30)

	// close
	lis1.Close()
	lis2.Close()

	wg.Wait()

}
