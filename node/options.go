package node

import (
	"time"

	"github.com/liyiysng/scatter/logger"
)

const (
	defaultWriteBufSize = 32 * 1024
	defaultReadBufSize  = 32 * 1024
)

// Options node 运行的所有选项
type Options struct {
	// 基础选项
	// ID 节点ID
	ID string
	// 节点名
	Name string
	// 日志前缀
	LogPrefix string
	// 日志实体
	Logger logger.DepthLogger

	// 缓冲设置
	writeBufferSize int
	readBufferSize  int

	// 超时设置
	connectionTimeout time.Duration

	// trace
	// 允许事件跟踪
	enableEventTrace bool

	// 指标
	metricsEnable bool
	// 当前链接数
	connCountEnable bool
}

func (o *Options) connCountEnabled() bool {
	if !o.metricsEnable {
		return false
	}
	return o.connCountEnable
}

var defaultOptions = Options{
	connectionTimeout: 120 * time.Second,
	writeBufferSize:   defaultWriteBufSize,
	readBufferSize:    defaultReadBufSize,
}

// IOption 设置 日志等级等....
type IOption interface {
	apply(*Options)
}

// funcOption wraps a function that modifies IOption into an
// implementation of the IOption interface.
type funcOption struct {
	f func(*Options)
}

func (fdo *funcOption) apply(do *Options) {
	fdo.f(do)
}

func newFuncServerOption(f func(*Options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

// WriteBufferSize determines how much data can be batched before doing a write on the wire.
// The corresponding memory allocation for this buffer will be twice the size to keep syscalls low.
// The default value for this buffer is 32KB.
func WriteBufferSize(s int) IOption {
	return newFuncServerOption(func(o *Options) {
		o.writeBufferSize = s
	})
}
