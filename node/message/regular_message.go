package message

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/liyiysng/scatter/constants"
)

// MsgRequest 消息请求
type MsgRequest struct {
	Sequence int32  // 请求时由客户端填充 , 回复时由服务器设置相同序列
	Service  string // 服务名称 如 Game.Foo
	Payload  []byte // 数据
}

// GetMsgType 实现Message接口
func (m *MsgRequest) GetMsgType() MsgType {
	return REQUEST
}

// FromBytes 实现Message接口
func (m *MsgRequest) FromBytes(b []byte) error {
	r := bytes.NewBuffer(b)
	err := binary.Read(r, binary.BigEndian, &m.Sequence)
	if err != nil {
		if err == io.EOF { // 消息不完整
			return io.ErrUnexpectedEOF
		}
		return err
	}
	str, err := ReadService(r)
	if err != nil {
		return err
	}
	m.Service = str

	// 读取剩余
	m.Payload = r.Bytes()

	return nil
}

// MsgResponse 消息回复
type MsgResponse struct {
	Sequence    int32  // 请求时由客户端填充 , 回复时由服务器设置相同序列
	Payload     []byte // 数据
	CustomError string //自定义错误描述 货币不足等逻辑错误. 当设置错误是,Payload无效
}

// GetMsgType 实现Message接口
func (m *MsgResponse) GetMsgType() MsgType {
	return RESPONSE
}

// MsgNotify 消息通知 客户端=>服务器
type MsgNotify struct {
	Payload []byte // 数据
}

// GetMsgType 实现Message接口
func (m *MsgNotify) GetMsgType() MsgType {
	return NOTIFY
}

// MsgPush 消息推送 服务器=>客户端
type MsgPush struct {
	Service [constants.MaxMessageServiceName]byte
	Payload []byte // 数据
}

// GetMsgType 实现Message接口
func (m *MsgPush) GetMsgType() MsgType {
	return PUSH
}

// MsgHeartbeat 心跳消息
type MsgHeartbeat struct {
	TimeStamp int64
}

// GetMsgType 实现Message接口
func (m *MsgHeartbeat) GetMsgType() MsgType {
	return HEARTBEAT
}

// MsgHeartbeatACK 心跳回复
type MsgHeartbeatACK struct {
	TimeStamp int64
}

// GetMsgType 实现Message接口
func (m *MsgHeartbeatACK) GetMsgType() MsgType {
	return HEARTBEATACK
}

// MsgHandShake 握手消息
type MsgHandShake struct {
	Platform      string
	ClientVersion string
	BuildVersion  string
}

// GetMsgType 实现Message接口
func (m *MsgHandShake) GetMsgType() MsgType {
	return HANDSHAKE
}

// MsgHandShakeACK 握手消息回复
type MsgHandShakeACK struct {
	ServerTimeStamp int64
}

// GetMsgType 实现Message接口
func (m *MsgHandShakeACK) GetMsgType() MsgType {
	return HANDSHAKEACK
}

// MsgError 错误消息
type MsgError struct {
	ErrorInfo string
}

// GetMsgType 实现Message接口
func (m *MsgError) GetMsgType() MsgType {
	return ERROR
}
