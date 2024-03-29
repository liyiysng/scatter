// Package node 服务节点
package node

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/liyiysng/scatter/cluster"
	"github.com/liyiysng/scatter/handle"
	"github.com/liyiysng/scatter/logger"
	"github.com/liyiysng/scatter/metrics"
	"github.com/liyiysng/scatter/node/acceptor"
	"github.com/liyiysng/scatter/node/conn"
	"github.com/liyiysng/scatter/node/session"
	"github.com/liyiysng/scatter/util"
	"golang.org/x/net/trace"
	"google.golang.org/grpc"
)

var (
	myLog = logger.Component("node")
)

// IStart 当节点启动时调用
type IStart interface {
	OnStart(n *Node)
}

// IBeforeStop 当节点停止前调用
type IStop interface {
	OnStop()
}

func init() {
	handle.AddIgnorMethod("OnStart")
	handle.AddIgnorMethod("OnStop")
}

var (
	// ErrNodeStopped 节点已停止
	ErrNodeStopped = errors.New("node: the node has been stopped")
	// ErrInvalidCertificates 证书配置错误
	ErrInvalidCertificates = errors.New("certificates must be exactly two")
	// ErrSubSrvNotSupport 不支持的子服务
	ErrSubSrvNotSupport = errors.New("sub service not support")
)

// SocketProtcol 协议类型
type SocketProtcol string

const (
	SocketProtcolUnknown = "unknown"
	// SocketProtcolTCP tcp 协议
	SocketProtcolTCP SocketProtcol = "tcp"
	// SocketProtcolWS websocket 协议
	SocketProtcolWS SocketProtcol = "ws"
)

// FrontRegister 注册前端服务
// 函数签名 func (context.Context,session session.Session,req proto.Message,optionalArg...) (res proto.Message , err error)
type FrontRegister interface {
	// RegisterFront 注册一个前端服务
	// 非协程安全
	RegisterFront(recv interface{}) error
	// RegisterFrontName 注册一个前端命名服务
	// 非协程安全
	RegisterFrontName(name string, recv interface{}) error
}

// GrpcRegister grpc服务注册
type GrpcRegister interface {
	grpc.ServiceRegistrar
	// RegisterSubService 注册子服务
	// 非协程安全
	RegisterSubService(recv interface{}) error
	// RegisterSubServiceName 注册子服务
	// 非协程安全
	RegisterSubServiceName(name string, recv interface{}) error
}

// INodeRegister node 注册
type INodeRegister interface {
	FrontRegister
	GrpcRegister
}

// INodeServe node 服务器
type INodeServe interface {
	// Serve 运行grpc服务,阻塞函数,一般运行在单独的goroutine
	ServeGrpc(lis net.Listener) error
	// Serve 启动一个前端Serve,阻塞函数,一般运行在单独的goroutine
	// Serve 除Stop或者 被调用之外,都返回一个非nil错误
	Serve(sp SocketProtcol, opts ...INodeServeOption) error
	// 停止该节点
	Stop()
}

// Info 节点信息
type Info struct {
	NID            int64
	FrontBindAddr  []string
	FrontOuterAddr []string
	FrontPattern   string
	InnerGrpcAddr  net.Addr
	NetProtocol    SocketProtcol
}

func (i *Info) OutAddr() string {
	addr := "unknown"
	if i.NetProtocol == SocketProtcolWS {
		addr = fmt.Sprintf("%s://%s/%s", string(i.NetProtocol), i.FrontOuterAddr[0], i.FrontPattern)
	} else if i.NetProtocol == SocketProtcolTCP {
		addr = i.FrontOuterAddr[0]
	}
	return addr
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
	sessions map[int64]session.Session

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

	// grpc node
	gnode *cluster.GrpcNode

	// 停止历程
	stopRoutine sync.Once

	// 标识是否服务
	serve bool

	// 前端协议
	sp      SocketProtcol
	pattern string
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
		opts.Logger = logger.GDepthLogger
	}

	if opts.LogPrefix != "" {
		opts.Logger = logger.NewPrefixLogger(opts.Logger, opts.LogPrefix)
	}

	if opts.Name == "" {
		opts.Name = strconv.FormatInt(nid, 10)
	}

	n = &Node{
		accs:      make(map[acceptor.Acceptor]bool),
		opts:      opts,
		startTime: time.Now(),
		sessions:  make(map[int64]session.Session),
		quit:      util.NewEvent(),
		sp:        SocketProtcolUnknown,
		pattern:   "",
		gnode: cluster.NewGrpcNode(strconv.FormatInt(opts.ID, 10),
			cluster.OptWithLogger(opts.Logger),
			cluster.OptWithGrpcOption(opts.grpcOpts...)),
	}

	n.idGen, err = util.NewNode(nid)
	if err != nil {
		return nil, err
	}

	if n.opts.subSrvEnable {
		n.srvHandle = newSrvHandleProxy(handle.NewServiceHandle(&handle.Option{
			Codec:      n.opts.getCodec(),
			HookCall:   n.onCall,
			HookNofify: n.onNotify,
			OptArgs:    n.opts.optArgs,
		}),
			n.opts.subSrvClientBuilder,
			n.opts.subSrvValidator)
	} else {
		n.srvHandle = handle.NewServiceHandle(&handle.Option{
			Codec:      n.opts.getCodec(),
			HookCall:   n.onCall,
			HookNofify: n.onNotify,
			OptArgs:    n.opts.optArgs,
		})
	}

	if opts.enableEventTrace {
		_, file, line, _ := runtime.Caller(1)
		n.trEvents = trace.NewEventLog("scatter.Node", fmt.Sprintf("%d-%s:%d", opts.ID, file, line))
	}

	opts.Logger.Infof("create node %d", opts.ID)

	return
}

// GetNID 获取node id
func (n *Node) GetNID() int64 {
	return n.opts.ID
}

// GetInfo 获取节点信息
func (n *Node) GetInfo() Info {
	n.mu.Lock()
	defer n.mu.Unlock()

	bindAddrs := []string{}
	outerAddrs := []string{}
	for k := range n.accs {
		bindAddrs = append(bindAddrs, k.GetAddr())
		outerAddrs = append(outerAddrs, k.GetOuterAddr())
	}

	return Info{
		NID:            n.opts.ID,
		FrontBindAddr:  bindAddrs,
		FrontOuterAddr: outerAddrs,
		InnerGrpcAddr:  n.gnode.GetAddr(),
		NetProtocol:    n.sp,
		FrontPattern:   n.pattern,
	}
}

// RegisterFront 注册前端服务
func (n *Node) RegisterFront(recv interface{}) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.accs == nil {
		return ErrNodeStopped
	}

	if n.serve {
		return fmt.Errorf("[Node.Register] register service after Node.Serve")
	}
	return n.srvHandle.Register(recv,
		handle.OptWithReqTypeValidator(n.opts.reqTypeValidator),
		handle.OptWithResTypeValidator(n.opts.resTypeValidator),
		handle.OptWithSessionTypeValidator(n.opts.sessionTypeValidator),
	)
}

// RegisterFrontName 注册命名前端服务
func (n *Node) RegisterFrontName(name string, recv interface{}) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.accs == nil {
		return ErrNodeStopped
	}

	if n.serve {
		return fmt.Errorf("[Node.Register] register service after Node.Serve %q", name)
	}

	return n.srvHandle.RegisterName(name, recv,
		handle.OptWithReqTypeValidator(n.opts.reqTypeValidator),
		handle.OptWithResTypeValidator(n.opts.resTypeValidator),
		handle.OptWithSessionTypeValidator(n.opts.sessionTypeValidator),
	)
}

func (n *Node) RegisterSubService(recv interface{}) error {
	return n.gnode.RegisterSubService(recv)
}

func (n *Node) RegisterSubServiceName(name string, recv interface{}) error {
	return n.gnode.RegisterSubServiceName(name, recv)
}

// RegisterService implement grpc.ServiceRegistrar
func (n *Node) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.gnode == nil {
		panic(ErrNodeStopped)
	}
	n.gnode.RegisterService(desc, impl)
}

// Subscribe 订阅
// 非协程安全
func (n *Node) Subscribe(topic string, recv interface{}) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.gnode == nil {
		panic(ErrNodeStopped)
	}
	return n.gnode.Subscribe(topic, recv)
}

// ServeGrpc run grpc server
func (n *Node) ServeGrpc(lis net.Listener) error {
	n.waitGroup.Add(1)
	defer n.waitGroup.Done()

	n.mu.Lock()
	if n.gnode == nil {
		n.mu.Unlock()
		return ErrNodeStopped
	}
	n.mu.Unlock()

	// 调用onStart
	allRecv := []interface{}{}
	n.mu.Lock()
	h := n.gnode.GetSubServiceHandle()
	if h != nil {
		allRecv = n.gnode.GetSubServiceHandle().AllRecv()
	}
	n.mu.Unlock()

	for _, v := range allRecv {
		if s, ok := v.(IStart); ok {
			go s.OnStart(n)
		}
	}

	defer func() {
		// OnStop
		for _, v := range allRecv {
			if s, ok := v.(IStop); ok {
				go s.OnStop()
			}
		}
	}()

	return n.gnode.Serve(lis)
}

// AddAfterStop 添加停止后需执行的函数
func (n *Node) AddAfterStop(f ...func()) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.accs == nil {
		return ErrNodeStopped
	}

	n.opts.afterStop = append(n.opts.afterStop, f...)
	return nil
}

func (n *Node) SessionCount() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return len(n.sessions)
}

// Serve 启动一个Serve
// Serve 除Stop或者 被调用之外,都返回一个非nil错误
// arg[0] = certfile
// arg[1] = keyfile
func (n *Node) Serve(sp SocketProtcol, addr string, opts ...INodeServeOption) error {

	opt := NodeServeOption{}
	for _, v := range opts {
		v.apply(&opt)
	}

	if opt.outerAddr == "" {
		opt.outerAddr = addr
	}

	n.mu.Lock()
	n.sp = sp
	if opt.pattern != "" {
		n.pattern = opt.pattern
	}
	n.mu.Unlock()

	n.waitGroup.Add(1)
	defer n.waitGroup.Done()

	n.mu.Lock()
	n.trEventLogf("starting")
	n.serve = true

	if n.accs == nil {
		// Start called after Stop
		n.mu.Unlock()
		return ErrNodeStopped
	}

	accOpt := acceptor.Option{
		Addr:      addr,
		OuterAddr: opt.outerAddr,
		CertFile:  opt.certFile,
		KeyFile:   opt.keyFile,
		Logger:    n.opts.Logger,
		GetConnOpt: func() conn.MsgConnOption {
			ret := conn.MsgConnOption{
				SID:                 n.idGen.Generate().Int64(),
				MaxLength:           n.opts.maxPayloadLength,
				ReadTimeout:         n.opts.readTimeout,
				WriteTimeout:        n.opts.writeTimeout,
				ReadBufferSize:      n.opts.readBufferSize,
				WriteBufferSize:     n.opts.writeBufferSize,
				Compresser:          n.opts.getCompressor(),
				EnableWriteLimit:    n.opts.enableWriteLimit,
				EnableReadLimit:     n.opts.enableReadLimit,
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
			if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
				n.opts.Logger.Errorf("[Node.Serve] node %v stoped cause by %v", n.opts.ID, err)
			}
		}, n.opts.Logger.Errorf)

	// 调用onStart
	n.mu.Lock()
	allRecv := n.srvHandle.AllRecv()
	n.mu.Unlock()

	for _, v := range allRecv {
		if s, ok := v.(IStart); ok {
			go s.OnStart(n)
		}
	}

	defer func() {
		// OnStop
		for _, v := range allRecv {
			if s, ok := v.(IStop); ok {
				go s.OnStop()
			}
		}
	}()

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

	if n.opts.Logger.V(logger.VDEBUG) {
		n.opts.Logger.Infof("stopping node %d", n.opts.ID)
	}

	n.quit.Fire()

	defer func() {
		n.waitGroup.Wait()
		// 清理历程
		n.stopRoutine.Do(func() {
			if n.opts.textLogWriter != nil && len(n.opts.textLogWriter) > 0 {
				for _, v := range n.opts.textLogWriter {
					err := v.Close()
					if err != nil {
						n.opts.Logger.Errorf("close text log failed %v", err)
					}
				}

			}
			if n.opts.Logger.V(logger.VIMPORTENT) {
				n.opts.Logger.Infof("node %d stoped", n.opts.ID)
			}
			if len(n.opts.afterStop) > 0 {
				for _, f := range n.opts.afterStop {
					f()
				}
			}
		})
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

	n.mu.Lock()
	if n.gnode != nil {
		n.gnode.Stop()
		n.gnode = nil
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
		MsgMaxLiveTime:    n.opts.msgMaxLiveTime,
		EnableTraceDetail: n.opts.enableTraceDetail,
		KeepAlive:         time.Minute * 5,
		MaxMsgCacheNum:    3,
		KickTimeout:       time.Second,
		TimerResolution:   time.Second,
		ShowMsgRWBytes:    n.opts.showMsgRWBytes,
	})

	if n.opts.Logger.V(logger.VDEBUG) {
		n.opts.Logger.Infof("new connection %v comming", s.PeerAddr())
	}

	n.mu.Lock()
	if n.sessions == nil { // stoped
		n.mu.Unlock()
		return
	}
	n.sessions[s.GetSID()] = s
	n.mu.Unlock()

	defer func() {
		if n.opts.Logger.V(logger.VDEBUG) {
			n.opts.Logger.Infof("connection %v leave", s.PeerAddr())
		}
		n.mu.Lock()
		if n.sessions == nil { // stoped
			n.mu.Unlock()
			return
		}
		delete(n.sessions, s.GetSID())
		n.mu.Unlock()
	}()

	s.Handle(n.srvHandle)
}

func (n *Node) onCall(ctx context.Context, s interface{}, srv interface{}, srvName string, methodName string, req interface{}, callee func(req interface{}) (res interface{}, err error)) (cres interface{}, err error) {

	beg := time.Now()

	if n.opts.callHook != nil {
		cres, err = n.opts.callHook(ctx, s, srv, srvName, methodName, req, callee)
	} else {
		cres, err = callee(req)
	}

	if n.opts.showHandleLog {
		n.opts.Logger.Infof("%s.%s(req:{%v}) (res:{%v},err:%v) => %v", srvName, methodName, req, cres, err, time.Since(beg))
	}

	// trace
	if n.opts.enableTraceDetail {
		session.SetReadPayloadObj(ctx, req)
		session.SetWritePayloadObj(ctx, cres)
	}

	if n.opts.metricsMsgProcDelayEnabled() {
		metrics.ReportMsgProcDelay(n.opts.metricsReporters, n.opts.ID, n.opts.Name, srvName, methodName, time.Since(beg))
	}

	return
}

func (n *Node) onNotify(ctx context.Context, s interface{}, srv interface{}, srvName string, methodName string, req interface{}, callee func(req interface{}) (err error)) error {

	beg := time.Now()

	var err error = nil

	if n.opts.notifyHook != nil {
		err = n.opts.notifyHook(ctx, s, srv, srvName, methodName, req, callee)
	} else {
		err = callee(req)
	}

	if n.opts.showHandleLog {
		n.opts.Logger.Infof("%s.%s(req:%v) (err:%v) => %v", srvName, methodName, req, err, time.Since(beg))
	}

	// trace
	if n.opts.enableTraceDetail {
		session.SetReadPayloadObj(ctx, req)
	}

	if n.opts.metricsMsgProcDelayEnabled() {
		metrics.ReportMsgProcDelay(n.opts.metricsReporters, n.opts.ID, n.opts.Name, srvName, methodName, time.Since(beg))
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

func (n *Node) AddJob(numPerSec int, foreachOnlineSession func(session session.Session)) (done <-chan struct{}, err error) {

	mdone := make(chan struct{})

	SIDs := make([]int64, 0, 128)
	// copy sid
	n.mu.Lock()
	for k := range n.sessions {
		SIDs = append(SIDs, k)
	}
	n.mu.Unlock()

	go func() {
		defer close(mdone)
		for index, value := range SIDs {
			if (index+1)%numPerSec == 0 {
				time.Sleep(time.Second)
			}
			// get session
			n.mu.Lock()
			s, ok := n.sessions[value]
			n.mu.Unlock()
			if ok {
				foreachOnlineSession(s)
			}
		}
	}()
	return mdone, nil
}
