// Package session 实现session相关
package session

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/liyiysng/scatter/handle"
	"github.com/liyiysng/scatter/logger"
	"github.com/liyiysng/scatter/node/message"
)

var (
	myLog = logger.Component("session")
)

var (
	// ErrSessionClosed session已关闭
	ErrSessionClosed = errors.New("session closed")
	// ErrorPushBufferFull session push缓冲已满
	ErrorPushBufferFull = errors.New("session push full")
	// ErrorMsgDiscard 消息被丢弃
	ErrorMsgDiscard = errors.New("message discard")
)

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
	SetAttr(key string, v interface{})
	GetAttr(key string) (v interface{}, ok bool)

	// uid
	//!
	GetUID() interface{}
	//!
	BindUID(uid interface{})
	//!
	IsUIDBind() bool
}

// IFrontendSession 前端session
type IFrontendSession interface {
	Session
	// 消息处理
	Handle(srvHandler handle.IHandler)
}
