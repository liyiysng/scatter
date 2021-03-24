package acceptor

import (
	"crypto/tls"
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

// GetOuterAddr outer addr
func (a *TCPAcceptor) GetOuterAddr() string {
	return a.opt.OuterAddr
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
			if a.listener != nil {
				a.listener.Close()
			}
			a.opt.Logger.Infof("stop listen: %s", a.opt.Addr)
		},
	)
}

func (a *TCPAcceptor) hasTLSCertificates() bool {
	return a.opt.CertFile != "" && a.opt.KeyFile != ""
}

// ListenAndServe using tcp acceptor
func (a *TCPAcceptor) ListenAndServe() error {
	return a.serve()
}

// ListenAndServeTLS listens using tls
func (a *TCPAcceptor) ListenAndServeTLS(cert, key string) error {
	a.opt.CertFile = cert
	a.opt.KeyFile = key
	return a.serve()
}

func (a *TCPAcceptor) serve() error {
	defer a.Stop()
	defer close(a.connChan)

	if a.hasTLSCertificates() {
		crt, err := tls.LoadX509KeyPair(a.opt.CertFile, a.opt.KeyFile)
		if err != nil {
			return err
		}

		tlsCfg := &tls.Config{Certificates: []tls.Certificate{crt}}

		listener, err := tls.Listen("tcp", a.opt.Addr, tlsCfg)
		if err != nil {
			return err
		}
		a.listener = listener
	} else {
		listener, err := net.Listen("tcp", a.opt.Addr)
		if err != nil {
			return err
		}
		a.listener = listener
	}

	a.opt.Logger.Infof("tcp listen: %s", a.opt.Addr)

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
					return nil
				}
				continue
			}
			// 不是临时错误/尝试失败
			if a.quit.HasFired() {
				return nil
			}
			return err
		}

		if a.quit.HasFired() {
			return nil
		}
		a.connChan <- conn.NewTCPMsgConn(c, a.opt.GetConnOpt())
	}
}
