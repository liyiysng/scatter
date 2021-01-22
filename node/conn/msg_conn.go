// Package conn 客户端链接
package conn

import (
	"net"
	"time"

	"github.com/liyiysng/scatter/encoding"
	"github.com/liyiysng/scatter/logger"
	"github.com/liyiysng/scatter/node/message"
)

var (
	myLog = logger.Component("node.acceptor")
)

// MsgConnOption 消息链接选项
type MsgConnOption struct {
	// 消息最大长度(不包含头长度)
	MaxLength int
	// 读超时 zero表示永不超时
	ReadTimeout time.Duration
	// 写超时 zero表示永不超时
	WriteTimeout time.Duration
	// 读缓冲大小
	ReadBufferSize int
	// 写缓冲大小
	WriteBufferSize int
	// 解压/压缩
	Compresser encoding.Compressor
	// 限流
	EnableLimit bool
	// 每秒读取字节数
	RateLimitReadBytes int64
	// 每秒写取字节数
	RateLimitWriteBytes int64
	// 报告读字节数
	ReadCountReport func(byteCount int)
	// 报告写字节数
	WriteCountReport func(byteCount int)
}

// MsgConn 表示面向消息的链接
type MsgConn interface {
	// 对端地址
	RemoteAddr() net.Addr
	// 本地地址
	LocalAddr() net.Addr
	// ReadNextMessage 获取下一个消息请求
	ReadNextMessage() (msg message.Message, err error)
	// WriteNextMessage 发送一个消息
	WriteNextMessage(msg message.Message, msgOpt message.MsgOpt) error
	// 关闭链接
	Close() error
}
