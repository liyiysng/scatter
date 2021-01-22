package message

import (
	"github.com/golang/protobuf/proto"
	"github.com/liyiysng/scatter/protobuf/node"
)

// ProtobufMsgRequest 消息请求
type ProtobufMsgRequest struct {
	node.MsgRequest
}

// GetMsgType 实现Message接口
func (m *ProtobufMsgRequest) GetMsgType() MsgType {
	return REQUEST
}

// FromBytes 实现Message接口
func (m *ProtobufMsgRequest) FromBytes(b []byte) error {
	return proto.Unmarshal(b, &m.MsgRequest)
}

// ToBytes 实现Message接口
func (m *ProtobufMsgRequest) ToBytes() (b []byte, err error) {
	return proto.Marshal(&m.MsgRequest)
}

// ProtobufMsgResponse 消息回复
type ProtobufMsgResponse struct {
	node.MsgResponse
}

// GetMsgType 实现Message接口
func (m *ProtobufMsgResponse) GetMsgType() MsgType {
	return RESPONSE
}

// FromBytes 实现Message接口
func (m *ProtobufMsgResponse) FromBytes(b []byte) error {
	return proto.Unmarshal(b, &m.MsgResponse)
}

// ToBytes 实现Message接口
func (m *ProtobufMsgResponse) ToBytes() (b []byte, err error) {
	return proto.Marshal(&m.MsgResponse)
}

// ProtobufMsgNotify 消息通知 客户端=>服务器
type ProtobufMsgNotify struct {
	node.MsgNotify
}

// GetMsgType 实现Message接口
func (m *ProtobufMsgNotify) GetMsgType() MsgType {
	return NOTIFY
}

// FromBytes 实现Message接口
func (m *ProtobufMsgNotify) FromBytes(b []byte) error {
	return proto.Unmarshal(b, &m.MsgNotify)
}

// ToBytes 实现Message接口
func (m *ProtobufMsgNotify) ToBytes() (b []byte, err error) {
	return proto.Marshal(&m.MsgNotify)
}

// ProtobufMsgPush 消息推送 服务器=>客户端
type ProtobufMsgPush struct {
	node.MsgPush
}

// GetMsgType 实现Message接口
func (m *ProtobufMsgPush) GetMsgType() MsgType {
	return PUSH
}

// FromBytes 实现Message接口
func (m *ProtobufMsgPush) FromBytes(b []byte) error {
	return proto.Unmarshal(b, &m.MsgPush)
}

// ToBytes 实现Message接口
func (m *ProtobufMsgPush) ToBytes() (b []byte, err error) {
	return proto.Marshal(&m.MsgPush)
}

// ProtobufMsgHeartbeat 心跳消息
type ProtobufMsgHeartbeat struct {
	node.MsgHeartbeat
}

// GetMsgType 实现Message接口
func (m *ProtobufMsgHeartbeat) GetMsgType() MsgType {
	return HEARTBEAT
}

// FromBytes 实现Message接口
func (m *ProtobufMsgHeartbeat) FromBytes(b []byte) error {
	return proto.Unmarshal(b, &m.MsgHeartbeat)
}

// ToBytes 实现Message接口
func (m *ProtobufMsgHeartbeat) ToBytes() (b []byte, err error) {
	return proto.Marshal(&m.MsgHeartbeat)
}

// ProtobufMsgHeartbeatACK 心跳回复
type ProtobufMsgHeartbeatACK struct {
	node.MsgHeartbeatACK
}

// GetMsgType 实现Message接口
func (m *ProtobufMsgHeartbeatACK) GetMsgType() MsgType {
	return HEARTBEATACK
}

// FromBytes 实现Message接口
func (m *ProtobufMsgHeartbeatACK) FromBytes(b []byte) error {
	return proto.Unmarshal(b, &m.MsgHeartbeatACK)
}

// ToBytes 实现Message接口
func (m *ProtobufMsgHeartbeatACK) ToBytes() (b []byte, err error) {
	return proto.Marshal(&m.MsgHeartbeatACK)
}

// ProtobufMsgHandShake 握手消息
type ProtobufMsgHandShake struct {
	node.MsgHandShake
}

// GetMsgType 实现Message接口
func (m *ProtobufMsgHandShake) GetMsgType() MsgType {
	return HANDSHAKE
}

// FromBytes 实现Message接口
func (m *ProtobufMsgHandShake) FromBytes(b []byte) error {
	return proto.Unmarshal(b, &m.MsgHandShake)
}

// ToBytes 实现Message接口
func (m *ProtobufMsgHandShake) ToBytes() (b []byte, err error) {
	return proto.Marshal(&m.MsgHandShake)
}

// ProtobufMsgHandShakeACK 握手消息回复
type ProtobufMsgHandShakeACK struct {
	node.MsgHandShakeACK
}

// GetMsgType 实现Message接口
func (m *ProtobufMsgHandShakeACK) GetMsgType() MsgType {
	return HANDSHAKEACK
}

// FromBytes 实现Message接口
func (m *ProtobufMsgHandShakeACK) FromBytes(b []byte) error {
	return proto.Unmarshal(b, &m.MsgHandShakeACK)
}

// ToBytes 实现Message接口
func (m *ProtobufMsgHandShakeACK) ToBytes() (b []byte, err error) {
	return proto.Marshal(&m.MsgHandShakeACK)
}

// ProtobufMsgError 错误消息
type ProtobufMsgError struct {
	node.MsgError
}

// GetMsgType 实现Message接口
func (m *ProtobufMsgError) GetMsgType() MsgType {
	return ERROR
}

// FromBytes 实现Message接口
func (m *ProtobufMsgError) FromBytes(b []byte) error {
	return proto.Unmarshal(b, &m.MsgError)
}

// ToBytes 实现Message接口
func (m *ProtobufMsgError) ToBytes() (b []byte, err error) {
	return proto.Marshal(&m.MsgError)
}
