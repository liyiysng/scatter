package session

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/liyiysng/scatter/node/message"
)

// MsgInfo 消息信息
type MsgInfo struct {
	// 消息处理完成时间
	FinishTime time.Time
	// 错误信息
	Error error

	// 回复/push
	MsgWrite message.Message
	// 回复数据类型
	WritePayload interface{}
	// 消息写入chan时间
	SendToChanTime time.Time

	// 请求/notify
	MsgRead message.Message
	// 回复数据类型
	ReadPayload interface{}
	// 消息读取时间
	ReadTime time.Time

	// 开始处理时间
	BeginHandleTime time.Time
	// 处理完成时间
	EndHandleTime time.Time
}

type jsonFormatInfo struct {

	// 请求/notify
	MsgRead string `json:"msg_read"`
	// 消息读取时间
	ReadTime time.Time `json:"read_time"`

	// 开始处理时间
	BeginHandleTime time.Time `json:"begin_handle_time"`
	// 处理完成时间
	EndHandleTime time.Time `json:"end_handle_time"`

	// 回复/push
	MsgWrite string `json:"msg_write"`
	// 消息写入chan时间
	SendToChanTime time.Time `json:"send_to_chan_time"`

	// 消息处理完成时间
	FinishTime time.Time `json:"finish_time"`
	// 错误信息
	Error string `json:"error"`
}

func (info *MsgInfo) String() string {

	sb := strings.Builder{}

	errStr := ""
	if info.Error != nil {
		errStr = info.Error.Error()
	}

	msgWriteStr := ""
	if info.MsgWrite != nil {
		sb.Reset()
		sb.WriteString(fmt.Sprintf("MsgType:%v ", info.MsgWrite.GetMsgType()))
		sb.WriteString(fmt.Sprintf("Service:%v ", info.MsgWrite.GetService()))
		sb.WriteString(fmt.Sprintf("Sequence:%v ", info.MsgWrite.GetSequence()))
		if info.WritePayload != nil {
			sb.WriteString(fmt.Sprintf("payload:%v ", info.WritePayload))
		}
		msgWriteStr = sb.String()
	}

	msgReadStr := ""
	if info.MsgRead != nil {
		sb.Reset()
		sb.WriteString(fmt.Sprintf("MsgType:%v ", info.MsgRead.GetMsgType()))
		sb.WriteString(fmt.Sprintf("Service:%v ", info.MsgRead.GetService()))
		sb.WriteString(fmt.Sprintf("Sequence:%v ", info.MsgRead.GetSequence()))
		if info.ReadPayload != nil {
			sb.WriteString(fmt.Sprintf("payload:%v ", info.ReadPayload))
		}
		msgReadStr = sb.String()
	}

	js := &jsonFormatInfo{
		FinishTime:      info.FinishTime,
		Error:           errStr,
		MsgWrite:        msgWriteStr,
		SendToChanTime:  info.SendToChanTime,
		MsgRead:         msgReadStr,
		ReadTime:        info.ReadTime,
		BeginHandleTime: info.BeginHandleTime,
		EndHandleTime:   info.EndHandleTime,
	}

	str, err := json.MarshalIndent(js, "", "\t")
	if err != nil {
		myLog.Errorf("[MsgInfo.String] %v", err)
		return "marshal failed"
	}
	return string(str)
}

type sessionContextKey string

const (
	// 消息信息
	_msgInfo sessionContextKey = "_msgInfo"
)

func withInfo(ctx context.Context) context.Context {
	info := &MsgInfo{}
	return context.WithValue(ctx, _msgInfo, info)
}

// MsgInfoFromContext 从context中获取消息信息
func MsgInfoFromContext(ctx context.Context) (info *MsgInfo, ok bool) {
	info, ok = ctx.Value(_msgInfo).(*MsgInfo)
	return
}

// 当读取到消息
func onMsgRead(ctx context.Context, msgRead message.Message) {
	info, ok := MsgInfoFromContext(ctx)
	if !ok {
		myLog.Warning("[profile.onMsgRead] message info not found in context")
		return
	}
	info.ReadTime = time.Now()
	info.MsgRead = msgRead
}

// 开始处理
func onMsgBegingHandle(ctx context.Context) {
	info, ok := MsgInfoFromContext(ctx)
	if !ok {
		myLog.Warning("[profile.onMsgBegingHandle] message info not found in context")
		return
	}
	info.BeginHandleTime = time.Now()
}

// 结束处理
func onMsgEndHandle(ctx context.Context) {
	info, ok := MsgInfoFromContext(ctx)
	if !ok {
		myLog.Warning("[profile.onMsgEndHandle] message info not found in context")
		return
	}
	info.EndHandleTime = time.Now()
}

// 回复/push
func onMsgWrite(ctx context.Context, msgWrite message.Message) {
	info, ok := MsgInfoFromContext(ctx)
	if !ok {
		myLog.Warning("[profile.onMsgWrite] message info not found in context")
		return
	}
	info.SendToChanTime = time.Now()
	info.MsgWrite = msgWrite
}

// 消息结束
func onMsgFinished(ctx context.Context, err error) {
	info, ok := MsgInfoFromContext(ctx)
	if !ok {
		myLog.Warning("[profile.onMsgFinished] message info not found in context")
		return
	}
	info.FinishTime = time.Now()
	info.Error = err
}

// SetReadPayloadObj 设置读取的对象
func SetReadPayloadObj(ctx context.Context, obj interface{}) {
	info, ok := MsgInfoFromContext(ctx)
	if !ok {
		myLog.Warning("[profile.SetReadPayloadObj] message info not found in context")
		return
	}
	info.ReadPayload = obj
}

// SetWritePayloadObj 设置回复的对象
func SetWritePayloadObj(ctx context.Context, obj interface{}) {
	info, ok := MsgInfoFromContext(ctx)
	if !ok {
		myLog.Warning("[profile.SetWritePayloadObj] message info not found in context")
		return
	}
	info.WritePayload = obj
}
