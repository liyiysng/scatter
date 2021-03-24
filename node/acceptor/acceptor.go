package acceptor

import (
	"github.com/liyiysng/scatter/logger"
	"github.com/liyiysng/scatter/node/conn"
)

// Acceptor type interface
type Acceptor interface {
	ListenAndServe() error
	Stop()
	GetAddr() string
	GetOuterAddr() string
	GetConnChan() <-chan conn.MsgConn
}

// Option 选项
type Option struct {
	Addr       string
	OuterAddr  string
	CertFile   string
	KeyFile    string
	GetConnOpt func() conn.MsgConnOption
	Logger     logger.Logger
}
