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

	// 开始时间
	StartTime time.Time

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

	// 开始时间
	StartTime time.Time `json:"@timestamp"`

	// 消息类型
	MsgType message.MsgType `json:"msg_type"`

	// 服务
	Srv string `json:"srv"`

	// 方法
	Method string `json:"method"`

	// 请求/notify
	MsgRead string `json:"msg_read"`
	// 消息读取时间
	ReadTime float32 `json:"read_time_us"`

	// 开始处理时间
	BeginHandleTime float32 `json:"begin_handle_time_us"`
	// 处理完成时间
	EndHandleTime float32 `json:"end_handle_time_us"`

	// 回复/push
	MsgWrite string `json:"msg_write"`
	// 消息写入chan时间
	SendToChanTime float32 `json:"send_to_chan_time_us"`

	// 消息处理完成时间
	FinishTime float32 `json:"finish_time_us"`
	// 错误信息
	Error string `json:"error"`
}

// MarshalJSON json marshal
func (info *MsgInfo) MarshalJSON() ([]byte, error) {
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
		StartTime: info.StartTime,
		Error:     errStr,
		MsgWrite:  msgWriteStr,
		MsgRead:   msgReadStr,
	}

	// time
	if !info.FinishTime.IsZero() {
		js.FinishTime = float32(info.FinishTime.Sub(info.StartTime).Microseconds())
	}
	if !info.SendToChanTime.IsZero() {
		js.SendToChanTime = float32(info.SendToChanTime.Sub(info.StartTime).Microseconds())
	}
	if !info.ReadTime.IsZero() {
		js.ReadTime = float32(info.ReadTime.Sub(info.StartTime).Microseconds())
	}
	if !info.BeginHandleTime.IsZero() {
		js.BeginHandleTime = float32(info.BeginHandleTime.Sub(info.StartTime).Microseconds())
	}
	if !info.EndHandleTime.IsZero() {
		js.EndHandleTime = float32(info.EndHandleTime.Sub(info.StartTime).Microseconds())
	}

	if info.MsgRead != nil {
		js.MsgType = info.MsgRead.GetMsgType()
		if info.MsgRead.GetService() != "" {
			serviceName, methodName, err := message.GetSrvMethod(info.MsgRead.GetService())
			if err != nil {
				myLog.Errorf("[MsgInfo.MarshalJSON] get srv failed %v", err)
			} else {
				js.Srv = serviceName
				js.Method = methodName
			}
		}
	} else {
		js.Method = info.MsgWrite.GetService()
		js.MsgType = message.PUSH
	}

	str, err := json.MarshalIndent(js, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("[MsgInfo.MarshalJSON] %v", err)
	}
	return []byte(str), err
}

func (info *MsgInfo) String() string {
	b, err := info.MarshalJSON()
	if err != nil {
		return err.Error()
	}
	return string(b)
}

type sessionContextKey string

const (
	// 消息信息
	_msgInfo sessionContextKey = "_msgInfo"
)

func withInfo(ctx context.Context) context.Context {
	info := &MsgInfo{}
	info.StartTime = time.Now()
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
