// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: cluster/sessionpb/session_service.proto

package sessionpb

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

type SessionType int32

const (
	SessionType_FrontEnd SessionType = 0 // 前端
	SessionType_Transfer SessionType = 1 // 转发
	SessionType_Backend  SessionType = 2 // 后端
	SessionType_Pub      SessionType = 3 // 由发布者
)

// Enum value maps for SessionType.
var (
	SessionType_name = map[int32]string{
		0: "FrontEnd",
		1: "Transfer",
		2: "Backend",
		3: "Pub",
	}
	SessionType_value = map[string]int32{
		"FrontEnd": 0,
		"Transfer": 1,
		"Backend":  2,
		"Pub":      3,
	}
)

func (x SessionType) Enum() *SessionType {
	p := new(SessionType)
	*p = x
	return p
}

func (x SessionType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SessionType) Descriptor() protoreflect.EnumDescriptor {
	return file_cluster_sessionpb_session_service_proto_enumTypes[0].Descriptor()
}

func (SessionType) Type() protoreflect.EnumType {
	return &file_cluster_sessionpb_session_service_proto_enumTypes[0]
}

func (x SessionType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SessionType.Descriptor instead.
func (SessionType) EnumDescriptor() ([]byte, []int) {
	return file_cluster_sessionpb_session_service_proto_rawDescGZIP(), []int{0}
}

type TransferInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UID int64 `protobuf:"varint,1,opt,name=UID,proto3" json:"UID,omitempty"`
}

func (x *TransferInfo) Reset() {
	*x = TransferInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_sessionpb_session_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransferInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransferInfo) ProtoMessage() {}

func (x *TransferInfo) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_sessionpb_session_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransferInfo.ProtoReflect.Descriptor instead.
func (*TransferInfo) Descriptor() ([]byte, []int) {
	return file_cluster_sessionpb_session_service_proto_rawDescGZIP(), []int{0}
}

func (x *TransferInfo) GetUID() int64 {
	if x != nil {
		return x.UID
	}
	return 0
}

type FrontEndInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UID int64 `protobuf:"varint,1,opt,name=UID,proto3" json:"UID,omitempty"`
	SID int64 `protobuf:"varint,2,opt,name=SID,proto3" json:"SID,omitempty"`
}

func (x *FrontEndInfo) Reset() {
	*x = FrontEndInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_sessionpb_session_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FrontEndInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FrontEndInfo) ProtoMessage() {}

func (x *FrontEndInfo) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_sessionpb_session_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FrontEndInfo.ProtoReflect.Descriptor instead.
func (*FrontEndInfo) Descriptor() ([]byte, []int) {
	return file_cluster_sessionpb_session_service_proto_rawDescGZIP(), []int{1}
}

func (x *FrontEndInfo) GetUID() int64 {
	if x != nil {
		return x.UID
	}
	return 0
}

func (x *FrontEndInfo) GetSID() int64 {
	if x != nil {
		return x.SID
	}
	return 0
}

type BackendInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *BackendInfo) Reset() {
	*x = BackendInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_sessionpb_session_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BackendInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BackendInfo) ProtoMessage() {}

func (x *BackendInfo) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_sessionpb_session_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BackendInfo.ProtoReflect.Descriptor instead.
func (*BackendInfo) Descriptor() ([]byte, []int) {
	return file_cluster_sessionpb_session_service_proto_rawDescGZIP(), []int{2}
}

type PubInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PubInfo) Reset() {
	*x = PubInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_sessionpb_session_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PubInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PubInfo) ProtoMessage() {}

func (x *PubInfo) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_sessionpb_session_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PubInfo.ProtoReflect.Descriptor instead.
func (*PubInfo) Descriptor() ([]byte, []int) {
	return file_cluster_sessionpb_session_service_proto_rawDescGZIP(), []int{3}
}

type SessionInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SType        SessionType   `protobuf:"varint,1,opt,name=SType,proto3,enum=scatter.service.SessionType" json:"SType,omitempty"`
	TransferInfo *TransferInfo `protobuf:"bytes,2,opt,name=TransferInfo,proto3" json:"TransferInfo,omitempty"`
	FrontEndInfo *FrontEndInfo `protobuf:"bytes,3,opt,name=FrontEndInfo,proto3" json:"FrontEndInfo,omitempty"`
	BackendInfo  *BackendInfo  `protobuf:"bytes,4,opt,name=BackendInfo,proto3" json:"BackendInfo,omitempty"`
	PubInfo      *PubInfo      `protobuf:"bytes,5,opt,name=PubInfo,proto3" json:"PubInfo,omitempty"`
}

func (x *SessionInfo) Reset() {
	*x = SessionInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_sessionpb_session_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionInfo) ProtoMessage() {}

func (x *SessionInfo) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_sessionpb_session_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionInfo.ProtoReflect.Descriptor instead.
func (*SessionInfo) Descriptor() ([]byte, []int) {
	return file_cluster_sessionpb_session_service_proto_rawDescGZIP(), []int{4}
}

func (x *SessionInfo) GetSType() SessionType {
	if x != nil {
		return x.SType
	}
	return SessionType_FrontEnd
}

func (x *SessionInfo) GetTransferInfo() *TransferInfo {
	if x != nil {
		return x.TransferInfo
	}
	return nil
}

func (x *SessionInfo) GetFrontEndInfo() *FrontEndInfo {
	if x != nil {
		return x.FrontEndInfo
	}
	return nil
}

func (x *SessionInfo) GetBackendInfo() *BackendInfo {
	if x != nil {
		return x.BackendInfo
	}
	return nil
}

func (x *SessionInfo) GetPubInfo() *PubInfo {
	if x != nil {
		return x.PubInfo
	}
	return nil
}

type CreateSessionReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Info *SessionInfo `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty"`
}

func (x *CreateSessionReq) Reset() {
	*x = CreateSessionReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_sessionpb_session_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateSessionReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSessionReq) ProtoMessage() {}

func (x *CreateSessionReq) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_sessionpb_session_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSessionReq.ProtoReflect.Descriptor instead.
func (*CreateSessionReq) Descriptor() ([]byte, []int) {
	return file_cluster_sessionpb_session_service_proto_rawDescGZIP(), []int{5}
}

func (x *CreateSessionReq) GetInfo() *SessionInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

type CreateSessionRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Info *SessionInfo `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty"`
}

func (x *CreateSessionRes) Reset() {
	*x = CreateSessionRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_sessionpb_session_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateSessionRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSessionRes) ProtoMessage() {}

func (x *CreateSessionRes) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_sessionpb_session_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSessionRes.ProtoReflect.Descriptor instead.
func (*CreateSessionRes) Descriptor() ([]byte, []int) {
	return file_cluster_sessionpb_session_service_proto_rawDescGZIP(), []int{6}
}

func (x *CreateSessionRes) GetInfo() *SessionInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

var File_cluster_sessionpb_session_service_proto protoreflect.FileDescriptor

var file_cluster_sessionpb_session_service_proto_rawDesc = []byte{
	0x0a, 0x27, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x70, 0x62, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x73, 0x63, 0x61, 0x74, 0x74,
	0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0x20, 0x0a, 0x0c, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x10, 0x0a, 0x03, 0x55, 0x49,
	0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x55, 0x49, 0x44, 0x22, 0x32, 0x0a, 0x0c,
	0x46, 0x72, 0x6f, 0x6e, 0x74, 0x45, 0x6e, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x10, 0x0a, 0x03,
	0x55, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x55, 0x49, 0x44, 0x12, 0x10,
	0x0a, 0x03, 0x53, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x53, 0x49, 0x44,
	0x22, 0x0d, 0x0a, 0x0b, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x22,
	0x09, 0x0a, 0x07, 0x50, 0x75, 0x62, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0xbb, 0x02, 0x0a, 0x0b, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x32, 0x0a, 0x05, 0x53, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x73, 0x63, 0x61, 0x74,
	0x74, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x05, 0x53, 0x54, 0x79, 0x70, 0x65, 0x12, 0x41,
	0x0a, 0x0c, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x0c, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x41, 0x0a, 0x0c, 0x46, 0x72, 0x6f, 0x6e, 0x74, 0x45, 0x6e, 0x64, 0x49, 0x6e, 0x66,
	0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65,
	0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x46, 0x72, 0x6f, 0x6e, 0x74, 0x45,
	0x6e, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0c, 0x46, 0x72, 0x6f, 0x6e, 0x74, 0x45, 0x6e, 0x64,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x3e, 0x0a, 0x0b, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x49,
	0x6e, 0x66, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x73, 0x63, 0x61, 0x74,
	0x74, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x42, 0x61, 0x63, 0x6b,
	0x65, 0x6e, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0b, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x32, 0x0a, 0x07, 0x50, 0x75, 0x62, 0x49, 0x6e, 0x66, 0x6f, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x50, 0x75, 0x62, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x07, 0x50, 0x75, 0x62, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x44, 0x0a, 0x10, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x12, 0x30, 0x0a, 0x04,
	0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x73, 0x63, 0x61,
	0x74, 0x74, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x22, 0x44,
	0x0a, 0x10, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x73, 0x12, 0x30, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1c, 0x2e, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x04,
	0x69, 0x6e, 0x66, 0x6f, 0x2a, 0x3f, 0x0a, 0x0b, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x0c, 0x0a, 0x08, 0x46, 0x72, 0x6f, 0x6e, 0x74, 0x45, 0x6e, 0x64, 0x10,
	0x00, 0x12, 0x0c, 0x0a, 0x08, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x10, 0x01, 0x12,
	0x0b, 0x0a, 0x07, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x10, 0x02, 0x12, 0x07, 0x0a, 0x03,
	0x50, 0x75, 0x62, 0x10, 0x03, 0x32, 0x55, 0x0a, 0x08, 0x53, 0x65, 0x73, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x49, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x1c, 0x2e, 0x73, 0x63,
	0x61, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x21, 0x2e, 0x73, 0x63, 0x61, 0x74,
	0x74, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x42, 0x2f, 0x5a, 0x2d,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x69, 0x79, 0x69, 0x79,
	0x73, 0x6e, 0x67, 0x2f, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2f, 0x63, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cluster_sessionpb_session_service_proto_rawDescOnce sync.Once
	file_cluster_sessionpb_session_service_proto_rawDescData = file_cluster_sessionpb_session_service_proto_rawDesc
)

func file_cluster_sessionpb_session_service_proto_rawDescGZIP() []byte {
	file_cluster_sessionpb_session_service_proto_rawDescOnce.Do(func() {
		file_cluster_sessionpb_session_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_cluster_sessionpb_session_service_proto_rawDescData)
	})
	return file_cluster_sessionpb_session_service_proto_rawDescData
}

var file_cluster_sessionpb_session_service_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_cluster_sessionpb_session_service_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_cluster_sessionpb_session_service_proto_goTypes = []interface{}{
	(SessionType)(0),         // 0: scatter.service.SessionType
	(*TransferInfo)(nil),     // 1: scatter.service.TransferInfo
	(*FrontEndInfo)(nil),     // 2: scatter.service.FrontEndInfo
	(*BackendInfo)(nil),      // 3: scatter.service.BackendInfo
	(*PubInfo)(nil),          // 4: scatter.service.PubInfo
	(*SessionInfo)(nil),      // 5: scatter.service.SessionInfo
	(*CreateSessionReq)(nil), // 6: scatter.service.CreateSessionReq
	(*CreateSessionRes)(nil), // 7: scatter.service.CreateSessionRes
}
var file_cluster_sessionpb_session_service_proto_depIdxs = []int32{
	0, // 0: scatter.service.SessionInfo.SType:type_name -> scatter.service.SessionType
	1, // 1: scatter.service.SessionInfo.TransferInfo:type_name -> scatter.service.TransferInfo
	2, // 2: scatter.service.SessionInfo.FrontEndInfo:type_name -> scatter.service.FrontEndInfo
	3, // 3: scatter.service.SessionInfo.BackendInfo:type_name -> scatter.service.BackendInfo
	4, // 4: scatter.service.SessionInfo.PubInfo:type_name -> scatter.service.PubInfo
	5, // 5: scatter.service.CreateSessionReq.info:type_name -> scatter.service.SessionInfo
	5, // 6: scatter.service.CreateSessionRes.info:type_name -> scatter.service.SessionInfo
	5, // 7: scatter.service.Sesssion.Create:input_type -> scatter.service.SessionInfo
	7, // 8: scatter.service.Sesssion.Create:output_type -> scatter.service.CreateSessionRes
	8, // [8:9] is the sub-list for method output_type
	7, // [7:8] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_cluster_sessionpb_session_service_proto_init() }
func file_cluster_sessionpb_session_service_proto_init() {
	if File_cluster_sessionpb_session_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cluster_sessionpb_session_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TransferInfo); i {
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
		file_cluster_sessionpb_session_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FrontEndInfo); i {
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
		file_cluster_sessionpb_session_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BackendInfo); i {
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
		file_cluster_sessionpb_session_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PubInfo); i {
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
		file_cluster_sessionpb_session_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SessionInfo); i {
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
		file_cluster_sessionpb_session_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateSessionReq); i {
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
		file_cluster_sessionpb_session_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateSessionRes); i {
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
			RawDescriptor: file_cluster_sessionpb_session_service_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_cluster_sessionpb_session_service_proto_goTypes,
		DependencyIndexes: file_cluster_sessionpb_session_service_proto_depIdxs,
		EnumInfos:         file_cluster_sessionpb_session_service_proto_enumTypes,
		MessageInfos:      file_cluster_sessionpb_session_service_proto_msgTypes,
	}.Build()
	File_cluster_sessionpb_session_service_proto = out.File
	file_cluster_sessionpb_session_service_proto_rawDesc = nil
	file_cluster_sessionpb_session_service_proto_goTypes = nil
	file_cluster_sessionpb_session_service_proto_depIdxs = nil
}