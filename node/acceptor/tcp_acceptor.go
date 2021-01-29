package acceptor

import (
	"crypto/tls"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/liyiysng/scatter/node/conn"
	"github.com/liyiysng/scatter/util"
)

// TCPAcceptor struct
type TCPAcceptor struct {
	connChan chan conn.MsgConn
	listener net.Listener
	// 退出chan
	quit      *util.Event
	opt       Option
	closeOnce sync.Once
}

// NewTCPAcceptor creates a new instance of tcp acceptor
func NewTCPAcceptor(opt Option) Acceptor {

	return &TCPAcceptor{
		connChan: make(chan conn.MsgConn),
		quit:     util.NewEvent(),
		opt:      opt,
	}
}

// GetAddr returns the addr the acceptor will listen on
func (a *TCPAcceptor) GetAddr() string {
	if a.listener != nil {
		return a.listener.Addr().String()
	}
	return ""
}

// GetConnChan gets a connection channel
func (a *TCPAcceptor) GetConnChan() <-chan conn.MsgConn {
	return a.connChan
}

// Stop stops the acceptor
func (a *TCPAcceptor) Stop() {
	a.closeOnce.Do(
		func() {
			a.quit.Fire()
			a.listener.Close()
			a.opt.Logger.Infof("stop listen: %s", a.opt.Addr)
		},
	)
}

func (a *TCPAcceptor) hasTLSCertificates() bool {
	return a.opt.CertFile != "" && a.opt.KeyFile != ""
}

// ListenAndServe using tcp acceptor
func (a *TCPAcceptor) ListenAndServe() error {
	if a.hasTLSCertificates() {
		return a.ListenAndServeTLS(a.opt.CertFile, a.opt.KeyFile)
	}

	listener, err := net.Listen("tcp", a.opt.Addr)
	if err != nil {
		return fmt.Errorf("[TCPAcceptor.ListenAndServe] Failed to listen: %s", err.Error())
	}
	a.listener = listener
	a.opt.Logger.Infof("tcp listen: %s", a.opt.Addr)

	a.serve()
	return nil
}

// ListenAndServeTLS listens using tls
func (a *TCPAcceptor) ListenAndServeTLS(cert, key string) error {
	crt, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return fmt.Errorf("[TCPAcceptor.ListenAndServeTLS] Failed to listen: %s", err.Error())
	}

	tlsCfg := &tls.Config{Certificates: []tls.Certificate{crt}}

	listener, err := tls.Listen("tcp", a.opt.Addr, tlsCfg)
	a.listener = listener

	a.opt.Logger.Infof("tcp tls listen: %s", a.opt.Addr)

	a.serve()
	return nil
}

func (a *TCPAcceptor) serve() {
	defer a.Stop()
	defer close(a.connChan)

	var tempDelay time.Duration // how long to sleep on accept failure

	for {
		c, err := a.listener.Accept()
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
				a.opt.Logger.Warningf("Accept error: %v; retrying in %v", err, tempDelay)
				timer := time.NewTimer(tempDelay)
				select {
				case <-timer.C:
				case <-a.quit.Done():
					timer.Stop()
					return
				}
				continue
			}
			// 不是临时错误/尝试失败
			if a.quit.HasFired() {
				return
			}
			a.opt.Logger.Errorf("done; Accept = %v", err)
			return
		}

		if a.quit.HasFired() {
			return
		}
		a.connChan <- conn.NewTCPMsgConn(c, a.opt.GetConnOpt())
	}
}
