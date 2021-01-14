// Package node 服务节点
package node

import (
	"errors"
	"fmt"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/liyiysng/scatter/logger"
	"github.com/liyiysng/scatter/metrics"
	"github.com/liyiysng/scatter/util"
	"golang.org/x/net/trace"
)

// ErrNodeStopped 节点已经停止
var ErrNodeStopped = errors.New("node: the node has been stopped")

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
	opts Options
	// 当前错误(最后一个错误)
	errValue atomic.Value
	// 开启时间
	startTime time.Time

	// 所有链接
	conns map[net.Conn]bool

	// 监听对象
	lis map[net.Listener]bool

	// 退出chan
	quit *util.Event
	// 检测所有内部goroutine
	waitGroup util.WaitGroupWrapper
	// trace
	trEvents trace.EventLog

	// 标识是否启动
	started bool
}

// NewNode 新建节点
func NewNode(opt ...IOption) (n *Node, err error) {

	opts := defaultOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	// 缺省日志
	if opts.Logger == nil {
		opts.Logger = logger.GDepthLogger
	}

	n = &Node{
		lis:       make(map[net.Listener]bool),
		opts:      opts,
		startTime: time.Now(),
		conns:     make(map[net.Conn]bool),
		quit:      util.NewEvent(),
	}

	if opts.enableEventTrace {
		_, file, line, _ := runtime.Caller(1)
		n.trEvents = trace.NewEventLog("scatter.Node", fmt.Sprintf("%s-%s:%d", opts.ID, file, line))
	}

	n.errValue.Store(&errStore{})

	opts.Logger.Infof("start node %s", opts.ID)

	return
}

// Serve 在该lis上启动一个Serve
// Serve 除Stop或者 被调用之外,都返回一个非nil错误,
func (n *Node) Serve(lis net.Listener) error {
	n.Lock()

	n.trEventLogf("starting")
	n.started = true

	if n.lis == nil {
		// Start called after Stop
		n.Unlock()
		lis.Close()
		return ErrNodeStopped
	}
	n.lis[lis] = true

	n.Unlock()

	defer func() {
		n.Lock()
		if n.lis != nil && n.lis[lis] {
			lis.Close()
			delete(n.lis, lis)
		}
		n.Unlock()
	}()

	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		rawConn, err := lis.Accept()
		if err != nil {
			if ne, ok := err.(interface {
				Temporary() bool
			}); ok && ne.Temporary() { // 临时错误
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				n.Lock()
				n.trEventLogf("Accept error: %v; retrying in %v", err, tempDelay)
				n.Unlock()
				timer := time.NewTimer(tempDelay)
				select {
				case <-timer.C:
				case <-n.quit.Done():
					timer.Stop()
					return nil
				}
				continue
			}
			// 不是临时错误/尝试失败
			n.Lock()
			n.trEventLogf("done; Accept = %v", err)
			n.Unlock()

			if n.quit.HasFired() {
				return nil
			}
			return err
		}
		tempDelay = 0

		n.waitGroup.Wrap(func() {
			n.handleRawConn(rawConn)
		}, nil)
	}
}

// Stop 停止/关闭该节点
func (n *Node) Stop() {
	n.quit.Fire()

	defer func() {
		n.waitGroup.Wait()
	}()

	n.Lock()
	listeners := n.lis
	n.lis = nil
	conns := n.conns
	n.conns = nil
	n.Unlock()

	for lis := range listeners {
		lis.Close()
	}
	for c := range conns {
		c.Close()
	}

	n.Lock()
	if n.trEvents != nil {
		n.trEvents.Finish()
		n.trEvents = nil
	}
	n.Unlock()
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

// 记录事件日志
// 需求:已经持有锁
func (n *Node) trEventLogf(format string, a ...interface{}) {
	if n.trEvents != nil {
		n.trEvents.Printf(format, a...)
	}
}

// 记录事件错误日志
// 需求:已经持有锁
func (n *Node) trEventErrorf(format string, a ...interface{}) {
	if n.trEvents != nil {
		n.trEvents.Errorf(format, a...)
	}
}

func (n *Node) handleRawConn(rawConn net.Conn) {
	if n.quit.HasFired() {
		rawConn.Close()
		return
	}

	if n.opts.connCountEnabled() { // 链接数
		metrics.ReportNodeNewCommingConn(n.opts.ID)
		defer metrics.ReportNodeNewCommingConn(n.opts.ID)
	}
}

func (n *Node) addConn(conn net.Conn) bool {
	n.Lock()
	defer n.Unlock()

	if n.conns == nil { //called after Stop
		conn.Close()
		return false
	}

	n.conns[conn] = true
	return true
}

func (n *Node) removeConn(conn net.Conn) {
	n.Lock()
	defer n.Unlock()
	if n.conns != nil {
		delete(n.conns, conn)
	}
}
