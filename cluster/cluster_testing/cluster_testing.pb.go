// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: testing/cluster_testing.proto

package cluster_testing

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

type String struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Str string `protobuf:"bytes,1,opt,name=str,proto3" json:"str,omitempty"`
}

func (x *String) Reset() {
	*x = String{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testing_cluster_testing_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *String) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*String) ProtoMessage() {}

func (x *String) ProtoReflect() protoreflect.Message {
	mi := &file_testing_cluster_testing_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use String.ProtoReflect.Descriptor instead.
func (*String) Descriptor() ([]byte, []int) {
	return file_testing_cluster_testing_proto_rawDescGZIP(), []int{0}
}

func (x *String) GetStr() string {
	if x != nil {
		return x.Str
	}
	return ""
}

type StringS struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Strs []string `protobuf:"bytes,1,rep,name=strs,proto3" json:"strs,omitempty"`
}

func (x *StringS) Reset() {
	*x = StringS{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testing_cluster_testing_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StringS) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringS) ProtoMessage() {}

func (x *StringS) ProtoReflect() protoreflect.Message {
	mi := &file_testing_cluster_testing_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringS.ProtoReflect.Descriptor instead.
func (*StringS) Descriptor() ([]byte, []int) {
	return file_testing_cluster_testing_proto_rawDescGZIP(), []int{1}
}

func (x *StringS) GetStrs() []string {
	if x != nil {
		return x.Strs
	}
	return nil
}

type Int struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	I int64 `protobuf:"varint,1,opt,name=i,proto3" json:"i,omitempty"`
}

func (x *Int) Reset() {
	*x = Int{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testing_cluster_testing_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Int) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Int) ProtoMessage() {}

func (x *Int) ProtoReflect() protoreflect.Message {
	mi := &file_testing_cluster_testing_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Int.ProtoReflect.Descriptor instead.
func (*Int) Descriptor() ([]byte, []int) {
	return file_testing_cluster_testing_proto_rawDescGZIP(), []int{2}
}

func (x *Int) GetI() int64 {
	if x != nil {
		return x.I
	}
	return 0
}

type Ints struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	I []int64 `protobuf:"varint,1,rep,packed,name=i,proto3" json:"i,omitempty"`
}

func (x *Ints) Reset() {
	*x = Ints{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testing_cluster_testing_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ints) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ints) ProtoMessage() {}

func (x *Ints) ProtoReflect() protoreflect.Message {
	mi := &file_testing_cluster_testing_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ints.ProtoReflect.Descriptor instead.
func (*Ints) Descriptor() ([]byte, []int) {
	return file_testing_cluster_testing_proto_rawDescGZIP(), []int{3}
}

func (x *Ints) GetI() []int64 {
	if x != nil {
		return x.I
	}
	return nil
}

var File_testing_cluster_testing_proto protoreflect.FileDescriptor

var file_testing_cluster_testing_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x2f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x5f, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0f, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x22, 0x1a, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x74,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x74, 0x72, 0x22, 0x1d, 0x0a, 0x07,
	0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x74, 0x72, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x73, 0x74, 0x72, 0x73, 0x22, 0x13, 0x0a, 0x03, 0x49,
	0x6e, 0x74, 0x12, 0x0c, 0x0a, 0x01, 0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x01, 0x69,
	0x22, 0x14, 0x0a, 0x04, 0x49, 0x6e, 0x74, 0x73, 0x12, 0x0c, 0x0a, 0x01, 0x69, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x03, 0x52, 0x01, 0x69, 0x32, 0xc2, 0x01, 0x0a, 0x0a, 0x53, 0x72, 0x76, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x3b, 0x0a, 0x07, 0x54, 0x6f, 0x4c, 0x6f, 0x77, 0x65, 0x72,
	0x12, 0x17, 0x2e, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x1a, 0x17, 0x2e, 0x73, 0x63, 0x61, 0x74,
	0x74, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x12, 0x3b, 0x0a, 0x07, 0x54, 0x6f, 0x55, 0x70, 0x70, 0x65, 0x72, 0x12, 0x17, 0x2e,
	0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x1a, 0x17, 0x2e, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x12,
	0x3a, 0x0a, 0x05, 0x53, 0x70, 0x6c, 0x69, 0x74, 0x12, 0x17, 0x2e, 0x73, 0x63, 0x61, 0x74, 0x74,
	0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x1a, 0x18, 0x2e, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x32, 0x73, 0x0a, 0x07, 0x53,
	0x72, 0x76, 0x49, 0x6e, 0x74, 0x73, 0x12, 0x32, 0x0a, 0x03, 0x53, 0x75, 0x6d, 0x12, 0x15, 0x2e,
	0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x49, 0x6e, 0x74, 0x73, 0x1a, 0x14, 0x2e, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x49, 0x6e, 0x74, 0x12, 0x34, 0x0a, 0x05, 0x4d, 0x75,
	0x6c, 0x74, 0x69, 0x12, 0x15, 0x2e, 0x73, 0x63, 0x61, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x49, 0x6e, 0x74, 0x73, 0x1a, 0x14, 0x2e, 0x73, 0x63, 0x61,
	0x74, 0x74, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x49, 0x6e, 0x74,
	0x42, 0x19, 0x5a, 0x17, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x63, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x5f, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_testing_cluster_testing_proto_rawDescOnce sync.Once
	file_testing_cluster_testing_proto_rawDescData = file_testing_cluster_testing_proto_rawDesc
)

func file_testing_cluster_testing_proto_rawDescGZIP() []byte {
	file_testing_cluster_testing_proto_rawDescOnce.Do(func() {
		file_testing_cluster_testing_proto_rawDescData = protoimpl.X.CompressGZIP(file_testing_cluster_testing_proto_rawDescData)
	})
	return file_testing_cluster_testing_proto_rawDescData
}

var file_testing_cluster_testing_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_testing_cluster_testing_proto_goTypes = []interface{}{
	(*String)(nil),  // 0: scatter.service.String
	(*StringS)(nil), // 1: scatter.service.StringS
	(*Int)(nil),     // 2: scatter.service.Int
	(*Ints)(nil),    // 3: scatter.service.Ints
}
var file_testing_cluster_testing_proto_depIdxs = []int32{
	0, // 0: scatter.service.SrvStrings.ToLower:input_type -> scatter.service.String
	0, // 1: scatter.service.SrvStrings.ToUpper:input_type -> scatter.service.String
	0, // 2: scatter.service.SrvStrings.Split:input_type -> scatter.service.String
	3, // 3: scatter.service.SrvInts.Sum:input_type -> scatter.service.Ints
	3, // 4: scatter.service.SrvInts.Multi:input_type -> scatter.service.Ints
	0, // 5: scatter.service.SrvStrings.ToLower:output_type -> scatter.service.String
	0, // 6: scatter.service.SrvStrings.ToUpper:output_type -> scatter.service.String
	1, // 7: scatter.service.SrvStrings.Split:output_type -> scatter.service.StringS
	2, // 8: scatter.service.SrvInts.Sum:output_type -> scatter.service.Int
	2, // 9: scatter.service.SrvInts.Multi:output_type -> scatter.service.Int
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_testing_cluster_testing_proto_init() }
func file_testing_cluster_testing_proto_init() {
	if File_testing_cluster_testing_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_testing_cluster_testing_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*String); i {
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
		file_testing_cluster_testing_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StringS); i {
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
		file_testing_cluster_testing_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Int); i {
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
		file_testing_cluster_testing_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ints); i {
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
			RawDescriptor: file_testing_cluster_testing_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_testing_cluster_testing_proto_goTypes,
		DependencyIndexes: file_testing_cluster_testing_proto_depIdxs,
		MessageInfos:      file_testing_cluster_testing_proto_msgTypes,
	}.Build()
	File_testing_cluster_testing_proto = out.File
	file_testing_cluster_testing_proto_rawDesc = nil
	file_testing_cluster_testing_proto_goTypes = nil
	file_testing_cluster_testing_proto_depIdxs = nil
}
