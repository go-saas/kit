// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.1
// - protoc             (unknown)
// source: payment/api/gateway/v1/gateway.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationPaymentGatewayServiceGetPaymentMethod = "/payment.api.gateway.v1.PaymentGatewayService/GetPaymentMethod"

type PaymentGatewayServiceHTTPServer interface {
	GetPaymentMethod(context.Context, *GetPaymentMethodRequest) (*GetPaymentMethodReply, error)
}

func RegisterPaymentGatewayServiceHTTPServer(s *http.Server, srv PaymentGatewayServiceHTTPServer) {
	r := s.Route("/")
	r.GET("/v1/payment/methods", _PaymentGatewayService_GetPaymentMethod0_HTTP_Handler(srv))
}

func _PaymentGatewayService_GetPaymentMethod0_HTTP_Handler(srv PaymentGatewayServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetPaymentMethodRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPaymentGatewayServiceGetPaymentMethod)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetPaymentMethod(ctx, req.(*GetPaymentMethodRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetPaymentMethodReply)
		return ctx.Result(200, reply)
	}
}

type PaymentGatewayServiceHTTPClient interface {
	GetPaymentMethod(ctx context.Context, req *GetPaymentMethodRequest, opts ...http.CallOption) (rsp *GetPaymentMethodReply, err error)
}

type PaymentGatewayServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewPaymentGatewayServiceHTTPClient(client *http.Client) PaymentGatewayServiceHTTPClient {
	return &PaymentGatewayServiceHTTPClientImpl{client}
}

func (c *PaymentGatewayServiceHTTPClientImpl) GetPaymentMethod(ctx context.Context, in *GetPaymentMethodRequest, opts ...http.CallOption) (*GetPaymentMethodReply, error) {
	var out GetPaymentMethodReply
	pattern := "/v1/payment/methods"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationPaymentGatewayServiceGetPaymentMethod))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

const OperationStripePaymentGatewayServiceCreateStripePaymentIntent = "/payment.api.gateway.v1.StripePaymentGatewayService/CreateStripePaymentIntent"
const OperationStripePaymentGatewayServiceStripeWebhook = "/payment.api.gateway.v1.StripePaymentGatewayService/StripeWebhook"

type StripePaymentGatewayServiceHTTPServer interface {
	CreateStripePaymentIntent(context.Context, *CreateStripePaymentIntentRequest) (*CreateStripePaymentIntentReply, error)
	StripeWebhook(context.Context, *emptypb.Empty) (*StripeWebhookReply, error)
}

func RegisterStripePaymentGatewayServiceHTTPServer(s *http.Server, srv StripePaymentGatewayServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/v1/payment/stripe/intent", _StripePaymentGatewayService_CreateStripePaymentIntent0_HTTP_Handler(srv))
	r.POST("/v1/payment/stripe/webhook", _StripePaymentGatewayService_StripeWebhook0_HTTP_Handler(srv))
}

func _StripePaymentGatewayService_CreateStripePaymentIntent0_HTTP_Handler(srv StripePaymentGatewayServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateStripePaymentIntentRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationStripePaymentGatewayServiceCreateStripePaymentIntent)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateStripePaymentIntent(ctx, req.(*CreateStripePaymentIntentRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateStripePaymentIntentReply)
		return ctx.Result(200, reply)
	}
}

func _StripePaymentGatewayService_StripeWebhook0_HTTP_Handler(srv StripePaymentGatewayServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in emptypb.Empty
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationStripePaymentGatewayServiceStripeWebhook)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.StripeWebhook(ctx, req.(*emptypb.Empty))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*StripeWebhookReply)
		return ctx.Result(200, reply)
	}
}

type StripePaymentGatewayServiceHTTPClient interface {
	CreateStripePaymentIntent(ctx context.Context, req *CreateStripePaymentIntentRequest, opts ...http.CallOption) (rsp *CreateStripePaymentIntentReply, err error)
	StripeWebhook(ctx context.Context, req *emptypb.Empty, opts ...http.CallOption) (rsp *StripeWebhookReply, err error)
}

type StripePaymentGatewayServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewStripePaymentGatewayServiceHTTPClient(client *http.Client) StripePaymentGatewayServiceHTTPClient {
	return &StripePaymentGatewayServiceHTTPClientImpl{client}
}

func (c *StripePaymentGatewayServiceHTTPClientImpl) CreateStripePaymentIntent(ctx context.Context, in *CreateStripePaymentIntentRequest, opts ...http.CallOption) (*CreateStripePaymentIntentReply, error) {
	var out CreateStripePaymentIntentReply
	pattern := "/v1/payment/stripe/intent"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationStripePaymentGatewayServiceCreateStripePaymentIntent))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *StripePaymentGatewayServiceHTTPClientImpl) StripeWebhook(ctx context.Context, in *emptypb.Empty, opts ...http.CallOption) (*StripeWebhookReply, error) {
	var out StripeWebhookReply
	pattern := "/v1/payment/stripe/webhook"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationStripePaymentGatewayServiceStripeWebhook))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
