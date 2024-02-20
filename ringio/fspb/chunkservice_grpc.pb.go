// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: ringio/fspb/chunkservice.proto

package fspb

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
	HashChunkSystemService_Get_FullMethodName             = "/fspb.HashChunkSystemService/Get"
	HashChunkSystemService_Put_FullMethodName             = "/fspb.HashChunkSystemService/Put"
	HashChunkSystemService_Delete_FullMethodName          = "/fspb.HashChunkSystemService/Delete"
	HashChunkSystemService_PutReplica_FullMethodName      = "/fspb.HashChunkSystemService/PutReplica"
	HashChunkSystemService_GetReplica_FullMethodName      = "/fspb.HashChunkSystemService/GetReplica"
	HashChunkSystemService_DeleteReplica_FullMethodName   = "/fspb.HashChunkSystemService/DeleteReplica"
	HashChunkSystemService_CheckReplica_FullMethodName    = "/fspb.HashChunkSystemService/CheckReplica"
	HashChunkSystemService_SyncReplicaInfo_FullMethodName = "/fspb.HashChunkSystemService/SyncReplicaInfo"
)

// HashChunkSystemServiceClient is the client API for HashChunkSystemService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HashChunkSystemServiceClient interface {
	Get(ctx context.Context, in *Key, opts ...grpc.CallOption) (HashChunkSystemService_GetClient, error)
	Put(ctx context.Context, opts ...grpc.CallOption) (HashChunkSystemService_PutClient, error)
	Delete(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Error, error)
	PutReplica(ctx context.Context, opts ...grpc.CallOption) (HashChunkSystemService_PutReplicaClient, error)
	GetReplica(ctx context.Context, in *Key, opts ...grpc.CallOption) (HashChunkSystemService_GetReplicaClient, error)
	DeleteReplica(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Error, error)
	CheckReplica(ctx context.Context, in *CheckReplicaRequest, opts ...grpc.CallOption) (*Error, error)
	SyncReplicaInfo(ctx context.Context, in *ReplicaChunkInfo, opts ...grpc.CallOption) (*Error, error)
}

type hashChunkSystemServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewHashChunkSystemServiceClient(cc grpc.ClientConnInterface) HashChunkSystemServiceClient {
	return &hashChunkSystemServiceClient{cc}
}

func (c *hashChunkSystemServiceClient) Get(ctx context.Context, in *Key, opts ...grpc.CallOption) (HashChunkSystemService_GetClient, error) {
	stream, err := c.cc.NewStream(ctx, &HashChunkSystemService_ServiceDesc.Streams[0], HashChunkSystemService_Get_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &hashChunkSystemServiceGetClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type HashChunkSystemService_GetClient interface {
	Recv() (*GetResponse, error)
	grpc.ClientStream
}

type hashChunkSystemServiceGetClient struct {
	grpc.ClientStream
}

func (x *hashChunkSystemServiceGetClient) Recv() (*GetResponse, error) {
	m := new(GetResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *hashChunkSystemServiceClient) Put(ctx context.Context, opts ...grpc.CallOption) (HashChunkSystemService_PutClient, error) {
	stream, err := c.cc.NewStream(ctx, &HashChunkSystemService_ServiceDesc.Streams[1], HashChunkSystemService_Put_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &hashChunkSystemServicePutClient{stream}
	return x, nil
}

type HashChunkSystemService_PutClient interface {
	Send(*PutRequest) error
	CloseAndRecv() (*Error, error)
	grpc.ClientStream
}

type hashChunkSystemServicePutClient struct {
	grpc.ClientStream
}

func (x *hashChunkSystemServicePutClient) Send(m *PutRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *hashChunkSystemServicePutClient) CloseAndRecv() (*Error, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Error)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *hashChunkSystemServiceClient) Delete(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Error, error) {
	out := new(Error)
	err := c.cc.Invoke(ctx, HashChunkSystemService_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hashChunkSystemServiceClient) PutReplica(ctx context.Context, opts ...grpc.CallOption) (HashChunkSystemService_PutReplicaClient, error) {
	stream, err := c.cc.NewStream(ctx, &HashChunkSystemService_ServiceDesc.Streams[2], HashChunkSystemService_PutReplica_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &hashChunkSystemServicePutReplicaClient{stream}
	return x, nil
}

type HashChunkSystemService_PutReplicaClient interface {
	Send(*PutReplicaRequest) error
	CloseAndRecv() (*Error, error)
	grpc.ClientStream
}

type hashChunkSystemServicePutReplicaClient struct {
	grpc.ClientStream
}

func (x *hashChunkSystemServicePutReplicaClient) Send(m *PutReplicaRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *hashChunkSystemServicePutReplicaClient) CloseAndRecv() (*Error, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Error)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *hashChunkSystemServiceClient) GetReplica(ctx context.Context, in *Key, opts ...grpc.CallOption) (HashChunkSystemService_GetReplicaClient, error) {
	stream, err := c.cc.NewStream(ctx, &HashChunkSystemService_ServiceDesc.Streams[3], HashChunkSystemService_GetReplica_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &hashChunkSystemServiceGetReplicaClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type HashChunkSystemService_GetReplicaClient interface {
	Recv() (*GetReplicaResponse, error)
	grpc.ClientStream
}

type hashChunkSystemServiceGetReplicaClient struct {
	grpc.ClientStream
}

func (x *hashChunkSystemServiceGetReplicaClient) Recv() (*GetReplicaResponse, error) {
	m := new(GetReplicaResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *hashChunkSystemServiceClient) DeleteReplica(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Error, error) {
	out := new(Error)
	err := c.cc.Invoke(ctx, HashChunkSystemService_DeleteReplica_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hashChunkSystemServiceClient) CheckReplica(ctx context.Context, in *CheckReplicaRequest, opts ...grpc.CallOption) (*Error, error) {
	out := new(Error)
	err := c.cc.Invoke(ctx, HashChunkSystemService_CheckReplica_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hashChunkSystemServiceClient) SyncReplicaInfo(ctx context.Context, in *ReplicaChunkInfo, opts ...grpc.CallOption) (*Error, error) {
	out := new(Error)
	err := c.cc.Invoke(ctx, HashChunkSystemService_SyncReplicaInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HashChunkSystemServiceServer is the server API for HashChunkSystemService service.
// All implementations must embed UnimplementedHashChunkSystemServiceServer
// for forward compatibility
type HashChunkSystemServiceServer interface {
	Get(*Key, HashChunkSystemService_GetServer) error
	Put(HashChunkSystemService_PutServer) error
	Delete(context.Context, *Key) (*Error, error)
	PutReplica(HashChunkSystemService_PutReplicaServer) error
	GetReplica(*Key, HashChunkSystemService_GetReplicaServer) error
	DeleteReplica(context.Context, *Key) (*Error, error)
	CheckReplica(context.Context, *CheckReplicaRequest) (*Error, error)
	SyncReplicaInfo(context.Context, *ReplicaChunkInfo) (*Error, error)
	mustEmbedUnimplementedHashChunkSystemServiceServer()
}

// UnimplementedHashChunkSystemServiceServer must be embedded to have forward compatible implementations.
type UnimplementedHashChunkSystemServiceServer struct {
}

func (UnimplementedHashChunkSystemServiceServer) Get(*Key, HashChunkSystemService_GetServer) error {
	return status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedHashChunkSystemServiceServer) Put(HashChunkSystemService_PutServer) error {
	return status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (UnimplementedHashChunkSystemServiceServer) Delete(context.Context, *Key) (*Error, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedHashChunkSystemServiceServer) PutReplica(HashChunkSystemService_PutReplicaServer) error {
	return status.Errorf(codes.Unimplemented, "method PutReplica not implemented")
}
func (UnimplementedHashChunkSystemServiceServer) GetReplica(*Key, HashChunkSystemService_GetReplicaServer) error {
	return status.Errorf(codes.Unimplemented, "method GetReplica not implemented")
}
func (UnimplementedHashChunkSystemServiceServer) DeleteReplica(context.Context, *Key) (*Error, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteReplica not implemented")
}
func (UnimplementedHashChunkSystemServiceServer) CheckReplica(context.Context, *CheckReplicaRequest) (*Error, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckReplica not implemented")
}
func (UnimplementedHashChunkSystemServiceServer) SyncReplicaInfo(context.Context, *ReplicaChunkInfo) (*Error, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncReplicaInfo not implemented")
}
func (UnimplementedHashChunkSystemServiceServer) mustEmbedUnimplementedHashChunkSystemServiceServer() {
}

// UnsafeHashChunkSystemServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HashChunkSystemServiceServer will
// result in compilation errors.
type UnsafeHashChunkSystemServiceServer interface {
	mustEmbedUnimplementedHashChunkSystemServiceServer()
}

func RegisterHashChunkSystemServiceServer(s grpc.ServiceRegistrar, srv HashChunkSystemServiceServer) {
	s.RegisterService(&HashChunkSystemService_ServiceDesc, srv)
}

func _HashChunkSystemService_Get_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Key)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(HashChunkSystemServiceServer).Get(m, &hashChunkSystemServiceGetServer{stream})
}

type HashChunkSystemService_GetServer interface {
	Send(*GetResponse) error
	grpc.ServerStream
}

type hashChunkSystemServiceGetServer struct {
	grpc.ServerStream
}

func (x *hashChunkSystemServiceGetServer) Send(m *GetResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _HashChunkSystemService_Put_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(HashChunkSystemServiceServer).Put(&hashChunkSystemServicePutServer{stream})
}

type HashChunkSystemService_PutServer interface {
	SendAndClose(*Error) error
	Recv() (*PutRequest, error)
	grpc.ServerStream
}

type hashChunkSystemServicePutServer struct {
	grpc.ServerStream
}

func (x *hashChunkSystemServicePutServer) SendAndClose(m *Error) error {
	return x.ServerStream.SendMsg(m)
}

func (x *hashChunkSystemServicePutServer) Recv() (*PutRequest, error) {
	m := new(PutRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _HashChunkSystemService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Key)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HashChunkSystemServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HashChunkSystemService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HashChunkSystemServiceServer).Delete(ctx, req.(*Key))
	}
	return interceptor(ctx, in, info, handler)
}

func _HashChunkSystemService_PutReplica_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(HashChunkSystemServiceServer).PutReplica(&hashChunkSystemServicePutReplicaServer{stream})
}

type HashChunkSystemService_PutReplicaServer interface {
	SendAndClose(*Error) error
	Recv() (*PutReplicaRequest, error)
	grpc.ServerStream
}

type hashChunkSystemServicePutReplicaServer struct {
	grpc.ServerStream
}

func (x *hashChunkSystemServicePutReplicaServer) SendAndClose(m *Error) error {
	return x.ServerStream.SendMsg(m)
}

func (x *hashChunkSystemServicePutReplicaServer) Recv() (*PutReplicaRequest, error) {
	m := new(PutReplicaRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _HashChunkSystemService_GetReplica_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Key)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(HashChunkSystemServiceServer).GetReplica(m, &hashChunkSystemServiceGetReplicaServer{stream})
}

type HashChunkSystemService_GetReplicaServer interface {
	Send(*GetReplicaResponse) error
	grpc.ServerStream
}

type hashChunkSystemServiceGetReplicaServer struct {
	grpc.ServerStream
}

func (x *hashChunkSystemServiceGetReplicaServer) Send(m *GetReplicaResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _HashChunkSystemService_DeleteReplica_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Key)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HashChunkSystemServiceServer).DeleteReplica(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HashChunkSystemService_DeleteReplica_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HashChunkSystemServiceServer).DeleteReplica(ctx, req.(*Key))
	}
	return interceptor(ctx, in, info, handler)
}

func _HashChunkSystemService_CheckReplica_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckReplicaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HashChunkSystemServiceServer).CheckReplica(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HashChunkSystemService_CheckReplica_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HashChunkSystemServiceServer).CheckReplica(ctx, req.(*CheckReplicaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HashChunkSystemService_SyncReplicaInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReplicaChunkInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HashChunkSystemServiceServer).SyncReplicaInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HashChunkSystemService_SyncReplicaInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HashChunkSystemServiceServer).SyncReplicaInfo(ctx, req.(*ReplicaChunkInfo))
	}
	return interceptor(ctx, in, info, handler)
}

// HashChunkSystemService_ServiceDesc is the grpc.ServiceDesc for HashChunkSystemService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var HashChunkSystemService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fspb.HashChunkSystemService",
	HandlerType: (*HashChunkSystemServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Delete",
			Handler:    _HashChunkSystemService_Delete_Handler,
		},
		{
			MethodName: "DeleteReplica",
			Handler:    _HashChunkSystemService_DeleteReplica_Handler,
		},
		{
			MethodName: "CheckReplica",
			Handler:    _HashChunkSystemService_CheckReplica_Handler,
		},
		{
			MethodName: "SyncReplicaInfo",
			Handler:    _HashChunkSystemService_SyncReplicaInfo_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Get",
			Handler:       _HashChunkSystemService_Get_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Put",
			Handler:       _HashChunkSystemService_Put_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "PutReplica",
			Handler:       _HashChunkSystemService_PutReplica_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "GetReplica",
			Handler:       _HashChunkSystemService_GetReplica_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "ringio/fspb/chunkservice.proto",
}
