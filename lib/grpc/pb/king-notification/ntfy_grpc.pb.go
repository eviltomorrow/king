// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: ntfy.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Ntfy_Send_FullMethodName = "/notification.Ntfy/Send"
)

// NtfyClient is the client API for Ntfy service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NtfyClient interface {
	Send(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*wrapperspb.StringValue, error)
}

type ntfyClient struct {
	cc grpc.ClientConnInterface
}

func NewNtfyClient(cc grpc.ClientConnInterface) NtfyClient {
	return &ntfyClient{cc}
}

func (c *ntfyClient) Send(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*wrapperspb.StringValue, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(wrapperspb.StringValue)
	err := c.cc.Invoke(ctx, Ntfy_Send_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NtfyServer is the server API for Ntfy service.
// All implementations must embed UnimplementedNtfyServer
// for forward compatibility.
type NtfyServer interface {
	Send(context.Context, *Msg) (*wrapperspb.StringValue, error)
	mustEmbedUnimplementedNtfyServer()
}

// UnimplementedNtfyServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedNtfyServer struct{}

func (UnimplementedNtfyServer) Send(context.Context, *Msg) (*wrapperspb.StringValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Send not implemented")
}
func (UnimplementedNtfyServer) mustEmbedUnimplementedNtfyServer() {}
func (UnimplementedNtfyServer) testEmbeddedByValue()              {}

// UnsafeNtfyServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NtfyServer will
// result in compilation errors.
type UnsafeNtfyServer interface {
	mustEmbedUnimplementedNtfyServer()
}

func RegisterNtfyServer(s grpc.ServiceRegistrar, srv NtfyServer) {
	// If the following call pancis, it indicates UnimplementedNtfyServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Ntfy_ServiceDesc, srv)
}

func _Ntfy_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Msg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NtfyServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Ntfy_Send_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NtfyServer).Send(ctx, req.(*Msg))
	}
	return interceptor(ctx, in, info, handler)
}

// Ntfy_ServiceDesc is the grpc.ServiceDesc for Ntfy service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Ntfy_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "notification.Ntfy",
	HandlerType: (*NtfyServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Send",
			Handler:    _Ntfy_Send_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ntfy.proto",
}