package acceptor

import (
	"crypto/tls"
	"fmt"
	"net"
	"sync"

	"github.com/liyiysng/scatter/node/conn"
)

// TCPAcceptor struct
type TCPAcceptor struct {
	connChan  chan conn.MsgConn
	listener  net.Listener
	running   bool
	opt       Option
	closeOnce sync.Once
}

// NewTCPAcceptor creates a new instance of tcp acceptor
func NewTCPAcceptor(opt Option) Acceptor {

	return &TCPAcceptor{
		connChan: make(chan conn.MsgConn),
		running:  false,
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
			a.running = false
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
	a.running = true
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
	a.running = true

	a.opt.Logger.Infof("tcp tls listen: %s", a.opt.Addr)

	a.serve()
	return nil
}

func (a *TCPAcceptor) serve() {
	defer a.Stop()
	defer close(a.connChan)
	for a.running {
		c, err := a.listener.Accept()
		if err != nil {
			a.opt.Logger.Errorf("Failed to accept TCP connection: %s", err.Error())
			break
		}
		a.connChan <- conn.NewTCPMsgConn(c, a.opt.GetConnOpt())
	}
}
