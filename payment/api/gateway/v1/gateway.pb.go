// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: payment/api/gateway/v1/gateway.proto

package v1

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "github.com/go-saas/kit/pkg/query"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	_ "google.golang.org/protobuf/types/known/fieldmaskpb"
	_ "google.golang.org/protobuf/types/known/structpb"
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

type GetPaymentMethodRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId string `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	IsTest  bool   `protobuf:"varint,2,opt,name=is_test,json=isTest,proto3" json:"is_test,omitempty"`
}

func (x *GetPaymentMethodRequest) Reset() {
	*x = GetPaymentMethodRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_payment_api_gateway_v1_gateway_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPaymentMethodRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPaymentMethodRequest) ProtoMessage() {}

func (x *GetPaymentMethodRequest) ProtoReflect() protoreflect.Message {
	mi := &file_payment_api_gateway_v1_gateway_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPaymentMethodRequest.ProtoReflect.Descriptor instead.
func (*GetPaymentMethodRequest) Descriptor() ([]byte, []int) {
	return file_payment_api_gateway_v1_gateway_proto_rawDescGZIP(), []int{0}
}

func (x *GetPaymentMethodRequest) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *GetPaymentMethodRequest) GetIsTest() bool {
	if x != nil {
		return x.IsTest
	}
	return false
}

type GetPaymentMethodReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Methods []*PaymentMethod `protobuf:"bytes,1,rep,name=methods,proto3" json:"methods,omitempty"`
}

func (x *GetPaymentMethodReply) Reset() {
	*x = GetPaymentMethodReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_payment_api_gateway_v1_gateway_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPaymentMethodReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPaymentMethodReply) ProtoMessage() {}

func (x *GetPaymentMethodReply) ProtoReflect() protoreflect.Message {
	mi := &file_payment_api_gateway_v1_gateway_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPaymentMethodReply.ProtoReflect.Descriptor instead.
func (*GetPaymentMethodReply) Descriptor() ([]byte, []int) {
	return file_payment_api_gateway_v1_gateway_proto_rawDescGZIP(), []int{1}
}

func (x *GetPaymentMethodReply) GetMethods() []*PaymentMethod {
	if x != nil {
		return x.Methods
	}
	return nil
}

type StripeWebhookRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StripeWebhookRequest) Reset() {
	*x = StripeWebhookRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_payment_api_gateway_v1_gateway_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StripeWebhookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StripeWebhookRequest) ProtoMessage() {}

func (x *StripeWebhookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_payment_api_gateway_v1_gateway_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StripeWebhookRequest.ProtoReflect.Descriptor instead.
func (*StripeWebhookRequest) Descriptor() ([]byte, []int) {
	return file_payment_api_gateway_v1_gateway_proto_rawDescGZIP(), []int{2}
}

type StripeWebhookReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StripeWebhookReply) Reset() {
	*x = StripeWebhookReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_payment_api_gateway_v1_gateway_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StripeWebhookReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StripeWebhookReply) ProtoMessage() {}

func (x *StripeWebhookReply) ProtoReflect() protoreflect.Message {
	mi := &file_payment_api_gateway_v1_gateway_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StripeWebhookReply.ProtoReflect.Descriptor instead.
func (*StripeWebhookReply) Descriptor() ([]byte, []int) {
	return file_payment_api_gateway_v1_gateway_proto_rawDescGZIP(), []int{3}
}

type PaymentMethod struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Logo   string `protobuf:"bytes,2,opt,name=logo,proto3" json:"logo,omitempty"`
	Desc   string `protobuf:"bytes,3,opt,name=desc,proto3" json:"desc,omitempty"`
	Notice string `protobuf:"bytes,4,opt,name=notice,proto3" json:"notice,omitempty"`
}

func (x *PaymentMethod) Reset() {
	*x = PaymentMethod{}
	if protoimpl.UnsafeEnabled {
		mi := &file_payment_api_gateway_v1_gateway_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PaymentMethod) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaymentMethod) ProtoMessage() {}

func (x *PaymentMethod) ProtoReflect() protoreflect.Message {
	mi := &file_payment_api_gateway_v1_gateway_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaymentMethod.ProtoReflect.Descriptor instead.
func (*PaymentMethod) Descriptor() ([]byte, []int) {
	return file_payment_api_gateway_v1_gateway_proto_rawDescGZIP(), []int{4}
}

func (x *PaymentMethod) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PaymentMethod) GetLogo() string {
	if x != nil {
		return x.Logo
	}
	return ""
}

func (x *PaymentMethod) GetDesc() string {
	if x != nil {
		return x.Desc
	}
	return ""
}

func (x *PaymentMethod) GetNotice() string {
	if x != nil {
		return x.Notice
	}
	return ""
}

type CreateStripePaymentIntentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId string `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
}

func (x *CreateStripePaymentIntentRequest) Reset() {
	*x = CreateStripePaymentIntentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_payment_api_gateway_v1_gateway_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateStripePaymentIntentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateStripePaymentIntentRequest) ProtoMessage() {}

func (x *CreateStripePaymentIntentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_payment_api_gateway_v1_gateway_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateStripePaymentIntentRequest.ProtoReflect.Descriptor instead.
func (*CreateStripePaymentIntentRequest) Descriptor() ([]byte, []int) {
	return file_payment_api_gateway_v1_gateway_proto_rawDescGZIP(), []int{5}
}

func (x *CreateStripePaymentIntentRequest) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

type CreateStripePaymentIntentReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PaymentIntent string `protobuf:"bytes,1,opt,name=payment_intent,json=paymentIntent,proto3" json:"payment_intent,omitempty"`
	EphemeralKey  string `protobuf:"bytes,2,opt,name=ephemeral_key,json=ephemeralKey,proto3" json:"ephemeral_key,omitempty"`
	CustomerId    string `protobuf:"bytes,3,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
}

func (x *CreateStripePaymentIntentReply) Reset() {
	*x = CreateStripePaymentIntentReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_payment_api_gateway_v1_gateway_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateStripePaymentIntentReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateStripePaymentIntentReply) ProtoMessage() {}

func (x *CreateStripePaymentIntentReply) ProtoReflect() protoreflect.Message {
	mi := &file_payment_api_gateway_v1_gateway_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateStripePaymentIntentReply.ProtoReflect.Descriptor instead.
func (*CreateStripePaymentIntentReply) Descriptor() ([]byte, []int) {
	return file_payment_api_gateway_v1_gateway_proto_rawDescGZIP(), []int{6}
}

func (x *CreateStripePaymentIntentReply) GetPaymentIntent() string {
	if x != nil {
		return x.PaymentIntent
	}
	return ""
}

func (x *CreateStripePaymentIntentReply) GetEphemeralKey() string {
	if x != nil {
		return x.EphemeralKey
	}
	return ""
}

func (x *CreateStripePaymentIntentReply) GetCustomerId() string {
	if x != nil {
		return x.CustomerId
	}
	return ""
}

var File_payment_api_gateway_v1_gateway_proto protoreflect.FileDescriptor

var file_payment_api_gateway_v1_gateway_proto_rawDesc = []byte{
	0x0a, 0x24, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x61,
	0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x76, 0x31, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74,
	0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65, 0x68,
	0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2f, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76,
	0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4d, 0x0a, 0x17, 0x47,
	0x65, 0x74, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x73, 0x5f, 0x74, 0x65, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x06, 0x69, 0x73, 0x54, 0x65, 0x73, 0x74, 0x22, 0x58, 0x0a, 0x15, 0x47, 0x65,
	0x74, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x3f, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x52, 0x07, 0x6d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x73, 0x22, 0x16, 0x0a, 0x14, 0x53, 0x74, 0x72, 0x69, 0x70, 0x65, 0x57, 0x65,
	0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x14, 0x0a, 0x12,
	0x53, 0x74, 0x72, 0x69, 0x70, 0x65, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x22, 0x63, 0x0a, 0x0d, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6c, 0x6f, 0x67, 0x6f, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6c, 0x6f, 0x67, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x64,
	0x65, 0x73, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x65, 0x73, 0x63, 0x12,
	0x16, 0x0a, 0x06, 0x6e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x6e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x22, 0x3d, 0x0a, 0x20, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x53, 0x74, 0x72, 0x69, 0x70, 0x65, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x22, 0x8d, 0x01, 0x0a, 0x1e, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x53, 0x74, 0x72, 0x69, 0x70, 0x65, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x25, 0x0a, 0x0e, 0x70, 0x61, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x12, 0x23, 0x0a, 0x0d, 0x65, 0x70, 0x68, 0x65, 0x6d, 0x65, 0x72, 0x61, 0x6c, 0x5f, 0x6b, 0x65,
	0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x70, 0x68, 0x65, 0x6d, 0x65, 0x72,
	0x61, 0x6c, 0x4b, 0x65, 0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x75, 0x73, 0x74,
	0x6f, 0x6d, 0x65, 0x72, 0x49, 0x64, 0x32, 0xa9, 0x01, 0x0a, 0x15, 0x50, 0x61, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x8f, 0x01, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x4d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x2f, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x47,
	0x65, 0x74, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x76, 0x31, 0x2e,
	0x47, 0x65, 0x74, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x12, 0x13, 0x2f,
	0x76, 0x31, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x6d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x73, 0x32, 0xcf, 0x02, 0x0a, 0x1b, 0x53, 0x74, 0x72, 0x69, 0x70, 0x65, 0x50, 0x61, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0xb3, 0x01, 0x0a, 0x19, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x74, 0x72,
	0x69, 0x70, 0x65, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x12, 0x38, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x67,
	0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x53, 0x74, 0x72, 0x69, 0x70, 0x65, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x36, 0x2e, 0x70, 0x61, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x74, 0x72, 0x69, 0x70, 0x65,
	0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x22, 0x24, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1e, 0x3a, 0x01, 0x2a, 0x22, 0x19, 0x2f,
	0x76, 0x31, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x73, 0x74, 0x72, 0x69, 0x70,
	0x65, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x7a, 0x0a, 0x0d, 0x53, 0x74, 0x72, 0x69,
	0x70, 0x65, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x1a, 0x2a, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x70,
	0x65, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x25, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x1f, 0x3a, 0x01, 0x2a, 0x22, 0x1a, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x61,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x73, 0x74, 0x72, 0x69, 0x70, 0x65, 0x2f, 0x77, 0x65, 0x62,
	0x68, 0x6f, 0x6f, 0x6b, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x2d, 0x73, 0x61, 0x61, 0x73, 0x2f, 0x6b, 0x69, 0x74, 0x2f, 0x70,
	0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77,
	0x61, 0x79, 0x2f, 0x76, 0x31, 0x3b, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_payment_api_gateway_v1_gateway_proto_rawDescOnce sync.Once
	file_payment_api_gateway_v1_gateway_proto_rawDescData = file_payment_api_gateway_v1_gateway_proto_rawDesc
)

func file_payment_api_gateway_v1_gateway_proto_rawDescGZIP() []byte {
	file_payment_api_gateway_v1_gateway_proto_rawDescOnce.Do(func() {
		file_payment_api_gateway_v1_gateway_proto_rawDescData = protoimpl.X.CompressGZIP(file_payment_api_gateway_v1_gateway_proto_rawDescData)
	})
	return file_payment_api_gateway_v1_gateway_proto_rawDescData
}

var file_payment_api_gateway_v1_gateway_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_payment_api_gateway_v1_gateway_proto_goTypes = []interface{}{
	(*GetPaymentMethodRequest)(nil),          // 0: payment.api.gateway.v1.GetPaymentMethodRequest
	(*GetPaymentMethodReply)(nil),            // 1: payment.api.gateway.v1.GetPaymentMethodReply
	(*StripeWebhookRequest)(nil),             // 2: payment.api.gateway.v1.StripeWebhookRequest
	(*StripeWebhookReply)(nil),               // 3: payment.api.gateway.v1.StripeWebhookReply
	(*PaymentMethod)(nil),                    // 4: payment.api.gateway.v1.PaymentMethod
	(*CreateStripePaymentIntentRequest)(nil), // 5: payment.api.gateway.v1.CreateStripePaymentIntentRequest
	(*CreateStripePaymentIntentReply)(nil),   // 6: payment.api.gateway.v1.CreateStripePaymentIntentReply
	(*emptypb.Empty)(nil),                    // 7: google.protobuf.Empty
}
var file_payment_api_gateway_v1_gateway_proto_depIdxs = []int32{
	4, // 0: payment.api.gateway.v1.GetPaymentMethodReply.methods:type_name -> payment.api.gateway.v1.PaymentMethod
	0, // 1: payment.api.gateway.v1.PaymentGatewayService.GetPaymentMethod:input_type -> payment.api.gateway.v1.GetPaymentMethodRequest
	5, // 2: payment.api.gateway.v1.StripePaymentGatewayService.CreateStripePaymentIntent:input_type -> payment.api.gateway.v1.CreateStripePaymentIntentRequest
	7, // 3: payment.api.gateway.v1.StripePaymentGatewayService.StripeWebhook:input_type -> google.protobuf.Empty
	1, // 4: payment.api.gateway.v1.PaymentGatewayService.GetPaymentMethod:output_type -> payment.api.gateway.v1.GetPaymentMethodReply
	6, // 5: payment.api.gateway.v1.StripePaymentGatewayService.CreateStripePaymentIntent:output_type -> payment.api.gateway.v1.CreateStripePaymentIntentReply
	3, // 6: payment.api.gateway.v1.StripePaymentGatewayService.StripeWebhook:output_type -> payment.api.gateway.v1.StripeWebhookReply
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_payment_api_gateway_v1_gateway_proto_init() }
func file_payment_api_gateway_v1_gateway_proto_init() {
	if File_payment_api_gateway_v1_gateway_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_payment_api_gateway_v1_gateway_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPaymentMethodRequest); i {
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
		file_payment_api_gateway_v1_gateway_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPaymentMethodReply); i {
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
		file_payment_api_gateway_v1_gateway_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StripeWebhookRequest); i {
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
		file_payment_api_gateway_v1_gateway_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StripeWebhookReply); i {
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
		file_payment_api_gateway_v1_gateway_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PaymentMethod); i {
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
		file_payment_api_gateway_v1_gateway_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateStripePaymentIntentRequest); i {
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
		file_payment_api_gateway_v1_gateway_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateStripePaymentIntentReply); i {
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
			RawDescriptor: file_payment_api_gateway_v1_gateway_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_payment_api_gateway_v1_gateway_proto_goTypes,
		DependencyIndexes: file_payment_api_gateway_v1_gateway_proto_depIdxs,
		MessageInfos:      file_payment_api_gateway_v1_gateway_proto_msgTypes,
	}.Build()
	File_payment_api_gateway_v1_gateway_proto = out.File
	file_payment_api_gateway_v1_gateway_proto_rawDesc = nil
	file_payment_api_gateway_v1_gateway_proto_goTypes = nil
	file_payment_api_gateway_v1_gateway_proto_depIdxs = nil
}
