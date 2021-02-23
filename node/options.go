package node

import (
	"time"

	"github.com/liyiysng/scatter/encoding"
	"github.com/liyiysng/scatter/logger"
	"github.com/liyiysng/scatter/node/textlog"

	//json编码
	_ "github.com/liyiysng/scatter/encoding/json"
	//proto编码
	_ "github.com/liyiysng/scatter/encoding/proto"
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
	// 显示处理日志
	showHandleLog bool

	// 文件日志
	// 是否记录消息详情
	needTextLog bool
	// 消息详情文件路径
	textLogWriter textlog.Sink

	// 链接设置
	// 缓冲设置
	writeBufferSize int
	readBufferSize  int
	// 数据最大长度
	maxPayloadLength int
	// 超时设置
	connectionTimeout time.Duration
	readTimeout       time.Duration
	writeTimeout      time.Duration
	// 压缩
	compresser string
	// 编码
	codec string
	// 限流
	enableLimit bool
	// 读限流
	rateLimitReadBytes int64
	// 写限流
	rateLimitWriteBytes int64

	// trace
	// 允许事件跟踪
	enableEventTrace bool
	// 监视详情
	enableTraceDetail bool

	// 指标
	metricsEnable bool
	// 当前链接数
	metricsConnCountEnable bool
	// 读/写字节数
	metricsReadWriteBytesCountEnable bool
}

func (o *Options) metricsConnCountEnabled() bool {
	if !o.metricsEnable {
		return false
	}
	return o.metricsConnCountEnable
}

func (o *Options) metricsReadWriteBytesCountEnabled() bool {
	if !o.metricsEnable {
		return false
	}
	return o.metricsReadWriteBytesCountEnable
}

func (o *Options) getCompressor() encoding.Compressor {
	if o.compresser != "" {
		return encoding.GetCompressor(o.compresser)
	}
	return nil
}

func (o *Options) getCodec() encoding.Codec {
	if o.codec != "" {
		return encoding.GetCodec(o.codec)
	}
	return encoding.GetCodec("proto")
}

var defaultOptions = Options{
	writeBufferSize:   defaultWriteBufSize,
	readBufferSize:    defaultReadBufSize,
	maxPayloadLength:  32 * 1024,
	connectionTimeout: 120 * time.Second,
	readTimeout:       0,
	writeTimeout:      time.Second * 5,
	compresser:        "gzip",
	enableLimit:       false,
	enableTraceDetail: true,
	showHandleLog:     true,
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

// EnableTextLog 开启文本日志
func EnableTextLog(sink textlog.Sink) IOption {
	if sink == nil {
		panic("[EnableTextLog] nil sink")
	}
	return newFuncServerOption(func(o *Options) {
		o.needTextLog = true
		o.textLogWriter = sink
	})
}
