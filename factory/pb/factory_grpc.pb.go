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
	GetFactoryData(ctx context.Context, in *GetFactoryDataRequest, opts ...grpc.CallOption) (*GetFactoryDataResponse, error)
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
}

type factoryClient struct {
	cc grpc.ClientConnInterface
}

func NewFactoryClient(cc grpc.ClientConnInterface) FactoryClient {
	return &factoryClient{cc}
}

func (c *factoryClient) GetFactoryData(ctx context.Context, in *GetFactoryDataRequest, opts ...grpc.CallOption) (*GetFactoryDataResponse, error) {
	out := new(GetFactoryDataResponse)
	err := c.cc.Invoke(ctx, "/Factory/GetFactoryData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *factoryClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, "/Factory/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FactoryServer is the server API for Factory service.
// All implementations must embed UnimplementedFactoryServer
// for forward compatibility
type FactoryServer interface {
	GetFactoryData(context.Context, *GetFactoryDataRequest) (*GetFactoryDataResponse, error)
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	mustEmbedUnimplementedFactoryServer()
}

// UnimplementedFactoryServer must be embedded to have forward compatible implementations.
type UnimplementedFactoryServer struct {
}

func (UnimplementedFactoryServer) GetFactoryData(context.Context, *GetFactoryDataRequest) (*GetFactoryDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFactoryData not implemented")
}
func (UnimplementedFactoryServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
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

func _Factory_GetFactoryData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFactoryDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FactoryServer).GetFactoryData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Factory/GetFactoryData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FactoryServer).GetFactoryData(ctx, req.(*GetFactoryDataRequest))
	}
	return interceptor(ctx, in, info, handler)
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

// Factory_ServiceDesc is the grpc.ServiceDesc for Factory service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Factory_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Factory",
	HandlerType: (*FactoryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFactoryData",
			Handler:    _Factory_GetFactoryData_Handler,
		},
		{
			MethodName: "CreateUser",
			Handler:    _Factory_CreateUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "factory.proto",
}
