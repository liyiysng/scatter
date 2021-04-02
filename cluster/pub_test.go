package cluster

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/liyiysng/scatter/cluster/cluster_testing"
	"github.com/liyiysng/scatter/cluster/session"
)

type Subscribe1 struct {
}

func (s *Subscribe1) Foo1(ctx context.Context, session session.ISession, data *cluster_testing.String) (err error) {
	myLog.Infof("[Subscribe1.Foo1] %v", data)
	return
}

func TestPub(t *testing.T) {

	wg := sync.WaitGroup{}

	n1, err := createGrpcNode(&wg, "100", "127.0.0.1:1414", func(gn *GrpcNode) {
		if serr := gn.Subscribe("foo", &Subscribe1{}); serr != nil {
			t.Fatal(serr)
		}
	})
	if err != nil {
		t.Fatal(err)
	}

	n2, err := createGrpcNode(&wg, "101", "127.0.0.1:1415", func(gn *GrpcNode) {
		if serr := gn.Subscribe("foo", &Subscribe1{}); serr != nil {
			t.Fatal(serr)
		}
	})
	if err != nil {
		t.Fatal(err)
	}

	n3, err := createGrpcNode(&wg, "102", "127.0.0.1:1469")
	if err != nil {
		t.Fatal(err)
	}

	myLog.Info("---------------------------------begin dial----------------------------------")

	pub := NewPublisher()

	go func() {
		for {
			err = pub.Publish(context.Background(), "foo", "Foo1", &cluster_testing.String{
				Str: "pub test",
			})
			if err != nil {
				myLog.Error(err)
				return
			}
			time.Sleep(time.Second * 1)
		}
	}()

	time.Sleep(time.Second * 50)

	n1.Stop()
	n2.Stop()
	n3.Stop()
	pub.Close()

	wg.Wait()
}
