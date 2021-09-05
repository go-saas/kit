// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.6.1
// source: internal_/conf/conf.proto

package conf

import (
	conf "github.com/goxiaoy/go-saas-kit/pkg/conf"
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

type Bootstrap struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data     *Data          `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Security *conf.Security `protobuf:"bytes,3,opt,name=security,proto3" json:"security,omitempty"`
	Services *conf.Services `protobuf:"bytes,4,opt,name=services,proto3" json:"services,omitempty"`
	User     *UserConf      `protobuf:"bytes,5,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *Bootstrap) Reset() {
	*x = Bootstrap{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal__conf_conf_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Bootstrap) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Bootstrap) ProtoMessage() {}

func (x *Bootstrap) ProtoReflect() protoreflect.Message {
	mi := &file_internal__conf_conf_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Bootstrap.ProtoReflect.Descriptor instead.
func (*Bootstrap) Descriptor() ([]byte, []int) {
	return file_internal__conf_conf_proto_rawDescGZIP(), []int{0}
}

func (x *Bootstrap) GetData() *Data {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Bootstrap) GetSecurity() *conf.Security {
	if x != nil {
		return x.Security
	}
	return nil
}

func (x *Bootstrap) GetServices() *conf.Services {
	if x != nil {
		return x.Services
	}
	return nil
}

func (x *Bootstrap) GetUser() *UserConf {
	if x != nil {
		return x.User
	}
	return nil
}

type Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Databases *conf.Databases `protobuf:"bytes,1,opt,name=databases,proto3" json:"databases,omitempty"`
	Redis     *conf.Redis     `protobuf:"bytes,2,opt,name=redis,proto3" json:"redis,omitempty"`
}

func (x *Data) Reset() {
	*x = Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal__conf_conf_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data) ProtoMessage() {}

func (x *Data) ProtoReflect() protoreflect.Message {
	mi := &file_internal__conf_conf_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data.ProtoReflect.Descriptor instead.
func (*Data) Descriptor() ([]byte, []int) {
	return file_internal__conf_conf_proto_rawDescGZIP(), []int{1}
}

func (x *Data) GetDatabases() *conf.Databases {
	if x != nil {
		return x.Databases
	}
	return nil
}

func (x *Data) GetRedis() *conf.Redis {
	if x != nil {
		return x.Redis
	}
	return nil
}

type Admin struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *Admin) Reset() {
	*x = Admin{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal__conf_conf_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Admin) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Admin) ProtoMessage() {}

func (x *Admin) ProtoReflect() protoreflect.Message {
	mi := &file_internal__conf_conf_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Admin.ProtoReflect.Descriptor instead.
func (*Admin) Descriptor() ([]byte, []int) {
	return file_internal__conf_conf_proto_rawDescGZIP(), []int{2}
}

func (x *Admin) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Admin) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type UserConf struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Minimum password score. [0-5]
	PasswordScoreMin int32  `protobuf:"varint,1,opt,name=password_score_min,json=passwordScoreMin,proto3" json:"password_score_min,omitempty"`
	Admin            *Admin `protobuf:"bytes,2,opt,name=admin,proto3" json:"admin,omitempty"`
}

func (x *UserConf) Reset() {
	*x = UserConf{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal__conf_conf_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserConf) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserConf) ProtoMessage() {}

func (x *UserConf) ProtoReflect() protoreflect.Message {
	mi := &file_internal__conf_conf_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserConf.ProtoReflect.Descriptor instead.
func (*UserConf) Descriptor() ([]byte, []int) {
	return file_internal__conf_conf_proto_rawDescGZIP(), []int{3}
}

func (x *UserConf) GetPasswordScoreMin() int32 {
	if x != nil {
		return x.PasswordScoreMin
	}
	return 0
}

func (x *UserConf) GetAdmin() *Admin {
	if x != nil {
		return x.Admin
	}
	return nil
}

var File_internal__conf_conf_proto protoreflect.FileDescriptor

var file_internal__conf_conf_proto_rawDesc = []byte{
	0x0a, 0x19, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x2f, 0x63, 0x6f, 0x6e, 0x66,
	0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x6b, 0x72, 0x61,
	0x74, 0x6f, 0x73, 0x2e, 0x61, 0x70, 0x69, 0x1a, 0x0f, 0x63, 0x6f, 0x6e, 0x66, 0x2f, 0x63, 0x6f,
	0x6e, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb3, 0x01, 0x0a, 0x09, 0x42, 0x6f, 0x6f,
	0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x12, 0x24, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6b, 0x72, 0x61, 0x74, 0x6f, 0x73, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x2a, 0x0a, 0x08,
	0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e,
	0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x53, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x52, 0x08,
	0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x12, 0x2a, 0x0a, 0x08, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x63, 0x6f, 0x6e,
	0x66, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x52, 0x08, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x12, 0x28, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6b, 0x72, 0x61, 0x74, 0x6f, 0x73, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x22, 0x58,
	0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x2d, 0x0a, 0x09, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61,
	0x73, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x63, 0x6f, 0x6e, 0x66,
	0x2e, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x73, 0x52, 0x09, 0x64, 0x61, 0x74, 0x61,
	0x62, 0x61, 0x73, 0x65, 0x73, 0x12, 0x21, 0x0a, 0x05, 0x72, 0x65, 0x64, 0x69, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x52, 0x65, 0x64, 0x69,
	0x73, 0x52, 0x05, 0x72, 0x65, 0x64, 0x69, 0x73, 0x22, 0x3f, 0x0a, 0x05, 0x41, 0x64, 0x6d, 0x69,
	0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a,
	0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x61, 0x0a, 0x08, 0x55, 0x73, 0x65,
	0x72, 0x43, 0x6f, 0x6e, 0x66, 0x12, 0x2c, 0x0a, 0x12, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x5f, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x5f, 0x6d, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x10, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x53, 0x63, 0x6f, 0x72, 0x65,
	0x4d, 0x69, 0x6e, 0x12, 0x27, 0x0a, 0x05, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6b, 0x72, 0x61, 0x74, 0x6f, 0x73, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x41, 0x64, 0x6d, 0x69, 0x6e, 0x52, 0x05, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x42, 0x39, 0x5a, 0x37,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x78, 0x69, 0x61,
	0x6f, 0x79, 0x2f, 0x67, 0x6f, 0x2d, 0x73, 0x61, 0x61, 0x73, 0x2d, 0x6b, 0x69, 0x74, 0x2f, 0x75,
	0x73, 0x65, 0x72, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x2f, 0x63, 0x6f,
	0x6e, 0x66, 0x3b, 0x63, 0x6f, 0x6e, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal__conf_conf_proto_rawDescOnce sync.Once
	file_internal__conf_conf_proto_rawDescData = file_internal__conf_conf_proto_rawDesc
)

func file_internal__conf_conf_proto_rawDescGZIP() []byte {
	file_internal__conf_conf_proto_rawDescOnce.Do(func() {
		file_internal__conf_conf_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal__conf_conf_proto_rawDescData)
	})
	return file_internal__conf_conf_proto_rawDescData
}

var file_internal__conf_conf_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_internal__conf_conf_proto_goTypes = []interface{}{
	(*Bootstrap)(nil),      // 0: kratos.api.Bootstrap
	(*Data)(nil),           // 1: kratos.api.Data
	(*Admin)(nil),          // 2: kratos.api.Admin
	(*UserConf)(nil),       // 3: kratos.api.UserConf
	(*conf.Security)(nil),  // 4: conf.Security
	(*conf.Services)(nil),  // 5: conf.Services
	(*conf.Databases)(nil), // 6: conf.Databases
	(*conf.Redis)(nil),     // 7: conf.Redis
}
var file_internal__conf_conf_proto_depIdxs = []int32{
	1, // 0: kratos.api.Bootstrap.data:type_name -> kratos.api.Data
	4, // 1: kratos.api.Bootstrap.security:type_name -> conf.Security
	5, // 2: kratos.api.Bootstrap.services:type_name -> conf.Services
	3, // 3: kratos.api.Bootstrap.user:type_name -> kratos.api.UserConf
	6, // 4: kratos.api.Data.databases:type_name -> conf.Databases
	7, // 5: kratos.api.Data.redis:type_name -> conf.Redis
	2, // 6: kratos.api.UserConf.admin:type_name -> kratos.api.Admin
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_internal__conf_conf_proto_init() }
func file_internal__conf_conf_proto_init() {
	if File_internal__conf_conf_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal__conf_conf_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Bootstrap); i {
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
		file_internal__conf_conf_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Data); i {
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
		file_internal__conf_conf_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Admin); i {
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
		file_internal__conf_conf_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserConf); i {
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
			RawDescriptor: file_internal__conf_conf_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal__conf_conf_proto_goTypes,
		DependencyIndexes: file_internal__conf_conf_proto_depIdxs,
		MessageInfos:      file_internal__conf_conf_proto_msgTypes,
	}.Build()
	File_internal__conf_conf_proto = out.File
	file_internal__conf_conf_proto_rawDesc = nil
	file_internal__conf_conf_proto_goTypes = nil
	file_internal__conf_conf_proto_depIdxs = nil
}
