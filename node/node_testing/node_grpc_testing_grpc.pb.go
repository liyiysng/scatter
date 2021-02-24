// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package node_testing

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// GrpcSrvTestClient is the client API for GrpcSrvTest service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GrpcSrvTestClient interface {
	GetSum(ctx context.Context, in *GRPCSumReq, opts ...grpc.CallOption) (*GRPCSumRes, error)
}

type grpcSrvTestClient struct {
	cc grpc.ClientConnInterface
}

func NewGrpcSrvTestClient(cc grpc.ClientConnInterface) GrpcSrvTestClient {
	return &grpcSrvTestClient{cc}
}

func (c *grpcSrvTestClient) GetSum(ctx context.Context, in *GRPCSumReq, opts ...grpc.CallOption) (*GRPCSumRes, error) {
	out := new(GRPCSumRes)
	err := c.cc.Invoke(ctx, "/GrpcSrvTest/GetSum", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GrpcSrvTestServer is the server API for GrpcSrvTest service.
// All implementations must embed UnimplementedGrpcSrvTestServer
// for forward compatibility
type GrpcSrvTestServer interface {
	GetSum(context.Context, *GRPCSumReq) (*GRPCSumRes, error)
	mustEmbedUnimplementedGrpcSrvTestServer()
}

// UnimplementedGrpcSrvTestServer must be embedded to have forward compatible implementations.
type UnimplementedGrpcSrvTestServer struct {
}

func (UnimplementedGrpcSrvTestServer) GetSum(context.Context, *GRPCSumReq) (*GRPCSumRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSum not implemented")
}
func (UnimplementedGrpcSrvTestServer) mustEmbedUnimplementedGrpcSrvTestServer() {}

// UnsafeGrpcSrvTestServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GrpcSrvTestServer will
// result in compilation errors.
type UnsafeGrpcSrvTestServer interface {
	mustEmbedUnimplementedGrpcSrvTestServer()
}

func RegisterGrpcSrvTestServer(s grpc.ServiceRegistrar, srv GrpcSrvTestServer) {
	s.RegisterService(&_GrpcSrvTest_serviceDesc, srv)
}

func _GrpcSrvTest_GetSum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GRPCSumReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcSrvTestServer).GetSum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/GrpcSrvTest/GetSum",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcSrvTestServer).GetSum(ctx, req.(*GRPCSumReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _GrpcSrvTest_serviceDesc = grpc.ServiceDesc{
	ServiceName: "GrpcSrvTest",
	HandlerType: (*GrpcSrvTestServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSum",
			Handler:    _GrpcSrvTest_GetSum_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "testing/node_grpc_testing.proto",
}
