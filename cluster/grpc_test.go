package cluster

import (
	"encoding/json"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/liyiysng/scatter/cluster/registry"
	"github.com/liyiysng/scatter/cluster/registry/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

func TestAddSrv(t *testing.T) {

	const (
		consulAddr = "127.0.0.1:8500"
		srvAddr    = "127.0.0.1:1155"
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
					Request: &registry.Value{
						Name: "req",
						Type: "proto.SumReq",
					},
					Response: &registry.Value{
						Name: "res",
						Type: "proto.SumRes",
					},
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
