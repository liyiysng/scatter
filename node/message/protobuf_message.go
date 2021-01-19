package message

import "github.com/liyiysng/scatter/protobuf/node"

// ProtobufMsgRequest 消息请求
type ProtobufMsgRequest struct {
	node.MsgError
}

// GetMsgType 实现Message接口
func (m *ProtobufMsgRequest) GetMsgType() MsgType {
	return REQUEST
}

// ProtobufMsgResponse 消息回复
type ProtobufMsgResponse struct {
	node.MsgResponse
}

// GetMsgType 实现Message接口
func (m *ProtobufMsgResponse) GetMsgType() MsgType {
	return RESPONSE
}

// ProtobufMsgNotify 消息通知 客户端=>服务器
type ProtobufMsgNotify struct {
	node.MsgNotify
}

// GetMsgType 实现Message接口
func (m *ProtobufMsgNotify) GetMsgType() MsgType {
	return NOTIFY
}

// ProtobufMsgPush 消息推送 服务器=>客户端
type ProtobufMsgPush struct {
	node.MsgPush
}

// GetMsgType 实现Message接口
func (m *ProtobufMsgPush) GetMsgType() MsgType {
	return PUSH
}

// ProtobufMsgHeartbeat 心跳消息
type ProtobufMsgHeartbeat struct {
	node.MsgHeartbeat
}

// GetMsgType 实现Message接口
func (m *ProtobufMsgHeartbeat) GetMsgType() MsgType {
	return HEARTBEAT
}

// ProtobufMsgHeartbeatACK 心跳回复
type ProtobufMsgHeartbeatACK struct {
	node.MsgHeartbeatACK
}

// GetMsgType 实现Message接口
func (m *ProtobufMsgHeartbeatACK) GetMsgType() MsgType {
	return HEARTBEATACK
}

// ProtobufMsgHandShake 握手消息
type ProtobufMsgHandShake struct {
	node.MsgHandShake
}

// GetMsgType 实现Message接口
func (m *ProtobufMsgHandShake) GetMsgType() MsgType {
	return HANDSHAKE
}

// ProtobufMsgHandShakeACK 握手消息回复
type ProtobufMsgHandShakeACK struct {
	node.MsgHandShakeACK
}

// GetMsgType 实现Message接口
func (m *ProtobufMsgHandShakeACK) GetMsgType() MsgType {
	return HANDSHAKEACK
}

// ProtobufMsgError 错误消息
type ProtobufMsgError struct {
	node.MsgError
}

// GetMsgType 实现Message接口
func (m *ProtobufMsgError) GetMsgType() MsgType {
	return ERROR
}
