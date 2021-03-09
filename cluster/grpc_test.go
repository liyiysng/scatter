package cluster

import (
	"encoding/json"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/liyiysng/scatter/cluster/cluster_testing"
	"github.com/liyiysng/scatter/cluster/registry"
	"github.com/liyiysng/scatter/cluster/registry/consul"
	"github.com/liyiysng/scatter/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	consulAddr = "127.0.0.1:8500"
)

func TestAddSrv(t *testing.T) {

	const (
		srvAddr = "127.0.0.1:1155"
	)

	wg := sync.WaitGroup{}

	//closeEvent := util.NewEvent()

	// create a grpc server
	s := grpc.NewServer()

	lis, err := net.Listen("tcp", srvAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer lis.Close()

	hs := health.NewServer()
	//hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthgrpc.RegisterHealthServer(s, hs)

	go func() {
		err := s.Serve(lis)
		if err != nil {
			myLog.Error(err)
		}
	}()

	// 确保grpc服务器启动
	time.Sleep(time.Second)

	regCreate := registry.GetRegistry("consul")

	reg := regCreate(
		registry.Addrs(consulAddr),
	)

	reg.Init(consul.WithGrpcCheck(time.Second * 2))

	err = reg.Register(
		&registry.Service{
			Name:    "reg_test",
			Version: "0.0.1",
			Endpoints: []*registry.Endpoint{
				{
					Name: "Foo1",
				},
			},
			Nodes: []*registry.Node{
				{
					ID:      "110",
					Address: srvAddr,
				},
			},
		},
		//registry.RegisterTTL(time.Second*10),
	)

	if err != nil {
		t.Fatal(err)
	}

	// 监视注册活动
	watcher, err := reg.Watch()
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		wg.Add(1)
		defer wg.Done()
		for {
			res, err := watcher.Next()
			if err != nil {
				myLog.Error(err)
				return
			}
			buf, _ := json.Marshal(res)
			myLog.Info(string(buf))
		}

	}()

	time.Sleep(time.Second * 50)

	watcher.Stop()

	wg.Wait()

}

func TestGrpc(t *testing.T) {
	regCreate := registry.GetRegistry("consul")

	reg := regCreate(
		registry.Addrs(consulAddr),
	)

	reg.Init(consul.WithGrpcCheck(time.Second * 2))

	s1 := NewGrpcServer("11010", reg)

	// register to grpc server
	cluster_testing.RegisterSrvStringsServer(s1, &srvStringsImp{})

	lis, err := net.Listen("tcp", "127.0.0.1:1155")
	if err != nil {
		t.Fatal(err)
	}
	defer lis.Close()
	defer s1.Stop()

	reg2 := regCreate(
		registry.Addrs(consulAddr),
	)

	reg2.Init(consul.WithGrpcCheck(time.Second * 2))

	s2 := NewGrpcServer("11011", reg2)

	// register to grpc server
	cluster_testing.RegisterSrvStringsServer(s2, &srvStringsImp{})

	lis2, err := net.Listen("tcp", "127.0.0.1:1156")
	if err != nil {
		t.Fatal(err)
	}
	defer lis2.Close()

	defer s2.Stop()

	reg3 := regCreate(
		registry.Addrs(consulAddr),
	)

	reg3.Init(consul.WithGrpcCheck(time.Second * 2))

	s3 := NewGrpcServer("11012", reg3)

	// register to grpc server
	cluster_testing.RegisterSrvStringsServer(s3, &srvStringsImp{})

	lis3, err := net.Listen("tcp", "127.0.0.1:1157")
	if err != nil {
		t.Fatal(err)
	}
	defer lis3.Close()

	defer s3.Stop()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s1.Serve(lis, config.NewConfig())
		if err != nil {
			myLog.Error(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s2.Serve(lis2, config.NewConfig())
		if err != nil {
			myLog.Error(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s3.Serve(lis3, config.NewConfig())
		if err != nil {
			myLog.Error(err)
		}
	}()

	go func() {
		time.Sleep(time.Second * 10)
		myLog.Info("close server 11010")
		lis.Close()
	}()

	// go func() {
	// 	time.Sleep(time.Second * 50)
	// 	myLog.Info("close server 11011")
	// 	lis2.Close()
	// }()

	wg.Wait()
	time.Sleep(time.Second * 2)
}
