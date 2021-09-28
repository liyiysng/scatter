package conn

import (
	"bytes"
	"io"
	"net/http"
	"time"
	"github.com/gorilla/websocket"
	"github.com/liyiysng/scatter/node/message"
	"github.com/liyiysng/scatter/ratelimit"
	"github.com/liyiysng/scatter/util"
)

type wsConn struct {
	*websocket.Conn
	opt     MsgConnOption
	rdBuket *ratelimit.Bucket
	wrBuket *ratelimit.Bucket
	readTotal int64
	writeTotal int64
}

// NewWSConn 创建ws链接
func NewWSConn(w http.ResponseWriter, r *http.Request, opt MsgConnOption) (MsgConn, error) {

	upgrader := websocket.Upgrader{
		HandshakeTimeout: time.Second * 30,
		ReadBufferSize:   opt.ReadBufferSize,
		WriteBufferSize:  opt.WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	if opt.MaxLength > 0 {
		conn.SetReadLimit(int64(opt.MaxLength))
	}

	ret := &wsConn{
		Conn: conn,
		opt:  opt,
	}

	if opt.EnableLimit {
		ret.rdBuket = ratelimit.NewBucketWithQuantum(time.Second, opt.RateLimitReadBytes, opt.RateLimitReadBytes)
		ret.wrBuket = ratelimit.NewBucketWithQuantum(time.Second, opt.RateLimitWriteBytes, opt.RateLimitWriteBytes)
	}
	return ret, nil
}

func (c *wsConn) GetSID() int64 {
	return c.opt.SID
}

// 当前读取字节数总量
func (c *wsConn)GetCurrentReadTotalBytes() int64{
	return c.readTotal;
}
// 当前写字节数总量
func (c *wsConn)GetCurrentWirteTotalBytes() int64{
	return c.writeTotal;
}

func (c *wsConn) Flush() error {
	return nil
}

func (c *wsConn) ReadNextMessage() (msg message.Message, popt message.PacketOpt, err error) {
	//记录读取字节数
	rdCount := 0
	defer func() {
		if c.opt.ReadCountReport != nil {
			c.opt.ReadCountReport(c, rdCount)
		}
		c.readTotal += int64(rdCount)
	}()

	p := message.PackagePoolGet()
	defer message.PackagePoolPut(p)

	if c.opt.ReadTimeout != 0 {
		c.SetReadDeadline(time.Now().Add(c.opt.ReadTimeout))
	}

	_, buf, err := c.Conn.ReadMessage()

	if err != nil {
		if werr, ok := err.(*websocket.CloseError); ok {
			if werr.Code == websocket.CloseNormalClosure {
				return nil, message.DEFAULTPOPT, io.EOF
			}
			if werr.Code == websocket.CloseAbnormalClosure {
				return nil, message.DEFAULTPOPT, io.ErrUnexpectedEOF
			}
		}
		return nil, message.DEFAULTPOPT, err
	}
	if c.rdBuket != nil {
		c.rdBuket.Wait(int64(len(buf)))
	}

	rdCount, err = p.ReadFrom(bytes.NewBuffer(buf), c.opt.Compresser, c.opt.MaxLength)
	if err != nil {
		return nil, message.DEFAULTPOPT, err
	}

	msg, err = message.MsgFactory.BuildMessage(p.Data)
	if err != nil {
		return nil, message.DEFAULTPOPT, err
	}

	popt = p.PacketOpt

	return
}

// WriteNextMessage 发送一个消息
func (c *wsConn) WriteNextMessage(msg message.Message, popt message.PacketOpt) error {
	buf, err := msg.ToBytes()
	if err != nil {
		return err
	}
	if len(buf) > c.opt.MaxLength {
		return message.ErrMsgTooLager
	}
	p := message.PackagePoolGet()
	defer message.PackagePoolPut(p)
	p.PacketOpt = popt
	p.Data = buf

	writeBuffer := util.BufferPoolGet()
	defer util.BufferPoolPut(writeBuffer)

	_, err = p.WriteTo(writeBuffer, c.opt.Compresser, c.opt.MaxLength)
	if err != nil {
		return err
	}
	if c.opt.WriteTimeout != 0 {
		c.SetWriteDeadline(time.Now().Add(c.opt.WriteTimeout))
	}
	if c.wrBuket != nil {
		c.wrBuket.Wait(int64(len(writeBuffer.Bytes())))
	}
	err = c.WriteMessage(websocket.BinaryMessage, writeBuffer.Bytes())
	if err != nil {
		return err
	}

	wcount := len(writeBuffer.Bytes())
	if c.opt.WriteCountReport != nil {
		c.opt.WriteCountReport(c, wcount)
	}
	c.writeTotal += int64(wcount)

	return nil
}

func (c *wsConn) Close() error {
	return c.Conn.Close()
}
