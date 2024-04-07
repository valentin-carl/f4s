// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: pkg/controller/functionstate.proto

package controller

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// FunctionControlServiceClient is the client API for FunctionControlService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FunctionControlServiceClient interface {
	GetInstance(ctx context.Context, in *FunctionInstanceRequest, opts ...grpc.CallOption) (*FunctionInstanceResponse, error)
	NotifyStreamClosed(ctx context.Context, in *StreamClosedNotification, opts ...grpc.CallOption) (*StreamClosedNotificationResponse, error)
}

type functionControlServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFunctionControlServiceClient(cc grpc.ClientConnInterface) FunctionControlServiceClient {
	return &functionControlServiceClient{cc}
}

func (c *functionControlServiceClient) GetInstance(ctx context.Context, in *FunctionInstanceRequest, opts ...grpc.CallOption) (*FunctionInstanceResponse, error) {
	out := new(FunctionInstanceResponse)
	err := c.cc.Invoke(ctx, "/controller.FunctionControlService/GetInstance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *functionControlServiceClient) NotifyStreamClosed(ctx context.Context, in *StreamClosedNotification, opts ...grpc.CallOption) (*StreamClosedNotificationResponse, error) {
	out := new(StreamClosedNotificationResponse)
	err := c.cc.Invoke(ctx, "/controller.FunctionControlService/NotifyStreamClosed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FunctionControlServiceServer is the server API for FunctionControlService service.
// All implementations must embed UnimplementedFunctionControlServiceServer
// for forward compatibility
type FunctionControlServiceServer interface {
	GetInstance(context.Context, *FunctionInstanceRequest) (*FunctionInstanceResponse, error)
	NotifyStreamClosed(context.Context, *StreamClosedNotification) (*StreamClosedNotificationResponse, error)
	mustEmbedUnimplementedFunctionControlServiceServer()
}

// UnimplementedFunctionControlServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFunctionControlServiceServer struct {
}

func (UnimplementedFunctionControlServiceServer) GetInstance(context.Context, *FunctionInstanceRequest) (*FunctionInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInstance not implemented")
}
func (UnimplementedFunctionControlServiceServer) NotifyStreamClosed(context.Context, *StreamClosedNotification) (*StreamClosedNotificationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyStreamClosed not implemented")
}
func (UnimplementedFunctionControlServiceServer) mustEmbedUnimplementedFunctionControlServiceServer() {
}

// UnsafeFunctionControlServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FunctionControlServiceServer will
// result in compilation errors.
type UnsafeFunctionControlServiceServer interface {
	mustEmbedUnimplementedFunctionControlServiceServer()
}

func RegisterFunctionControlServiceServer(s grpc.ServiceRegistrar, srv FunctionControlServiceServer) {
	s.RegisterService(&FunctionControlService_ServiceDesc, srv)
}

func _FunctionControlService_GetInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FunctionInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FunctionControlServiceServer).GetInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/controller.FunctionControlService/GetInstance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FunctionControlServiceServer).GetInstance(ctx, req.(*FunctionInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FunctionControlService_NotifyStreamClosed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StreamClosedNotification)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FunctionControlServiceServer).NotifyStreamClosed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/controller.FunctionControlService/NotifyStreamClosed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FunctionControlServiceServer).NotifyStreamClosed(ctx, req.(*StreamClosedNotification))
	}
	return interceptor(ctx, in, info, handler)
}

// FunctionControlService_ServiceDesc is the grpc.ServiceDesc for FunctionControlService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FunctionControlService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "controller.FunctionControlService",
	HandlerType: (*FunctionControlServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetInstance",
			Handler:    _FunctionControlService_GetInstance_Handler,
		},
		{
			MethodName: "NotifyStreamClosed",
			Handler:    _FunctionControlService_NotifyStreamClosed_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/controller/functionstate.proto",
}
