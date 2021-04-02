// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package cluster_testing

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// SrvStringsClient is the client API for SrvStrings service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SrvStringsClient interface {
	ToLower(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error)
	ToUpper(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error)
	Split(ctx context.Context, in *String, opts ...grpc.CallOption) (*StringS, error)
}

type srvStringsClient struct {
	cc grpc.ClientConnInterface
}

func NewSrvStringsClient(cc grpc.ClientConnInterface) SrvStringsClient {
	return &srvStringsClient{cc}
}

func (c *srvStringsClient) ToLower(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error) {
	out := new(String)
	err := c.cc.Invoke(ctx, "/scatter.service.SrvStrings/ToLower", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *srvStringsClient) ToUpper(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error) {
	out := new(String)
	err := c.cc.Invoke(ctx, "/scatter.service.SrvStrings/ToUpper", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *srvStringsClient) Split(ctx context.Context, in *String, opts ...grpc.CallOption) (*StringS, error) {
	out := new(StringS)
	err := c.cc.Invoke(ctx, "/scatter.service.SrvStrings/Split", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SrvStringsServer is the server API for SrvStrings service.
// All implementations must embed UnimplementedSrvStringsServer
// for forward compatibility
type SrvStringsServer interface {
	ToLower(context.Context, *String) (*String, error)
	ToUpper(context.Context, *String) (*String, error)
	Split(context.Context, *String) (*StringS, error)
	mustEmbedUnimplementedSrvStringsServer()
}

// UnimplementedSrvStringsServer must be embedded to have forward compatible implementations.
type UnimplementedSrvStringsServer struct {
}

func (UnimplementedSrvStringsServer) ToLower(context.Context, *String) (*String, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ToLower not implemented")
}
func (UnimplementedSrvStringsServer) ToUpper(context.Context, *String) (*String, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ToUpper not implemented")
}
func (UnimplementedSrvStringsServer) Split(context.Context, *String) (*StringS, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Split not implemented")
}
func (UnimplementedSrvStringsServer) mustEmbedUnimplementedSrvStringsServer() {}

// UnsafeSrvStringsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SrvStringsServer will
// result in compilation errors.
type UnsafeSrvStringsServer interface {
	mustEmbedUnimplementedSrvStringsServer()
}

func RegisterSrvStringsServer(s grpc.ServiceRegistrar, srv SrvStringsServer) {
	s.RegisterService(&_SrvStrings_serviceDesc, srv)
}

func _SrvStrings_ToLower_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SrvStringsServer).ToLower(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scatter.service.SrvStrings/ToLower",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SrvStringsServer).ToLower(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

func _SrvStrings_ToUpper_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SrvStringsServer).ToUpper(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scatter.service.SrvStrings/ToUpper",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SrvStringsServer).ToUpper(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

func _SrvStrings_Split_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SrvStringsServer).Split(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scatter.service.SrvStrings/Split",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SrvStringsServer).Split(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

var _SrvStrings_serviceDesc = grpc.ServiceDesc{
	ServiceName: "scatter.service.SrvStrings",
	HandlerType: (*SrvStringsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ToLower",
			Handler:    _SrvStrings_ToLower_Handler,
		},
		{
			MethodName: "ToUpper",
			Handler:    _SrvStrings_ToUpper_Handler,
		},
		{
			MethodName: "Split",
			Handler:    _SrvStrings_Split_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "testing/cluster_testing.proto",
}

// SrvIntsClient is the client API for SrvInts service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SrvIntsClient interface {
	Sum(ctx context.Context, in *Ints, opts ...grpc.CallOption) (*Int, error)
	Multi(ctx context.Context, in *Ints, opts ...grpc.CallOption) (*Int, error)
}

type srvIntsClient struct {
	cc grpc.ClientConnInterface
}

func NewSrvIntsClient(cc grpc.ClientConnInterface) SrvIntsClient {
	return &srvIntsClient{cc}
}

func (c *srvIntsClient) Sum(ctx context.Context, in *Ints, opts ...grpc.CallOption) (*Int, error) {
	out := new(Int)
	err := c.cc.Invoke(ctx, "/scatter.service.SrvInts/Sum", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *srvIntsClient) Multi(ctx context.Context, in *Ints, opts ...grpc.CallOption) (*Int, error) {
	out := new(Int)
	err := c.cc.Invoke(ctx, "/scatter.service.SrvInts/Multi", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SrvIntsServer is the server API for SrvInts service.
// All implementations must embed UnimplementedSrvIntsServer
// for forward compatibility
type SrvIntsServer interface {
	Sum(context.Context, *Ints) (*Int, error)
	Multi(context.Context, *Ints) (*Int, error)
	mustEmbedUnimplementedSrvIntsServer()
}

// UnimplementedSrvIntsServer must be embedded to have forward compatible implementations.
type UnimplementedSrvIntsServer struct {
}

func (UnimplementedSrvIntsServer) Sum(context.Context, *Ints) (*Int, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sum not implemented")
}
func (UnimplementedSrvIntsServer) Multi(context.Context, *Ints) (*Int, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Multi not implemented")
}
func (UnimplementedSrvIntsServer) mustEmbedUnimplementedSrvIntsServer() {}

// UnsafeSrvIntsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SrvIntsServer will
// result in compilation errors.
type UnsafeSrvIntsServer interface {
	mustEmbedUnimplementedSrvIntsServer()
}

func RegisterSrvIntsServer(s grpc.ServiceRegistrar, srv SrvIntsServer) {
	s.RegisterService(&_SrvInts_serviceDesc, srv)
}

func _SrvInts_Sum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Ints)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SrvIntsServer).Sum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scatter.service.SrvInts/Sum",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SrvIntsServer).Sum(ctx, req.(*Ints))
	}
	return interceptor(ctx, in, info, handler)
}

func _SrvInts_Multi_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Ints)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SrvIntsServer).Multi(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scatter.service.SrvInts/Multi",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SrvIntsServer).Multi(ctx, req.(*Ints))
	}
	return interceptor(ctx, in, info, handler)
}

var _SrvInts_serviceDesc = grpc.ServiceDesc{
	ServiceName: "scatter.service.SrvInts",
	HandlerType: (*SrvIntsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Sum",
			Handler:    _SrvInts_Sum_Handler,
		},
		{
			MethodName: "Multi",
			Handler:    _SrvInts_Multi_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "testing/cluster_testing.proto",
}