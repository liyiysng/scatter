// Package node 服务节点
package node

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/liyiysng/scatter/constants"
	"github.com/liyiysng/scatter/logger"
	"github.com/liyiysng/scatter/metrics"
	"github.com/liyiysng/scatter/node/acceptor"
	"github.com/liyiysng/scatter/node/conn"
	"github.com/liyiysng/scatter/node/handle"
	"github.com/liyiysng/scatter/node/session"
	"github.com/liyiysng/scatter/util"
	"golang.org/x/net/trace"
)

var (
	myLog = logger.Component("node")
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
	mu sync.RWMutex
	// session ID
	idGen *util.Node
	// 选项
	opts Options
	// 开启时间
	startTime time.Time

	// 所有链接
	sessions map[int64]session.FrontendSession

	// 监听对象
	accs map[acceptor.Acceptor]bool

	// 退出chan
	quit *util.Event
	// 检测所有内部goroutine
	waitGroup util.WaitGroupWrapper
	// trace
	trEvents trace.EventLog
	// 处理
	srvHandle handle.IHandler

	// 标识是否服务
	serve bool
}

// NewNode 新建节点
func NewNode(nid int64, opt ...IOption) (n *Node, err error) {

	opts := defaultOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	if opts.lastError != nil {
		return nil, opts.lastError
	}

	if err = opts.validate(); err != nil {
		return nil, err
	}

	opts.ID = nid

	// 缺省日志
	if opts.Logger == nil {
		if opts.LogPrefix != "" {
			opts.Logger = logger.NewPrefixLogger(logger.GDepthLogger, opts.LogPrefix)
		} else {
			opts.Logger = logger.GDepthLogger
		}
	}

	if opts.Name == "" {
		opts.Name = strconv.FormatInt(nid, 10)
	}

	n = &Node{
		accs:      make(map[acceptor.Acceptor]bool),
		opts:      opts,
		startTime: time.Now(),
		sessions:  make(map[int64]session.FrontendSession),
		quit:      util.NewEvent(),
	}

	n.idGen, err = util.NewNode(nid)
	if err != nil {
		return nil, err
	}

	n.srvHandle = handle.NewServiceHandle(&handle.Option{
		Codec:            n.opts.getCodec(),
		ReqTypeValidator: n.opts.reqTypeValidator,
		ResTypeValidator: n.opts.resTypeValidator,
		SessionType:      reflect.TypeOf((*session.Session)(nil)).Elem(),
		HookCall:         n.onCall,
		HookNofify:       n.onNotify,
	})

	if opts.enableEventTrace {
		_, file, line, _ := runtime.Caller(1)
		n.trEvents = trace.NewEventLog("scatter.Node", fmt.Sprintf("%d-%s:%d", opts.ID, file, line))
	}

	opts.Logger.Infof("start node %d", opts.ID)

	return
}

// Register 注册服务
func (n *Node) Register(recv interface{}) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.accs == nil {
		return constants.ErrNodeStopped
	}

	if n.serve {
		return fmt.Errorf("[Node.Register] register service after Node.Serve")
	}

	return n.srvHandle.Register(recv)
}

// RegisterName 注册命名服务
func (n *Node) RegisterName(name string, recv interface{}) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.accs == nil {
		return constants.ErrNodeStopped
	}

	if n.serve {
		return fmt.Errorf("[Node.Register] register service after Node.Serve %q", name)
	}

	return n.srvHandle.RegisterName(name, recv)
}

// Serve 启动一个Serve
// Serve 除Stop或者 被调用之外,都返回一个非nil错误
// arg[0] = certfile
// arg[1] = keyfile
func (n *Node) Serve(sp SocketProtcol, addr string, cert ...string) error {

	n.waitGroup.Add(1)
	defer n.waitGroup.Done()

	n.mu.Lock()
	n.trEventLogf("starting")
	n.serve = true

	if n.accs == nil {
		// Start called after Stop
		n.mu.Unlock()
		return constants.ErrNodeStopped
	}

	certFile := ""
	keyFile := ""

	if len(cert) > 0 {
		if len(cert) != 2 {
			n.mu.Unlock()
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
				SID:                 n.idGen.Generate().Int64(),
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
					metrics.ReportNodeReadBytes(n.opts.metricsReporters, n.opts.ID, n.opts.Name, byteCount)
				}
				ret.WriteCountReport = func(info conn.MsgConnInfo, byteCount int) {
					metrics.ReportNodeWriteBytes(n.opts.metricsReporters, n.opts.ID, n.opts.Name, byteCount)
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
		n.mu.Unlock()
		return fmt.Errorf("[Node.Serve] unsupport socket protcol %s", sp)
	}
	n.accs[acc] = true
	n.mu.Unlock()

	defer func() {
		n.mu.Lock()
		if n.accs != nil && n.accs[acc] {
			delete(n.accs, acc)
		}
		n.mu.Unlock()
	}()

	n.mu.Lock()
	n.trEventLogf("[Node.Serve] protcol:%s addr %s", sp, addr)
	n.mu.Unlock()

	// accept go
	n.waitGroup.Wrap(
		func() {
			err := acc.ListenAndServe()
			n.mu.Lock()
			n.trEventLogf("[Node.Serve] done serve protcol:%s addr %s", sp, addr)
			n.mu.Unlock()
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
					n.waitGroup.Wrap(func() { n.handleConn(c) }, n.opts.Logger.Errorf)
				} else {
					return nil
				}
			}
		}
	}
}

// Stop 停止/关闭该节点
func (n *Node) Stop() {

	if n.quit.HasFired() {
		return
	}

	n.opts.Logger.Infof("stopping node %d", n.opts.ID)

	n.quit.Fire()

	defer func() {
		n.waitGroup.Wait()
		if n.opts.textLogWriter != nil && len(n.opts.textLogWriter) > 0 {
			for _, v := range n.opts.textLogWriter {
				err := v.Close()
				if err != nil {
					n.opts.Logger.Errorf("close text log failed %v", err)
				}
			}

		}
		n.opts.Logger.Infof("node %d stoped", n.opts.ID)
	}()

	n.mu.Lock()
	accs := n.accs
	n.accs = nil
	ss := n.sessions
	n.sessions = nil
	n.mu.Unlock()

	for acc := range accs {
		acc.Stop()
	}
	for _, c := range ss {
		c.Close()
	}

	n.mu.Lock()
	if n.trEvents != nil {
		n.trEvents.Finish()
		n.trEvents = nil
	}
	n.mu.Unlock()
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

	n.mu.Lock()
	n.trEventLogf("[Node.handleConn] new connection local:%v remote:%v", conn.LocalAddr(), conn.RemoteAddr())
	n.mu.Unlock()

	//指标:链接数
	if n.opts.metricsConnCountEnabled() {
		metrics.ReportNodeConnectionInc(n.opts.metricsReporters, n.opts.ID, n.opts.Name)
		defer metrics.ReportNodeConnectionDec(n.opts.metricsReporters, n.opts.ID, n.opts.Name)
	}

	// 创建session
	s := session.NewFrontendSession(n.opts.ID, conn, &session.Option{
		Logger:            n.opts.Logger,
		ConnectTimeout:    n.opts.connectionTimeout,
		ReadChanSize:      n.opts.readChanBufSize,
		WriteChanSize:     n.opts.writeChanBufSize,
		Codec:             n.opts.getCodec(),
		RateLimitMsgProc:  n.opts.rateLimitMsgProcNum,
		PushInterceptor:   nil,
		OnMsgFinish:       n.onMessageFinished,
		MsgHandleTimeOut:  0,
		MsgMaxLiveTime:    time.Second * 5,
		EnableTraceDetail: n.opts.enableTraceDetail,
		KeepAlive:         time.Minute * 5,
		MaxMsgCacheNum:    3,
	})

	n.opts.Logger.Infof("new connect %v comming", s.Stats().RemoteAddress)

	n.mu.Lock()
	if n.sessions == nil { // stoped
		n.mu.Unlock()
		return
	}
	n.sessions[s.Stats().SID] = s
	n.mu.Unlock()

	defer func() {
		n.opts.Logger.Infof("connect %v leave", s.Stats().RemoteAddress)
		n.mu.Lock()
		if n.sessions == nil { // stoped
			n.mu.Unlock()
			return
		}
		delete(n.sessions, s.Stats().SID)
		n.mu.Unlock()
	}()

	s.Handle(n.srvHandle)
}

func (n *Node) onCall(ctx context.Context, s interface{}, srv interface{}, srvName string, methodName string, req interface{}, callee func(req interface{}) (res interface{}, err error)) error {

	beg := time.Now()

	res, err := callee(req)

	if n.opts.showHandleLog {
		n.opts.Logger.Infof("%s.%s(req:{%v}) (res:{%v},err:%v) => %v", srvName, methodName, req, res, err, time.Now().Sub(beg))
	}

	// trace
	if n.opts.enableTraceDetail {
		session.SetReadPayloadObj(ctx, req)
		session.SetWritePayloadObj(ctx, res)
	}

	if n.opts.metricsMsgProcDelayEnabled() {
		metrics.ReportMsgProcDelay(n.opts.metricsReporters, n.opts.ID, n.opts.Name, srvName, methodName, time.Now().Sub(beg))
	}

	return err
}

func (n *Node) onNotify(ctx context.Context, s interface{}, srv interface{}, srvName string, methodName string, req interface{}, callee func(req interface{}) (err error)) error {

	beg := time.Now()

	err := callee(req)

	if n.opts.showHandleLog {
		n.opts.Logger.Infof("%s.%s(req:%v) (err:%v) => %v", srvName, methodName, req, err, time.Now().Sub(beg))
	}

	// trace
	if n.opts.enableTraceDetail {
		session.SetReadPayloadObj(ctx, req)
	}

	if n.opts.metricsMsgProcDelayEnabled() {
		metrics.ReportMsgProcDelay(n.opts.metricsReporters, n.opts.ID, n.opts.Name, srvName, methodName, time.Now().Sub(beg))
	}

	return err
}

func (n *Node) onMessageFinished(ctx context.Context) {
	if len(n.opts.textLogWriter) > 0 {
		if info, ok := session.MsgInfoFromContext(ctx); ok {
			for _, v := range n.opts.textLogWriter {
				err := v.Write(info)
				if err != nil {
					n.opts.Logger.Errorf("[Node.onMessageFinished] write text log faild %v", err)
				}
			}
		}
	}
}
