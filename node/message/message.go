// Package message 消息格式定义
package message

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/liyiysng/scatter/protobuf/node"
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

func (mt *MsgType) String() string {
	return msgTypes[*mt]
}

// MsgOpt 消息选项
type MsgOpt byte

const (
	// COMPRESS 是否压缩
	COMPRESS MsgOpt = 0x01
)

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
}

// BuildPushMessage 创建一个push Message
func BuildPushMessage(cmd string, data []byte) (msg Message) {
	msg = &ProtobufMsgPush{
		node.MsgPush{
			Command: cmd,
			Payload: data,
		},
	}
	return
}

// BuildHeatAckMessage 创建一个心跳回复 Message
func BuildHeatAckMessage() (msg Message) {
	msg = &ProtobufMsgHeartbeatACK{
		node.MsgHeartbeatACK{
			TimeStamp: time.Now().Unix(),
		},
	}
	return
}

// BuildHandShakeAckMessage 创建一个握手回复 Message
func BuildHandShakeAckMessage() (msg Message) {
	msg = &ProtobufMsgHandShakeACK{
		node.MsgHandShakeACK{
			TimeStamp: time.Now().Unix(),
		},
	}
	return
}

// BuildResponseMessage 创建一个回复
func BuildResponseMessage(sequence int32, payload []byte) (msg Message) {
	msg = &ProtobufMsgResponse{
		node.MsgResponse{
			Sequence: sequence,
			Payload:  payload,
		},
	}
	return
}

// BuildResponseCustomErrorMessage 创建一个自定义错误回复
func BuildResponseCustomErrorMessage(sequence int32, customError string) (msg Message) {
	msg = &ProtobufMsgResponse{
		node.MsgResponse{
			Sequence:    sequence,
			CustomError: customError,
		},
	}
	return
}

// BuildMessage 根据给定类型和buf构建消息
func BuildMessage(mtype MsgType, buf []byte) (msg Message, err error) {
	switch mtype {

	case REQUEST: // REQUEST 请求消息
		msg = &ProtobufMsgRequest{}
		break
	case RESPONSE: // RESPONSE 回复消息
		msg = &ProtobufMsgResponse{}
		break
	case NOTIFY: // NOTIFY 通知消息 客户端=>服务器
		msg = &ProtobufMsgNotify{}
		break
	case PUSH: // PUSH 推送消息 服务器=>客户端
		msg = &ProtobufMsgPush{}
		break
	case HEARTBEAT: // HEARTBEAT 心跳/ping 消息
		msg = &ProtobufMsgHeartbeat{}
		break
	case HEARTBEATACK: // HEARTBEATACK 心跳/ping 回复
		msg = &ProtobufMsgHeartbeatACK{}
		break
	case HANDSHAKE: // HANDSHAKE 握手消息
		msg = &ProtobufMsgHandShake{}
		break
	case HANDSHAKEACK: // HANDSHAKEACK 握手消息回复
		msg = &ProtobufMsgHandShakeACK{}
		break
	case ERROR: // ERROR 错误消息 服务器=>客户端
		msg = &ProtobufMsgError{}
	default:
		return nil, fmt.Errorf("[BuildMessage] invalid message type [%s]", mtype.String())
	}

	err = msg.FromBytes(buf)
	if err != nil {
		return nil, err
	}
	return
}

// GetSrvMethod 获取服务名和方法名
func GetSrvMethod(service string) (srvName, methodName string, err error) {
	dot := strings.LastIndex(srvName, ".")
	if dot < 0 {
		return "", "", errors.New("[frontendSession.handleMsg]: service/method ill-formed: " + service)
	}
	return service[:dot], service[dot+1:], nil
}
