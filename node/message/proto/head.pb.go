// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: node/message/proto/head.proto

package proto

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

//消息类型
type MsgType int32

const (
	// NULL 非法消息
	MsgType_NULL MsgType = 0
	// REQUEST 请求消息
	MsgType_REQUEST MsgType = 1
	// RESPONSE 回复消息
	MsgType_RESPONSE MsgType = 2
	// NOTIFY 通知消息 客户端=>服务器
	MsgType_NOTIFY MsgType = 3
	// PUSH 推送消息 服务器=>客户端
	MsgType_PUSH MsgType = 4
	// HEARTBEAT 心跳/ping 消息
	MsgType_HEARTBEAT MsgType = 5
	// HEARTBEATACK 心跳/ping 回复
	MsgType_HEARTBEATACK MsgType = 6
	// HANDSHAKE 握手消息
	MsgType_HANDSHAKE MsgType = 7
	// HANDSHAKEACK 握手消息回复
	MsgType_HANDSHAKEACK MsgType = 8
	// ERROR 错误消息 服务器=>客户端
	MsgType_ERROR MsgType = 9
	// KICK 剔除 , 客户端收到该消息应断开链接(服务器器主动断开会有 TIME_WAIT 问题)
	MsgType_KICK               MsgType = 10
	MsgType_LOCAL_CONNECTED    MsgType = 20 //本地连接成功
	MsgType_LOCAL_DISCONNECTED MsgType = 21 //本地连接断开
	MsgType_LOCAL_TIMEOUT      MsgType = 22 //本地超时
)

// Enum value maps for MsgType.
var (
	MsgType_name = map[int32]string{
		0:  "NULL",
		1:  "REQUEST",
		2:  "RESPONSE",
		3:  "NOTIFY",
		4:  "PUSH",
		5:  "HEARTBEAT",
		6:  "HEARTBEATACK",
		7:  "HANDSHAKE",
		8:  "HANDSHAKEACK",
		9:  "ERROR",
		10: "KICK",
		20: "LOCAL_CONNECTED",
		21: "LOCAL_DISCONNECTED",
		22: "LOCAL_TIMEOUT",
	}
	MsgType_value = map[string]int32{
		"NULL":               0,
		"REQUEST":            1,
		"RESPONSE":           2,
		"NOTIFY":             3,
		"PUSH":               4,
		"HEARTBEAT":          5,
		"HEARTBEATACK":       6,
		"HANDSHAKE":          7,
		"HANDSHAKEACK":       8,
		"ERROR":              9,
		"KICK":               10,
		"LOCAL_CONNECTED":    20,
		"LOCAL_DISCONNECTED": 21,
		"LOCAL_TIMEOUT":      22,
	}
)

func (x MsgType) Enum() *MsgType {
	p := new(MsgType)
	*p = x
	return p
}

func (x MsgType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MsgType) Descriptor() protoreflect.EnumDescriptor {
	return file_node_message_proto_head_proto_enumTypes[0].Descriptor()
}

func (MsgType) Type() protoreflect.EnumType {
	return &file_node_message_proto_head_proto_enumTypes[0]
}

func (x MsgType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MsgType.Descriptor instead.
func (MsgType) EnumDescriptor() ([]byte, []int) {
	return file_node_message_proto_head_proto_rawDescGZIP(), []int{0}
}

type MsgError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code     int32  `protobuf:"varint,1,opt,name=Code,proto3" json:"Code,omitempty"`
	Describe string `protobuf:"bytes,2,opt,name=Describe,proto3" json:"Describe,omitempty"`
}

func (x *MsgError) Reset() {
	*x = MsgError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_message_proto_head_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgError) ProtoMessage() {}

func (x *MsgError) ProtoReflect() protoreflect.Message {
	mi := &file_node_message_proto_head_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgError.ProtoReflect.Descriptor instead.
func (*MsgError) Descriptor() ([]byte, []int) {
	return file_node_message_proto_head_proto_rawDescGZIP(), []int{0}
}

func (x *MsgError) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *MsgError) GetDescribe() string {
	if x != nil {
		return x.Describe
	}
	return ""
}

type Head struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MsgType  MsgType   `protobuf:"varint,1,opt,name=MsgType,proto3,enum=scatter.node.proto.MsgType" json:"MsgType,omitempty"` // 消息类型
	Sequence int32     `protobuf:"varint,2,opt,name=Sequence,proto3" json:"Sequence,omitempty"`                               // 请求时由客户端填充 , 回复时由服务器设置相同序列
	Service  string    `protobuf:"bytes,3,opt,name=Service,proto3" json:"Service,omitempty"`                                  // 服务名称 如 Game.Foo
	Error    *MsgError `protobuf:"bytes,4,opt,name=Error,proto3" json:"Error,omitempty"`                                      // 自定义错误描述 货币不足等逻辑错误
	Payload  []byte    `protobuf:"bytes,5,opt,name=Payload,proto3" json:"Payload,omitempty"`                                  // 数据
}

func (x *Head) Reset() {
	*x = Head{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_message_proto_head_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Head) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Head) ProtoMessage() {}

func (x *Head) ProtoReflect() protoreflect.Message {
	mi := &file_node_message_proto_head_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Head.ProtoReflect.Descriptor instead.
func (*Head) Descriptor() ([]byte, []int) {
	return file_node_message_proto_head_proto_rawDescGZIP(), []int{1}
}

func (x *Head) GetMsgType() MsgType {
	if x != nil {
		return x.MsgType
	}
	return MsgType_NULL
}

func (x *Head) GetSequence() int32 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

func (x *Head) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

func (x *Head) GetError() *MsgError {
	if x != nil {
		return x.Error
	}
	return nil
}

func (x *Head) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type MsgHandShake struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Platform      string `protobuf:"bytes,1,opt,name=Platform,proto3" json:"Platform,omitempty"`
	ClientVersion string `protobuf:"bytes,2,opt,name=ClientVersion,proto3" json:"ClientVersion,omitempty"`
	BuildVersion  string `protobuf:"bytes,3,opt,name=BuildVersion,proto3" json:"BuildVersion,omitempty"`
	UDID          string `protobuf:"bytes,4,opt,name=UDID,proto3" json:"UDID,omitempty"`              // 用户硬件设备号（用户硬件设备号--Android和iOS都用的uuid，32位通用唯一识别码）
	Model         string `protobuf:"bytes,5,opt,name=Model,proto3" json:"Model,omitempty"`            // 设备的机型，例如Samsung GT-I9208
	OsVersion     string `protobuf:"bytes,6,opt,name=OsVersion,proto3" json:"OsVersion,omitempty"`    // 操作系统版本，例如13.0.2
	Network       string `protobuf:"bytes,7,opt,name=Network,proto3" json:"Network,omitempty"`        // 网络信息（4G/3G/WIFI/2G）。ps：不能区分移动新号就统一填2G吧
	BSdkUDid      string `protobuf:"bytes,10,opt,name=BSdkUDid,proto3" json:"BSdkUDid,omitempty"`     // 用户硬件设备号（b服SDK udid，客户端SDK登录事件接口获取，32位通用唯一识别码）
	BGameID       string `protobuf:"bytes,11,opt,name=BGameID,proto3" json:"BGameID,omitempty"`       // 游戏id（一款游戏的ID，类似app_id，b服SDK获得）
	BChannelID    string `protobuf:"bytes,12,opt,name=BChannelID,proto3" json:"BChannelID,omitempty"` // 游戏的渠道ID（游戏安装包的渠道ID）
	BSkdTypeID    string `protobuf:"bytes,13,opt,name=BSkdTypeID,proto3" json:"BSkdTypeID,omitempty"` // 客户端获取到的sdk_type（如果客户端拿不到BGameID，就填充该值，后端映射填充BGameID）
}

func (x *MsgHandShake) Reset() {
	*x = MsgHandShake{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_message_proto_head_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgHandShake) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgHandShake) ProtoMessage() {}

func (x *MsgHandShake) ProtoReflect() protoreflect.Message {
	mi := &file_node_message_proto_head_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgHandShake.ProtoReflect.Descriptor instead.
func (*MsgHandShake) Descriptor() ([]byte, []int) {
	return file_node_message_proto_head_proto_rawDescGZIP(), []int{2}
}

func (x *MsgHandShake) GetPlatform() string {
	if x != nil {
		return x.Platform
	}
	return ""
}

func (x *MsgHandShake) GetClientVersion() string {
	if x != nil {
		return x.ClientVersion
	}
	return ""
}

func (x *MsgHandShake) GetBuildVersion() string {
	if x != nil {
		return x.BuildVersion
	}
	return ""
}

func (x *MsgHandShake) GetUDID() string {
	if x != nil {
		return x.UDID
	}
	return ""
}

func (x *MsgHandShake) GetModel() string {
	if x != nil {
		return x.Model
	}
	return ""
}

func (x *MsgHandShake) GetOsVersion() string {
	if x != nil {
		return x.OsVersion
	}
	return ""
}

func (x *MsgHandShake) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
}

func (x *MsgHandShake) GetBSdkUDid() string {
	if x != nil {
		return x.BSdkUDid
	}
	return ""
}

func (x *MsgHandShake) GetBGameID() string {
	if x != nil {
		return x.BGameID
	}
	return ""
}

func (x *MsgHandShake) GetBChannelID() string {
	if x != nil {
		return x.BChannelID
	}
	return ""
}

func (x *MsgHandShake) GetBSkdTypeID() string {
	if x != nil {
		return x.BSkdTypeID
	}
	return ""
}

var File_node_message_proto_head_proto protoreflect.FileDescriptor

var file_node_message_proto_head_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x68, 0x65, 0x61, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x12, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x3a, 0x0a, 0x08, 0x4d, 0x73, 0x67, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12,
	0x12, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x43,
	0x6f, 0x64, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x22,
	0xc1, 0x01, 0x0a, 0x04, 0x48, 0x65, 0x61, 0x64, 0x12, 0x35, 0x0a, 0x07, 0x4d, 0x73, 0x67, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x73, 0x63, 0x61, 0x74,
	0x74, 0x65, 0x72, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d,
	0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x07, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x08, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x32, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x6e,
	0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x73, 0x67, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x52, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x50, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x50, 0x61, 0x79, 0x6c,
	0x6f, 0x61, 0x64, 0x22, 0xcc, 0x02, 0x0a, 0x0c, 0x4d, 0x73, 0x67, 0x48, 0x61, 0x6e, 0x64, 0x53,
	0x68, 0x61, 0x6b, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x50, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d,
	0x12, 0x24, 0x0a, 0x0d, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x0c, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x42, 0x75,
	0x69, 0x6c, 0x64, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x55, 0x44,
	0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x55, 0x44, 0x49, 0x44, 0x12, 0x14,
	0x0a, 0x05, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x4d,
	0x6f, 0x64, 0x65, 0x6c, 0x12, 0x1c, 0x0a, 0x09, 0x4f, 0x73, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x4f, 0x73, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12, 0x1a, 0x0a, 0x08,
	0x42, 0x53, 0x64, 0x6b, 0x55, 0x44, 0x69, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x42, 0x53, 0x64, 0x6b, 0x55, 0x44, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x42, 0x47, 0x61, 0x6d,
	0x65, 0x49, 0x44, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x42, 0x47, 0x61, 0x6d, 0x65,
	0x49, 0x44, 0x12, 0x1e, 0x0a, 0x0a, 0x42, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x49, 0x44,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x42, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
	0x49, 0x44, 0x12, 0x1e, 0x0a, 0x0a, 0x42, 0x53, 0x6b, 0x64, 0x54, 0x79, 0x70, 0x65, 0x49, 0x44,
	0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x42, 0x53, 0x6b, 0x64, 0x54, 0x79, 0x70, 0x65,
	0x49, 0x44, 0x2a, 0xdb, 0x01, 0x0a, 0x07, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08,
	0x0a, 0x04, 0x4e, 0x55, 0x4c, 0x4c, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x45, 0x51, 0x55,
	0x45, 0x53, 0x54, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x45, 0x53, 0x50, 0x4f, 0x4e, 0x53,
	0x45, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06, 0x4e, 0x4f, 0x54, 0x49, 0x46, 0x59, 0x10, 0x03, 0x12,
	0x08, 0x0a, 0x04, 0x50, 0x55, 0x53, 0x48, 0x10, 0x04, 0x12, 0x0d, 0x0a, 0x09, 0x48, 0x45, 0x41,
	0x52, 0x54, 0x42, 0x45, 0x41, 0x54, 0x10, 0x05, 0x12, 0x10, 0x0a, 0x0c, 0x48, 0x45, 0x41, 0x52,
	0x54, 0x42, 0x45, 0x41, 0x54, 0x41, 0x43, 0x4b, 0x10, 0x06, 0x12, 0x0d, 0x0a, 0x09, 0x48, 0x41,
	0x4e, 0x44, 0x53, 0x48, 0x41, 0x4b, 0x45, 0x10, 0x07, 0x12, 0x10, 0x0a, 0x0c, 0x48, 0x41, 0x4e,
	0x44, 0x53, 0x48, 0x41, 0x4b, 0x45, 0x41, 0x43, 0x4b, 0x10, 0x08, 0x12, 0x09, 0x0a, 0x05, 0x45,
	0x52, 0x52, 0x4f, 0x52, 0x10, 0x09, 0x12, 0x08, 0x0a, 0x04, 0x4b, 0x49, 0x43, 0x4b, 0x10, 0x0a,
	0x12, 0x13, 0x0a, 0x0f, 0x4c, 0x4f, 0x43, 0x41, 0x4c, 0x5f, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43,
	0x54, 0x45, 0x44, 0x10, 0x14, 0x12, 0x16, 0x0a, 0x12, 0x4c, 0x4f, 0x43, 0x41, 0x4c, 0x5f, 0x44,
	0x49, 0x53, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54, 0x45, 0x44, 0x10, 0x15, 0x12, 0x11, 0x0a,
	0x0d, 0x4c, 0x4f, 0x43, 0x41, 0x4c, 0x5f, 0x54, 0x49, 0x4d, 0x45, 0x4f, 0x55, 0x54, 0x10, 0x16,
	0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c,
	0x69, 0x79, 0x69, 0x79, 0x73, 0x6e, 0x67, 0x2f, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2f,
	0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_node_message_proto_head_proto_rawDescOnce sync.Once
	file_node_message_proto_head_proto_rawDescData = file_node_message_proto_head_proto_rawDesc
)

func file_node_message_proto_head_proto_rawDescGZIP() []byte {
	file_node_message_proto_head_proto_rawDescOnce.Do(func() {
		file_node_message_proto_head_proto_rawDescData = protoimpl.X.CompressGZIP(file_node_message_proto_head_proto_rawDescData)
	})
	return file_node_message_proto_head_proto_rawDescData
}

var file_node_message_proto_head_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_node_message_proto_head_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_node_message_proto_head_proto_goTypes = []interface{}{
	(MsgType)(0),         // 0: scatter.node.proto.MsgType
	(*MsgError)(nil),     // 1: scatter.node.proto.MsgError
	(*Head)(nil),         // 2: scatter.node.proto.Head
	(*MsgHandShake)(nil), // 3: scatter.node.proto.MsgHandShake
}
var file_node_message_proto_head_proto_depIdxs = []int32{
	0, // 0: scatter.node.proto.Head.MsgType:type_name -> scatter.node.proto.MsgType
	1, // 1: scatter.node.proto.Head.Error:type_name -> scatter.node.proto.MsgError
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_node_message_proto_head_proto_init() }
func file_node_message_proto_head_proto_init() {
	if File_node_message_proto_head_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_node_message_proto_head_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgError); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_node_message_proto_head_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Head); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_node_message_proto_head_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgHandShake); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_node_message_proto_head_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_node_message_proto_head_proto_goTypes,
		DependencyIndexes: file_node_message_proto_head_proto_depIdxs,
		EnumInfos:         file_node_message_proto_head_proto_enumTypes,
		MessageInfos:      file_node_message_proto_head_proto_msgTypes,
	}.Build()
	File_node_message_proto_head_proto = out.File
	file_node_message_proto_head_proto_rawDesc = nil
	file_node_message_proto_head_proto_goTypes = nil
	file_node_message_proto_head_proto_depIdxs = nil
}
