// Package message 消息格式定义
package message

import (
	"io"
)

// MsgType 消息类型
type MsgType byte

const (
	// NULL 非法消息
	NULL MsgType = iota
	// REQUEST 请求消息
	REQUEST
	// RESPONSE 回复消息
	RESPONSE
	// NOTIFY 通知消息 客户端=>服务器
	NOTIFY
	// PUSH 推送消息 服务器=>客户端
	PUSH
	// HEARTBEAT 心跳/ping 消息
	HEARTBEAT
	// HEARTBEATACK 心跳/ping 回复
	HEARTBEATACK
	// HANDSHAKE 握手消息
	HANDSHAKE
	// HANDSHAKEACK 握手消息回复
	HANDSHAKEACK
	// ERROR 错误消息 服务器=>客户端
	ERROR
)

var msgTypes = map[MsgType]string{
	NULL:         "NULL",
	REQUEST:      "REQUEST",
	RESPONSE:     "RESPONSE",
	NOTIFY:       "NOTIFY",
	PUSH:         "PUSH",
	HEARTBEAT:    "HEARTBEAT",
	HEARTBEATACK: "HEARTBEATACK",
}

func (mt *MsgType) String() string {
	return msgTypes[*mt]
}

// Message 表示一条消息
type Message interface {
	// 消息类型
	GetMsgType() MsgType
	// 从io.Reader序列化出一个消息
	FromReader(r io.Reader) error
	// 写入writer
	ToWriter(w io.Writer) error
}
