// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.4
// source: ringio/fspb/chunkinfo.proto

package fspb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type HashChunkInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChunkName  string                 `protobuf:"bytes,1,opt,name=chunk_name,json=chunkName,proto3" json:"chunk_name,omitempty"`
	ChunkHash  []byte                 `protobuf:"bytes,2,opt,name=chunk_hash,json=chunkHash,proto3" json:"chunk_hash,omitempty"`
	BasePath   string                 `protobuf:"bytes,3,opt,name=base_path,json=basePath,proto3" json:"base_path,omitempty"`
	Size       int64                  `protobuf:"varint,4,opt,name=size,proto3" json:"size,omitempty"`
	ChunkCount int64                  `protobuf:"varint,5,opt,name=chunk_count,json=chunkCount,proto3" json:"chunk_count,omitempty"`
	ModTime    *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=mod_time,json=modTime,proto3" json:"mod_time,omitempty"`
	CreateTime *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
}

func (x *HashChunkInfo) Reset() {
	*x = HashChunkInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ringio_fspb_chunkinfo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HashChunkInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HashChunkInfo) ProtoMessage() {}

func (x *HashChunkInfo) ProtoReflect() protoreflect.Message {
	mi := &file_ringio_fspb_chunkinfo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HashChunkInfo.ProtoReflect.Descriptor instead.
func (*HashChunkInfo) Descriptor() ([]byte, []int) {
	return file_ringio_fspb_chunkinfo_proto_rawDescGZIP(), []int{0}
}

func (x *HashChunkInfo) GetChunkName() string {
	if x != nil {
		return x.ChunkName
	}
	return ""
}

func (x *HashChunkInfo) GetChunkHash() []byte {
	if x != nil {
		return x.ChunkHash
	}
	return nil
}

func (x *HashChunkInfo) GetBasePath() string {
	if x != nil {
		return x.BasePath
	}
	return ""
}

func (x *HashChunkInfo) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *HashChunkInfo) GetChunkCount() int64 {
	if x != nil {
		return x.ChunkCount
	}
	return 0
}

func (x *HashChunkInfo) GetModTime() *timestamppb.Timestamp {
	if x != nil {
		return x.ModTime
	}
	return nil
}

func (x *HashChunkInfo) GetCreateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CreateTime
	}
	return nil
}

var File_ringio_fspb_chunkinfo_proto protoreflect.FileDescriptor

var file_ringio_fspb_chunkinfo_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x72, 0x69, 0x6e, 0x67, 0x69, 0x6f, 0x2f, 0x66, 0x73, 0x70, 0x62, 0x2f, 0x63, 0x68,
	0x75, 0x6e, 0x6b, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x66,
	0x73, 0x70, 0x62, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x93, 0x02, 0x0a, 0x0d, 0x48, 0x61, 0x73, 0x68, 0x43, 0x68, 0x75,
	0x6e, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x68, 0x75, 0x6e,
	0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x5f, 0x68,
	0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x63, 0x68, 0x75, 0x6e, 0x6b,
	0x48, 0x61, 0x73, 0x68, 0x12, 0x1b, 0x0a, 0x09, 0x62, 0x61, 0x73, 0x65, 0x5f, 0x70, 0x61, 0x74,
	0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x62, 0x61, 0x73, 0x65, 0x50, 0x61, 0x74,
	0x68, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x5f, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x63, 0x68, 0x75, 0x6e,
	0x6b, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x35, 0x0a, 0x08, 0x6d, 0x6f, 0x64, 0x5f, 0x74, 0x69,
	0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x6d, 0x6f, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x3b, 0x0a,
	0x0b, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x18, 0x5a, 0x16, 0x63, 0x6c,
	0x6f, 0x75, 0x64, 0x62, 0x6f, 0x72, 0x61, 0x64, 0x2f, 0x72, 0x69, 0x6e, 0x67, 0x69, 0x6f, 0x2f,
	0x66, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ringio_fspb_chunkinfo_proto_rawDescOnce sync.Once
	file_ringio_fspb_chunkinfo_proto_rawDescData = file_ringio_fspb_chunkinfo_proto_rawDesc
)

func file_ringio_fspb_chunkinfo_proto_rawDescGZIP() []byte {
	file_ringio_fspb_chunkinfo_proto_rawDescOnce.Do(func() {
		file_ringio_fspb_chunkinfo_proto_rawDescData = protoimpl.X.CompressGZIP(file_ringio_fspb_chunkinfo_proto_rawDescData)
	})
	return file_ringio_fspb_chunkinfo_proto_rawDescData
}

var file_ringio_fspb_chunkinfo_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_ringio_fspb_chunkinfo_proto_goTypes = []interface{}{
	(*HashChunkInfo)(nil),         // 0: fspb.HashChunkInfo
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_ringio_fspb_chunkinfo_proto_depIdxs = []int32{
	1, // 0: fspb.HashChunkInfo.mod_time:type_name -> google.protobuf.Timestamp
	1, // 1: fspb.HashChunkInfo.create_time:type_name -> google.protobuf.Timestamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_ringio_fspb_chunkinfo_proto_init() }
func file_ringio_fspb_chunkinfo_proto_init() {
	if File_ringio_fspb_chunkinfo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ringio_fspb_chunkinfo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HashChunkInfo); i {
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
			RawDescriptor: file_ringio_fspb_chunkinfo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ringio_fspb_chunkinfo_proto_goTypes,
		DependencyIndexes: file_ringio_fspb_chunkinfo_proto_depIdxs,
		MessageInfos:      file_ringio_fspb_chunkinfo_proto_msgTypes,
	}.Build()
	File_ringio_fspb_chunkinfo_proto = out.File
	file_ringio_fspb_chunkinfo_proto_rawDesc = nil
	file_ringio_fspb_chunkinfo_proto_goTypes = nil
	file_ringio_fspb_chunkinfo_proto_depIdxs = nil
}
