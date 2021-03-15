package message

import (
	"github.com/golang/protobuf/proto"
	phead "github.com/liyiysng/scatter/node/message/proto"
)

// ProtobufMsg protobuf类型消息
type ProtobufMsg struct {
	phead.Head
}

// GetMsgType 实现Message接口
func (m *ProtobufMsg) GetMsgType() MsgType {
	return MsgType(m.MsgType)
}

// FromBytes 实现Message接口
func (m *ProtobufMsg) FromBytes(b []byte) error {
	return proto.Unmarshal(b, &m.Head)
}

// ToBytes 实现Message接口
func (m *ProtobufMsg) ToBytes() (b []byte, err error) {
	//若时回复消息,service字段无需序列化
	if m.GetMsgType() == phead.MsgType_RESPONSE {
		srv := m.Head.Service
		m.Head.Service = ""
		b, err = proto.Marshal(&m.Head)
		m.Head.Service = srv
		return
	}
	return proto.Marshal(&m.Head)
}

// protoBufFactory probuf类型的工厂
type protoBufFactory struct {
}

// BuildMessage 根据buf构建消息
func (f *protoBufFactory) BuildMessage(buf []byte) (msg Message, err error) {
	msg = &ProtobufMsg{}
	err = msg.FromBytes(buf)
	if err != nil {
		return nil, err
	}

	return
}

// BuildPushMessage 创建一个push Message
func (f *protoBufFactory) BuildPushMessage(cmd string, data []byte) (msg Message, err error) {
	msg = &ProtobufMsg{
		Head: phead.Head{
			MsgType: phead.MsgType_PUSH,
			Service: cmd,
			Payload: data,
		},
	}
	return
}

// BuildHeatAckMessage 创建一个心跳回复 Message
func (f *protoBufFactory) BuildHeatAckMessage() (msg Message, err error) {
	msg = &ProtobufMsg{
		Head: phead.Head{
			MsgType: phead.MsgType_HEARTBEATACK,
		},
	}
	return
}

// BuildHandShakeMessage 创建一个握手消息 Message
func (f *protoBufFactory) BuildHandShakeMessage(platform, clientVersion, buildVersion string) (msg Message, err error) {

	handShake := &phead.MsgHandShake{
		Platform:      platform,
		ClientVersion: clientVersion,
		BuildVersion:  buildVersion,
	}

	buf, err := proto.Marshal(handShake)
	if err != nil {
		return nil, err
	}

	msg = &ProtobufMsg{
		Head: phead.Head{
			MsgType: phead.MsgType_HANDSHAKE,
			Payload: buf,
		},
	}
	return
}

// BuildHandShakeAckMessage 创建一个握手回复 Message
func (f *protoBufFactory) BuildHandShakeAckMessage() (msg Message, err error) {
	msg = &ProtobufMsg{
		Head: phead.Head{
			MsgType: phead.MsgType_HANDSHAKEACK,
		},
	}
	return
}

// BuildResponseMessage 创建一个回复
func (f *protoBufFactory) BuildResponseMessage(sequence int32, srv string, payload []byte) (msg Message, err error) {
	msg = &ProtobufMsg{
		Head: phead.Head{
			MsgType:  phead.MsgType_RESPONSE,
			Service:  srv,
			Sequence: sequence,
			Payload:  payload,
		},
	}
	return
}

//BuildResponseCustomErrorMessage 创建一个自定义错误回复
func (f *protoBufFactory) BuildResponseCustomErrorMessage(sequence int32, srv string, customError string) (msg Message, err error) {
	msg = &ProtobufMsg{
		Head: phead.Head{
			MsgType:     phead.MsgType_RESPONSE,
			Service:     srv,
			Sequence:    sequence,
			CustomError: customError,
		},
	}
	return
}

// ParseHandShake 解析握手数据
func (f *protoBufFactory) ParseHandShake(buf []byte) (h interface{}, err error) {
	handshake := &phead.MsgHandShake{}
	err = proto.Unmarshal(buf, handshake)
	if err != nil {
		return
	}
	h = handshake
	return
}

// BuildRequestMessage 创建一个请求
func (f *protoBufFactory) BuildRequestMessage(sequence int32, srv string, payload []byte) (msg Message, err error) {
	msg = &ProtobufMsg{
		Head: phead.Head{
			MsgType:  phead.MsgType_REQUEST,
			Service:  srv,
			Sequence: sequence,
			Payload:  payload,
		},
	}
	return
}

// BuildNotifyMessage 创建一个通知
func (f *protoBufFactory) BuildNotifyMessage(srv string, payload []byte) (msg Message, err error) {
	msg = &ProtobufMsg{
		Head: phead.Head{
			MsgType: phead.MsgType_NOTIFY,
			Service: srv,
			Payload: payload,
		},
	}
	return
}
