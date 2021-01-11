// Package node 服务节点
package node

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/liyiysng/scatter/internal/statsd"
	"github.com/liyiysng/scatter/internal/util"

	"github.com/liyiysng/scatter/internal/lg"
)

// 内部错误
type errStore struct {
	err error
}

// Node represent
type Node struct {
	sync.RWMutex
	// 客户端自增ID
	clientIDSequence int64
	// 选项
	opts atomic.Value
	// 当前错误(最后一个错误)
	errValue atomic.Value
	// 开启时间
	startTime time.Time

	// 所有客户端
	clientLock sync.RWMutex
	clients    map[int64]Client

	// tcp 监听
	tcpListener net.Listener

	// 退出chan
	exitChan chan int
	// 检测所有内部goroutine
	waitGroup util.WaitGroupWrapper
}

// NewNode 新建节点
func NewNode(opts *Options) (n *Node, err error) {
	// 缺省日志
	if opts.Logger == nil {
		opts.Logger = log.New(os.Stderr, opts.LogPrefix, log.Ldate|log.Ltime|log.Lmicroseconds)
	}

	n = &Node{
		startTime: time.Now(),
		clients:   make(map[int64]Client),
		exitChan:  make(chan int),
	}

	n.swapOpts(opts)
	n.errValue.Store(&errStore{})

	if opts.StatsdPrefix != "" {
		var port string
		_, port, err = net.SplitHostPort(opts.HTTPAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to parse HTTP address (%s) - %s", opts.HTTPAddress, err)
		}
		statsdHostKey := statsd.HostKey(net.JoinHostPort(opts.BroadcastAddress, port))
		prefixWithHost := strings.Replace(opts.StatsdPrefix, "%s", statsdHostKey, -1)
		if prefixWithHost[len(prefixWithHost)-1] != '.' {
			prefixWithHost += "."
		}
		opts.StatsdPrefix = prefixWithHost
	}

	n.logf(lg.INFO, "ID: %d", opts.ID)

	n.tcpListener, err = net.Listen("tcp", opts.TCPAddress)
	if err != nil {
		return nil, fmt.Errorf("listen (%s) failed - %s", opts.TCPAddress, err)
	}

	return
}

// SetHealth 设置当前节点最后一个错误 , 可重制错误
func (n *Node) SetHealth(err error) {
	n.errValue.Store(errStore{err: err})
}

// IsHealthy 当前节点是否健康
func (n *Node) IsHealthy() bool {
	return n.GetError() == nil
}

// GetError 返回当前节点的最后一个错误
func (n *Node) GetError() error {
	errValue := n.errValue.Load()
	return errValue.(errStore).err
}

// GetStartTime 获取创建时间
func (n *Node) GetStartTime() time.Time {
	return n.startTime
}

// GetHealth 获取当前节点健康状态信息
func (n *Node) GetHealth() string {
	err := n.GetError()
	if err != nil {
		return fmt.Sprintf("NOK - %s", err)
	}
	return "OK"
}

func (n *Node) swapOpts(opts *Options) {
	n.opts.Store(opts)
}

func (n *Node) getOpts() *Options {
	return n.opts.Load().(*Options)
}

func (n *Node) logf(level lg.LogLevel, f string, args ...interface{}) {
	opts := n.getOpts()
	lg.Logf(opts.Logger, opts.LogLevel, level, f, args...)
}
