package session

import (
	"context"
	"fmt"
	"time"

	"github.com/liyiysng/scatter/node/message"
)

// MsgInfo 消息信息
type MsgInfo struct {
	// 消息处理完成时间
	FinishTime time.Time
	// 错误信息
	Error error

	// push 消息
	Push message.Message
	// 回复
	Response message.Message
	// 消息写入chan时间
	SendToChanTime time.Time

	// 请求
	ReqOrNotify message.Message
	// 消息读取时间
	ReadTime time.Time

	// 开始处理时间
	BeginHandleTime time.Time
	// 处理完成时间
	EndHandleTime time.Time
}

type sessionContextKey string

const (
	// 消息信息
	_msgInfo sessionContextKey = "_msgInfo"
)

func withInfo(ctx context.Context, msg message.Message) (context.Context, *MsgInfo) {
	info := &MsgInfo{}

	if msg.GetMsgType() == message.REQUEST || msg.GetMsgType() == message.NOTIFY {
		info.ReqOrNotify = msg
	} else if msg.GetMsgType() == message.PUSH {
		info.Push = msg
	} else {
		panic(fmt.Sprintf("invalid message type %v", msg.GetMsgType()))
	}

	return ctx, info
}

// MsgInfoFromContext 从context中获取消息信息
func MsgInfoFromContext(ctx context.Context) (info *MsgInfo, ok bool) {
	info, ok = ctx.Value(_msgInfo).(*MsgInfo)
	return
}
