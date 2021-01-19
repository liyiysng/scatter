package conn

import "net"

type tcpConn struct {
	net.Conn
}

func (c *tcpConn) GetNextMessage() (data []byte, err error) {
	return nil, nil
}
