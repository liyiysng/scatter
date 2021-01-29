package session

import "github.com/liyiysng/scatter/node/conn"

type frontendSession struct {
	sid  string
	nid  string
	conn conn.MsgConn
}

// NewFrontendSession 创建一个session
func NewFrontendSession(c conn.MsgConn) {

}

func (s *frontendSession) Stats() State {
	return State{
		SID:           s.sid,
		NID:           s.nid,
		RemoteAddress: s.conn.RemoteAddr().String(),
	}
}
