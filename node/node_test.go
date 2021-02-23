package node

import (
	"context"
	"net"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"

	"github.com/liyiysng/scatter/node/conn"
	"github.com/liyiysng/scatter/node/message"
	"github.com/liyiysng/scatter/node/session"
	"github.com/liyiysng/scatter/node/textlog"
)

func startServer(t *testing.T) (n *Node, err error) {

	sink, err := textlog.NewTempFileSink()
	if err != nil {
		return
	}

	n, err = NewNode(EnableTextLog(sink))
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
		SID:                 1,
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

	handShakeAck, err := msgConn.ReadNextMessage()
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

type ServiceTest struct {
}

func (srv *ServiceTest) Sum(ctx context.Context, session session.Session, req *TestReq) {

}

func TestNodeService(t *testing.T) {

	// start node
	sink, err := textlog.NewTempFileSink()
	if err != nil {
		t.Fatal(err)
	}

	n, err := NewNode(EnableTextLog(sink))
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		err = n.Serve(SocketProtcolTCP, "127.0.0.1:7788")
		if err != nil {
			myLog.Error(err)
			return
		}
		myLog.Info("node Serve finished")
	}()
}
