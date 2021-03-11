package cluster

import (
	"encoding/json"
	"net"
	"sync"
	"testing"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/liyiysng/scatter/cluster/registry"
	"github.com/liyiysng/scatter/cluster/registry/consul"
	_ "github.com/liyiysng/scatter/cluster/registry/consul"
	"github.com/liyiysng/scatter/cluster/registry/publisher"
	"github.com/liyiysng/scatter/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

func TestRegistyService(t *testing.T) {

	cfg := consulapi.DefaultConfig()
	cfg.Address = "127.0.0.1:8500"
	client, err := consulapi.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}

	closeEvent := util.NewEvent()

	regCreate := registry.GetRegistry("consul")

	reg := regCreate(
		registry.Addrs("127.0.0.1:8500"),
	)

	reg.Init()

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
					SrvNodeID: "110",
					Address:   "127.0.0.1:1155",
				},
			},
		},
		registry.RegisterTTL(time.Second*10),
	)

	if err != nil {
		t.Fatal(err)
	}

	go func() {

		wg.Add(1)
		defer wg.Done()

		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()
		for {
			select {
			case <-closeEvent.Done():
				{
					return
				}
			case <-ticker.C:
				{
					err := client.Agent().UpdateTTL("service:110", "", "pass")
					if err != nil {
						myLog.Error(err)
					}
				}
			}
		}

	}()

	time.Sleep(time.Second * 3)

	go func() {
		time.Sleep(time.Second * 20)
		closeEvent.Fire()
	}()

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

func TestRegistyMultiService(t *testing.T) {
	regCreate := registry.GetRegistry("consul")

	reg := regCreate(
		registry.Addrs("127.0.0.1:8500"),
	)

	reg.Init()

	err := reg.Register(
		&registry.Service{
			Name:    "reg_test",
			Version: "0.0.1",
			Endpoints: []*registry.Endpoint{
				{
					Name:     "Foo1",
					Metadata: map[string]string{"foo": "bar", "fooo": "barr"},
				},
			},
			Nodes: []*registry.Node{
				{
					SrvNodeID: "110",
					Address:   "127.0.0.1:1155",
					Metadata:  map[string]string{"nm": "bar", "nm1": "barr"},
				},
			},
			Metadata: map[string]string{"sm": "bar", "sm1": "barr"},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	err = reg.Register(
		&registry.Service{
			Name:    "reg_test",
			Version: "0.0.2",
			Endpoints: []*registry.Endpoint{
				{
					Name: "Foo1",
				},
			},
			Nodes: []*registry.Node{
				{
					SrvNodeID: "220",
					Address:   "127.0.0.1:1154",
				},
			},
		},
		consul.WithRegistryTags([]string{"chat"}),
	)

	if err != nil {
		t.Fatal(err)
	}

	err = reg.Register(
		&registry.Service{
			Name:    "reg_test",
			Version: "0.0.2",
			Endpoints: []*registry.Endpoint{
				{
					Name: "Foo1",
				},
			},
			Nodes: []*registry.Node{
				{
					SrvNodeID: "330",
					Address:   "127.0.0.1:1158",
				},
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 2)

}

var (
	grpcAddr = []string{"127.0.0.1:5544", "127.0.0.1:5533"}
)

func TestRegistrySingleton(t *testing.T) {

	s := grpc.NewServer()

	lis, err := net.Listen("tcp", grpcAddr[0])
	if err != nil {
		t.Fatal(err)
	}
	defer lis.Close()

	hs := health.NewServer()
	healthgrpc.RegisterHealthServer(s, hs)

	go func() {
		defer s.Stop()
		err := s.Serve(lis)
		if err != nil {
			myLog.Error(err)
		}
	}()

	registry.Register(&registry.Service{
		Name:    "singleton_test_1",
		Version: "0.0.2",
		Nodes: []*registry.Node{
			{
				SrvNodeID: "330",
				Address:   grpcAddr[0],
			},
		},
	}, registry.RegisterGrpcTTL(time.Second*5))

	// registry.Register(&registry.Service{
	// 	Name:    "singleton_test_2",
	// 	Version: "0.0.2",
	// 	Nodes: []*registry.Node{
	// 		{
	// 			ID:      "330",
	// 			Address: grpcAddr[0],
	// 		},
	// 	},
	// }, registry.RegisterGrpcTTL(time.Second*5))

	go func() {
		chanRes := publisher.GetPublisher().Subscribe(func(srvName string, node *registry.Node) bool {
			return true
		})
		for {
			_, cok := <-chanRes
			if !cok {
				return
			}
			allNode := publisher.GetPublisher().FindAllNodes(func(srv *registry.Service, node *registry.Node) bool {
				return true
			})
			myLog.Info("--------------current nodes-------------------------------------")
			for _, v := range allNode {
				myLog.Info(v)
			}
			myLog.Info("----------end of current nodes-----------------------------------\n")
		}
	}()

	time.Sleep(time.Second * 50)

	lis.Close()

	time.Sleep(time.Second * 20)

}
