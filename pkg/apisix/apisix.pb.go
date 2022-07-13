// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.15.8
// source: apisix/apisix.proto

package apisix

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

type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host   string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	Port   uint64 `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	Weight int64  `protobuf:"varint,3,opt,name=weight,proto3" json:"weight,omitempty"`
}

func (x *Node) Reset() {
	*x = Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apisix_apisix_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_apisix_apisix_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_apisix_apisix_proto_rawDescGZIP(), []int{0}
}

func (x *Node) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *Node) GetPort() uint64 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *Node) GetWeight() int64 {
	if x != nil {
		return x.Weight
	}
	return 0
}

type Upstream struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nodes  []*Node `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
	Type   string  `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	Scheme string  `protobuf:"bytes,3,opt,name=scheme,proto3" json:"scheme,omitempty"`
}

func (x *Upstream) Reset() {
	*x = Upstream{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apisix_apisix_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Upstream) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Upstream) ProtoMessage() {}

func (x *Upstream) ProtoReflect() protoreflect.Message {
	mi := &file_apisix_apisix_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Upstream.ProtoReflect.Descriptor instead.
func (*Upstream) Descriptor() ([]byte, []int) {
	return file_apisix_apisix_proto_rawDescGZIP(), []int{1}
}

func (x *Upstream) GetNodes() []*Node {
	if x != nil {
		return x.Nodes
	}
	return nil
}

func (x *Upstream) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Upstream) GetScheme() string {
	if x != nil {
		return x.Scheme
	}
	return ""
}

var File_apisix_apisix_proto protoreflect.FileDescriptor

var file_apisix_apisix_proto_rawDesc = []byte{
	0x0a, 0x13, 0x61, 0x70, 0x69, 0x73, 0x69, 0x78, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x69, 0x78, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x69, 0x73, 0x69, 0x78, 0x22, 0x46, 0x0a,
	0x04, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x77,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x22, 0x5a, 0x0a, 0x08, 0x55, 0x70, 0x73, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x12, 0x22, 0x0a, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0c, 0x2e, 0x61, 0x70, 0x69, 0x73, 0x69, 0x78, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x05,
	0x6e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x63, 0x68,
	0x65, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d,
	0x65, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x67, 0x6f, 0x2d, 0x73, 0x61, 0x61, 0x73, 0x2f, 0x6b, 0x69, 0x74, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x61, 0x70, 0x69, 0x73, 0x69, 0x78, 0x3b, 0x61, 0x70, 0x69, 0x73, 0x69, 0x78, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_apisix_apisix_proto_rawDescOnce sync.Once
	file_apisix_apisix_proto_rawDescData = file_apisix_apisix_proto_rawDesc
)

func file_apisix_apisix_proto_rawDescGZIP() []byte {
	file_apisix_apisix_proto_rawDescOnce.Do(func() {
		file_apisix_apisix_proto_rawDescData = protoimpl.X.CompressGZIP(file_apisix_apisix_proto_rawDescData)
	})
	return file_apisix_apisix_proto_rawDescData
}

var file_apisix_apisix_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_apisix_apisix_proto_goTypes = []interface{}{
	(*Node)(nil),     // 0: apisix.Node
	(*Upstream)(nil), // 1: apisix.Upstream
}
var file_apisix_apisix_proto_depIdxs = []int32{
	0, // 0: apisix.Upstream.nodes:type_name -> apisix.Node
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_apisix_apisix_proto_init() }
func file_apisix_apisix_proto_init() {
	if File_apisix_apisix_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_apisix_apisix_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Node); i {
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
		file_apisix_apisix_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Upstream); i {
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
			RawDescriptor: file_apisix_apisix_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_apisix_apisix_proto_goTypes,
		DependencyIndexes: file_apisix_apisix_proto_depIdxs,
		MessageInfos:      file_apisix_apisix_proto_msgTypes,
	}.Build()
	File_apisix_apisix_proto = out.File
	file_apisix_apisix_proto_rawDesc = nil
	file_apisix_apisix_proto_goTypes = nil
	file_apisix_apisix_proto_depIdxs = nil
}
