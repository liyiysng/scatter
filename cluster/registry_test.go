package cluster

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/liyiysng/scatter/cluster/registry"
	"github.com/liyiysng/scatter/cluster/registry/consul"
	_ "github.com/liyiysng/scatter/cluster/registry/consul"
	"github.com/liyiysng/scatter/util"
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
					Address: "127.0.0.1:1155",
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
					Address: "127.0.0.1:1155",
				},
			},
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
					ID:      "220",
					Address: "127.0.0.1:1154",
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
					ID:      "330",
					Address: "127.0.0.1:1158",
				},
			},
		},
		consul.WithRegistryCompress(),
	)

	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 2)

}
