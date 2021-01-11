package node

import "github.com/liyiysng/scatter/internal/lg"

// Options node 运行的所有选项
type Options struct {
	// 基础选项
	// ID 节点ID
	ID int64
	// 日志等级
	LogLevel lg.LogLevel
	// 日志前缀
	LogPrefix string
	// 日志实体
	Logger lg.Logger

	// 地址
	// http 地址
	HTTPAddress string
	// tcp 地址
	TCPAddress             string
	BroadcastAddress       string
	NSQLookupdTCPAddresses []string

	// 状态信息前缀 , 被当作key发送至statsd (%s for host replacement)
	StatsdPrefix string
}
