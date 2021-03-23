// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package sessionpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// SesssionClient is the client API for Sesssion service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SesssionClient interface {
	Create(ctx context.Context, in *SessionInfo, opts ...grpc.CallOption) (*CreateSessionRes, error)
}

type sesssionClient struct {
	cc grpc.ClientConnInterface
}

func NewSesssionClient(cc grpc.ClientConnInterface) SesssionClient {
	return &sesssionClient{cc}
}

func (c *sesssionClient) Create(ctx context.Context, in *SessionInfo, opts ...grpc.CallOption) (*CreateSessionRes, error) {
	out := new(CreateSessionRes)
	err := c.cc.Invoke(ctx, "/scatter.service.Sesssion/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SesssionServer is the server API for Sesssion service.
// All implementations must embed UnimplementedSesssionServer
// for forward compatibility
type SesssionServer interface {
	Create(context.Context, *SessionInfo) (*CreateSessionRes, error)
	mustEmbedUnimplementedSesssionServer()
}

// UnimplementedSesssionServer must be embedded to have forward compatible implementations.
type UnimplementedSesssionServer struct {
}

func (UnimplementedSesssionServer) Create(context.Context, *SessionInfo) (*CreateSessionRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedSesssionServer) mustEmbedUnimplementedSesssionServer() {}

// UnsafeSesssionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SesssionServer will
// result in compilation errors.
type UnsafeSesssionServer interface {
	mustEmbedUnimplementedSesssionServer()
}

func RegisterSesssionServer(s grpc.ServiceRegistrar, srv SesssionServer) {
	s.RegisterService(&_Sesssion_serviceDesc, srv)
}

func _Sesssion_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SesssionServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scatter.service.Sesssion/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SesssionServer).Create(ctx, req.(*SessionInfo))
	}
	return interceptor(ctx, in, info, handler)
}

var _Sesssion_serviceDesc = grpc.ServiceDesc{
	ServiceName: "scatter.service.Sesssion",
	HandlerType: (*SesssionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Sesssion_Create_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cluster/session_service.proto",
}
