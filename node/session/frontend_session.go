package session

import (
	"context"

	"github.com/liyiysng/scatter/node/conn"
	"google.golang.org/protobuf/proto"
)

type frontendSession struct {
	sid  string
	nid  string
	conn conn.MsgConn
}

// NewFrontendSession 创建一个session
func NewFrontendSession(nid string, c conn.MsgConn) Session {
	return &frontendSession{
		sid:  "",
		nid:  nid,
		conn: c,
	}
}

func (s *frontendSession) Stats() State {
	return State{
		SID:           s.sid,
		NID:           s.nid,
		RemoteAddress: s.conn.RemoteAddr().String(),
	}
}

func (s *frontendSession) Push(ctx context.Context, cmd string, msg proto.Message) error {
	panic("unimplement")
}

func (s *frontendSession) OnClose(onClose OnClose) {
	panic("unimplement")
}

func (s *frontendSession) Closed() bool {
	panic("unimplement")
}
