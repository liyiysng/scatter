// Package node 服务节点
package node

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/liyiysng/scatter/constants"
	"github.com/liyiysng/scatter/logger"
	"github.com/liyiysng/scatter/metrics"
	"github.com/liyiysng/scatter/node/acceptor"
	"github.com/liyiysng/scatter/node/conn"
	"github.com/liyiysng/scatter/util"
	"golang.org/x/net/trace"
)

// SocketProtcol 协议类型
type SocketProtcol string

const (
	// SocketProtcolTCP tcp 协议
	SocketProtcolTCP SocketProtcol = "tcp"
	// SocketProtcolWS websocket 协议
	SocketProtcolWS SocketProtcol = "ws"
)

// 内部错误
type errStore struct {
	err error
}

// Node represent
type Node struct {
	sync.RWMutex
	// session ID
	sIDSequence int64
	// 选项
	opts Options
	// 当前错误(最后一个错误)
	errValue atomic.Value
	// 开启时间
	startTime time.Time

	// 所有链接
	conns map[conn.MsgConn]bool

	// 监听对象
	accs map[acceptor.Acceptor]bool

	// 退出chan
	quit *util.Event
	// 检测所有内部goroutine
	waitGroup util.WaitGroupWrapper
	// trace
	trEvents trace.EventLog

	// 标识是否启动
	started bool
}

// NewNode 新建节点
func NewNode(opt ...IOption) (n *Node, err error) {

	opts := defaultOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	// 缺省日志
	if opts.Logger == nil {
		opts.Logger = logger.GDepthLogger
	}

	n = &Node{
		accs:      make(map[acceptor.Acceptor]bool),
		opts:      opts,
		startTime: time.Now(),
		conns:     make(map[conn.MsgConn]bool),
		quit:      util.NewEvent(),
	}

	if opts.enableEventTrace {
		_, file, line, _ := runtime.Caller(1)
		n.trEvents = trace.NewEventLog("scatter.Node", fmt.Sprintf("%s-%s:%d", opts.ID, file, line))
	}

	n.errValue.Store(&errStore{})

	opts.Logger.Infof("start node %s", opts.ID)

	return
}

// Serve 启动一个Serve
// Serve 除Stop或者 被调用之外,都返回一个非nil错误
// arg[0] = certfile
// arg[1] = keyfile
func (n *Node) Serve(sp SocketProtcol, addr string, cert ...string) error {

	n.waitGroup.Add(1)
	defer n.waitGroup.Done()

	n.Lock()
	n.trEventLogf("starting")
	n.started = true

	if n.accs == nil {
		// Start called after Stop
		n.Unlock()
		return constants.ErrNodeStopped
	}

	certFile := ""
	keyFile := ""

	if len(cert) > 0 {
		if len(cert) != 2 {
			n.Unlock()
			return constants.ErrInvalidCertificates
		}
		certFile = cert[0]
		keyFile = cert[1]
	}

	accOpt := acceptor.Option{
		Addr:     addr,
		CertFile: certFile,
		KeyFile:  keyFile,
		Logger:   n.opts.Logger,
		GetConnOpt: func() conn.MsgConnOption {
			ret := conn.MsgConnOption{
				SID:                 atomic.AddInt64(&n.sIDSequence, 1),
				MaxLength:           n.opts.maxPayloadLength,
				ReadTimeout:         n.opts.readTimeout,
				WriteTimeout:        n.opts.writeTimeout,
				ReadBufferSize:      n.opts.readBufferSize,
				WriteBufferSize:     n.opts.writeBufferSize,
				Compresser:          n.opts.getCompressor(),
				EnableLimit:         n.opts.enableLimit,
				RateLimitReadBytes:  n.opts.rateLimitReadBytes,
				RateLimitWriteBytes: n.opts.rateLimitWriteBytes,
			}

			// 读写字节数指标
			if n.opts.metricsReadWriteBytesCountEnabled() {
				ret.ReadCountReport = func(info conn.MsgConnInfo, byteCount int) {

				}
				ret.WriteCountReport = func(info conn.MsgConnInfo, byteCount int) {

				}
			}

			return ret
		},
	}

	var acc acceptor.Acceptor
	if sp == SocketProtcolTCP {
		acc = acceptor.NewTCPAcceptor(accOpt)
	} else if sp == SocketProtcolWS {
		acc = acceptor.NewWSAcceptor(accOpt)
	} else {
		n.Unlock()
		return fmt.Errorf("[Node.Serve] unsupport socket protcol %s", sp)
	}
	n.accs[acc] = true
	n.Unlock()

	defer func() {
		n.Lock()
		if n.accs != nil && n.accs[acc] {
			delete(n.accs, acc)
		}
		n.Unlock()
	}()

	n.Lock()
	n.trEventLogf("[Node.Serve] protcol:%s addr %s", sp, addr)
	n.Unlock()

	// accept go
	n.waitGroup.Wrap(
		func() {
			err := acc.ListenAndServe()
			n.Lock()
			n.trEventLogf("[Node.Serve] done serve protcol:%s addr %s", sp, addr)
			n.Unlock()
			if err != nil {
				n.opts.Logger.Errorf("[Node.Serve] node %v stop %v", n.opts.ID, err)
			}
		}, n.opts.Logger.Errorf)

	// 此处无需检查node是否退出
	// node 停止时,会关闭所有的acceptor
	for {
		select {
		case c, ok := <-acc.GetConnChan():
			{
				if ok {
					//handel conn go
					n.waitGroup.Wrap(func() {
						n.handleConn(c)
					}, n.opts.Logger.Errorf)
				} else {
					return nil
				}
			}
		}
	}
}

// Stop 停止/关闭该节点
func (n *Node) Stop() {
	n.quit.Fire()

	defer func() {
		n.waitGroup.Wait()
	}()

	n.Lock()
	accs := n.accs
	n.accs = nil
	conns := n.conns
	n.conns = nil
	n.Unlock()

	for acc := range accs {
		acc.Stop()
	}
	for c := range conns {
		c.Close()
	}

	n.Lock()
	if n.trEvents != nil {
		n.trEvents.Finish()
		n.trEvents = nil
	}
	n.Unlock()
}

// 记录事件日志
// 需求:已经持有锁
func (n *Node) trEventLogf(format string, a ...interface{}) {
	if n.trEvents != nil {
		n.trEvents.Printf(format, a...)
	}
}

// 记录事件错误日志
// 需求:已经持有锁
func (n *Node) trEventErrorf(format string, a ...interface{}) {
	if n.trEvents != nil {
		n.trEvents.Errorf(format, a...)
	}
}

func (n *Node) handleConn(conn conn.MsgConn) {
	if n.quit.HasFired() {
		conn.Close()
		return
	}

	n.Lock()
	n.trEventLogf("[Node.handleConn] new connection local:%v remote:%v", conn.LocalAddr(), conn.RemoteAddr())
	n.Unlock()

	//指标:链接数
	if n.opts.metricsConnCountEnabled() {
		metrics.ReportNodeNewCommingConn(n.opts.ID)
		defer metrics.ReportNodeNewCommingConn(n.opts.ID)
	}

	// 创建session
}
