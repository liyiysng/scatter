// Package session 实现session相关
package session

import (
	"context"

	"google.golang.org/protobuf/proto"
)

// State 表示session当前状态
type State struct {
	SID           string `json:"sid"`
	Hostname      string `json:"hostname"`
	Version       string `json:"version"`
	RemoteAddress string `json:"remote_address"`
}

// OnClose 关闭回调类型
type OnClose func(s Session)

// Session 表示一个客户端,可能是TCP/UDP/ws(s)/http(s)/
type Session interface {
	// Stats 获取客户端状态
	Stats() State
	// 向客户端推送消息
	Push(ctx context.Context, cmd string, msg proto.Message) error
	// 关闭回调
	OnClose(onClose OnClose)
	// session是否关闭
	Closed() bool
}

// BackendSession 标识在内部节点session
// BackendSession 会在各个节点同步Context
type BackendSession interface {
	Session
	// 绑定上下文
	BindContext(ctx context.Context, key string, data proto.Message) error
	// 获取上下文
	GetContext(key string) (exists bool, data proto.Message)
}
