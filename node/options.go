package node

import (
	"fmt"
	"reflect"
	"time"

	"github.com/liyiysng/scatter/encoding"
	"github.com/liyiysng/scatter/handle"
	"github.com/liyiysng/scatter/logger"
	"github.com/liyiysng/scatter/metrics"
	"github.com/liyiysng/scatter/node/textlog"
	"github.com/olivere/elastic/v7"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	//json编码
	_ "github.com/liyiysng/scatter/encoding/json"
	//proto编码
	_ "github.com/liyiysng/scatter/encoding/proto"
	//gzip压缩
	_ "github.com/liyiysng/scatter/encoding/gzip"
	//snappy压缩
	_ "github.com/liyiysng/scatter/encoding/snappy"
)

var typeProtoMessage = reflect.TypeOf((*proto.Message)(nil)).Elem()

const (
	defaultWriteBufSize = 32 * 1024
	defaultReadBufSize  = 32 * 1024
)

// Options node 运行的所有选项
type Options struct {
	// 基础选项
	// ID 节点ID
	ID int64
	// 节点名
	Name string
	// 日志前缀
	LogPrefix string
	// 日志实体
	Logger logger.DepthLogger
	// 显示处理日志
	showHandleLog bool

	// 文件日志
	// 文件详情记录
	needFileTextLog bool
	// es详情记录
	needEsTextLog bool
	// 消息详情日志
	textLogWriter []textlog.Sink

	// 链接设置
	// 缓冲设置
	writeBufferSize int
	readBufferSize  int
	// 数据最大长度
	maxPayloadLength int
	// 超时设置
	// 链接超时(当链接创建后多久事件未接受到handshake消息)
	// 次值不能为0
	connectionTimeout time.Duration
	// 读写超时
	readTimeout  time.Duration
	writeTimeout time.Duration
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
	// 限制单个链接每秒消息处理个数
	// <=0 不做限制
	rateLimitMsgProcNum int64
	// 读chan缓冲
	readChanBufSize int
	// 写chan缓冲
	writeChanBufSize int

	// 消息设置

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
	// 消息处理延时
	metricsMsgProcDelayEnable bool
	// repoters
	metricsReporters []metrics.Reporter

	// rpc
	// 請求類型驗證
	reqTypeValidator func(reqType reflect.Type) error
	// 回复類型驗證
	resTypeValidator func(reqType reflect.Type) error
	// 可选参数
	optArgs *handle.OptionalArgs

	// 配置错误
	lastError error

	// 节点停止之后执行
	afterStop []func()

	// sub servive validator
	subSrvValidator func(srvName string) bool

	// grpc 配置
	grpcOpts []grpc.ServerOption
}

func (o *Options) validate() error {
	if (o.needFileTextLog || o.needEsTextLog) && !o.enableTraceDetail {
		return fmt.Errorf("want text log , but trace detail was disabled")
	}

	if c := encoding.GetCompressor(o.compresser); c == nil {
		return fmt.Errorf("cmpressor %q not support", o.compresser)
	}

	return nil
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

func (o *Options) metricsMsgProcDelayEnabled() bool {
	if !o.metricsEnable {
		return false
	}
	return o.metricsMsgProcDelayEnable
}

func (o *Options) getCompressor() encoding.Compressor {
	if o.compresser != "" {
		return encoding.GetCompressor(o.compresser)
	}
	return encoding.GetCompressor("gzip")
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
	codec:             "proto",
	compresser:        "gzip",
	enableLimit:       false,
	enableTraceDetail: true,
	showHandleLog:     true,
	reqTypeValidator: func(reqType reflect.Type) error {
		if reqType.Implements(typeProtoMessage) {
			return nil
		}
		return handle.ErrRequstTypeError
	},
	resTypeValidator: func(reqType reflect.Type) error {
		if reqType.Implements(typeProtoMessage) {
			return nil
		}
		return handle.ErrResponseTypeError
	},
	readChanBufSize:  1024,
	writeChanBufSize: 1024,
	subSrvValidator:  func(srvName string) bool { return true },
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

// NOptWriteBufferSize determines how much data can be batched before doing a write on the wire.
// The corresponding memory allocation for this buffer will be twice the size to keep syscalls low.
// The default value for this buffer is 32KB.
func NOptWriteBufferSize(s int) IOption {
	return newFuncServerOption(func(o *Options) {
		if o.lastError != nil {
			return
		}
		o.writeBufferSize = s
	})
}

// NOptEnableEnableFileTextLog 开启文件文本日志
func NOptEnableEnableFileTextLog() IOption {
	return newFuncServerOption(func(o *Options) {
		if o.needFileTextLog {
			return
		}
		if o.lastError != nil {
			return
		}
		sink, err := textlog.NewTempFileSink()
		if err != nil {
			o.lastError = err
			return
		}
		o.needFileTextLog = true
		o.enableTraceDetail = true
		o.textLogWriter = append(o.textLogWriter, sink)
	})
}

// NOptEnableEnableEsTextLog 开启es文本日志
func NOptEnableEnableEsTextLog(index string, numBulk int, flushInverval time.Duration, options ...elastic.ClientOptionFunc) IOption {
	return newFuncServerOption(func(o *Options) {
		if o.needEsTextLog {
			return
		}
		if o.lastError != nil {
			return
		}
		sink, err := textlog.NewEsSink(index, numBulk, flushInverval, options...)
		if err != nil {
			o.lastError = err
			return
		}
		o.needEsTextLog = true
		o.enableTraceDetail = true
		o.textLogWriter = append(o.textLogWriter, sink)
	})
}

// NOptWithOptArgs rpc可选参数配置
func NOptWithOptArgs(optArgs *handle.OptionalArgs) IOption {
	if optArgs == nil {
		panic("[WithOptArgs] nil param")
	}
	return newFuncServerOption(func(o *Options) {
		if o.lastError != nil {
			return
		}
		o.optArgs = optArgs
	})
}

// NOptShowHandleLog 是否显示处理日志
func NOptShowHandleLog(show bool) IOption {
	return newFuncServerOption(func(o *Options) {
		if o.lastError != nil {
			return
		}
		o.showHandleLog = show
	})
}

// NOptTraceDetail 是否监视细节
func NOptTraceDetail(trace bool) IOption {
	return newFuncServerOption(func(o *Options) {
		if o.lastError != nil {
			return
		}
		o.enableTraceDetail = trace
	})
}

// NOptCompress 压缩算法
func NOptCompress(c string) IOption {
	return newFuncServerOption(func(o *Options) {
		if o.lastError != nil {
			return
		}
		o.compresser = c
	})
}

// NOptMetricsReporter 指标提交
func NOptMetricsReporter(r metrics.Reporter) IOption {
	return newFuncServerOption(func(o *Options) {
		if o.lastError != nil {
			return
		}
		o.metricsReporters = append(o.metricsReporters, r)
	})
}

// NOptEnableMetrics 指标监视
func NOptEnableMetrics(enable bool) IOption {
	return newFuncServerOption(func(o *Options) {
		if o.lastError != nil {
			return
		}
		o.metricsEnable = enable
		o.metricsConnCountEnable = enable
		o.metricsReadWriteBytesCountEnable = enable
		o.metricsMsgProcDelayEnable = enable
	})
}

// NOptNodeName 指标监视
func NOptNodeName(nname string) IOption {
	return newFuncServerOption(func(o *Options) {
		if o.lastError != nil {
			return
		}
		o.Name = nname
	})
}

// NOptProcMsgNumRateLimit 单个链接每秒消息处理个数限制
func NOptProcMsgNumRateLimit(num int64) IOption {
	return newFuncServerOption(func(o *Options) {
		if o.lastError != nil {
			return
		}
		o.rateLimitMsgProcNum = num
	})
}

// NOptAfterStop 当节点停止后调用
func NOptAfterStop(f ...func()) IOption {
	return newFuncServerOption(func(o *Options) {
		if o.lastError != nil {
			return
		}
		o.afterStop = append(o.afterStop, f...)
	})
}

// NOptWithSubSrvValidator 子服务验证
func NOptWithSubSrvValidator(f func(srvName string) bool) IOption {
	return newFuncServerOption(func(o *Options) {
		if o.lastError != nil {
			return
		}
		o.subSrvValidator = f
	})
}

// NOptWithGrpcOpts grpc配置
func NOptWithGrpcOpts(gopt ...grpc.ServerOption) IOption {
	return newFuncServerOption(func(o *Options) {
		o.grpcOpts = append(o.grpcOpts, gopt...)
	})
}

// NodeServeOption 开启服务选项
type NodeServeOption struct {
	outerAddr string
	certFile  string
	keyFile   string
}

// IGrpcClientOpt 客户端选项
type INodeServeOption interface {
	apply(*NodeServeOption)
}

type funcNodeServeOption struct {
	f func(*NodeServeOption)
}

func (fdo *funcNodeServeOption) apply(do *NodeServeOption) {
	fdo.f(do)
}

func newFuncNodeServeOption(f func(*NodeServeOption)) INodeServeOption {
	return &funcNodeServeOption{
		f: f,
	}
}

// OptNodeServeOptionWithOuterAddr 外部地址
func OptNodeServeOptionWithOuterAddr(outerAddr string) INodeServeOption {
	return newFuncNodeServeOption(func(o *NodeServeOption) {
		o.outerAddr = outerAddr
	})
}

// OptNodeServeOptionWithOuterAddr 外部地址
func OptNodeServeOptionWithCert(certFile, keyFile string) INodeServeOption {
	return newFuncNodeServeOption(func(o *NodeServeOption) {
		o.certFile = certFile
		o.keyFile = keyFile
	})
}
