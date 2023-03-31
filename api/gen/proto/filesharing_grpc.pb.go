// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: proto/filesharing.proto

package __proto

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

// FileSharingServiceClient is the client API for FileSharingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FileSharingServiceClient interface {
	SendFile(ctx context.Context, in *SendFileRequest, opts ...grpc.CallOption) (*SendFileResponse, error)
	DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (*DownloadFileResponse, error)
}

type fileSharingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFileSharingServiceClient(cc grpc.ClientConnInterface) FileSharingServiceClient {
	return &fileSharingServiceClient{cc}
}

func (c *fileSharingServiceClient) SendFile(ctx context.Context, in *SendFileRequest, opts ...grpc.CallOption) (*SendFileResponse, error) {
	out := new(SendFileResponse)
	err := c.cc.Invoke(ctx, "/proto.FileSharingService/SendFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileSharingServiceClient) DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (*DownloadFileResponse, error) {
	out := new(DownloadFileResponse)
	err := c.cc.Invoke(ctx, "/proto.FileSharingService/DownloadFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FileSharingServiceServer is the server API for FileSharingService service.
// All implementations must embed UnimplementedFileSharingServiceServer
// for forward compatibility
type FileSharingServiceServer interface {
	SendFile(context.Context, *SendFileRequest) (*SendFileResponse, error)
	DownloadFile(context.Context, *DownloadFileRequest) (*DownloadFileResponse, error)
	mustEmbedUnimplementedFileSharingServiceServer()
}

// UnimplementedFileSharingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFileSharingServiceServer struct {
}

func (UnimplementedFileSharingServiceServer) SendFile(context.Context, *SendFileRequest) (*SendFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendFile not implemented")
}
func (UnimplementedFileSharingServiceServer) DownloadFile(context.Context, *DownloadFileRequest) (*DownloadFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DownloadFile not implemented")
}
func (UnimplementedFileSharingServiceServer) mustEmbedUnimplementedFileSharingServiceServer() {}

// UnsafeFileSharingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileSharingServiceServer will
// result in compilation errors.
type UnsafeFileSharingServiceServer interface {
	mustEmbedUnimplementedFileSharingServiceServer()
}

func RegisterFileSharingServiceServer(s grpc.ServiceRegistrar, srv FileSharingServiceServer) {
	s.RegisterService(&FileSharingService_ServiceDesc, srv)
}

func _FileSharingService_SendFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileSharingServiceServer).SendFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.FileSharingService/SendFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileSharingServiceServer).SendFile(ctx, req.(*SendFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileSharingService_DownloadFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileSharingServiceServer).DownloadFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.FileSharingService/DownloadFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileSharingServiceServer).DownloadFile(ctx, req.(*DownloadFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FileSharingService_ServiceDesc is the grpc.ServiceDesc for FileSharingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FileSharingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.FileSharingService",
	HandlerType: (*FileSharingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendFile",
			Handler:    _FileSharingService_SendFile_Handler,
		},
		{
			MethodName: "DownloadFile",
			Handler:    _FileSharingService_DownloadFile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/filesharing.proto",
}
