package cluster

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/liyiysng/scatter/cluster/cluster_testing"
	"github.com/liyiysng/scatter/cluster/selector/policy"
	"github.com/liyiysng/scatter/cluster/subsrv"
	"github.com/liyiysng/scatter/cluster/subsrvpb"
	"google.golang.org/protobuf/proto"
)

type SubService0 struct {
}

func (s *SubService0) Foo1(ctx context.Context, session *subsrv.DummySession, req *cluster_testing.String) (res *cluster_testing.String, err error) {
	res = &cluster_testing.String{
		Str: req.Str,
	}
	return
}

type SubService1 struct {
}

func (s *SubService1) Foo1(ctx context.Context, session *subsrv.DummySession, req *cluster_testing.String) (res *cluster_testing.String, err error) {
	res = &cluster_testing.String{
		Str: req.Str,
	}
	return
}

func TestSubService(t *testing.T) {

	wg := sync.WaitGroup{}

	n1, err := createGrpcNode(&wg, "100", "127.0.0.1:1414", func(gn *GrpcNode) {
		gn.RegisterSubServiceName("foo", &SubService0{})
	})
	if err != nil {
		t.Fatal(err)
	}

	n2, err := createGrpcNode(&wg, "101", "127.0.0.1:1415", func(gn *GrpcNode) {
		gn.RegisterSubServiceName("foo1", &SubService1{})
	})
	if err != nil {
		t.Fatal(err)
	}

	n3, err := createGrpcNode(&wg, "102", "127.0.0.1:1469")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	var connBind interface{}

	ctx = policy.WithSessionAffinity(ctx, func(conn interface{}) {
		connBind = conn
	}, func() (conn interface{}) {
		return connBind
	})

	ctx = policy.WithConsistentHashID(ctx, "88576")

	myLog.Info("---------------------------------begin dial----------------------------------")

	client, err := n1.GetSubSrvClient("foo")
	if err != nil {
		t.Fatal(err)
	}

	client1, err := n1.GetSubSrvClient("foo1")
	if err != nil {
		t.Fatal(err)
	}

	req := &cluster_testing.String{
		Str: "----------------------",
	}

	payload, err := proto.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Call(ctx, &subsrvpb.CallReq{
		ServiceName: "foo",
		MethodName:  "Foo1",
		Payload:     payload,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)
	res, err = client1.Call(ctx, &subsrvpb.CallReq{
		ServiceName: "foo1",
		MethodName:  "Foo1",
		Payload:     payload,
	})

	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)

	time.Sleep(time.Second * 50)

	n1.Stop()
	n2.Stop()
	n3.Stop()

	wg.Wait()
}
