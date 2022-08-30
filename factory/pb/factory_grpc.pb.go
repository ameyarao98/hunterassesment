// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.4
// source: factory.proto

package pb

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

// FactoryClient is the client API for Factory service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FactoryClient interface {
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
	GetResourceData(ctx context.Context, in *GetResourceDataRequest, opts ...grpc.CallOption) (*GetResourceDataResponse, error)
}

type factoryClient struct {
	cc grpc.ClientConnInterface
}

func NewFactoryClient(cc grpc.ClientConnInterface) FactoryClient {
	return &factoryClient{cc}
}

func (c *factoryClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, "/Factory/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *factoryClient) GetResourceData(ctx context.Context, in *GetResourceDataRequest, opts ...grpc.CallOption) (*GetResourceDataResponse, error) {
	out := new(GetResourceDataResponse)
	err := c.cc.Invoke(ctx, "/Factory/GetResourceData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FactoryServer is the server API for Factory service.
// All implementations must embed UnimplementedFactoryServer
// for forward compatibility
type FactoryServer interface {
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	GetResourceData(context.Context, *GetResourceDataRequest) (*GetResourceDataResponse, error)
	mustEmbedUnimplementedFactoryServer()
}

// UnimplementedFactoryServer must be embedded to have forward compatible implementations.
type UnimplementedFactoryServer struct {
}

func (UnimplementedFactoryServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedFactoryServer) GetResourceData(context.Context, *GetResourceDataRequest) (*GetResourceDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetResourceData not implemented")
}
func (UnimplementedFactoryServer) mustEmbedUnimplementedFactoryServer() {}

// UnsafeFactoryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FactoryServer will
// result in compilation errors.
type UnsafeFactoryServer interface {
	mustEmbedUnimplementedFactoryServer()
}

func RegisterFactoryServer(s grpc.ServiceRegistrar, srv FactoryServer) {
	s.RegisterService(&Factory_ServiceDesc, srv)
}

func _Factory_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FactoryServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Factory/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FactoryServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Factory_GetResourceData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetResourceDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FactoryServer).GetResourceData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Factory/GetResourceData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FactoryServer).GetResourceData(ctx, req.(*GetResourceDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Factory_ServiceDesc is the grpc.ServiceDesc for Factory service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Factory_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Factory",
	HandlerType: (*FactoryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _Factory_CreateUser_Handler,
		},
		{
			MethodName: "GetResourceData",
			Handler:    _Factory_GetResourceData_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "factory.proto",
}
