// Package message 消息格式定义
package message

import (
	"errors"
	"fmt"
	"strings"

	"github.com/liyiysng/scatter/logger"
)

var (
	myLog = logger.Component("message")
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
	HANDSHAKE:    "HANDSHAKE",
	HANDSHAKEACK: "HANDSHAKEACK",
	ERROR:        "ERROR",
}

func (mt MsgType) String() string {
	return msgTypes[mt]
}

// MarshalJSON json marshal
func (mt *MsgType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", msgTypes[*mt])), nil
}

// MsgFactory 消息工厂
var MsgFactory Factor = &protoBufFactory{}

// Message 表示一条消息
type Message interface {
	// 消息类型
	GetMsgType() MsgType
	// 从[]byte序列化出一个消息
	// 返回 erro : io.EOF:说明消息完整,并且是最后一个
	//	io.ErrUnexpectedEOF 说明消息不完整,不予处理
	FromBytes(b []byte) error
	// 写入b
	ToBytes() (b []byte, err error)
	// 获取服务Game.Foo
	GetService() string
	// 获取序号
	GetSequence() int32
	// 获取数据
	GetPayload() []byte
	// 字符表达
	String() string
}

// Factor 消息工厂,更具数据构建消息
type Factor interface {
	// BuildMessage 根据buf构建消息
	BuildMessage(buf []byte) (msg Message, err error)
	// BuildPushMessage 创建一个push Message
	BuildPushMessage(cmd string, data []byte) (msg Message, err error)
	// BuildHeatAckMessage 创建一个心跳回复 Message
	BuildHeatAckMessage() (msg Message, err error)
	// BuildHandShakeMessage 创建一个握手消息 Message
	BuildHandShakeMessage(platform, clientVersion, buildVersion string) (msg Message, err error)
	// BuildHandShakeAckMessage 创建一个握手回复 Message
	BuildHandShakeAckMessage() (msg Message, err error)
	// BuildResponseMessage 创建一个回复
	BuildResponseMessage(sequence int32, srv string, payload []byte) (msg Message, err error)
	//BuildResponseCustomErrorMessage 创建一个自定义错误回复
	BuildResponseCustomErrorMessage(sequence int32, srv string, customError string) (msg Message, err error)

	// ParseHandShake 解析握手数据
	ParseHandShake(buf []byte) (h interface{}, err error)
	// BuildRequestMessage 创建一个请求
	BuildRequestMessage(sequence int32, srv string, payload []byte) (msg Message, err error)
	// BuildNotifyMessage 创建一个通知
	BuildNotifyMessage(srv string, payload []byte) (msg Message, err error)
}

// GetSrvMethod 获取服务名和方法名
func GetSrvMethod(service string) (srvName, methodName string, err error) {
	dot := strings.LastIndex(service, ".")
	if dot < 0 {
		return "", "", errors.New("[GetSrvMethod]: service/method ill-formed: " + service)
	}
	return service[:dot], service[dot+1:], nil
}
