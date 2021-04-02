// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: cluster/subsrvpb/sub_service.proto

package subsrvpb

import (
	proto "github.com/golang/protobuf/proto"
	sessionpb "github.com/liyiysng/scatter/cluster/sessionpb"
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

type ErrorInfo_ErrorType int32

const (
	ErrorInfo_ErrorTypeCustom   ErrorInfo_ErrorType = 0
	ErrorInfo_ErrorTypeCritical ErrorInfo_ErrorType = 1
	ErrorInfo_ErrorTypeCommon   ErrorInfo_ErrorType = 2
)

// Enum value maps for ErrorInfo_ErrorType.
var (
	ErrorInfo_ErrorType_name = map[int32]string{
		0: "ErrorTypeCustom",
		1: "ErrorTypeCritical",
		2: "ErrorTypeCommon",
	}
	ErrorInfo_ErrorType_value = map[string]int32{
		"ErrorTypeCustom":   0,
		"ErrorTypeCritical": 1,
		"ErrorTypeCommon":   2,
	}
)

func (x ErrorInfo_ErrorType) Enum() *ErrorInfo_ErrorType {
	p := new(ErrorInfo_ErrorType)
	*p = x
	return p
}

func (x ErrorInfo_ErrorType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorInfo_ErrorType) Descriptor() protoreflect.EnumDescriptor {
	return file_cluster_subsrvpb_sub_service_proto_enumTypes[0].Descriptor()
}

func (ErrorInfo_ErrorType) Type() protoreflect.EnumType {
	return &file_cluster_subsrvpb_sub_service_proto_enumTypes[0]
}

func (x ErrorInfo_ErrorType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorInfo_ErrorType.Descriptor instead.
func (ErrorInfo_ErrorType) EnumDescriptor() ([]byte, []int) {
	return file_cluster_subsrvpb_sub_service_proto_rawDescGZIP(), []int{0, 0}
}

type ErrorInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrType ErrorInfo_ErrorType `protobuf:"varint,1,opt,name=errType,proto3,enum=scatter.service.ErrorInfo_ErrorType" json:"errType,omitempty"`
	Err     string              `protobuf:"bytes,2,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *ErrorInfo) Reset() {
	*x = ErrorInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_subsrvpb_sub_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ErrorInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrorInfo) ProtoMessage() {}

func (x *ErrorInfo) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_subsrvpb_sub_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ErrorInfo.ProtoReflect.Descriptor instead.
func (*ErrorInfo) Descriptor() ([]byte, []int) {
	return file_cluster_subsrvpb_sub_service_proto_rawDescGZIP(), []int{0}
}

func (x *ErrorInfo) GetErrType() ErrorInfo_ErrorType {
	if x != nil {
		return x.ErrType
	}
	return ErrorInfo_ErrorTypeCustom
}

func (x *ErrorInfo) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

type CallReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sinfo       *sessionpb.SessionInfo `protobuf:"bytes,1,opt,name=sinfo,proto3" json:"sinfo,omitempty"`
	ServiceName string                 `protobuf:"bytes,2,opt,name=serviceName,proto3" json:"serviceName,omitempty"`
	MethodName  string                 `protobuf:"bytes,3,opt,name=methodName,proto3" json:"methodName,omitempty"`
	Payload     []byte                 `protobuf:"bytes,4,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *CallReq) Reset() {
	*x = CallReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_subsrvpb_sub_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CallReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CallReq) ProtoMessage() {}

func (x *CallReq) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_subsrvpb_sub_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CallReq.ProtoReflect.Descriptor instead.
func (*CallReq) Descriptor() ([]byte, []int) {
	return file_cluster_subsrvpb_sub_service_proto_rawDescGZIP(), []int{1}
}

func (x *CallReq) GetSinfo() *sessionpb.SessionInfo {
	if x != nil {
		return x.Sinfo
	}
	return nil
}

func (x *CallReq) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

func (x *CallReq) GetMethodName() string {
	if x != nil {
		return x.MethodName
	}
	return ""
}

func (x *CallReq) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type CallRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrInfo *ErrorInfo `protobuf:"bytes,1,opt,name=errInfo,proto3" json:"errInfo,omitempty"`
	Payload []byte     `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *CallRes) Reset() {
	*x = CallRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_subsrvpb_sub_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CallRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CallRes) ProtoMessage() {}

func (x *CallRes) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_subsrvpb_sub_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CallRes.ProtoReflect.Descriptor instead.
func (*CallRes) Descriptor() ([]byte, []int) {
	return file_cluster_subsrvpb_sub_service_proto_rawDescGZIP(), []int{2}
}

func (x *CallRes) GetErrInfo() *ErrorInfo {
	if x != nil {
		return x.ErrInfo
	}
	return nil
}

func (x *CallRes) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type NotifyReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sinfo       *sessionpb.SessionInfo `protobuf:"bytes,1,opt,name=sinfo,proto3" json:"sinfo,omitempty"`
	ServiceName string                 `protobuf:"bytes,2,opt,name=serviceName,proto3" json:"serviceName,omitempty"`
	MethodName  string                 `protobuf:"bytes,3,opt,name=methodName,proto3" json:"methodName,omitempty"`
	Payload     []byte                 `protobuf:"bytes,4,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *NotifyReq) Reset() {
	*x = NotifyReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_subsrvpb_sub_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyReq) ProtoMessage() {}

func (x *NotifyReq) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_subsrvpb_sub_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyReq.ProtoReflect.Descriptor instead.
func (*NotifyReq) Descriptor() ([]byte, []int) {
	return file_cluster_subsrvpb_sub_service_proto_rawDescGZIP(), []int{3}
}

func (x *NotifyReq) GetSinfo() *sessionpb.SessionInfo {
	if x != nil {
		return x.Sinfo
	}
	return nil
}

func (x *NotifyReq) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

func (x *NotifyReq) GetMethodName() string {
	if x != nil {
		return x.MethodName
	}
	return ""
}

func (x *NotifyReq) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type NotifyRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrInfo *ErrorInfo `protobuf:"bytes,1,opt,name=errInfo,proto3" json:"errInfo,omitempty"`
}

func (x *NotifyRes) Reset() {
	*x = NotifyRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_subsrvpb_sub_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyRes) ProtoMessage() {}

func (x *NotifyRes) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_subsrvpb_sub_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyRes.ProtoReflect.Descriptor instead.
func (*NotifyRes) Descriptor() ([]byte, []int) {
	return file_cluster_subsrvpb_sub_service_proto_rawDescGZIP(), []int{4}
}

func (x *NotifyRes) GetErrInfo() *ErrorInfo {
	if x != nil {
		return x.ErrInfo
	}
	return nil
}

type PubReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sinfo   *sessionpb.SessionInfo `protobuf:"bytes,1,opt,name=sinfo,proto3" json:"sinfo,omitempty"`
	Topic   string                 `protobuf:"bytes,2,opt,name=topic,proto3" json:"topic,omitempty"`
	Cmd     string                 `protobuf:"bytes,3,opt,name=cmd,proto3" json:"cmd,omitempty"`
	Payload []byte                 `protobuf:"bytes,4,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *PubReq) Reset() {
	*x = PubReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_subsrvpb_sub_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PubReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PubReq) ProtoMessage() {}

func (x *PubReq) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_subsrvpb_sub_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PubReq.ProtoReflect.Descriptor instead.
func (*PubReq) Descriptor() ([]byte, []int) {
	return file_cluster_subsrvpb_sub_service_proto_rawDescGZIP(), []int{5}
}

func (x *PubReq) GetSinfo() *sessionpb.SessionInfo {
	if x != nil {
		return x.Sinfo
	}
	return nil
}

func (x *PubReq) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *PubReq) GetCmd() string {
	if x != nil {
		return x.Cmd
	}
	return ""
}

func (x *PubReq) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type PubRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrInfo *ErrorInfo `protobuf:"bytes,1,opt,name=errInfo,proto3" json:"errInfo,omitempty"`
}

func (x *PubRes) Reset() {
	*x = PubRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_subsrvpb_sub_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PubRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PubRes) ProtoMessage() {}

func (x *PubRes) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_subsrvpb_sub_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PubRes.ProtoReflect.Descriptor instead.
func (*PubRes) Descriptor() ([]byte, []int) {
	return file_cluster_subsrvpb_sub_service_proto_rawDescGZIP(), []int{6}
}

func (x *PubRes) GetErrInfo() *ErrorInfo {
	if x != nil {
		return x.ErrInfo
	}
	return nil
}

var File_cluster_subsrvpb_sub_service_proto protoreflect.FileDescriptor

var file_cluster_subsrvpb_sub_service_proto_rawDesc = []byte{
	0x0a, 0x22, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x73, 0x75, 0x62, 0x73, 0x72, 0x76,
	0x70, 0x62, 0x2f, 0x73, 0x75, 0x62, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x27, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x73,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xab,
	0x01, 0x0a, 0x09, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x3e, 0x0a, 0x07,
	0x65, 0x72, 0x72, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x24, 0x2e,
	0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x07, 0x65, 0x72, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x03,
	0x65, 0x72, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x72, 0x72, 0x22, 0x4c,
	0x0a, 0x09, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12, 0x13, 0x0a, 0x0f, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x54, 0x79, 0x70, 0x65, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x10, 0x00,
	0x12, 0x15, 0x0a, 0x11, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x54, 0x79, 0x70, 0x65, 0x43, 0x72, 0x69,
	0x74, 0x69, 0x63, 0x61, 0x6c, 0x10, 0x01, 0x12, 0x13, 0x0a, 0x0f, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x54, 0x79, 0x70, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x10, 0x02, 0x22, 0x99, 0x01, 0x0a,
	0x07, 0x43, 0x61, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x12, 0x32, 0x0a, 0x05, 0x73, 0x69, 0x6e, 0x66,
	0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65,
	0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x73, 0x69, 0x6e, 0x66, 0x6f, 0x12, 0x20, 0x0a, 0x0b,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1e,
	0x0a, 0x0a, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x59, 0x0a, 0x07, 0x43, 0x61, 0x6c, 0x6c,
	0x52, 0x65, 0x73, 0x12, 0x34, 0x0a, 0x07, 0x65, 0x72, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x07, 0x65, 0x72, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c,
	0x6f, 0x61, 0x64, 0x22, 0x9b, 0x01, 0x0a, 0x09, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x52, 0x65,
	0x71, 0x12, 0x32, 0x0a, 0x05, 0x73, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1c, 0x2e, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05,
	0x73, 0x69, 0x6e, 0x66, 0x6f, 0x12, 0x20, 0x0a, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x6d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f,
	0x61, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61,
	0x64, 0x22, 0x41, 0x0a, 0x09, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x52, 0x65, 0x73, 0x12, 0x34,
	0x0a, 0x07, 0x65, 0x72, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x07, 0x65, 0x72, 0x72,
	0x49, 0x6e, 0x66, 0x6f, 0x22, 0x7e, 0x0a, 0x06, 0x50, 0x75, 0x62, 0x52, 0x65, 0x71, 0x12, 0x32,
	0x0a, 0x05, 0x73, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x73, 0x69, 0x6e,
	0x66, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61,
	0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x22, 0x3e, 0x0a, 0x06, 0x50, 0x75, 0x62, 0x52, 0x65, 0x73, 0x12, 0x34,
	0x0a, 0x07, 0x65, 0x72, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x07, 0x65, 0x72, 0x72,
	0x49, 0x6e, 0x66, 0x6f, 0x32, 0xc3, 0x01, 0x0a, 0x0a, 0x53, 0x75, 0x62, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x3a, 0x0a, 0x04, 0x43, 0x61, 0x6c, 0x6c, 0x12, 0x18, 0x2e, 0x73, 0x63,
	0x61, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x61,
	0x6c, 0x6c, 0x52, 0x65, 0x71, 0x1a, 0x18, 0x2e, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x61, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x12,
	0x40, 0x0a, 0x06, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x12, 0x1a, 0x2e, 0x73, 0x63, 0x61, 0x74,
	0x74, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4e, 0x6f, 0x74, 0x69,
	0x66, 0x79, 0x52, 0x65, 0x71, 0x1a, 0x1a, 0x2e, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x52, 0x65,
	0x73, 0x12, 0x37, 0x0a, 0x03, 0x50, 0x75, 0x62, 0x12, 0x17, 0x2e, 0x73, 0x63, 0x61, 0x74, 0x74,
	0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x50, 0x75, 0x62, 0x52, 0x65,
	0x71, 0x1a, 0x17, 0x2e, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x50, 0x75, 0x62, 0x52, 0x65, 0x73, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x69, 0x79, 0x69, 0x79, 0x73, 0x6e,
	0x67, 0x2f, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x2f, 0x73, 0x75, 0x62, 0x73, 0x72, 0x76, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_cluster_subsrvpb_sub_service_proto_rawDescOnce sync.Once
	file_cluster_subsrvpb_sub_service_proto_rawDescData = file_cluster_subsrvpb_sub_service_proto_rawDesc
)

func file_cluster_subsrvpb_sub_service_proto_rawDescGZIP() []byte {
	file_cluster_subsrvpb_sub_service_proto_rawDescOnce.Do(func() {
		file_cluster_subsrvpb_sub_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_cluster_subsrvpb_sub_service_proto_rawDescData)
	})
	return file_cluster_subsrvpb_sub_service_proto_rawDescData
}

var file_cluster_subsrvpb_sub_service_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_cluster_subsrvpb_sub_service_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_cluster_subsrvpb_sub_service_proto_goTypes = []interface{}{
	(ErrorInfo_ErrorType)(0),      // 0: scatter.service.ErrorInfo.ErrorType
	(*ErrorInfo)(nil),             // 1: scatter.service.ErrorInfo
	(*CallReq)(nil),               // 2: scatter.service.CallReq
	(*CallRes)(nil),               // 3: scatter.service.CallRes
	(*NotifyReq)(nil),             // 4: scatter.service.NotifyReq
	(*NotifyRes)(nil),             // 5: scatter.service.NotifyRes
	(*PubReq)(nil),                // 6: scatter.service.PubReq
	(*PubRes)(nil),                // 7: scatter.service.PubRes
	(*sessionpb.SessionInfo)(nil), // 8: scatter.service.SessionInfo
}
var file_cluster_subsrvpb_sub_service_proto_depIdxs = []int32{
	0,  // 0: scatter.service.ErrorInfo.errType:type_name -> scatter.service.ErrorInfo.ErrorType
	8,  // 1: scatter.service.CallReq.sinfo:type_name -> scatter.service.SessionInfo
	1,  // 2: scatter.service.CallRes.errInfo:type_name -> scatter.service.ErrorInfo
	8,  // 3: scatter.service.NotifyReq.sinfo:type_name -> scatter.service.SessionInfo
	1,  // 4: scatter.service.NotifyRes.errInfo:type_name -> scatter.service.ErrorInfo
	8,  // 5: scatter.service.PubReq.sinfo:type_name -> scatter.service.SessionInfo
	1,  // 6: scatter.service.PubRes.errInfo:type_name -> scatter.service.ErrorInfo
	2,  // 7: scatter.service.SubService.Call:input_type -> scatter.service.CallReq
	4,  // 8: scatter.service.SubService.Notify:input_type -> scatter.service.NotifyReq
	6,  // 9: scatter.service.SubService.Pub:input_type -> scatter.service.PubReq
	3,  // 10: scatter.service.SubService.Call:output_type -> scatter.service.CallRes
	5,  // 11: scatter.service.SubService.Notify:output_type -> scatter.service.NotifyRes
	7,  // 12: scatter.service.SubService.Pub:output_type -> scatter.service.PubRes
	10, // [10:13] is the sub-list for method output_type
	7,  // [7:10] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_cluster_subsrvpb_sub_service_proto_init() }
func file_cluster_subsrvpb_sub_service_proto_init() {
	if File_cluster_subsrvpb_sub_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cluster_subsrvpb_sub_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ErrorInfo); i {
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
		file_cluster_subsrvpb_sub_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CallReq); i {
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
		file_cluster_subsrvpb_sub_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CallRes); i {
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
		file_cluster_subsrvpb_sub_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyReq); i {
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
		file_cluster_subsrvpb_sub_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyRes); i {
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
		file_cluster_subsrvpb_sub_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PubReq); i {
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
		file_cluster_subsrvpb_sub_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PubRes); i {
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
			RawDescriptor: file_cluster_subsrvpb_sub_service_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_cluster_subsrvpb_sub_service_proto_goTypes,
		DependencyIndexes: file_cluster_subsrvpb_sub_service_proto_depIdxs,
		EnumInfos:         file_cluster_subsrvpb_sub_service_proto_enumTypes,
		MessageInfos:      file_cluster_subsrvpb_sub_service_proto_msgTypes,
	}.Build()
	File_cluster_subsrvpb_sub_service_proto = out.File
	file_cluster_subsrvpb_sub_service_proto_rawDesc = nil
	file_cluster_subsrvpb_sub_service_proto_goTypes = nil
	file_cluster_subsrvpb_sub_service_proto_depIdxs = nil
}