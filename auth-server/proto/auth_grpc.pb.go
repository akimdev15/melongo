// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.3
// source: auth.proto

package proto

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

const (
	AuthService_AuthorizeUser_FullMethodName    = "/proto.AuthService/AuthorizeUser"
	AuthService_AuthenticateUser_FullMethodName = "/proto.AuthService/AuthenticateUser"
)

// AuthServiceClient is the client API for AuthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthServiceClient interface {
	// When user login through spotify, will receive callback.
	// This cotains a code which can be used to get an auth token.
	// With the auth token, will get user info and create a new user.
	// Also save the token, refresh token, expire date in a database.
	// Return userID and name that user can use to make a subsequent request.
	AuthorizeUser(ctx context.Context, in *AuthCallbackRequest, opts ...grpc.CallOption) (*AuthCallbackResponse, error)
	// When user is requesting for something like playlist, will always have
	// to go through this authentication method
	AuthenticateUser(ctx context.Context, in *AuthenticateRequest, opts ...grpc.CallOption) (*AuthenticateResponse, error)
}

type authServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthServiceClient(cc grpc.ClientConnInterface) AuthServiceClient {
	return &authServiceClient{cc}
}

func (c *authServiceClient) AuthorizeUser(ctx context.Context, in *AuthCallbackRequest, opts ...grpc.CallOption) (*AuthCallbackResponse, error) {
	out := new(AuthCallbackResponse)
	err := c.cc.Invoke(ctx, AuthService_AuthorizeUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) AuthenticateUser(ctx context.Context, in *AuthenticateRequest, opts ...grpc.CallOption) (*AuthenticateResponse, error) {
	out := new(AuthenticateResponse)
	err := c.cc.Invoke(ctx, AuthService_AuthenticateUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServiceServer is the server API for AuthService service.
// All implementations must embed UnimplementedAuthServiceServer
// for forward compatibility
type AuthServiceServer interface {
	// When user login through spotify, will receive callback.
	// This cotains a code which can be used to get an auth token.
	// With the auth token, will get user info and create a new user.
	// Also save the token, refresh token, expire date in a database.
	// Return userID and name that user can use to make a subsequent request.
	AuthorizeUser(context.Context, *AuthCallbackRequest) (*AuthCallbackResponse, error)
	// When user is requesting for something like playlist, will always have
	// to go through this authentication method
	AuthenticateUser(context.Context, *AuthenticateRequest) (*AuthenticateResponse, error)
	mustEmbedUnimplementedAuthServiceServer()
}

// UnimplementedAuthServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuthServiceServer struct {
}

func (UnimplementedAuthServiceServer) AuthorizeUser(context.Context, *AuthCallbackRequest) (*AuthCallbackResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthorizeUser not implemented")
}
func (UnimplementedAuthServiceServer) AuthenticateUser(context.Context, *AuthenticateRequest) (*AuthenticateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthenticateUser not implemented")
}
func (UnimplementedAuthServiceServer) mustEmbedUnimplementedAuthServiceServer() {}

// UnsafeAuthServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServiceServer will
// result in compilation errors.
type UnsafeAuthServiceServer interface {
	mustEmbedUnimplementedAuthServiceServer()
}

func RegisterAuthServiceServer(s grpc.ServiceRegistrar, srv AuthServiceServer) {
	s.RegisterService(&AuthService_ServiceDesc, srv)
}

func _AuthService_AuthorizeUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthCallbackRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).AuthorizeUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_AuthorizeUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).AuthorizeUser(ctx, req.(*AuthCallbackRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_AuthenticateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthenticateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).AuthenticateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_AuthenticateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).AuthenticateUser(ctx, req.(*AuthenticateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthService_ServiceDesc is the grpc.ServiceDesc for AuthService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AuthorizeUser",
			Handler:    _AuthService_AuthorizeUser_Handler,
		},
		{
			MethodName: "AuthenticateUser",
			Handler:    _AuthService_AuthenticateUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth.proto",
}