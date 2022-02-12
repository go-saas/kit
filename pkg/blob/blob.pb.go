// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.1
// source: blob/blob.proto

package blob

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BlobConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Provider    string          `protobuf:"bytes,1,opt,name=provider,proto3" json:"provider,omitempty"`
	BasePath    string          `protobuf:"bytes,2,opt,name=base_path,json=basePath,proto3" json:"base_path,omitempty"`
	ReadOnly    bool            `protobuf:"varint,3,opt,name=read_only,json=readOnly,proto3" json:"read_only,omitempty"`
	RegexFilter string          `protobuf:"bytes,4,opt,name=regex_filter,json=regexFilter,proto3" json:"regex_filter,omitempty"`
	PublicUrl   string          `protobuf:"bytes,5,opt,name=public_url,json=publicUrl,proto3" json:"public_url,omitempty"`
	InternalUrl string          `protobuf:"bytes,6,opt,name=internal_url,json=internalUrl,proto3" json:"internal_url,omitempty"`
	S3          *BlobProviderS3 `protobuf:"bytes,100,opt,name=s3,proto3" json:"s3,omitempty"`
}

func (x *BlobConfig) Reset() {
	*x = BlobConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blob_blob_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlobConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlobConfig) ProtoMessage() {}

func (x *BlobConfig) ProtoReflect() protoreflect.Message {
	mi := &file_blob_blob_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlobConfig.ProtoReflect.Descriptor instead.
func (*BlobConfig) Descriptor() ([]byte, []int) {
	return file_blob_blob_proto_rawDescGZIP(), []int{0}
}

func (x *BlobConfig) GetProvider() string {
	if x != nil {
		return x.Provider
	}
	return ""
}

func (x *BlobConfig) GetBasePath() string {
	if x != nil {
		return x.BasePath
	}
	return ""
}

func (x *BlobConfig) GetReadOnly() bool {
	if x != nil {
		return x.ReadOnly
	}
	return false
}

func (x *BlobConfig) GetRegexFilter() string {
	if x != nil {
		return x.RegexFilter
	}
	return ""
}

func (x *BlobConfig) GetPublicUrl() string {
	if x != nil {
		return x.PublicUrl
	}
	return ""
}

func (x *BlobConfig) GetInternalUrl() string {
	if x != nil {
		return x.InternalUrl
	}
	return ""
}

func (x *BlobConfig) GetS3() *BlobProviderS3 {
	if x != nil {
		return x.S3
	}
	return nil
}

type BlobProviderS3 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Region string `protobuf:"bytes,1,opt,name=region,proto3" json:"region,omitempty"`
	Key    string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Secret string `protobuf:"bytes,3,opt,name=secret,proto3" json:"secret,omitempty"`
	Bucket string `protobuf:"bytes,4,opt,name=bucket,proto3" json:"bucket,omitempty"`
}

func (x *BlobProviderS3) Reset() {
	*x = BlobProviderS3{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blob_blob_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlobProviderS3) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlobProviderS3) ProtoMessage() {}

func (x *BlobProviderS3) ProtoReflect() protoreflect.Message {
	mi := &file_blob_blob_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlobProviderS3.ProtoReflect.Descriptor instead.
func (*BlobProviderS3) Descriptor() ([]byte, []int) {
	return file_blob_blob_proto_rawDescGZIP(), []int{1}
}

func (x *BlobProviderS3) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *BlobProviderS3) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *BlobProviderS3) GetSecret() string {
	if x != nil {
		return x.Secret
	}
	return ""
}

func (x *BlobProviderS3) GetBucket() string {
	if x != nil {
		return x.Bucket
	}
	return ""
}

type BlobFile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string           `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name     string           `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Mime     string           `protobuf:"bytes,3,opt,name=mime,proto3" json:"mime,omitempty"`
	Size     int64            `protobuf:"varint,4,opt,name=size,proto3" json:"size,omitempty"`
	Url      string           `protobuf:"bytes,5,opt,name=url,proto3" json:"url,omitempty"`
	Metadata *structpb.Struct `protobuf:"bytes,6,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *BlobFile) Reset() {
	*x = BlobFile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blob_blob_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlobFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlobFile) ProtoMessage() {}

func (x *BlobFile) ProtoReflect() protoreflect.Message {
	mi := &file_blob_blob_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlobFile.ProtoReflect.Descriptor instead.
func (*BlobFile) Descriptor() ([]byte, []int) {
	return file_blob_blob_proto_rawDescGZIP(), []int{2}
}

func (x *BlobFile) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *BlobFile) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *BlobFile) GetMime() string {
	if x != nil {
		return x.Mime
	}
	return ""
}

func (x *BlobFile) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *BlobFile) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *BlobFile) GetMetadata() *structpb.Struct {
	if x != nil {
		return x.Metadata
	}
	return nil
}

var File_blob_blob_proto protoreflect.FileDescriptor

var file_blob_blob_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x62, 0x6c, 0x6f, 0x62, 0x2f, 0x62, 0x6c, 0x6f, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x04, 0x62, 0x6c, 0x6f, 0x62, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xed, 0x01, 0x0a, 0x0a, 0x42, 0x6c, 0x6f, 0x62, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
	0x12, 0x1b, 0x0a, 0x09, 0x62, 0x61, 0x73, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x62, 0x61, 0x73, 0x65, 0x50, 0x61, 0x74, 0x68, 0x12, 0x1b, 0x0a,
	0x09, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x08, 0x72, 0x65, 0x61, 0x64, 0x4f, 0x6e, 0x6c, 0x79, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65,
	0x67, 0x65, 0x78, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x72, 0x65, 0x67, 0x65, 0x78, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x1d, 0x0a,
	0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x55, 0x72, 0x6c, 0x12, 0x21, 0x0a, 0x0c,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x55, 0x72, 0x6c, 0x12,
	0x24, 0x0a, 0x02, 0x73, 0x33, 0x18, 0x64, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x62, 0x6c,
	0x6f, 0x62, 0x2e, 0x42, 0x6c, 0x6f, 0x62, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x53,
	0x33, 0x52, 0x02, 0x73, 0x33, 0x22, 0x6a, 0x0a, 0x0e, 0x42, 0x6c, 0x6f, 0x62, 0x50, 0x72, 0x6f,
	0x76, 0x69, 0x64, 0x65, 0x72, 0x53, 0x33, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x75, 0x63,
	0x6b, 0x65, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x62, 0x75, 0x63, 0x6b, 0x65,
	0x74, 0x22, 0x9d, 0x01, 0x0a, 0x08, 0x42, 0x6c, 0x6f, 0x62, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6d, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6d, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72,
	0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x33, 0x0a, 0x08,
	0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x67, 0x6f, 0x78, 0x69, 0x61, 0x6f, 0x79, 0x2f, 0x67, 0x6f, 0x2d, 0x73, 0x61, 0x61, 0x73, 0x2d,
	0x6b, 0x69, 0x74, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x62, 0x6c, 0x6f, 0x62, 0x3b, 0x62, 0x6c, 0x6f,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_blob_blob_proto_rawDescOnce sync.Once
	file_blob_blob_proto_rawDescData = file_blob_blob_proto_rawDesc
)

func file_blob_blob_proto_rawDescGZIP() []byte {
	file_blob_blob_proto_rawDescOnce.Do(func() {
		file_blob_blob_proto_rawDescData = protoimpl.X.CompressGZIP(file_blob_blob_proto_rawDescData)
	})
	return file_blob_blob_proto_rawDescData
}

var file_blob_blob_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_blob_blob_proto_goTypes = []interface{}{
	(*BlobConfig)(nil),      // 0: blob.BlobConfig
	(*BlobProviderS3)(nil),  // 1: blob.BlobProviderS3
	(*BlobFile)(nil),        // 2: blob.BlobFile
	(*structpb.Struct)(nil), // 3: google.protobuf.Struct
}
var file_blob_blob_proto_depIdxs = []int32{
	1, // 0: blob.BlobConfig.s3:type_name -> blob.BlobProviderS3
	3, // 1: blob.BlobFile.metadata:type_name -> google.protobuf.Struct
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_blob_blob_proto_init() }
func file_blob_blob_proto_init() {
	if File_blob_blob_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_blob_blob_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlobConfig); i {
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
		file_blob_blob_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlobProviderS3); i {
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
		file_blob_blob_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlobFile); i {
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
			RawDescriptor: file_blob_blob_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_blob_blob_proto_goTypes,
		DependencyIndexes: file_blob_blob_proto_depIdxs,
		MessageInfos:      file_blob_blob_proto_msgTypes,
	}.Build()
	File_blob_blob_proto = out.File
	file_blob_blob_proto_rawDesc = nil
	file_blob_blob_proto_goTypes = nil
	file_blob_blob_proto_depIdxs = nil
}