// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: product/api/product/v1/brand.proto

package v1

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	query "github.com/go-saas/kit/pkg/query"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CreateBrandRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code string `protobuf:"bytes,2,opt,name=code,proto3" json:"code,omitempty"`
	Name string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Logo string `protobuf:"bytes,4,opt,name=logo,proto3" json:"logo,omitempty"`
	Url  string `protobuf:"bytes,6,opt,name=url,proto3" json:"url,omitempty"`
	Desc string `protobuf:"bytes,7,opt,name=desc,proto3" json:"desc,omitempty"`
}

func (x *CreateBrandRequest) Reset() {
	*x = CreateBrandRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_api_product_v1_brand_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateBrandRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateBrandRequest) ProtoMessage() {}

func (x *CreateBrandRequest) ProtoReflect() protoreflect.Message {
	mi := &file_product_api_product_v1_brand_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateBrandRequest.ProtoReflect.Descriptor instead.
func (*CreateBrandRequest) Descriptor() ([]byte, []int) {
	return file_product_api_product_v1_brand_proto_rawDescGZIP(), []int{0}
}

func (x *CreateBrandRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *CreateBrandRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateBrandRequest) GetLogo() string {
	if x != nil {
		return x.Logo
	}
	return ""
}

func (x *CreateBrandRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *CreateBrandRequest) GetDesc() string {
	if x != nil {
		return x.Desc
	}
	return ""
}

type UpdateBrandRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Brand      *UpdateBrand           `protobuf:"bytes,1,opt,name=brand,proto3" json:"brand,omitempty"`
	UpdateMask *fieldmaskpb.FieldMask `protobuf:"bytes,2,opt,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
}

func (x *UpdateBrandRequest) Reset() {
	*x = UpdateBrandRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_api_product_v1_brand_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateBrandRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateBrandRequest) ProtoMessage() {}

func (x *UpdateBrandRequest) ProtoReflect() protoreflect.Message {
	mi := &file_product_api_product_v1_brand_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateBrandRequest.ProtoReflect.Descriptor instead.
func (*UpdateBrandRequest) Descriptor() ([]byte, []int) {
	return file_product_api_product_v1_brand_proto_rawDescGZIP(), []int{1}
}

func (x *UpdateBrandRequest) GetBrand() *UpdateBrand {
	if x != nil {
		return x.Brand
	}
	return nil
}

func (x *UpdateBrandRequest) GetUpdateMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.UpdateMask
	}
	return nil
}

type UpdateBrand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Code string `protobuf:"bytes,2,opt,name=code,proto3" json:"code,omitempty"`
	Name string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Logo string `protobuf:"bytes,4,opt,name=logo,proto3" json:"logo,omitempty"`
	Url  string `protobuf:"bytes,6,opt,name=url,proto3" json:"url,omitempty"`
	Desc string `protobuf:"bytes,7,opt,name=desc,proto3" json:"desc,omitempty"`
}

func (x *UpdateBrand) Reset() {
	*x = UpdateBrand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_api_product_v1_brand_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateBrand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateBrand) ProtoMessage() {}

func (x *UpdateBrand) ProtoReflect() protoreflect.Message {
	mi := &file_product_api_product_v1_brand_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateBrand.ProtoReflect.Descriptor instead.
func (*UpdateBrand) Descriptor() ([]byte, []int) {
	return file_product_api_product_v1_brand_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateBrand) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateBrand) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *UpdateBrand) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateBrand) GetLogo() string {
	if x != nil {
		return x.Logo
	}
	return ""
}

func (x *UpdateBrand) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *UpdateBrand) GetDesc() string {
	if x != nil {
		return x.Desc
	}
	return ""
}

type DeleteBrandRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteBrandRequest) Reset() {
	*x = DeleteBrandRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_api_product_v1_brand_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteBrandRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteBrandRequest) ProtoMessage() {}

func (x *DeleteBrandRequest) ProtoReflect() protoreflect.Message {
	mi := &file_product_api_product_v1_brand_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteBrandRequest.ProtoReflect.Descriptor instead.
func (*DeleteBrandRequest) Descriptor() ([]byte, []int) {
	return file_product_api_product_v1_brand_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteBrandRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type DeleteBrandReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Code string `protobuf:"bytes,2,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *DeleteBrandReply) Reset() {
	*x = DeleteBrandReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_api_product_v1_brand_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteBrandReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteBrandReply) ProtoMessage() {}

func (x *DeleteBrandReply) ProtoReflect() protoreflect.Message {
	mi := &file_product_api_product_v1_brand_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteBrandReply.ProtoReflect.Descriptor instead.
func (*DeleteBrandReply) Descriptor() ([]byte, []int) {
	return file_product_api_product_v1_brand_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteBrandReply) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *DeleteBrandReply) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

type GetBrandRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetBrandRequest) Reset() {
	*x = GetBrandRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_api_product_v1_brand_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBrandRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBrandRequest) ProtoMessage() {}

func (x *GetBrandRequest) ProtoReflect() protoreflect.Message {
	mi := &file_product_api_product_v1_brand_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBrandRequest.ProtoReflect.Descriptor instead.
func (*GetBrandRequest) Descriptor() ([]byte, []int) {
	return file_product_api_product_v1_brand_proto_rawDescGZIP(), []int{5}
}

func (x *GetBrandRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type BrandFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   *query.StringFilterOperation `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Code *query.StringFilterOperation `protobuf:"bytes,2,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *BrandFilter) Reset() {
	*x = BrandFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_api_product_v1_brand_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BrandFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BrandFilter) ProtoMessage() {}

func (x *BrandFilter) ProtoReflect() protoreflect.Message {
	mi := &file_product_api_product_v1_brand_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BrandFilter.ProtoReflect.Descriptor instead.
func (*BrandFilter) Descriptor() ([]byte, []int) {
	return file_product_api_product_v1_brand_proto_rawDescGZIP(), []int{6}
}

func (x *BrandFilter) GetId() *query.StringFilterOperation {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *BrandFilter) GetCode() *query.StringFilterOperation {
	if x != nil {
		return x.Code
	}
	return nil
}

type ListBrandRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageOffset int32                  `protobuf:"varint,1,opt,name=page_offset,json=pageOffset,proto3" json:"page_offset,omitempty"`
	PageSize   int32                  `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	Search     string                 `protobuf:"bytes,3,opt,name=search,proto3" json:"search,omitempty"`
	Sort       []string               `protobuf:"bytes,4,rep,name=sort,proto3" json:"sort,omitempty"`
	Fields     *fieldmaskpb.FieldMask `protobuf:"bytes,5,opt,name=fields,proto3" json:"fields,omitempty"`
	Filter     *BrandFilter           `protobuf:"bytes,6,opt,name=filter,proto3" json:"filter,omitempty"`
}

func (x *ListBrandRequest) Reset() {
	*x = ListBrandRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_api_product_v1_brand_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListBrandRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListBrandRequest) ProtoMessage() {}

func (x *ListBrandRequest) ProtoReflect() protoreflect.Message {
	mi := &file_product_api_product_v1_brand_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListBrandRequest.ProtoReflect.Descriptor instead.
func (*ListBrandRequest) Descriptor() ([]byte, []int) {
	return file_product_api_product_v1_brand_proto_rawDescGZIP(), []int{7}
}

func (x *ListBrandRequest) GetPageOffset() int32 {
	if x != nil {
		return x.PageOffset
	}
	return 0
}

func (x *ListBrandRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListBrandRequest) GetSearch() string {
	if x != nil {
		return x.Search
	}
	return ""
}

func (x *ListBrandRequest) GetSort() []string {
	if x != nil {
		return x.Sort
	}
	return nil
}

func (x *ListBrandRequest) GetFields() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ListBrandRequest) GetFilter() *BrandFilter {
	if x != nil {
		return x.Filter
	}
	return nil
}

type ListBrandReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TotalSize  int32    `protobuf:"varint,1,opt,name=total_size,json=totalSize,proto3" json:"total_size,omitempty"`
	FilterSize int32    `protobuf:"varint,2,opt,name=filter_size,json=filterSize,proto3" json:"filter_size,omitempty"`
	Items      []*Brand `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *ListBrandReply) Reset() {
	*x = ListBrandReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_api_product_v1_brand_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListBrandReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListBrandReply) ProtoMessage() {}

func (x *ListBrandReply) ProtoReflect() protoreflect.Message {
	mi := &file_product_api_product_v1_brand_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListBrandReply.ProtoReflect.Descriptor instead.
func (*ListBrandReply) Descriptor() ([]byte, []int) {
	return file_product_api_product_v1_brand_proto_rawDescGZIP(), []int{8}
}

func (x *ListBrandReply) GetTotalSize() int32 {
	if x != nil {
		return x.TotalSize
	}
	return 0
}

func (x *ListBrandReply) GetFilterSize() int32 {
	if x != nil {
		return x.FilterSize
	}
	return 0
}

func (x *ListBrandReply) GetItems() []*Brand {
	if x != nil {
		return x.Items
	}
	return nil
}

type Brand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Code string `protobuf:"bytes,2,opt,name=code,proto3" json:"code,omitempty"`
	Name string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Logo string `protobuf:"bytes,4,opt,name=logo,proto3" json:"logo,omitempty"`
	Url  string `protobuf:"bytes,6,opt,name=url,proto3" json:"url,omitempty"`
	Desc string `protobuf:"bytes,7,opt,name=desc,proto3" json:"desc,omitempty"`
}

func (x *Brand) Reset() {
	*x = Brand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_api_product_v1_brand_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Brand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Brand) ProtoMessage() {}

func (x *Brand) ProtoReflect() protoreflect.Message {
	mi := &file_product_api_product_v1_brand_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Brand.ProtoReflect.Descriptor instead.
func (*Brand) Descriptor() ([]byte, []int) {
	return file_product_api_product_v1_brand_proto_rawDescGZIP(), []int{9}
}

func (x *Brand) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Brand) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *Brand) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Brand) GetLogo() string {
	if x != nil {
		return x.Logo
	}
	return ""
}

func (x *Brand) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Brand) GetDesc() string {
	if x != nil {
		return x.Desc
	}
	return ""
}

var File_product_api_product_v1_brand_proto protoreflect.FileDescriptor

var file_product_api_product_v1_brand_proto_rawDesc = []byte{
	0x0a, 0x22, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f,
	0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2f, 0x6f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61,
	0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x82,
	0x01, 0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x0a, 0xe0, 0x41, 0x02, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6c, 0x6f, 0x67,
	0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6c, 0x6f, 0x67, 0x6f, 0x12, 0x10, 0x0a,
	0x03, 0x75, 0x72, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12,
	0x12, 0x0a, 0x04, 0x64, 0x65, 0x73, 0x63, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64,
	0x65, 0x73, 0x63, 0x22, 0x99, 0x01, 0x0a, 0x12, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x72,
	0x61, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x46, 0x0a, 0x05, 0x62, 0x72,
	0x61, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e,
	0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x42, 0x0b,
	0xe0, 0x41, 0x02, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x05, 0x62, 0x72, 0x61,
	0x6e, 0x64, 0x12, 0x3b, 0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x6d, 0x61, 0x73,
	0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4d,
	0x61, 0x73, 0x6b, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x73, 0x6b, 0x22,
	0x8b, 0x01, 0x0a, 0x0b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x12,
	0x1a, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xe0, 0x41, 0x02,
	0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x63,
	0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6c, 0x6f, 0x67, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6c, 0x6f, 0x67, 0x6f, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x65, 0x73,
	0x63, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x65, 0x73, 0x63, 0x22, 0x24, 0x0a,
	0x12, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x22, 0x36, 0x0a, 0x10, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x42, 0x72, 0x61,
	0x6e, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x2d, 0x0a, 0x0f, 0x47,
	0x65, 0x74, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xe0, 0x41, 0x02, 0xfa,
	0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x02, 0x69, 0x64, 0x22, 0x81, 0x01, 0x0a, 0x0b, 0x42,
	0x72, 0x61, 0x6e, 0x64, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x36, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2e, 0x6f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x46,
	0x69, 0x6c, 0x74, 0x65, 0x72, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x3a, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x26, 0x2e, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x4f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0xed,
	0x01, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x6f, 0x66, 0x66, 0x73,
	0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x4f, 0x66,
	0x66, 0x73, 0x65, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6f, 0x72,
	0x74, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x12, 0x32, 0x0a,
	0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x4d, 0x61, 0x73, 0x6b, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64,
	0x73, 0x12, 0x3b, 0x0a, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x72, 0x61, 0x6e, 0x64,
	0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x22, 0x85,
	0x01, 0x0a, 0x0e, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x53, 0x69, 0x7a, 0x65,
	0x12, 0x1f, 0x0a, 0x0b, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x53, 0x69, 0x7a,
	0x65, 0x12, 0x33, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x52,
	0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x79, 0x0a, 0x05, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63,
	0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6c, 0x6f, 0x67, 0x6f, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6c, 0x6f, 0x67, 0x6f, 0x12, 0x10, 0x0a, 0x03, 0x75,
	0x72, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x12, 0x0a,
	0x04, 0x64, 0x65, 0x73, 0x63, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x65, 0x73,
	0x63, 0x32, 0xc0, 0x05, 0x0a, 0x0c, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x96, 0x01, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x72, 0x61, 0x6e, 0x64,
	0x12, 0x28, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x72,
	0x61, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x22, 0x37, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x31, 0x5a, 0x1b, 0x3a, 0x01, 0x2a, 0x22,
	0x16, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2f, 0x62, 0x72, 0x61,
	0x6e, 0x64, 0x2f, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x12, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x2f, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x73, 0x12, 0x72, 0x0a, 0x08, 0x47,
	0x65, 0x74, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x12, 0x27, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x76, 0x31,
	0x2e, 0x47, 0x65, 0x74, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x22,
	0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x12, 0x16, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x2f, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12,
	0x76, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x12, 0x2a,
	0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x72,
	0x61, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x22, 0x1c, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x16, 0x3a, 0x01, 0x2a, 0x22, 0x11, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x74, 0x2f, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x12, 0xa4, 0x01, 0x0a, 0x0b, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x12, 0x2a, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x76, 0x31,
	0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x72, 0x61,
	0x6e, 0x64, 0x22, 0x4a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x44, 0x3a, 0x01, 0x2a, 0x5a, 0x21, 0x3a,
	0x01, 0x2a, 0x32, 0x1c, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2f,
	0x62, 0x72, 0x61, 0x6e, 0x64, 0x2f, 0x7b, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x2e, 0x69, 0x64, 0x7d,
	0x1a, 0x1c, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2f, 0x62, 0x72,
	0x61, 0x6e, 0x64, 0x2f, 0x7b, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x2e, 0x69, 0x64, 0x7d, 0x12, 0x83,
	0x01, 0x0a, 0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x12, 0x2a,
	0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x42, 0x72,
	0x61, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x2a, 0x16, 0x2f, 0x76,
	0x31, 0x2f, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2f, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x2f,
	0x7b, 0x69, 0x64, 0x7d, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x2d, 0x73, 0x61, 0x61, 0x73, 0x2f, 0x6b, 0x69, 0x74, 0x2f, 0x70,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x2f, 0x76, 0x31, 0x3b, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_product_api_product_v1_brand_proto_rawDescOnce sync.Once
	file_product_api_product_v1_brand_proto_rawDescData = file_product_api_product_v1_brand_proto_rawDesc
)

func file_product_api_product_v1_brand_proto_rawDescGZIP() []byte {
	file_product_api_product_v1_brand_proto_rawDescOnce.Do(func() {
		file_product_api_product_v1_brand_proto_rawDescData = protoimpl.X.CompressGZIP(file_product_api_product_v1_brand_proto_rawDescData)
	})
	return file_product_api_product_v1_brand_proto_rawDescData
}

var file_product_api_product_v1_brand_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_product_api_product_v1_brand_proto_goTypes = []interface{}{
	(*CreateBrandRequest)(nil),          // 0: product.api.product.v1.CreateBrandRequest
	(*UpdateBrandRequest)(nil),          // 1: product.api.product.v1.UpdateBrandRequest
	(*UpdateBrand)(nil),                 // 2: product.api.product.v1.UpdateBrand
	(*DeleteBrandRequest)(nil),          // 3: product.api.product.v1.DeleteBrandRequest
	(*DeleteBrandReply)(nil),            // 4: product.api.product.v1.DeleteBrandReply
	(*GetBrandRequest)(nil),             // 5: product.api.product.v1.GetBrandRequest
	(*BrandFilter)(nil),                 // 6: product.api.product.v1.BrandFilter
	(*ListBrandRequest)(nil),            // 7: product.api.product.v1.ListBrandRequest
	(*ListBrandReply)(nil),              // 8: product.api.product.v1.ListBrandReply
	(*Brand)(nil),                       // 9: product.api.product.v1.Brand
	(*fieldmaskpb.FieldMask)(nil),       // 10: google.protobuf.FieldMask
	(*query.StringFilterOperation)(nil), // 11: query.operation.StringFilterOperation
}
var file_product_api_product_v1_brand_proto_depIdxs = []int32{
	2,  // 0: product.api.product.v1.UpdateBrandRequest.brand:type_name -> product.api.product.v1.UpdateBrand
	10, // 1: product.api.product.v1.UpdateBrandRequest.update_mask:type_name -> google.protobuf.FieldMask
	11, // 2: product.api.product.v1.BrandFilter.id:type_name -> query.operation.StringFilterOperation
	11, // 3: product.api.product.v1.BrandFilter.code:type_name -> query.operation.StringFilterOperation
	10, // 4: product.api.product.v1.ListBrandRequest.fields:type_name -> google.protobuf.FieldMask
	6,  // 5: product.api.product.v1.ListBrandRequest.filter:type_name -> product.api.product.v1.BrandFilter
	9,  // 6: product.api.product.v1.ListBrandReply.items:type_name -> product.api.product.v1.Brand
	7,  // 7: product.api.product.v1.BrandService.ListBrand:input_type -> product.api.product.v1.ListBrandRequest
	5,  // 8: product.api.product.v1.BrandService.GetBrand:input_type -> product.api.product.v1.GetBrandRequest
	0,  // 9: product.api.product.v1.BrandService.CreateBrand:input_type -> product.api.product.v1.CreateBrandRequest
	1,  // 10: product.api.product.v1.BrandService.UpdateBrand:input_type -> product.api.product.v1.UpdateBrandRequest
	3,  // 11: product.api.product.v1.BrandService.DeleteBrand:input_type -> product.api.product.v1.DeleteBrandRequest
	8,  // 12: product.api.product.v1.BrandService.ListBrand:output_type -> product.api.product.v1.ListBrandReply
	9,  // 13: product.api.product.v1.BrandService.GetBrand:output_type -> product.api.product.v1.Brand
	9,  // 14: product.api.product.v1.BrandService.CreateBrand:output_type -> product.api.product.v1.Brand
	9,  // 15: product.api.product.v1.BrandService.UpdateBrand:output_type -> product.api.product.v1.Brand
	4,  // 16: product.api.product.v1.BrandService.DeleteBrand:output_type -> product.api.product.v1.DeleteBrandReply
	12, // [12:17] is the sub-list for method output_type
	7,  // [7:12] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_product_api_product_v1_brand_proto_init() }
func file_product_api_product_v1_brand_proto_init() {
	if File_product_api_product_v1_brand_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_product_api_product_v1_brand_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateBrandRequest); i {
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
		file_product_api_product_v1_brand_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateBrandRequest); i {
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
		file_product_api_product_v1_brand_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateBrand); i {
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
		file_product_api_product_v1_brand_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteBrandRequest); i {
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
		file_product_api_product_v1_brand_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteBrandReply); i {
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
		file_product_api_product_v1_brand_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetBrandRequest); i {
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
		file_product_api_product_v1_brand_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BrandFilter); i {
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
		file_product_api_product_v1_brand_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListBrandRequest); i {
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
		file_product_api_product_v1_brand_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListBrandReply); i {
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
		file_product_api_product_v1_brand_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Brand); i {
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
			RawDescriptor: file_product_api_product_v1_brand_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_product_api_product_v1_brand_proto_goTypes,
		DependencyIndexes: file_product_api_product_v1_brand_proto_depIdxs,
		MessageInfos:      file_product_api_product_v1_brand_proto_msgTypes,
	}.Build()
	File_product_api_product_v1_brand_proto = out.File
	file_product_api_product_v1_brand_proto_rawDesc = nil
	file_product_api_product_v1_brand_proto_goTypes = nil
	file_product_api_product_v1_brand_proto_depIdxs = nil
}
