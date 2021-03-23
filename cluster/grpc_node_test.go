package cluster

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/liyiysng/scatter/cluster/cluster_testing"
	"github.com/liyiysng/scatter/cluster/selector/policy"
	"google.golang.org/grpc"

	"net/http"
	_ "net/http/pprof"
)

func TestGrpcNode(t *testing.T) {

	go func() {
		http.ListenAndServe("127.0.0.1:11112", nil)
	}()

	n1 := NewGrpcNode("110")
	n2 := NewGrpcNode("111")

	lis1, err := net.Listen("tcp", srvAddr[0])
	if err != nil {
		t.Fatal(err)
	}

	lis2, err := net.Listen("tcp", srvAddr[1])
	if err != nil {
		t.Fatal(err)
	}

	cluster_testing.RegisterSrvStringsServer(n1, &srvStringsImp{id: "1"})
	cluster_testing.RegisterSrvStringsServer(n2, &srvStringsImp{id: "2"})
	cluster_testing.RegisterSrvIntsServer(n1, &srvIntsImp{id: "1"})
	cluster_testing.RegisterSrvIntsServer(n2, &srvIntsImp{id: "2"})

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer n1.Stop()
		err := n1.Serve(lis1)
		if err != nil {
			myLog.Error(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer n2.Stop()
		err := n2.Serve(lis2)
		if err != nil {
			myLog.Error(err)
		}
	}()

	client, err := grpc.Dial("scatter://110/scatter.cluster.SrvStrings", grpc.WithInsecure(), grpc.WithDefaultServiceConfig(
		`{
			"loadBalancingConfig":[ { "session_affinity": {} } ]
		}
		`,
	), grpc.WithBlock())

	stringsClient := cluster_testing.NewSrvStringsClient(client)

	ctx := context.Background()

	var connBind interface{}

	ctx = policy.WithSessionAffinity(ctx, func(conn interface{}) {
		connBind = conn
	}, func() (conn interface{}) {
		return connBind
	})

	go func() {
		time.Sleep(time.Second * 10)
		lis1.Close()
	}()

	for i := 0; i < 50; i++ {
		res, err := stringsClient.ToUpper(ctx, &cluster_testing.String{
			Str: "hello world.",
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(res)
		time.Sleep(time.Second * 1)
	}

	time.Sleep(time.Second * 10)

	lis1.Close()
	lis2.Close()

	wg.Wait()
}

func createGrpcNode(wg *sync.WaitGroup, id, addr string, cb ...func(gn *GrpcNode)) (*GrpcNode, error) {
	n := NewGrpcNode(id)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	cluster_testing.RegisterSrvStringsServer(n, &srvStringsImp{id: id})
	cluster_testing.RegisterSrvIntsServer(n, &srvIntsImp{id: id})

	for _, v := range cb {
		v(n)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := n.Serve(lis)
		if err != nil {
			myLog.Error(err)
		}
	}()

	return n, nil
}

func TestGrpcNodeClient(t *testing.T) {

	wg := sync.WaitGroup{}

	n1, err := createGrpcNode(&wg, "100", "127.0.0.1:1414")
	if err != nil {
		t.Fatal(err)
	}

	n2, err := createGrpcNode(&wg, "101", "127.0.0.1:1415")
	if err != nil {
		t.Fatal(err)
	}

	n3, err := createGrpcNode(&wg, "102", "127.0.0.1:1469")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		time.Sleep(time.Second * 10)
		n2.Stop()
	}()

	ctx := context.Background()

	var connBind interface{}

	ctx = policy.WithSessionAffinity(ctx, func(conn interface{}) {
		connBind = conn
	}, func() (conn interface{}) {
		return connBind
	})

	ctx = policy.WithConsistentHashID(ctx, "88576")

	myLog.Info("---------------------------------begin dial----------------------------------")

	client, err := n1.GetClient("scatter.service.SrvStrings")
	if err != nil {
		t.Fatal(err)
	}

	iClient, err := n1.GetClient("scatter.service.SrvInts")
	if err != nil {
		t.Fatal(err)
	}

	stringsClient := cluster_testing.NewSrvStringsClient(client)
	intsClient := cluster_testing.NewSrvIntsClient(iClient)

	for i := 0; i < 50; i++ {
		res, err := stringsClient.ToUpper(ctx, &cluster_testing.String{
			Str: "hello world.",
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(res)
		time.Sleep(time.Second * 1)

		sres, err := intsClient.Sum(ctx, &cluster_testing.Ints{
			I: []int64{10, 32, 55},
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(sres)
	}

	time.Sleep(time.Second * 10)

	n1.Stop()
	n2.Stop()
	n3.Stop()

	wg.Wait()
}
