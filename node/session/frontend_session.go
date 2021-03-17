package session

import (
	"context"
	"fmt"
	"io"
	"net"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/liyiysng/scatter/cluster/selector/policy"
	"github.com/liyiysng/scatter/constants"
	"github.com/liyiysng/scatter/encoding"
	"github.com/liyiysng/scatter/logger"
	"github.com/liyiysng/scatter/node/conn"
	"github.com/liyiysng/scatter/node/handle"
	"github.com/liyiysng/scatter/node/message"
	phead "github.com/liyiysng/scatter/node/message/proto"
	"github.com/liyiysng/scatter/ratelimit"
	"github.com/liyiysng/scatter/util"
)

// Option session选项
type Option struct {
	Logger logger.Logger
	// 链接超时(当链接创建后多久事件未接受到handshake消息)
	ConnectTimeout time.Duration
	// 读写Chan缓冲大小
	ReadChanSize  int
	WriteChanSize int
	// 限制每秒消息处理
	// <=0 不受限制
	RateLimitMsgProc int64
	// push拦截
	PushInterceptor handle.SerrvicePushInterceptor
	// 编码
	Codec encoding.Codec
	// 当消息完成
	OnMsgFinish func(ctx context.Context)
	// 消息处理超时
	// 0:表示不超时
	MsgHandleTimeOut time.Duration
	// 消息存活时间
	// 0:表示一直存活
	MsgMaxLiveTime time.Duration
	// 是监视详情
	EnableTraceDetail bool
	// keepAlive
	KeepAlive time.Duration
	// 消息缓存量
	// flush调用时机:
	// 1.缓存消息数量大于该值
	// 2.wirteChan数量为0
	// session 会接连发送[1,MaxMsgCacheNum]条消息到bufio中
	MaxMsgCacheNum int
	// kick timeout ,剔出超时
	KickTimeout time.Duration
}

type msgCtx struct {
	// request , notify , heartbeat , handshake...
	msgRead message.Message
	// reponse , push , heatbeatact , handshakeact...
	msgWrite message.Message
	// 是否压缩等
	popt          message.PacketOpt
	ctx           context.Context
	cancel        context.CancelFunc
	timeoutCancel context.CancelFunc
}

func (mctx *msgCtx) WithTimeout(d time.Duration) {
	if mctx.ctx == nil || mctx.cancel == nil {
		panic("nil msgCtx")
	}
	mctx.ctx, mctx.timeoutCancel = context.WithTimeout(mctx.ctx, d)
}

func (mctx *msgCtx) WithDeadline(d time.Time) {
	if mctx.ctx == nil || mctx.cancel == nil {
		panic("nil msgCtx")
	}
	mctx.ctx, mctx.timeoutCancel = context.WithDeadline(mctx.ctx, d)
}

func (mctx *msgCtx) WithValue(key, val interface{}) {
	if mctx.ctx == nil || mctx.cancel == nil {
		panic("nil msgCtx")
	}
	mctx.ctx = context.WithValue(mctx.ctx, key, val)
}

var msgCtxPool = sync.Pool{
	New: func() interface{} {
		return &msgCtx{}
	},
}

func getMsgCtxWithContext(ctx context.Context, enableProfile bool, popt message.PacketOpt) *msgCtx {
	ret := msgCtxPool.Get().(*msgCtx)
	ret.ctx, ret.cancel = context.WithCancel(ctx)
	if enableProfile {
		ret.ctx = withInfo(ret.ctx)
	}
	ret.popt = popt
	return ret
}

func getMsgCtx(enableProfile bool, popt message.PacketOpt) *msgCtx {
	ret := msgCtxPool.Get().(*msgCtx)
	ret.ctx, ret.cancel = context.WithCancel(context.Background())
	if enableProfile {
		ret.ctx = withInfo(ret.ctx)
	}
	ret.popt = popt
	return ret
}

func putMsgCtx(mctx *msgCtx) {
	// reset
	mctx.msgRead = nil
	mctx.msgWrite = nil
	mctx.cancel = nil
	mctx.ctx = nil

	// put in pool
	msgCtxPool.Put(mctx)
}

type frontendSession struct {
	nid int64

	opt *Option

	rateLimt *ratelimit.Bucket

	conn       conn.MsgConn
	connClosed bool
	connMu     sync.Mutex

	readChan chan *msgCtx // request chan  ,  ensure request sequence

	writeChanClosed bool
	writeChan       chan *msgCtx // need response to client , include push

	closeEvent *util.Event

	onClose OnClose
	mu      sync.Mutex

	// session context
	ctx                 context.Context
	concel              context.CancelFunc
	sessionAffinityData interface{}

	// 握手数据
	handShake interface{}

	// 上次心跳时间
	lastHeartBeat time.Time

	attrs  map[string]interface{}
	attrMu sync.Mutex

	wg util.WaitGroupWrapper
}

// NewFrontendSession 创建一个session
func NewFrontendSession(nid int64, c conn.MsgConn, opt *Option) IFrontendSession {
	ret := &frontendSession{
		nid:        nid,
		conn:       c,
		opt:        opt,
		readChan:   make(chan *msgCtx, opt.ReadChanSize),
		writeChan:  make(chan *msgCtx, opt.WriteChanSize),
		closeEvent: util.NewEvent(),
		ctx:        context.Background(),
		attrs:      make(map[string]interface{}),
	}
	ret.ctx, ret.concel = context.WithCancel(ret.ctx)
	// session_affinity support
	ret.ctx = policy.WithSessionAffinity(ret.ctx, func(data interface{}) {
		ret.sessionAffinityData = data
	}, func() (conn interface{}) {
		return ret.sessionAffinityData
	})

	if ret.opt.RateLimitMsgProc > 0 {
		ret.rateLimt = ratelimit.NewBucketWithQuantum(time.Second, ret.opt.RateLimitMsgProc, ret.opt.RateLimitMsgProc)
	}

	// 读取消息
	ret.wg.Wrap(ret.runRead, ret.opt.Logger.Errorf)
	// 写消息
	ret.wg.Wrap(ret.runWrite, ret.opt.Logger.Errorf)

	return ret
}

func (s *frontendSession) SetAttr(key string, v interface{}) {
	s.attrMu.Lock()
	defer s.attrMu.Unlock()
	s.attrs[key] = v
}

func (s *frontendSession) GetAttr(key string) (v interface{}, ok bool) {
	s.attrMu.Lock()
	defer s.attrMu.Unlock()
	v, ok = s.attrs[key]
	return
}

func (s *frontendSession) GetCtx() context.Context {
	return s.ctx
}

func (s *frontendSession) GetSID() int64 {
	return s.conn.GetSID()
}

func (s *frontendSession) GetNID() int64 {
	return s.nid
}

func (s *frontendSession) PeerAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *frontendSession) Push(ctx context.Context, cmd string, v interface{}, popt ...message.IPacketOption) error {

	data, err := s.opt.Codec.Marshal(v)
	if err != nil {
		return err
	}

	if s.closeEvent.HasFired() {
		return constants.ErrSessionClosed
	}

	p := message.DEFAULTPOPT

	if len(popt) > 0 {
		for _, v := range popt {
			v.Apply(&p)
		}
	}

	mctx := getMsgCtxWithContext(ctx, s.opt.EnableTraceDetail, p)

	mctx.msgWrite, err = message.MsgFactory.BuildPushMessage(cmd, data)
	if err != nil {
		return err
	}

	if s.opt.EnableTraceDetail {
		SetWritePayloadObj(mctx.ctx, v)
	}

	return s.push(mctx)
}

func (s *frontendSession) push(mctx *msgCtx) error {

	if s.opt.PushInterceptor != nil {
		needPush, err := s.opt.PushInterceptor(mctx.ctx, mctx.msgWrite)
		if err != nil {
			s.finishMsg(mctx, err)
			return err
		}
		if !needPush {
			s.finishMsg(mctx, fmt.Errorf("push message intercepted"))
			return nil
		}
	}

	if s.opt.EnableTraceDetail {
		onMsgWrite(mctx.ctx, mctx.msgWrite)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.writeChanClosed {
		return constants.ErrSessionClosed
	}

	select {
	case s.writeChan <- mctx:
		{
			return nil
		}
	case <-mctx.ctx.Done(): // 取消或者超时
		{
			s.finishMsg(mctx, mctx.ctx.Err())
			return mctx.ctx.Err()
		}
	case <-s.closeEvent.Done():
		{
			s.finishMsg(mctx, constants.ErrSessionClosed)
			return constants.ErrSessionClosed
		}
	}
}

func (s *frontendSession) PushTimeout(ctx context.Context, cmd string, v interface{}, timeout time.Duration, popt ...message.IPacketOption) error {

	if s.closeEvent.HasFired() {
		return constants.ErrSessionClosed
	}

	data, err := s.opt.Codec.Marshal(v)
	if err != nil {
		return err
	}

	p := message.DEFAULTPOPT

	if len(popt) > 0 {
		for _, v := range popt {
			v.Apply(&p)
		}
	}

	mctx := getMsgCtxWithContext(ctx, s.opt.EnableTraceDetail, p)

	mctx.msgWrite, err = message.MsgFactory.BuildPushMessage(cmd, data)
	if err != nil {
		return err
	}

	mctx.WithTimeout(timeout)

	if s.opt.EnableTraceDetail {
		SetWritePayloadObj(mctx.ctx, v)
	}

	return s.push(mctx)
}

func (s *frontendSession) PushImmediately(ctx context.Context, cmd string, v interface{}, popt ...message.IPacketOption) error {

	if s.closeEvent.HasFired() {
		return constants.ErrSessionClosed
	}

	data, err := s.opt.Codec.Marshal(v)
	if err != nil {
		return err
	}

	p := message.DEFAULTPOPT

	if len(popt) > 0 {
		for _, v := range popt {
			v.Apply(&p)
		}
	}

	mctx := getMsgCtxWithContext(ctx, s.opt.EnableTraceDetail, p)

	mctx.msgWrite, err = message.MsgFactory.BuildPushMessage(cmd, data)
	if err != nil {
		return err
	}

	if s.opt.EnableTraceDetail {
		SetWritePayloadObj(mctx.ctx, v)
	}

	if s.opt.PushInterceptor != nil {
		needPush, err := s.opt.PushInterceptor(mctx.ctx, mctx.msgWrite)
		if err != nil {
			s.finishMsg(mctx, err)
			return err
		}
		if !needPush {
			s.finishMsg(mctx, fmt.Errorf("push message intercepted"))
			return nil
		}
	}

	if s.opt.EnableTraceDetail {
		onMsgWrite(mctx.ctx, mctx.msgWrite)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.writeChanClosed {
		return constants.ErrSessionClosed
	}

	select {
	case s.writeChan <- mctx:
		{
			return nil
		}
	case <-s.closeEvent.Done():
		{
			s.finishMsg(mctx, constants.ErrSessionClosed)
			return constants.ErrSessionClosed
		}
	default:
		{
			s.finishMsg(mctx, constants.ErrorPushBufferFull)
			return constants.ErrorPushBufferFull
		}
	}
}

func (s *frontendSession) SetOnClose(onClose OnClose) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onClose = onClose
}

func (s *frontendSession) Closed() bool {
	return s.closeEvent.HasFired()
}

func (s *frontendSession) Kick() error {
	if s.closeEvent.HasFired() {
		return constants.ErrSessionClosed
	}

	mctx := getMsgCtx(s.opt.EnableTraceDetail, message.DEFAULTPOPT)

	var err error
	mctx.msgWrite, err = message.MsgFactory.BuildKickMessage()
	if err != nil {
		return err
	}

	mctx.WithTimeout(s.opt.KickTimeout)

	return s.push(mctx)
}

// Close 关闭session 可多次调用
func (s *frontendSession) Close() {
	s.closeEvent.Fire()
}

func (s *frontendSession) Handle(srvHandler handle.IHandler) {

	tick := time.NewTicker(s.opt.KeepAlive)
	defer tick.Stop()

	checkConnectTimeout := time.NewTimer(s.opt.ConnectTimeout)
	defer checkConnectTimeout.Stop()

	defer func() {

		s.Close()

		s.mu.Lock()
		s.writeChanClosed = true
		// 关闭写chan
		close(s.writeChan)
		s.mu.Unlock()

		//丢弃剩余的消息
		s.drainRead()

		// 读写协程都结束
		s.wg.Wait()

		////////////////////////////////close routine/////////////////////////////////
		// copy onClose
		var onClose OnClose
		s.mu.Lock()
		onClose = s.onClose
		s.mu.Unlock()

		if onClose != nil {
			onClose(s)
		}

		// cancel context
		s.concel()
	}()

	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			s.opt.Logger.Errorf("goroutine panic %v\n%s", err, buf)
		}
	}()

	// 退出条件
	// 1.当处理完所有请求后退出(reqChan关闭是退出该协程)
	// 2.遇到关键错误
	// 3.关闭事件
	for {

		select {
		case <-s.closeEvent.Done():
			{
				return
			}
		case mctx, ok := <-s.readChan: // 处理消息
			{
				if !ok { // reqChan closed
					return
				}

				// 若收到握手数据
				if mctx.msgRead.GetMsgType() == phead.MsgType_HANDSHAKE {
					checkConnectTimeout.Stop()
				}

				// 处理消息
				err := s.handleMsg(mctx, srvHandler)
				if err != nil {
					// 关键错误
					if cErr, ok := err.(handle.ICriticalError); ok {
						s.finishMsg(mctx, cErr)
						s.opt.Logger.Errorf("handle message encount a critical error %v", cErr)
						return
					} else if customErr, ok := err.(handle.ICustomError); ok {
						// 非关键错误
						if s.opt.Logger.V(logger.VDEBUG) {
							s.opt.Logger.Warningf("handle message encount custom error %v", customErr)
						}
					} else {
						s.finishMsg(mctx, cErr)
						if s.opt.Logger.V(logger.VDEBUG) {
							s.opt.Logger.Warningf("handle message encount a error %v", err)
						}
						return
					}
				}
			}
		case <-tick.C: // 检查心跳
			{
				if !s.lastHeartBeat.IsZero() {

				}
			}
		case <-checkConnectTimeout.C:
			{
				if s.handShake == nil {
					s.opt.Logger.Warningf("connect timeout , handshake not recved in %v", s.opt.ConnectTimeout)
					return
				}
			}
		}
	}

}

func (s *frontendSession) handleMsg(mctx *msgCtx, srvHandler handle.IHandler) error {

	if mctx.ctx.Err() != nil { // 已经超时或取消
		return handle.NewCustomErrorWithError(mctx.ctx.Err())
	}

	if s.opt.MsgHandleTimeOut != 0 {
		mctx.WithTimeout(s.opt.MsgHandleTimeOut)
	}

	// 处理读取的消息
	switch mctx.msgRead.GetMsgType() {
	case phead.MsgType_NOTIFY:
		{
			err := s.handleNotify(mctx, srvHandler)
			if err != nil {
				return err
			}
			s.finishMsg(mctx, nil)
		}
	case phead.MsgType_REQUEST:
		{
			sequence := mctx.msgRead.GetSequence()
			srv := mctx.msgRead.GetService()
			resBuf, err := s.handleRequest(mctx, srvHandler)

			if err != nil {
				if customErr, ok := err.(handle.ICustomError); ok {
					// 非关键错误
					// 回复客户端
					mctx.msgWrite, err = message.MsgFactory.BuildResponseCustomErrorMessage(sequence, srv, customErr.Error())
					if err != nil {
						return err
					}
					err = s.sendWrite(mctx)
					if err != nil {
						return err
					}
					return customErr
				}
				return err
			}

			// 回复客户端
			mctx.msgWrite, err = message.MsgFactory.BuildResponseMessage(sequence, srv, resBuf)
			if err != nil {
				return err
			}
			err = s.sendWrite(mctx)
			if err != nil {
				return err
			}
		}
	case phead.MsgType_HEARTBEAT:
		{
			err := s.handleHeartbeat(mctx)
			if err != nil {
				return err
			}

			// write ack
			mctx.msgWrite, err = message.MsgFactory.BuildHeatAckMessage()
			if err != nil {
				return err
			}

			err = s.sendWrite(mctx)
			if err != nil {
				return err
			}
		}
	case phead.MsgType_HANDSHAKE:
		{
			err := s.handleHandshake(mctx)
			if err != nil {
				return err
			}

			// write ack
			mctx.msgWrite, err = message.MsgFactory.BuildHandShakeAckMessage()
			if err != nil {
				return err
			}
			err = s.sendWrite(mctx)
			if err != nil {
				return err
			}
		}
	default:
		return handle.NewCriticalErrorf("invalid message type %v", mctx.msgRead.GetMsgType())
	}

	return nil
}

func (s *frontendSession) handleRequest(mctx *msgCtx, srvHandler handle.IHandler) (resBuf []byte, err error) {
	// profile
	if s.opt.EnableTraceDetail {
		onMsgBegingHandle(mctx.ctx)
	}
	if s.opt.EnableTraceDetail {
		defer onMsgEndHandle(mctx.ctx)
	}

	reqMsg := mctx.msgRead
	serviceName, methodName, err := message.GetSrvMethod(reqMsg.GetService())
	if err != nil {
		return nil, handle.NewCriticalErrorWithError(err)
	}
	resBuf, err = srvHandler.Call(mctx.ctx, s, serviceName, methodName, reqMsg.GetPayload())
	return
}

func (s *frontendSession) handleNotify(mctx *msgCtx, srvHandler handle.IHandler) error {
	// profile
	if s.opt.EnableTraceDetail {
		onMsgBegingHandle(mctx.ctx)
	}
	if s.opt.EnableTraceDetail {
		defer onMsgEndHandle(mctx.ctx)
	}

	reqMsg := mctx.msgRead

	serviceName, methodName, err := message.GetSrvMethod(reqMsg.GetService())
	if err != nil {
		return handle.NewCriticalErrorWithError(err)
	}

	err = srvHandler.Notify(mctx.ctx, s, serviceName, methodName, reqMsg.GetPayload())
	if err != nil {
		return err
	}

	return nil
}

func (s *frontendSession) handleHeartbeat(mctx *msgCtx) error {
	// profile
	if s.opt.EnableTraceDetail {
		onMsgBegingHandle(mctx.ctx)
	}
	if s.opt.EnableTraceDetail {
		defer onMsgEndHandle(mctx.ctx)
	}

	s.lastHeartBeat = time.Now()

	return nil
}

func (s *frontendSession) handleHandshake(mctx *msgCtx) error {
	// profile
	if s.opt.EnableTraceDetail {
		onMsgBegingHandle(mctx.ctx)
	}
	if s.opt.EnableTraceDetail {
		defer onMsgEndHandle(mctx.ctx)
	}

	if s.handShake != nil {
		return handle.NewCriticalError("dupulicate handshake")
	}

	handShake := mctx.msgRead
	h, err := message.MsgFactory.ParseHandShake(handShake.GetPayload())
	if err != nil {
		return err
	}
	s.handShake = h

	if s.opt.EnableTraceDetail {
		SetReadPayloadObj(mctx.ctx, h)
	}

	return nil
}

// prevent block in write chan when closing
func (s *frontendSession) sendWrite(mctx *msgCtx) error {

	if s.opt.EnableTraceDetail {
		onMsgWrite(mctx.ctx, mctx.msgWrite)
	}

	select {
	case <-mctx.ctx.Done():
		{
			return handle.NewCustomErrorWithError(mctx.ctx.Err())
		}
	case s.writeChan <- mctx:
		{
			return nil
		}
	case <-s.closeEvent.Done():
		{
			return constants.ErrSessionClosed
		}
	}
}

func (s *frontendSession) runWrite() {
	defer func() {
		s.Close()
		// 关闭链接
		s.closeConn("sever close")
	}()

	cacheMsgNum := 0
	// 退出写协程条件:
	// 1.writeChan关闭
	// 2.conn写消息错误
	// 3.关闭事件
	for {
		select {
		case <-s.closeEvent.Done():
			{
				// 确保所有消息已被正确发送
				s.drainWrite()
				return
			}
		case mctx, ok := <-s.writeChan:
			{
				if !ok { // chan closed
					return
				}

				err := s.writeMsg(mctx)
				if err != nil {
					s.opt.Logger.Errorf("write message error %v", err)
					// io错误 关闭链接
					s.closeConn(fmt.Sprintf("write io error:%v", err))
					return
				}
				cacheMsgNum++
				// need flush?
				if len(s.writeChan) == 0 || cacheMsgNum >= s.opt.MaxMsgCacheNum {
					if err = s.flush(); err != nil {
						s.opt.Logger.Errorf("flush message error %v", err)
						// io错误 关闭链接
						s.closeConn(fmt.Sprintf("write io error:%v", err))
						return
					}
					cacheMsgNum = 0
				}
			}
		}
	}
}

func (s *frontendSession) flush() error {
	var err error
	s.connMu.Lock()
	if !s.connClosed {
		err = s.conn.Flush()
	} else {
		err = constants.ErrSessionClosed
	}
	s.connMu.Unlock()
	return err
}

func (s *frontendSession) writeMsg(mctx *msgCtx) error {
	if mctx.msgWrite == nil {
		panic("nil wirte msg")
	}

	var err error

	s.connMu.Lock()
	if !s.connClosed {
		err = s.conn.WriteNextMessage(mctx.msgWrite, mctx.popt)
	} else {
		err = fmt.Errorf("%v message:{%v}", constants.ErrSessionClosed, mctx.msgWrite)
	}
	s.connMu.Unlock()

	s.finishMsg(mctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *frontendSession) finishMsg(mctx *msgCtx, err error) {
	mctx.cancel() // release resource
	defer putMsgCtx(mctx)

	if s.opt.EnableTraceDetail {
		onMsgFinished(mctx.ctx, err)
	}

	if s.opt.OnMsgFinish != nil {
		s.opt.OnMsgFinish(mctx.ctx)
	}

}

func (s *frontendSession) runRead() {
	defer func() {
		s.Close()
		// 关闭req chan
		close(s.readChan)
	}()

	// 退出读协程条件:
	// 1.关闭事件
	// 2.conn读取消息错误,一般为EOF
	// 3.链接已关闭
	for {

		// 消息限流
		if s.rateLimt != nil {
			s.rateLimt.Wait(1)
		}

		msg, popt, err := s.conn.ReadNextMessage()
		if err != nil {
			if err != io.EOF && err != io.ErrUnexpectedEOF && !strings.Contains(err.Error(), "use of closed network connection") {
				s.opt.Logger.Errorf("read message error %v", err)
				s.closeConn(err.Error())
			} else {
				// io err 关闭链接
				s.closeConn("peer close")
			}
			return
		}

		mctx := getMsgCtx(s.opt.EnableTraceDetail, popt)
		mctx.msgRead = msg

		// 消息存活时间
		if s.opt.MsgMaxLiveTime != 0 {
			mctx.WithTimeout(s.opt.MsgMaxLiveTime)
		}

		// profile
		if s.opt.EnableTraceDetail {
			onMsgRead(mctx.ctx, mctx.msgRead)
		}

		select {
		case <-s.closeEvent.Done():
			{
				return
			}
		case s.readChan <- mctx:
			{
			}
		}
	}
}

func (s *frontendSession) drainWrite() {

	var err error

	for {
		mctx, ok := <-s.writeChan
		if !ok { // chan closed
			return
		}

		if err == nil { // 无错误,继续发送
			err = s.writeMsg(mctx)
			if err != nil {
				s.opt.Logger.Errorf("write message error %v", err)
				continue
			}
			err = s.flush()
			if err != nil {
				s.opt.Logger.Errorf("flush message error %v", err)
				continue
			}
		} else { // 写消息错误 , 丢弃剩余消息
			s.finishMsg(mctx, constants.ErrorMsgDiscard)
		}
	}
}

// 丢弃所有读到的消息
func (s *frontendSession) drainRead() {
	for {
		mctx, ok := <-s.readChan
		if !ok { // chan closed
			return
		}
		s.finishMsg(mctx, constants.ErrorMsgDiscard)
	}
}

func (s *frontendSession) closeConn(why string) {
	s.connMu.Lock()
	defer s.connMu.Unlock()

	if s.connClosed {
		return
	}

	if s.opt.Logger.V(logger.VDEBUG) {
		s.opt.Logger.Infof("close conn caused by %s", why)
	}

	// 关闭链接
	err := s.conn.Close()
	if err != nil {
		s.opt.Logger.Errorf("session close error %v", err)
	}
	s.connClosed = true
}
