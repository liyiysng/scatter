package conn

import (
	"bufio"
	"io"
	"net"
	"time"

	"github.com/liyiysng/scatter/constants"
	"github.com/liyiysng/scatter/node/message"
	"github.com/liyiysng/scatter/ratelimit"
)

type tcpConn struct {
	net.Conn
	bufIO *bufio.ReadWriter
	opt   MsgConnOption
}

// NewTCPMsgConn 创建一个tcp消息链接
func NewTCPMsgConn(conn net.Conn, opt MsgConnOption) MsgConn {

	var rd io.Reader = conn
	var wr io.Writer = conn

	if opt.EnableLimit {
		rdBuket := ratelimit.NewBucketWithQuantum(time.Second, opt.RateLimitReadBytes, opt.RateLimitReadBytes)
		wrBuket := ratelimit.NewBucketWithQuantum(time.Second, opt.RateLimitWriteBytes, opt.RateLimitWriteBytes)
		rd = ratelimit.Reader(rd, rdBuket)
		wr = ratelimit.Writer(wr, wrBuket)
	}

	return &tcpConn{
		Conn:  conn,
		bufIO: bufio.NewReadWriter(bufio.NewReaderSize(rd, opt.ReadBufferSize), bufio.NewWriterSize(wr, opt.WriteBufferSize)),
		opt:   opt,
	}
}

func (c *tcpConn) GetSID() int64 {
	return c.opt.SID
}

func (c *tcpConn) ReadNextMessage() (msg message.Message, err error) {

	//记录读取字节数
	rdCount := 0
	defer func() {
		if c.opt.ReadCountReport != nil {
			c.opt.ReadCountReport(c, rdCount)
		}
	}()

	p := message.PackagePoolGet()
	defer message.PackagePoolPut(p)

	if c.opt.ReadTimeout != 0 {
		c.SetReadDeadline(time.Now().Add(c.opt.ReadTimeout))
	}
	rdCount, err = p.ReadFrom(c.bufIO, c.opt.Compresser, c.opt.MaxLength)
	if err != nil {
		return nil, err
	}

	msg, err = message.BuildMessage(p.MsgType, p.Data)
	if err != nil {
		return nil, err
	}

	return
}

// WriteNextMessage 发送一个消息
func (c *tcpConn) WriteNextMessage(msg message.Message, msgOpt message.MsgOpt) error {
	buf, err := msg.ToBytes()
	if err != nil {
		return err
	}
	if len(buf) > c.opt.MaxLength {
		return constants.ErrMsgTooLager
	}
	p := message.PackagePoolGet()
	defer message.PackagePoolPut(p)
	p.MsgType = msg.GetMsgType()
	p.MsgOpt = msgOpt
	p.Data = buf

	if c.opt.WriteTimeout != 0 {
		c.SetWriteDeadline(time.Now().Add(c.opt.WriteTimeout))
	}
	wLen, err := p.WriteTo(c.bufIO, c.opt.Compresser, c.opt.MaxLength)
	if err != nil {
		return err
	}

	err = c.bufIO.Flush()
	if err != nil {
		return err
	}

	if c.opt.WriteCountReport != nil {
		c.opt.WriteCountReport(c, wLen)
	}

	return nil
}

func (c *tcpConn) Close() error {
	return c.Conn.Close()
}
