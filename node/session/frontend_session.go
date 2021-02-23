package session

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/liyiysng/scatter/constants"
	"github.com/liyiysng/scatter/encoding"
	"github.com/liyiysng/scatter/logger"
	"github.com/liyiysng/scatter/node/conn"
	"github.com/liyiysng/scatter/node/handle"
	"github.com/liyiysng/scatter/node/message"
	"github.com/liyiysng/scatter/util"
)

// Option session选项
type Option struct {
	Logger      logger.Logger
	ReqChanSize int
	ResChanSize int
	// 获取消息选项 , 如某些消息需要压缩等
	GetMessageOpt func(msg message.Message) message.PacketOpt
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
}

type msgCtx struct {
	// request , notify , heartbeat , handshake...
	msgRead message.Message
	// reponse , push , heatbeatact , handshakeact...
	msgWrite      message.Message
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

func getMsgCtxWithContext(ctx context.Context, enableProfile bool) *msgCtx {
	ret := msgCtxPool.Get().(*msgCtx)
	ret.ctx, ret.cancel = context.WithCancel(ctx)
	if enableProfile {
		ret.ctx = withInfo(ret.ctx)
	}
	return ret
}

func getMsgCtx(enableProfile bool) *msgCtx {
	ret := msgCtxPool.Get().(*msgCtx)
	ret.ctx, ret.cancel = context.WithCancel(context.Background())
	if enableProfile {
		ret.ctx = withInfo(ret.ctx)
	}
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
	nid string
	opt *Option

	conn       conn.MsgConn
	connClosed bool
	connMu     sync.Mutex

	readChan chan *msgCtx // request chan  ,  ensure request sequence

	writeChanClosed bool
	writeChan       chan *msgCtx // need response to client , include push

	closeEvent *util.Event

	onClose OnClose
	mu      sync.Mutex

	// 握手数据
	handShake interface{}

	// 上次心跳时间
	lastHeartBeat time.Time

	wg util.WaitGroupWrapper
}

// NewFrontendSession 创建一个session
func NewFrontendSession(nid string, c conn.MsgConn, opt *Option) FrontendSession {
	ret := &frontendSession{
		nid:        nid,
		conn:       c,
		opt:        opt,
		readChan:   make(chan *msgCtx, opt.ReqChanSize),
		writeChan:  make(chan *msgCtx, opt.ResChanSize),
		closeEvent: util.NewEvent(),
	}

	wg := util.WaitGroupWrapper{}

	// 读取消息
	wg.Wrap(ret.runRead, ret.opt.Logger.Errorf)
	// 写消息
	wg.Wrap(ret.runWrite, ret.opt.Logger.Errorf)

	return ret
}

func (s *frontendSession) Stats() State {
	return State{
		SID:           s.conn.GetSID(),
		NID:           s.nid,
		RemoteAddress: s.conn.RemoteAddr().String(),
	}
}

func (s *frontendSession) Push(ctx context.Context, cmd string, v interface{}) error {

	data, err := s.opt.Codec.Marshal(v)
	if err != nil {
		return err
	}

	if s.closeEvent.HasFired() {
		return constants.ErrSessionClosed
	}
	mctx := getMsgCtxWithContext(ctx, s.opt.EnableTraceDetail)

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

func (s *frontendSession) PushTimeout(ctx context.Context, cmd string, v interface{}, timeout time.Duration) error {

	if s.closeEvent.HasFired() {
		return constants.ErrSessionClosed
	}

	data, err := s.opt.Codec.Marshal(v)
	if err != nil {
		return err
	}

	mctx := getMsgCtxWithContext(ctx, s.opt.EnableTraceDetail)

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

func (s *frontendSession) PushImmediately(ctx context.Context, cmd string, v interface{}) error {

	if s.closeEvent.HasFired() {
		return constants.ErrSessionClosed
	}

	data, err := s.opt.Codec.Marshal(v)
	if err != nil {
		return err
	}

	mctx := getMsgCtxWithContext(ctx, s.opt.EnableTraceDetail)

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

func (s *frontendSession) OnClose(onClose OnClose) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onClose = onClose
}

func (s *frontendSession) Closed() bool {
	return s.closeEvent.HasFired()
}

// Close 关闭session 可多次调用
func (s *frontendSession) Close() {
	s.closeEvent.Fire()
}

func (s *frontendSession) Handle(srvHandler handle.IHandler) {

	tick := time.NewTicker(s.opt.KeepAlive)
	defer tick.Stop()

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

		// copy onClose
		var onClose OnClose
		s.mu.Lock()
		onClose = s.onClose
		s.mu.Unlock()

		if onClose != nil {
			onClose(s)
		}
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
						s.opt.Logger.Warningf("handle message encount custom error %v", customErr)
					} else {
						s.finishMsg(mctx, cErr)
						s.opt.Logger.Errorf("handle message encount a error %v", err)
						return
					}
				}
			}
		case <-tick.C: // 检查心跳
			{
				if !s.lastHeartBeat.IsZero() {

				}
			}
		}
	}

}

func (s *frontendSession) handleMsg(mctx *msgCtx, srvHandler handle.IHandler) error {

	// profile
	if s.opt.EnableTraceDetail {
		onMsgBegingHandle(mctx.ctx)
	}

	defer func() {
		if s.opt.EnableTraceDetail {
			onMsgEndHandle(mctx.ctx)
		}
	}()

	if mctx.ctx.Err() != nil { // 已经超时或取消
		return handle.NewCustomErrorWithError(mctx.ctx.Err())
	}

	if s.opt.MsgHandleTimeOut != 0 {
		mctx.WithTimeout(s.opt.MsgHandleTimeOut)
	}

	// 处理读取的消息
	switch mctx.msgRead.GetMsgType() {
	case message.NOTIFY:
		{
			reqMsg := mctx.msgRead

			serviceName, methodName, err := message.GetSrvMethod(reqMsg.GetService())
			if err != nil {
				return handle.NewCriticalErrorWithError(err)
			}

			err = srvHandler.Notify(mctx.ctx, s, serviceName, methodName, reqMsg.GetPayload())

			if err != nil {
				return err
			}
			s.finishMsg(mctx, nil)
		}
	case message.REQUEST:
		{
			reqMsg := mctx.msgRead

			serviceName, methodName, err := message.GetSrvMethod(reqMsg.GetService())
			if err != nil {
				return handle.NewCriticalErrorWithError(err)
			}

			resBuf, err := srvHandler.Call(mctx.ctx, s, serviceName, methodName, reqMsg.GetPayload())

			if err != nil {
				if customErr, ok := err.(handle.ICustomError); ok {
					// 非关键错误
					// 回复客户端
					mctx.msgWrite, err = message.MsgFactory.BuildResponseCustomErrorMessage(reqMsg.GetSequence(), customErr.Error())
					if err != nil {
						return err
					}
					err = s.sendWrite(mctx)
					if err != nil {
						return err
					}
				}
				return err
			}

			// 回复客户端
			mctx.msgWrite, err = message.MsgFactory.BuildResponseMessage(reqMsg.GetSequence(), resBuf)
			if err != nil {
				return err
			}
			err = s.sendWrite(mctx)
			if err != nil {
				return err
			}
		}
	case message.HEARTBEAT:
		{
			s.lastHeartBeat = time.Now()
			// write ack
			var err error
			mctx.msgWrite, err = message.MsgFactory.BuildHeatAckMessage()
			if err != nil {
				return err
			}

			err = s.sendWrite(mctx)
			if err != nil {
				return err
			}
		}
	case message.HANDSHAKE:
		{
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
		return handle.NewCriticalErrorf("invalid message type %v", mctx.msgRead)
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
		s.closeConn()
	}()

	cacheMsgNum := 0
	// 退出写协程条件:
	// 1.resChan关闭
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
					s.closeConn()
					return
				}
				cacheMsgNum++
				// need flush?
				if len(s.writeChan) == 0 || cacheMsgNum >= s.opt.MaxMsgCacheNum {
					if err = s.flush(); err != nil {
						s.opt.Logger.Errorf("flush message error %v", err)
						// io错误 关闭链接
						s.closeConn()
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
		err = s.conn.WriteNextMessage(mctx.msgWrite, s.opt.GetMessageOpt(mctx.msgWrite))
	} else {
		err = constants.ErrSessionClosed
	}
	s.connMu.Unlock()

	if err != nil {
		s.finishMsg(mctx, err)
		return err
	}
	s.finishMsg(mctx, nil)
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
	// 2.conn读取消息错误,一半为EOF
	// 3.链接已关闭
	for {
		msg, err := s.conn.ReadNextMessage()
		if err != nil {
			if err != io.EOF && err != io.ErrUnexpectedEOF && !strings.Contains(err.Error(), "use of closed network connection") {
				s.opt.Logger.Errorf("read message error %v", err)
			}
			// io err 关闭链接
			s.closeConn()
			return
		}

		mctx := getMsgCtx(s.opt.EnableTraceDetail)
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

func (s *frontendSession) closeConn() {
	s.connMu.Lock()
	defer s.connMu.Unlock()

	if s.connClosed {
		return
	}

	s.opt.Logger.Info("close conn")

	// 关闭链接
	err := s.conn.Close()
	if err != nil {
		s.opt.Logger.Errorf("session close error %v", err)
	}
	s.connClosed = true
}
