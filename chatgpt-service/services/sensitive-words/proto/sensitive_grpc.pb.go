// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.0
// source: proto/sensitive.proto

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

// SensitiveWordsClient is the client API for SensitiveWords service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SensitiveWordsClient interface {
	Validate(ctx context.Context, in *ValidateReq, opts ...grpc.CallOption) (*ValidateRes, error)
}

type sensitiveWordsClient struct {
	cc grpc.ClientConnInterface
}

func NewSensitiveWordsClient(cc grpc.ClientConnInterface) SensitiveWordsClient {
	return &sensitiveWordsClient{cc}
}

func (c *sensitiveWordsClient) Validate(ctx context.Context, in *ValidateReq, opts ...grpc.CallOption) (*ValidateRes, error) {
	out := new(ValidateRes)
	err := c.cc.Invoke(ctx, "/leoh_package.SensitiveWords/Validate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SensitiveWordsServer is the server API for SensitiveWords service.
// All implementations must embed UnimplementedSensitiveWordsServer
// for forward compatibility
type SensitiveWordsServer interface {
	Validate(context.Context, *ValidateReq) (*ValidateRes, error)
	mustEmbedUnimplementedSensitiveWordsServer()
}

// UnimplementedSensitiveWordsServer must be embedded to have forward compatible implementations.
type UnimplementedSensitiveWordsServer struct {
}

func (UnimplementedSensitiveWordsServer) Validate(context.Context, *ValidateReq) (*ValidateRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validate not implemented")
}
func (UnimplementedSensitiveWordsServer) mustEmbedUnimplementedSensitiveWordsServer() {}

// UnsafeSensitiveWordsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SensitiveWordsServer will
// result in compilation errors.
type UnsafeSensitiveWordsServer interface {
	mustEmbedUnimplementedSensitiveWordsServer()
}

func RegisterSensitiveWordsServer(s grpc.ServiceRegistrar, srv SensitiveWordsServer) {
	s.RegisterService(&SensitiveWords_ServiceDesc, srv)
}

func _SensitiveWords_Validate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SensitiveWordsServer).Validate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/leoh_package.SensitiveWords/Validate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SensitiveWordsServer).Validate(ctx, req.(*ValidateReq))
	}
	return interceptor(ctx, in, info, handler)
}

// SensitiveWords_ServiceDesc is the grpc.ServiceDesc for SensitiveWords service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SensitiveWords_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "leoh_package.SensitiveWords",
	HandlerType: (*SensitiveWordsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Validate",
			Handler:    _SensitiveWords_Validate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/sensitive.proto",
}
