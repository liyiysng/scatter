package conn

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/liyiysng/scatter/constants"
	"github.com/liyiysng/scatter/node/message"
	"github.com/liyiysng/scatter/ratelimit"
	"github.com/liyiysng/scatter/util"
)

type wsConn struct {
	*websocket.Conn
	opt     MsgConnOption
	rdBuket *ratelimit.Bucket
	wrBuket *ratelimit.Bucket
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

func (c *wsConn) Flush() error {
	return nil
}

func (c *wsConn) ReadNextMessage() (msg message.Message, err error) {
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

	_, buf, err := c.Conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	if c.rdBuket != nil {
		c.rdBuket.Wait(int64(len(buf)))
	}

	rdCount, err = p.ReadFrom(bytes.NewBuffer(buf), c.opt.Compresser, c.opt.MaxLength)
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
func (c *wsConn) WriteNextMessage(msg message.Message, msgOpt message.MsgOpt) error {
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

	writeBuffer := util.BufferPoolGet()
	defer util.BufferPoolPut(writeBuffer)

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

	if c.opt.WriteCountReport != nil {
		c.opt.WriteCountReport(c, len(writeBuffer.Bytes()))
	}

	return nil
}

func (c *wsConn) Close() error {
	return c.Conn.Close()
}
