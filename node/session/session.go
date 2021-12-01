// Package session 实现session相关
package session

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	csession "github.com/liyiysng/scatter/cluster/subsrv"
	"github.com/liyiysng/scatter/constants"
	"github.com/liyiysng/scatter/handle"
	"github.com/liyiysng/scatter/logger"
	"github.com/liyiysng/scatter/node/message"
)

var (
	myLog = logger.Component("session")
)

var (
	// ErrorPushBufferFull session push缓冲已满
	ErrorPushBufferFull = errors.New("session push full")
	// ErrorMsgDiscard 消息被丢弃
	ErrorMsgDiscard = errors.New("message discard")
)

type SessionAttr string

const (
	SessionPeerInfo SessionAttr = "__session_peer_info__"
)

func newSessionClosedError(sid int64, why string) error {
	return &sessionClosedError{
		sid: sid,
		why: why,
	}
}

type sessionClosedError struct {
	why string
	sid int64
}

func (s *sessionClosedError) Error() string {
	return fmt.Sprintf("session %d closed : %s", s.sid, s.why)
}

// OnClose 关闭回调类型
type OnClose func(s Session)

// ISessionInfo session信息
type ISessionInfo interface {
	GetSID() int64
	// GetNID 节点ID
	GetNID() int64
	// PeerAddr 对端地址
	PeerAddr() net.Addr
}

// Session 表示一个客户端,可能是TCP/UDP/ws(s)/http(s) 的一次会话
type Session interface {
	ISessionInfo
	csession.Session

	// 向客户端推送消息
	Push(ctx context.Context, cmd string, v interface{}, popt ...message.IPacketOption) error
	// 向客户端推送消息
	PushTimeout(ctx context.Context, cmd string, v interface{}, timeout time.Duration, popt ...message.IPacketOption) error
	// 向客户端推送,若发送缓冲已满则会返回 ErrorPushBufferFull
	PushImmediately(ctx context.Context, cmd string, v interface{}, popt ...message.IPacketOption) error
	// 关闭回调
	SetOnClose(onClose OnClose)
	// session是否关闭
	Closed() bool
	// Kick 提出 , 防止服务器出现过多TIME_WAIT状态,使用Kick断开链接
	Kick() error
	// 关闭session , 强制关闭session , 服务器主动断开链接
	Close()

	// context
	GetCtx() context.Context

	// attr 元数据
	SetAttr(key SessionAttr, v interface{})
	GetAttr(key SessionAttr) (v interface{}, ok bool)

	// uid
	BindUID(uid constants.UID)

	// 获取处理
	// cb协程安全
	// 切勿再cb内访问session 或 执行较长时间得阻塞任务
	SetOnHandle(cb func(ss Session, msgRead message.Message, h handle.IHandler, err error))

	// 设置push
	// cb协程安全
	// 切勿再cb内访问session 或 执行较长时间得阻塞任务
	SetOnPush(cb func(ss Session, msgPush message.Message, err error))

	// ticker job
	// 任务将在单独的协程中执行
	// 当session结束时,所有的定时任务自动清理
	// AddTicker 添加定时任务
	AddTicker(name string, duration time.Duration, delay time.Duration, cb func(t time.Time)) (err error)
	// RemoveTicker 移除定时任务
	RemoveTicker(name string) (err error)
}

// IFrontendSession 前端session
type IFrontendSession interface {
	Session
	// 消息处理
	Handle(srvHandler handle.IHandler)
}
