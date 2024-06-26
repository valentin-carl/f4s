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

// FunctionStateServiceClient is the client API for FunctionStateService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FunctionStateServiceClient interface {
	GetInstance(ctx context.Context, in *FunctionInstanceRequest, opts ...grpc.CallOption) (*FunctionInstanceResponse, error)
	NotifyStreamClosed(ctx context.Context, in *StreamClosedNotification, opts ...grpc.CallOption) (*StreamClosedNotificationResponse, error)
}

type functionStateServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFunctionStateServiceClient(cc grpc.ClientConnInterface) FunctionStateServiceClient {
	return &functionStateServiceClient{cc}
}

func (c *functionStateServiceClient) GetInstance(ctx context.Context, in *FunctionInstanceRequest, opts ...grpc.CallOption) (*FunctionInstanceResponse, error) {
	out := new(FunctionInstanceResponse)
	err := c.cc.Invoke(ctx, "/controller.FunctionStateService/GetInstance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *functionStateServiceClient) NotifyStreamClosed(ctx context.Context, in *StreamClosedNotification, opts ...grpc.CallOption) (*StreamClosedNotificationResponse, error) {
	out := new(StreamClosedNotificationResponse)
	err := c.cc.Invoke(ctx, "/controller.FunctionStateService/NotifyStreamClosed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FunctionStateServiceServer is the server API for FunctionStateService service.
// All implementations must embed UnimplementedFunctionStateServiceServer
// for forward compatibility
type FunctionStateServiceServer interface {
	GetInstance(context.Context, *FunctionInstanceRequest) (*FunctionInstanceResponse, error)
	NotifyStreamClosed(context.Context, *StreamClosedNotification) (*StreamClosedNotificationResponse, error)
	mustEmbedUnimplementedFunctionStateServiceServer()
}

// UnimplementedFunctionStateServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFunctionStateServiceServer struct {
}

func (UnimplementedFunctionStateServiceServer) GetInstance(context.Context, *FunctionInstanceRequest) (*FunctionInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInstance not implemented")
}
func (UnimplementedFunctionStateServiceServer) NotifyStreamClosed(context.Context, *StreamClosedNotification) (*StreamClosedNotificationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyStreamClosed not implemented")
}
func (UnimplementedFunctionStateServiceServer) mustEmbedUnimplementedFunctionStateServiceServer() {}

// UnsafeFunctionStateServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FunctionStateServiceServer will
// result in compilation errors.
type UnsafeFunctionStateServiceServer interface {
	mustEmbedUnimplementedFunctionStateServiceServer()
}

func RegisterFunctionStateServiceServer(s grpc.ServiceRegistrar, srv FunctionStateServiceServer) {
	s.RegisterService(&FunctionStateService_ServiceDesc, srv)
}

func _FunctionStateService_GetInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FunctionInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FunctionStateServiceServer).GetInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/controller.FunctionStateService/GetInstance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FunctionStateServiceServer).GetInstance(ctx, req.(*FunctionInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FunctionStateService_NotifyStreamClosed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StreamClosedNotification)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FunctionStateServiceServer).NotifyStreamClosed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/controller.FunctionStateService/NotifyStreamClosed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FunctionStateServiceServer).NotifyStreamClosed(ctx, req.(*StreamClosedNotification))
	}
	return interceptor(ctx, in, info, handler)
}

// FunctionStateService_ServiceDesc is the grpc.ServiceDesc for FunctionStateService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FunctionStateService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "controller.FunctionStateService",
	HandlerType: (*FunctionStateServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetInstance",
			Handler:    _FunctionStateService_GetInstance_Handler,
		},
		{
			MethodName: "NotifyStreamClosed",
			Handler:    _FunctionStateService_NotifyStreamClosed_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/controller/functionstate.proto",
}
