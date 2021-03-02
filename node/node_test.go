package node

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/liyiysng/scatter/config"
	"github.com/liyiysng/scatter/metrics"
	"github.com/liyiysng/scatter/node/conn"
	"github.com/liyiysng/scatter/node/message"
	"github.com/liyiysng/scatter/node/node_testing"
	"github.com/liyiysng/scatter/node/session"
	"github.com/liyiysng/scatter/util"
	"github.com/olivere/elastic/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
)

func startServer(t *testing.T) (n *Node, err error) {
	n, err = NewNode(1, NOptEnableEnableFileTextLog())
	if err != nil {
		return
	}
	go func() {
		err = n.Serve(SocketProtcolTCP, "127.0.0.1:7788")
		if err != nil {
			myLog.Error(err)
			return
		}
		myLog.Info("node Serve finished")
	}()

	return
}

func TestNodeClose(t *testing.T) {
	n, err := startServer(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	time.Sleep(time.Second * 4)

	n.Stop()

}

func TestNodeConnect(t *testing.T) {
	n, err := startServer(t)
	if err != nil {
		t.Fatal(err)
		return
	}

	defer n.Stop()

	time.Sleep(time.Second)

	_, err = net.Dial("tcp", "127.0.0.1:7788")
	if err != nil {
		t.Fatal(err)
		return
	}
}

func TestNodeHandShake(t *testing.T) {
	n, err := startServer(t)
	if err != nil {
		t.Fatal(err)
		return
	}

	go func() {
		err = http.ListenAndServe("127.0.0.1:8081", nil)
		if err != nil {
			myLog.Error(err)
			return
		}
	}()

	defer n.Stop()

	time.Sleep(time.Second)

	con, err := net.Dial("tcp", "127.0.0.1:7788")
	if err != nil {
		t.Fatal(err)
		return
	}

	msgConn := conn.NewTCPMsgConn(con, conn.MsgConnOption{
		MaxLength:           1024 * 1024,
		ReadTimeout:         time.Second * 2,
		WriteTimeout:        time.Second * 2,
		ReadBufferSize:      1024 * 1024 * 4,
		WriteBufferSize:     1024 * 1024 * 4,
		Compresser:          n.opts.getCompressor(),
		EnableLimit:         n.opts.enableLimit,
		RateLimitReadBytes:  n.opts.rateLimitReadBytes,
		RateLimitWriteBytes: n.opts.rateLimitWriteBytes,
	})
	defer msgConn.Close()

	handShake, err := message.MsgFactory.BuildHandShakeMessage("win64", "0.0.1", "0.1.0")
	if err != nil {
		t.Fatal(err)
		return
	}

	err = msgConn.WriteNextMessage(handShake, 0)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = msgConn.Flush()
	if err != nil {
		t.Fatal(err)
		return
	}
	// go func() {
	// 	handShakeAck, err := msgConn.ReadNextMessage()
	// 	if err != nil {
	// 		myLog.Error(err)
	// 		return
	// 	}
	// 	myLog.Info(handShakeAck)

	// 	err = msgConn.Close()
	// 	if err != nil {
	// 		myLog.Error(err)
	// 		return
	// 	}
	// }()

	time.Sleep(time.Second)

	n.Stop()

	handShakeAck, _, err := msgConn.ReadNextMessage()
	if err != nil {
		myLog.Error(err)
		return
	}
	myLog.Info(handShakeAck)

	err = msgConn.Close()
	if err != nil {
		myLog.Error(err)
		return
	}

	time.Sleep(time.Second * 2)
}

const (
	nodeBindAddr = "127.0.0.1:7788"
)

var callSequence int32 = 0

type ServiceTest struct {
}

func (srv *ServiceTest) Sum(ctx context.Context, session session.Session, req *node_testing.SumReq) (res *node_testing.SumRes, err error) {

	res = &node_testing.SumRes{
		Sum: req.LOP + req.ROP,
	}

	return
}

func createNode() (n *Node, err error) {
	// // start node
	// sink, err := textlog.NewTempFileSink()
	// if err != nil {
	// 	return nil, err
	// }

	// n, err = NewNode(NOptEnableTextLog(sink), NOptShowHandleLog(false))
	n, err = NewNode(1, NOptShowHandleLog(false), NOptTraceDetail(false))
	if err != nil {
		return nil, err
	}
	return
}

func srvNode(n *Node) {
	go func() {
		err := n.Serve(SocketProtcolTCP, nodeBindAddr)
		if err != nil {
			myLog.Error(err)
			return
		}
		myLog.Info("node Serve finished")
	}()
}

func createClient(n *Node) (c conn.MsgConn, err error) {
	con, err := net.Dial("tcp", nodeBindAddr)
	if err != nil {
		return nil, err
	}

	msgConn := conn.NewTCPMsgConn(con, conn.MsgConnOption{
		MaxLength:           1024 * 1024,
		ReadTimeout:         time.Second * 20,
		WriteTimeout:        time.Second * 2,
		ReadBufferSize:      1024 * 1024 * 4,
		WriteBufferSize:     1024 * 1024 * 4,
		Compresser:          n.opts.getCompressor(),
		EnableLimit:         n.opts.enableLimit,
		RateLimitReadBytes:  n.opts.rateLimitReadBytes,
		RateLimitWriteBytes: n.opts.rateLimitWriteBytes,
	})

	// 发送握手
	handShake, err := message.MsgFactory.BuildHandShakeMessage("win64", "0.0.1", "0.1.0")
	if err != nil {
		return nil, err
	}
	err = msgConn.WriteNextMessage(handShake, 0)
	if err != nil {
		return nil, err
	}
	err = msgConn.Flush()
	if err != nil {
		return nil, err
	}

	// 等待握手完成
	_, _, err = msgConn.ReadNextMessage()
	if err != nil {
		return nil, err
	}

	return msgConn, nil
}

type dummyLock struct{}

func (l *dummyLock) Lock()   {}
func (l *dummyLock) Unlock() {}

func call(c conn.MsgConn, wl sync.Locker, rl sync.Locker, wg *sync.WaitGroup, codec encoding.Codec, srv string, req interface{}, res interface{}) error {

	wg.Add(1)
	defer wg.Done()

	sequence := atomic.AddInt32(&callSequence, 1)

	buf, err := codec.Marshal(req)
	if err != nil {
		return err
	}

	msg, err := message.MsgFactory.BuildRequestMessage(sequence, srv, buf)
	if err != nil {
		return err
	}

	wl.Lock()
	err = c.WriteNextMessage(msg, 0)
	if err != nil {
		wl.Unlock()
		return err
	}

	err = c.Flush()
	if err != nil {
		wl.Unlock()
		return err
	}
	wl.Unlock()

	rl.Lock()
	recvMsg, _, err := c.ReadNextMessage()
	rl.Unlock()
	if err != nil {
		return err
	}

	err = codec.Unmarshal(recvMsg.GetPayload(), res)
	if err != nil {
		return err
	}

	return nil
}

func TestNodeService(t *testing.T) {

	n, err := createNode()
	if err != nil {
		t.Fatal(err)
	}
	defer n.Stop()

	n.Register(&ServiceTest{})

	srvNode(n)

	// 确保node启动
	time.Sleep(time.Second)

	// 创建客户端
	client, err := createClient(n)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	res := &node_testing.SumRes{}

	wg := &sync.WaitGroup{}

	err = call(
		client,
		&dummyLock{},
		&dummyLock{},
		wg,
		n.opts.getCodec(),
		"ServiceTest.Sum",
		&node_testing.SumReq{
			LOP: 10,
			ROP: 10,
		},
		res,
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)

	times := 10000

	beg := time.Now()

	for i := 0; i < times; i++ {
		err = call(
			client,
			&dummyLock{},
			&dummyLock{},
			wg,
			n.opts.getCodec(),
			"ServiceTest.Sum",
			&node_testing.SumReq{
				LOP: 10,
				ROP: 10,
			},
			res,
		)
		if err != nil {
			t.Fatal(err)
		}
	}
	wg.Wait()
	t.Logf("rpc times %d execution time %v , %v/op", times, time.Now().Sub(beg), time.Now().Sub(beg)/time.Duration(times))

}

func BenchmarkRPCSingleClient(b *testing.B) {

	n, err := createNode()
	if err != nil {
		b.Fatal(err)
	}
	defer n.Stop()

	n.Register(&ServiceTest{})

	srvNode(n)

	// 确保node启动
	time.Sleep(time.Second)

	// 创建客户端
	client, err := createClient(n)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	wwg := sync.WaitGroup{}
	wl := &sync.Mutex{}
	rl := &sync.Mutex{}
	b.RunParallel(func(tb *testing.PB) {
		for tb.Next() {
			wwg.Add(1)
			defer wwg.Done()

			res := &node_testing.SumRes{}

			err = call(
				client,
				wl,
				rl,
				&wwg,
				n.opts.getCodec(),
				"ServiceTest.Sum",
				&node_testing.SumReq{
					LOP: 10,
					ROP: 10,
				},
				res,
			)
			if err != nil {
				b.Fatal(err)
				return
			}
		}
	})

	wwg.Wait()
}

type ncall struct {
	res  interface{}
	done chan struct{}
}

type nodeClient struct {
	c     conn.MsgConn
	codec encoding.Codec
	wl    sync.Mutex
	rl    sync.Mutex
	wg    sync.WaitGroup
	seq   int32

	closeEvent *util.Event

	mu          sync.Mutex
	pendingCall map[int32] /*sequence*/ *ncall
}

func (c *nodeClient) runRead() {

	for {

		if c.closeEvent.HasFired() {
			return
		}

		err := c.readOne()
		if err != nil {
			myLog.Error(err)
			return
		}
	}

}

func (c *nodeClient) readOne() error {

	c.rl.Lock()
	recvMsg, _, err := c.c.ReadNextMessage()
	c.rl.Unlock()
	if err != nil {
		return err
	}

	var res interface{} = nil

	c.mu.Lock()
	if call, ok := c.pendingCall[recvMsg.GetSequence()]; ok {
		res = call.res
		delete(c.pendingCall, recvMsg.GetSequence())
		if call.done != nil {
			close(call.done)
		}

		c.wg.Done()
	} else {
		c.mu.Unlock()
		return fmt.Errorf("sequence %d not found , current sequence %d", recvMsg.GetSequence(), c.seq)
	}
	c.mu.Unlock()

	err = c.codec.Unmarshal(recvMsg.GetPayload(), res)
	if err != nil {
		return err
	}

	return nil
}

func (c *nodeClient) call(srv string, req interface{}, res interface{}, popt message.PacketOpt) (done <-chan struct{}, err error) {

	callDone := make(chan struct{})

	sequence := atomic.AddInt32(&c.seq, 1)

	buf, err := c.codec.Marshal(req)
	if err != nil {
		return nil, err
	}

	msg, err := message.MsgFactory.BuildRequestMessage(sequence, srv, buf)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	if c.pendingCall == nil {
		c.pendingCall = map[int32]*ncall{}
	}

	c.pendingCall[sequence] = &ncall{
		res:  res,
		done: callDone,
	}
	c.wg.Add(1)
	c.mu.Unlock()

	defer func() {
		if err != nil {
			c.mu.Lock()
			delete(c.pendingCall, sequence)
			c.wg.Done()
			c.mu.Unlock()
		}
	}()

	err = c.wirteOne(msg, popt)
	if err != nil {
		return nil, err
	}

	return callDone, nil
}

func (c *nodeClient) wirteOne(msg message.Message, popt message.PacketOpt) error {
	c.wl.Lock()
	defer c.wl.Unlock()

	err := c.c.WriteNextMessage(msg, popt)
	if err != nil {
		return err
	}

	err = c.c.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (c *nodeClient) close() {
	c.wg.Wait()
	c.closeEvent.Fire()
	c.c.Close()
}

func BenchmarkRPC(b *testing.B) {

	n, err := createNode()
	if err != nil {
		b.Fatal(err)
	}
	defer n.Stop()

	n.Register(&ServiceTest{})

	srvNode(n)

	// 确保node启动
	time.Sleep(time.Second)

	clientCount := 1000

	clients := make([]*nodeClient, clientCount)

	for i := 0; i < clientCount; i++ {
		// 创建客户端
		client, err := createClient(n)
		if err != nil {
			b.Fatal(err)
		}
		clients[i] = &nodeClient{
			codec:      n.opts.getCodec(),
			c:          client,
			closeEvent: util.NewEvent(),
		}
		defer clients[i].close()
		go clients[i].runRead()
	}

	b.ResetTimer()
	wwg := sync.WaitGroup{}

	count := int32(0)

	b.RunParallel(func(tb *testing.PB) {
		for tb.Next() {
			wwg.Add(1)
			defer wwg.Done()

			res := &node_testing.SumRes{}

			cur := atomic.AddInt32(&count, 1)

			cl := clients[int(cur)%clientCount]

			_, err = cl.call(
				"ServiceTest.Sum",
				&node_testing.SumReq{
					LOP: 10,
					ROP: 10,
				},
				res,
				message.DEFAULTPOPT,
			)
			if err != nil {
				b.Fatal(err)
				return
			}
		}
	})

	wwg.Wait()
}

type grpcTestImp struct {
	node_testing.UnimplementedGrpcSrvTestServer
}

func (srv *grpcTestImp) GetSum(ctx context.Context, req *node_testing.GRPCSumReq) (*node_testing.GRPCSumRes, error) {
	return &node_testing.GRPCSumRes{
		Sum: req.LOP + req.ROP,
	}, nil
}

func BenchmarkGRPC(b *testing.B) {
	s := grpc.NewServer()
	node_testing.RegisterGrpcSrvTestServer(s, &grpcTestImp{})

	lis, err := net.Listen("tcp", "127.0.0.1:8989")
	if err != nil {
		b.Fatal(err)
		return
	}

	go func() {
		err = s.Serve(lis)
		if err != nil {
			myLog.Error(err)
			return
		}
	}()

	// 确保服务器启动
	time.Sleep(time.Second)
	defer s.Stop()

	clientCount := 100

	clients := make([]node_testing.GrpcSrvTestClient, clientCount)

	for i := 0; i < clientCount; i++ {
		// 创建客户端
		conn, err := grpc.Dial("127.0.0.1:8989", grpc.WithInsecure())
		if err != nil {
			b.Fatal(err)
		}
		defer conn.Close()

		clients[i] = node_testing.NewGrpcSrvTestClient(conn)
	}

	b.ResetTimer()

	wwg := sync.WaitGroup{}

	count := int32(0)

	b.RunParallel(func(tb *testing.PB) {
		for tb.Next() {
			wwg.Add(1)
			defer wwg.Done()

			cur := atomic.AddInt32(&count, 1)

			client := clients[int(cur)%clientCount]

			_, err := client.GetSum(context.Background(),
				&node_testing.GRPCSumReq{
					LOP: 10,
					ROP: 10,
				},
			)

			if err != nil {
				b.Fatal(err)
				return
			}
		}
	})

	wwg.Wait()

}

// 消息选项测试
func TestNodeMsgOpt(t *testing.T) {

	n, err := NewNode(1, NOptShowHandleLog(false), NOptTraceDetail(false), NOptCompress("gzip"))
	if err != nil {
		t.Fatal(err)
	}

	defer n.Stop()

	n.Register(&ServiceTest{})

	srvNode(n)

	// 确保node启动
	time.Sleep(time.Second)

	// 创建客户端
	conn, err := createClient(n)
	if err != nil {
		t.Fatal(err)
	}
	client := &nodeClient{
		codec:      n.opts.getCodec(),
		c:          conn,
		closeEvent: util.NewEvent(),
	}
	defer client.close()
	go client.runRead()

	for i := 0; i < 10; i++ {
		res := &node_testing.SumRes{}

		done, err := client.call(
			"ServiceTest.Sum",
			&node_testing.SumReq{
				LOP: 10,
				ROP: 10,
				DataStr: `xxxxxxxxxxxxxxxffffffffffffffffffffffffffffff--------------------------
				--------------ssssssssssssssssssssssssssssssabfoabhfoasbfobaogfbhodhgosd
				hgoidshoghdshgoidshgoidshogihsdihgosdhgosdhgoihsdighoisdhgoishdgoihsdigh
				ohfoiahfohaofhoashfoashfoashfoahsofhaoshfoas`,
			},
			res,
			message.COMPRESS,
		)

		if err != nil {
			t.Fatal(err)
		}

		<-done

		t.Log(res)
	}

}

type ServiceMetricsTest struct {
}

func (srv *ServiceMetricsTest) Foo1(ctx context.Context, session session.Session, req *node_testing.SumReq) (res *node_testing.SumRes, err error) {

	res = &node_testing.SumRes{
		Sum: req.LOP + req.ROP,
	}

	time.Sleep(time.Millisecond * time.Duration(rand.Int31n(5)))

	return
}

func (srv *ServiceMetricsTest) Foo2(ctx context.Context, session session.Session, req *node_testing.SumReq) (res *node_testing.SumRes, err error) {

	res = &node_testing.SumRes{
		Sum: req.LOP + req.ROP,
	}

	time.Sleep(time.Millisecond * time.Duration(rand.Int31n(10)))

	return
}

func (srv *ServiceMetricsTest) Foo3(ctx context.Context, session session.Session, req *node_testing.SumReq) (res *node_testing.SumRes, err error) {

	res = &node_testing.SumRes{
		Sum: req.LOP + req.ROP,
	}

	time.Sleep(time.Millisecond * time.Duration(rand.Int31n(15)))

	return
}

func (srv *ServiceMetricsTest) Foo4(ctx context.Context, session session.Session, req *node_testing.SumReq) (res *node_testing.SumRes, err error) {

	res = &node_testing.SumRes{
		Sum: req.LOP + req.ROP,
	}

	time.Sleep(time.Millisecond * time.Duration(rand.Int31n(20)))

	return
}

// 指标测试
func TestNodeMetrics(t *testing.T) {
	reporter, err := metrics.NewPrometheusReporter("node_metrics", config.NewConfig(), map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	defer reporter.Close()

	go func() {
		for {
			metrics.ReportSysMetrics([]metrics.Reporter{reporter}, time.Second)
		}
	}()

	n, err := NewNode(
		1,
		NOptShowHandleLog(false),
		NOptTraceDetail(false),
		NOptCompress("gzip"),
		NOptMetricsReporter(reporter),
		NOptEnableMetrics(true),
		NOptNodeName("test_node"),
	)
	if err != nil {
		t.Fatal(err)
	}

	defer n.Stop()

	n.Register(&ServiceMetricsTest{})

	srvNode(n)

	// 确保node启动
	time.Sleep(time.Second)

	wwg := sync.WaitGroup{}

	srvs := []string{
		"ServiceMetricsTest.Foo1",
		"ServiceMetricsTest.Foo2",
		"ServiceMetricsTest.Foo3",
		"ServiceMetricsTest.Foo4",
	}

	finish := util.NewEvent()

	go func() {
		time.Sleep(time.Second * 400)
		finish.Fire()
	}()

	for i := 0; i < 10; i++ {

		go func() {

			wwg.Add(1)
			defer wwg.Done()

			// 创建客户端
			conn, err := createClient(n)
			if err != nil {
				myLog.Error(err)
				return
			}
			client := &nodeClient{
				codec:      n.opts.getCodec(),
				c:          conn,
				closeEvent: util.NewEvent(),
			}
			defer client.close()
			go client.runRead()

			for {

				if finish.HasFired() {
					return
				}

				res := &node_testing.SumRes{}

				done, err := client.call(
					srvs[rand.Intn(len(srvs))],
					&node_testing.SumReq{
						LOP:     10,
						ROP:     10,
						DataStr: `xxxxxxxxxxxxxxx`,
					},
					res,
					message.COMPRESS,
				)

				if err != nil {
					myLog.Error(err)
					return
				}

				<-done
				time.Sleep(time.Second * 10)

				t.Log(res)
			}
		}()

	}

	wwg.Wait()

}

func TestNodeEsSink(t *testing.T) {

	cfg := config.NewConfig()

	n, err := NewNode(
		1,
		NOptShowHandleLog(false),
		NOptTraceDetail(false),
		NOptCompress("gzip"),
		NOptEnableMetrics(true),
		NOptNodeName("test_node"),
		NOptEnableEnableEsTextLog(
			cfg.GetString("scatter.es.write_index"),
			100, time.Second,
			elastic.SetURL(cfg.GetStringSlice("scatter.es.url")...),
			elastic.SetSniff(false),
			//elastic.SetTraceLog(log.New(os.Stderr, "", log.LstdFlags))
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	defer n.Stop()

	n.Register(&ServiceMetricsTest{})

	srvNode(n)

	// 确保node启动
	time.Sleep(time.Second)

	wwg := sync.WaitGroup{}

	srvs := []string{
		"ServiceMetricsTest.Foo1",
		"ServiceMetricsTest.Foo2",
		"ServiceMetricsTest.Foo3",
		"ServiceMetricsTest.Foo4",
	}

	finish := util.NewEvent()

	go func() {
		time.Sleep(time.Second * 400)
		finish.Fire()
	}()

	for i := 0; i < 10; i++ {

		go func() {

			wwg.Add(1)
			defer wwg.Done()

			// 创建客户端
			conn, err := createClient(n)
			if err != nil {
				myLog.Error(err)
				return
			}
			client := &nodeClient{
				codec:      n.opts.getCodec(),
				c:          conn,
				closeEvent: util.NewEvent(),
			}
			defer client.close()
			go client.runRead()

			for {

				if finish.HasFired() {
					return
				}

				res := &node_testing.SumRes{}

				done, err := client.call(
					srvs[rand.Intn(len(srvs))],
					&node_testing.SumReq{
						LOP:     10,
						ROP:     10,
						DataStr: `xxxxxxxxxxxxxxx`,
					},
					res,
					message.COMPRESS,
				)

				if err != nil {
					myLog.Error(err)
					return
				}

				<-done
				time.Sleep(time.Second * 1)

				t.Log(res)
			}
		}()

	}

	wwg.Wait()

}
