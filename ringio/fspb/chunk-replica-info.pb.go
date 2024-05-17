// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.4
// source: ringio/fspb/chunk-replica-info.proto

package fspb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ReplicaChunkInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key          []byte         `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	ChunkInfo    *HashChunkInfo `protobuf:"bytes,2,opt,name=chunk_info,json=chunkInfo,proto3" json:"chunk_info,omitempty"`
	ReplicaCount int64          `protobuf:"varint,3,opt,name=replica_count,json=replicaCount,proto3" json:"replica_count,omitempty"`
	Checksum     []byte         `protobuf:"bytes,4,opt,name=checksum,proto3" json:"checksum,omitempty"`
	NodeIds      []string       `protobuf:"bytes,5,rep,name=node_ids,json=nodeIds,proto3" json:"node_ids,omitempty"`
}

func (x *ReplicaChunkInfo) Reset() {
	*x = ReplicaChunkInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ringio_fspb_chunk_replica_info_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReplicaChunkInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReplicaChunkInfo) ProtoMessage() {}

func (x *ReplicaChunkInfo) ProtoReflect() protoreflect.Message {
	mi := &file_ringio_fspb_chunk_replica_info_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReplicaChunkInfo.ProtoReflect.Descriptor instead.
func (*ReplicaChunkInfo) Descriptor() ([]byte, []int) {
	return file_ringio_fspb_chunk_replica_info_proto_rawDescGZIP(), []int{0}
}

func (x *ReplicaChunkInfo) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *ReplicaChunkInfo) GetChunkInfo() *HashChunkInfo {
	if x != nil {
		return x.ChunkInfo
	}
	return nil
}

func (x *ReplicaChunkInfo) GetReplicaCount() int64 {
	if x != nil {
		return x.ReplicaCount
	}
	return 0
}

func (x *ReplicaChunkInfo) GetChecksum() []byte {
	if x != nil {
		return x.Checksum
	}
	return nil
}

func (x *ReplicaChunkInfo) GetNodeIds() []string {
	if x != nil {
		return x.NodeIds
	}
	return nil
}

// Stream
type PutReplicaRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte            `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Info *ReplicaChunkInfo `protobuf:"bytes,2,opt,name=info,proto3" json:"info,omitempty"`
}

func (x *PutReplicaRequest) Reset() {
	*x = PutReplicaRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ringio_fspb_chunk_replica_info_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutReplicaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutReplicaRequest) ProtoMessage() {}

func (x *PutReplicaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ringio_fspb_chunk_replica_info_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutReplicaRequest.ProtoReflect.Descriptor instead.
func (*PutReplicaRequest) Descriptor() ([]byte, []int) {
	return file_ringio_fspb_chunk_replica_info_proto_rawDescGZIP(), []int{1}
}

func (x *PutReplicaRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *PutReplicaRequest) GetInfo() *ReplicaChunkInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

// Stream
type GetReplicaResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data  []byte            `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Info  *ReplicaChunkInfo `protobuf:"bytes,2,opt,name=info,proto3" json:"info,omitempty"`
	Error *Error            `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *GetReplicaResponse) Reset() {
	*x = GetReplicaResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ringio_fspb_chunk_replica_info_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetReplicaResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetReplicaResponse) ProtoMessage() {}

func (x *GetReplicaResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ringio_fspb_chunk_replica_info_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetReplicaResponse.ProtoReflect.Descriptor instead.
func (*GetReplicaResponse) Descriptor() ([]byte, []int) {
	return file_ringio_fspb_chunk_replica_info_proto_rawDescGZIP(), []int{2}
}

func (x *GetReplicaResponse) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *GetReplicaResponse) GetInfo() *ReplicaChunkInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

func (x *GetReplicaResponse) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

type CheckReplicaRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Info *ReplicaChunkInfo `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty"`
}

func (x *CheckReplicaRequest) Reset() {
	*x = CheckReplicaRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ringio_fspb_chunk_replica_info_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckReplicaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckReplicaRequest) ProtoMessage() {}

func (x *CheckReplicaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ringio_fspb_chunk_replica_info_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckReplicaRequest.ProtoReflect.Descriptor instead.
func (*CheckReplicaRequest) Descriptor() ([]byte, []int) {
	return file_ringio_fspb_chunk_replica_info_proto_rawDescGZIP(), []int{3}
}

func (x *CheckReplicaRequest) GetInfo() *ReplicaChunkInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

type UpdateReplicaInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Info *ReplicaChunkInfo `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty"`
}

func (x *UpdateReplicaInfoRequest) Reset() {
	*x = UpdateReplicaInfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ringio_fspb_chunk_replica_info_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateReplicaInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateReplicaInfoRequest) ProtoMessage() {}

func (x *UpdateReplicaInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ringio_fspb_chunk_replica_info_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateReplicaInfoRequest.ProtoReflect.Descriptor instead.
func (*UpdateReplicaInfoRequest) Descriptor() ([]byte, []int) {
	return file_ringio_fspb_chunk_replica_info_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateReplicaInfoRequest) GetInfo() *ReplicaChunkInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

var File_ringio_fspb_chunk_replica_info_proto protoreflect.FileDescriptor

var file_ringio_fspb_chunk_replica_info_proto_rawDesc = []byte{
	0x0a, 0x24, 0x72, 0x69, 0x6e, 0x67, 0x69, 0x6f, 0x2f, 0x66, 0x73, 0x70, 0x62, 0x2f, 0x63, 0x68,
	0x75, 0x6e, 0x6b, 0x2d, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x2d, 0x69, 0x6e, 0x66, 0x6f,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x66, 0x73, 0x70, 0x62, 0x1a, 0x1b, 0x72, 0x69,
	0x6e, 0x67, 0x69, 0x6f, 0x2f, 0x66, 0x73, 0x70, 0x62, 0x2f, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x69,
	0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x72, 0x69, 0x6e, 0x67, 0x69,
	0x6f, 0x2f, 0x66, 0x73, 0x70, 0x62, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xb4, 0x01, 0x0a, 0x10, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x43, 0x68,
	0x75, 0x6e, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x32, 0x0a, 0x0a, 0x63, 0x68, 0x75,
	0x6e, 0x6b, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e,
	0x66, 0x73, 0x70, 0x62, 0x2e, 0x48, 0x61, 0x73, 0x68, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x49, 0x6e,
	0x66, 0x6f, 0x52, 0x09, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x23, 0x0a,
	0x0d, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d, 0x12, 0x19,
	0x0a, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x73, 0x22, 0x53, 0x0a, 0x11, 0x50, 0x75, 0x74,
	0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x2a, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x16, 0x2e, 0x66, 0x73, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x43,
	0x68, 0x75, 0x6e, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x22, 0x77,
	0x0a, 0x12, 0x47, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x2a, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x66, 0x73, 0x70, 0x62, 0x2e, 0x52, 0x65,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x04,
	0x69, 0x6e, 0x66, 0x6f, 0x12, 0x21, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x66, 0x73, 0x70, 0x62, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x41, 0x0a, 0x13, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a,
	0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x66,
	0x73, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x43, 0x68, 0x75, 0x6e, 0x6b,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x22, 0x46, 0x0a, 0x18, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x66, 0x73, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x70, 0x6c,
	0x69, 0x63, 0x61, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x04, 0x69, 0x6e,
	0x66, 0x6f, 0x42, 0x18, 0x5a, 0x16, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x62, 0x6f, 0x72, 0x61, 0x64,
	0x2f, 0x72, 0x69, 0x6e, 0x67, 0x69, 0x6f, 0x2f, 0x66, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ringio_fspb_chunk_replica_info_proto_rawDescOnce sync.Once
	file_ringio_fspb_chunk_replica_info_proto_rawDescData = file_ringio_fspb_chunk_replica_info_proto_rawDesc
)

func file_ringio_fspb_chunk_replica_info_proto_rawDescGZIP() []byte {
	file_ringio_fspb_chunk_replica_info_proto_rawDescOnce.Do(func() {
		file_ringio_fspb_chunk_replica_info_proto_rawDescData = protoimpl.X.CompressGZIP(file_ringio_fspb_chunk_replica_info_proto_rawDescData)
	})
	return file_ringio_fspb_chunk_replica_info_proto_rawDescData
}

var file_ringio_fspb_chunk_replica_info_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_ringio_fspb_chunk_replica_info_proto_goTypes = []interface{}{
	(*ReplicaChunkInfo)(nil),         // 0: fspb.ReplicaChunkInfo
	(*PutReplicaRequest)(nil),        // 1: fspb.PutReplicaRequest
	(*GetReplicaResponse)(nil),       // 2: fspb.GetReplicaResponse
	(*CheckReplicaRequest)(nil),      // 3: fspb.CheckReplicaRequest
	(*UpdateReplicaInfoRequest)(nil), // 4: fspb.UpdateReplicaInfoRequest
	(*HashChunkInfo)(nil),            // 5: fspb.HashChunkInfo
	(*Error)(nil),                    // 6: fspb.Error
}
var file_ringio_fspb_chunk_replica_info_proto_depIdxs = []int32{
	5, // 0: fspb.ReplicaChunkInfo.chunk_info:type_name -> fspb.HashChunkInfo
	0, // 1: fspb.PutReplicaRequest.info:type_name -> fspb.ReplicaChunkInfo
	0, // 2: fspb.GetReplicaResponse.info:type_name -> fspb.ReplicaChunkInfo
	6, // 3: fspb.GetReplicaResponse.error:type_name -> fspb.Error
	0, // 4: fspb.CheckReplicaRequest.info:type_name -> fspb.ReplicaChunkInfo
	0, // 5: fspb.UpdateReplicaInfoRequest.info:type_name -> fspb.ReplicaChunkInfo
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_ringio_fspb_chunk_replica_info_proto_init() }
func file_ringio_fspb_chunk_replica_info_proto_init() {
	if File_ringio_fspb_chunk_replica_info_proto != nil {
		return
	}
	file_ringio_fspb_chunkinfo_proto_init()
	file_ringio_fspb_error_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_ringio_fspb_chunk_replica_info_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReplicaChunkInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ringio_fspb_chunk_replica_info_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutReplicaRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ringio_fspb_chunk_replica_info_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetReplicaResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ringio_fspb_chunk_replica_info_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckReplicaRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ringio_fspb_chunk_replica_info_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateReplicaInfoRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ringio_fspb_chunk_replica_info_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ringio_fspb_chunk_replica_info_proto_goTypes,
		DependencyIndexes: file_ringio_fspb_chunk_replica_info_proto_depIdxs,
		MessageInfos:      file_ringio_fspb_chunk_replica_info_proto_msgTypes,
	}.Build()
	File_ringio_fspb_chunk_replica_info_proto = out.File
	file_ringio_fspb_chunk_replica_info_proto_rawDesc = nil
	file_ringio_fspb_chunk_replica_info_proto_goTypes = nil
	file_ringio_fspb_chunk_replica_info_proto_depIdxs = nil
}