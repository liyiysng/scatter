package node

import (
	"context"
	"testing"

	"github.com/liyiysng/scatter/node/handle"
	"github.com/liyiysng/scatter/node/node_testing"
	"github.com/liyiysng/scatter/node/session"
)

type SrvTest struct {
	n *Node
}

func (s *SrvTest) FooTest(ctx context.Context, session session.Session, req *node_testing.FooReq) (res *node_testing.FooRes, err error) {
	myLog.Infof("[SrvTest.FooTest] %v", req)
	return &node_testing.FooRes{
		Res:        "---------------",
		PlayerUID:  112,
		PlayerUIDs: 3344,
	}, nil
}

func (s *SrvTest) FooErrTest(ctx context.Context, session session.Session, req *node_testing.FooReq) (res *node_testing.FooRes, err error) {
	myLog.Infof("[SrvTest.FooErrTest] %v", req)
	return nil, handle.NewCustomError("error info")
}

func (s *SrvTest) FooNotify(ctx context.Context, session session.Session, req *node_testing.FooNotiry) (err error) {
	myLog.Infof("[SrvTest.FooNotify] %v", req)
	return nil
}

func (s *SrvTest) MakeFooPush(ctx context.Context, session session.Session, req *node_testing.FooReq) (res *node_testing.FooRes, err error) {
	myLog.Infof("[SrvTest.MakeFooPush] %v", req)

	session.Push(context.Background(), "ClientFoo", &node_testing.FooPush{
		Push: "push data",
	})

	return &node_testing.FooRes{
		Res:        "---------------",
		PlayerUID:  112,
		PlayerUIDs: 3344,
	}, nil
}

func TestCSharpClient(t *testing.T) {
	//cfg := config.NewConfig()

	n, err := NewNode(
		1,
		NOptShowHandleLog(false),
		NOptTraceDetail(false),
		NOptCompress("gzip"),
		NOptEnableMetrics(false),
		NOptNodeName("test_node"),
		// NOptEnableEnableEsTextLog(
		// 	cfg.GetString("scatter.es.write_index"),
		// 	100, time.Second,
		// 	elastic.SetURL(cfg.GetStringSlice("scatter.es.url")...),
		// 	elastic.SetSniff(false),
		// 	//elastic.SetTraceLog(log.New(os.Stderr, "", log.LstdFlags)),
		// ),
	)
	if err != nil {
		t.Fatal(err)
	}

	defer n.Stop()

	n.Register(&SrvTest{})

	err = n.Serve(SocketProtcolWS, "127.0.0.1:4545")
	if err != nil {
		myLog.Error(err)
		return
	}
	myLog.Info("node Serve finished")

}
