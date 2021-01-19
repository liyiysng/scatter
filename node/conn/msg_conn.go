// Package conn 客户端链接
package conn

import (
	"net"

	"github.com/liyiysng/scatter/logger"
)

var (
	myLog = logger.Component("node.acceptor")
)

// MsgConn 表示面向消息的链接
type MsgConn interface {
	net.Conn
	// ReadNextMessage 获取下一个消息请求
	ReadNextMessage() (data []byte, err error)
	// SetReadLimit 设置单个消息最大长度 , 包含消息头长度
	SetReadLimit(limit int64)
	// WriteNextMessage 发送一个消息
	WriteNextMessage(data []byte) error
}
