// Package session 实现session相关
package session

import (
	"context"
	"time"

	"github.com/liyiysng/scatter/logger"
	"github.com/liyiysng/scatter/node/handle"
	"github.com/liyiysng/scatter/node/message"
)

var (
	myLog = logger.Component("session")
)

// State 表示session当前状态
type State struct {
	SID           int64  `json:"sid"`
	NID           int64  `json:"nid"`
	RemoteAddress string `json:"remote_address"`
}

// OnClose 关闭回调类型
type OnClose func(s Session)

// Session 表示一个客户端,可能是TCP/UDP/ws(s)/http(s) 的一次会话
type Session interface {
	GetSID() int64
	// 向客户端推送消息
	Push(ctx context.Context, cmd string, v interface{}, popt ...message.IPacketOption) error
	// 向客户端推送消息
	PushTimeout(ctx context.Context, cmd string, v interface{}, timeout time.Duration, popt ...message.IPacketOption) error
	// 向客户端推送,若发送缓冲已满则会返回 ErrorPushBufferFull
	PushImmediately(ctx context.Context, cmd string, v interface{}, popt ...message.IPacketOption) error
	// 关闭回调
	OnClose(onClose OnClose)
	// session是否关闭
	Closed() bool
}

// FrontendSession 前端Session
type FrontendSession interface {
	Session
	// GetNID 节点ID
	GetNID() int64
	// PeerAddr 对端地址
	PeerAddr() string
	// 消息处理
	Handle(srvHandler handle.IHandler)
	// 关闭session
	Close()
}
