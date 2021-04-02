package session

import (
	"errors"

	"github.com/liyiysng/scatter/cluster/sessionpb"
	"github.com/liyiysng/scatter/constants"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

// ISession 描述一个session
type ISession interface {
	Info() *sessionpb.SessionInfo
}

type ITransferSession interface {
	ISession
	// uid
	//!
	GetUID() constants.UID
	//!
	IsUIDBind() bool
}

type baseSession struct {
	*sessionpb.SessionInfo
}

func (s *baseSession) Info() *sessionpb.SessionInfo {
	return s.SessionInfo
}

type transferSession struct {
	*baseSession
}

func (s *transferSession) GetUID() constants.UID {
	return s.TransferInfo.UID
}

func (s *transferSession) IsUIDBind() bool {
	return s.TransferInfo.UID == constants.DefaultUID
}

type pubSession struct {
	*baseSession
}

func GetSession(info *sessionpb.SessionInfo) (ISession, error) {

	switch info.SType {
	case sessionpb.SessionType_Backend:
		{
			panic("not implement")
		}
	case sessionpb.SessionType_FrontEnd:
		{
			panic("not implement")
		}
	case sessionpb.SessionType_Pub:
		{
			return &pubSession{
				baseSession: &baseSession{
					SessionInfo: info,
				},
			}, nil
		}
	case sessionpb.SessionType_Transfer:
		{
			return &transferSession{
				baseSession: &baseSession{
					SessionInfo: info,
				},
			}, nil
		}
	}

	return nil, nil
}
