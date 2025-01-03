// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: storage.proto

package pb

import (
	context "context"
	entity "github.com/eviltomorrow/king/lib/grpc/pb/entity"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Storage_ShowMetadata_FullMethodName   = "/storage.Storage/ShowMetadata"
	Storage_PushMetadata_FullMethodName   = "/storage.Storage/PushMetadata"
	Storage_GetStockOne_FullMethodName    = "/storage.Storage/GetStockOne"
	Storage_GetStockAll_FullMethodName    = "/storage.Storage/GetStockAll"
	Storage_GetQuoteLatest_FullMethodName = "/storage.Storage/GetQuoteLatest"
)

// StorageClient is the client API for Storage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StorageClient interface {
	ShowMetadata(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*ShowResponse, error)
	PushMetadata(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[entity.Metadata, PushResponse], error)
	GetStockOne(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*Stock, error)
	GetStockAll(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Stock], error)
	GetQuoteLatest(ctx context.Context, in *GetQuoteLatestRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Quote], error)
}

type storageClient struct {
	cc grpc.ClientConnInterface
}

func NewStorageClient(cc grpc.ClientConnInterface) StorageClient {
	return &storageClient{cc}
}

func (c *storageClient) ShowMetadata(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*ShowResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShowResponse)
	err := c.cc.Invoke(ctx, Storage_ShowMetadata_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageClient) PushMetadata(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[entity.Metadata, PushResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Storage_ServiceDesc.Streams[0], Storage_PushMetadata_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[entity.Metadata, PushResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Storage_PushMetadataClient = grpc.ClientStreamingClient[entity.Metadata, PushResponse]

func (c *storageClient) GetStockOne(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*Stock, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Stock)
	err := c.cc.Invoke(ctx, Storage_GetStockOne_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageClient) GetStockAll(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Stock], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Storage_ServiceDesc.Streams[1], Storage_GetStockAll_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[emptypb.Empty, Stock]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Storage_GetStockAllClient = grpc.ServerStreamingClient[Stock]

func (c *storageClient) GetQuoteLatest(ctx context.Context, in *GetQuoteLatestRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Quote], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Storage_ServiceDesc.Streams[2], Storage_GetQuoteLatest_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[GetQuoteLatestRequest, Quote]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Storage_GetQuoteLatestClient = grpc.ServerStreamingClient[Quote]

// StorageServer is the server API for Storage service.
// All implementations must embed UnimplementedStorageServer
// for forward compatibility.
type StorageServer interface {
	ShowMetadata(context.Context, *wrapperspb.StringValue) (*ShowResponse, error)
	PushMetadata(grpc.ClientStreamingServer[entity.Metadata, PushResponse]) error
	GetStockOne(context.Context, *wrapperspb.StringValue) (*Stock, error)
	GetStockAll(*emptypb.Empty, grpc.ServerStreamingServer[Stock]) error
	GetQuoteLatest(*GetQuoteLatestRequest, grpc.ServerStreamingServer[Quote]) error
	mustEmbedUnimplementedStorageServer()
}

// UnimplementedStorageServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedStorageServer struct{}

func (UnimplementedStorageServer) ShowMetadata(context.Context, *wrapperspb.StringValue) (*ShowResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShowMetadata not implemented")
}
func (UnimplementedStorageServer) PushMetadata(grpc.ClientStreamingServer[entity.Metadata, PushResponse]) error {
	return status.Errorf(codes.Unimplemented, "method PushMetadata not implemented")
}
func (UnimplementedStorageServer) GetStockOne(context.Context, *wrapperspb.StringValue) (*Stock, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStockOne not implemented")
}
func (UnimplementedStorageServer) GetStockAll(*emptypb.Empty, grpc.ServerStreamingServer[Stock]) error {
	return status.Errorf(codes.Unimplemented, "method GetStockAll not implemented")
}
func (UnimplementedStorageServer) GetQuoteLatest(*GetQuoteLatestRequest, grpc.ServerStreamingServer[Quote]) error {
	return status.Errorf(codes.Unimplemented, "method GetQuoteLatest not implemented")
}
func (UnimplementedStorageServer) mustEmbedUnimplementedStorageServer() {}
func (UnimplementedStorageServer) testEmbeddedByValue()                 {}

// UnsafeStorageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StorageServer will
// result in compilation errors.
type UnsafeStorageServer interface {
	mustEmbedUnimplementedStorageServer()
}

func RegisterStorageServer(s grpc.ServiceRegistrar, srv StorageServer) {
	// If the following call pancis, it indicates UnimplementedStorageServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Storage_ServiceDesc, srv)
}

func _Storage_ShowMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).ShowMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Storage_ShowMetadata_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).ShowMetadata(ctx, req.(*wrapperspb.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _Storage_PushMetadata_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StorageServer).PushMetadata(&grpc.GenericServerStream[entity.Metadata, PushResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Storage_PushMetadataServer = grpc.ClientStreamingServer[entity.Metadata, PushResponse]

func _Storage_GetStockOne_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).GetStockOne(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Storage_GetStockOne_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).GetStockOne(ctx, req.(*wrapperspb.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _Storage_GetStockAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StorageServer).GetStockAll(m, &grpc.GenericServerStream[emptypb.Empty, Stock]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Storage_GetStockAllServer = grpc.ServerStreamingServer[Stock]

func _Storage_GetQuoteLatest_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetQuoteLatestRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StorageServer).GetQuoteLatest(m, &grpc.GenericServerStream[GetQuoteLatestRequest, Quote]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Storage_GetQuoteLatestServer = grpc.ServerStreamingServer[Quote]

// Storage_ServiceDesc is the grpc.ServiceDesc for Storage service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Storage_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "storage.Storage",
	HandlerType: (*StorageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ShowMetadata",
			Handler:    _Storage_ShowMetadata_Handler,
		},
		{
			MethodName: "GetStockOne",
			Handler:    _Storage_GetStockOne_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PushMetadata",
			Handler:       _Storage_PushMetadata_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "GetStockAll",
			Handler:       _Storage_GetStockAll_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetQuoteLatest",
			Handler:       _Storage_GetQuoteLatest_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "storage.proto",
}
