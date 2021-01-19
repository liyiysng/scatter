package conn

import (
	"io"
	"time"

	"github.com/gorilla/websocket"
)

type wsConn struct {
	*websocket.Conn
	typ    int // message type
	reader io.Reader
}

// NewWSConn 创建ws链接
func NewWSConn(conn *websocket.Conn) MsgConn {
	return &wsConn{
		Conn: conn,
	}
}

// Read reads data from the connection.
// Read can be made to time out and return an Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetReadDeadline.
func (c *wsConn) Read(b []byte) (int, error) {
	if c.reader == nil {
		t, r, err := c.NextReader()
		if err != nil {
			return 0, err
		}
		c.typ = t
		c.reader = r
	}
	n, err := c.reader.Read(b)
	if err != nil && err != io.EOF {
		return n, err
	} else if err == io.EOF {
		_, r, err := c.NextReader()
		if err != nil {
			return 0, err
		}
		c.reader = r
	}

	return n, nil
}

func (c *wsConn) Write(b []byte) (int, error) {
	err := c.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}

	return len(b), nil
}

func (c *wsConn) WriteNextMessage(data []byte) error {
	return c.WriteMessage(websocket.BinaryMessage, data)
}

func (c *wsConn) ReadNextMessage() (data []byte, err error) {
	return nil, nil
}

func (c *wsConn) SetDeadline(t time.Time) error {
	if err := c.SetReadDeadline(t); err != nil {
		return err
	}

	return c.SetWriteDeadline(t)
}
