package acceptor

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/liyiysng/scatter/node/conn"
)

// WSAcceptor struct
type WSAcceptor struct {
	connChan  chan conn.MsgConn
	listener  net.Listener
	opt       Option
	closeOnce sync.Once
	closed    bool
}

// NewWSAcceptor creates a new instance of websocket acceptor
func NewWSAcceptor(opt Option) Acceptor {
	return &WSAcceptor{
		connChan: make(chan conn.MsgConn),
		opt:      opt,
	}
}

// GetAddr returns the addr the acceptor will listen on
func (w *WSAcceptor) GetAddr() string {
	if w.listener != nil {
		return w.listener.Addr().String()
	}
	return ""
}

// GetConnChan gets a connection channel
func (w *WSAcceptor) GetConnChan() <-chan conn.MsgConn {
	return w.connChan
}

func (w *WSAcceptor) hasTLSCertificates() bool {
	return w.opt.CertFile != "" && w.opt.KeyFile != ""
}

// ListenAndServe listens and serve in the specified addr
func (w *WSAcceptor) ListenAndServe() error {

	listener, err := net.Listen("tcp", w.opt.Addr)
	if err != nil {
		return fmt.Errorf("Failed to listen: %s", err.Error())
	}
	w.listener = listener

	defer w.Stop()
	defer close(w.connChan)

	w.listener = listener

	if w.hasTLSCertificates() {
		return http.ServeTLS(w.listener, w, w.opt.CertFile, w.opt.KeyFile)
	}

	return http.Serve(w.listener, w)
}

// Stop stops the acceptor
func (w *WSAcceptor) Stop() {
	w.closeOnce.Do(func() {
		if w.listener != nil {
			err := w.listener.Close()
			if err != nil {
				w.opt.Logger.Errorf("Failed to stop: %s", err.Error())
			}
		}
		w.closed = true
		w.opt.Logger.Info("stop ws listen: %s", w.opt.Addr)
	})
}

func (w *WSAcceptor) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	conn, err := conn.NewWSConn(rw, r, w.opt.GetConnOpt())
	if err != nil {
		w.opt.Logger.Errorf("[WSAcceptor.ServeHTTP] new NewWSConn err %s", err.Error())
		return
	}

	if w.closed {
		w.opt.Logger.Warningf("[WSAcceptor.ServeHTTP] acceptor closed.")
		_ = conn.Close() // ignor error
		return
	}

	w.connChan <- conn
}
