// Code generated by protoc-gen-go-grpc-proxy. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc-proxy v1.2.0
// - protoc             (unknown)
// source: payment/api/subscription/v1/subscription.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

var _ SubscriptionServiceServer = (*subscriptionServiceClientProxy)(nil)

// subscriptionServiceClientProxy is the proxy to turn SubscriptionService client to server interface.
type subscriptionServiceClientProxy struct {
	cc SubscriptionServiceClient
}

func NewSubscriptionServiceClientProxy(cc SubscriptionServiceClient) SubscriptionServiceServer {
	return &subscriptionServiceClientProxy{cc}
}

func (c *subscriptionServiceClientProxy) UpdateSubscription(ctx context.Context, in *UpdateSubscriptionRequest) (*Subscription, error) {
	return c.cc.UpdateSubscription(ctx, in)
}
func (c *subscriptionServiceClientProxy) ListSubscription(ctx context.Context, in *ListSubscriptionRequest) (*ListSubscriptionReply, error) {
	return c.cc.ListSubscription(ctx, in)
}
func (c *subscriptionServiceClientProxy) GetSubscription(ctx context.Context, in *GetSubscriptionRequest) (*Subscription, error) {
	return c.cc.GetSubscription(ctx, in)
}
func (c *subscriptionServiceClientProxy) CancelSubscription(ctx context.Context, in *CancelSubscriptionRequest) (*Subscription, error) {
	return c.cc.CancelSubscription(ctx, in)
}
func (c *subscriptionServiceClientProxy) CreateMySubscription(ctx context.Context, in *CreateMySubscriptionRequest) (*Subscription, error) {
	return c.cc.CreateMySubscription(ctx, in)
}
func (c *subscriptionServiceClientProxy) CancelMySubscription(ctx context.Context, in *CancelSubscriptionRequest) (*Subscription, error) {
	return c.cc.CancelMySubscription(ctx, in)
}
func (c *subscriptionServiceClientProxy) UpdateMySubscription(ctx context.Context, in *UpdateMySubscriptionRequest) (*Subscription, error) {
	return c.cc.UpdateMySubscription(ctx, in)
}
func (c *subscriptionServiceClientProxy) ListMySubscription(ctx context.Context, in *ListMySubscriptionRequest) (*ListMySubscriptionReply, error) {
	return c.cc.ListMySubscription(ctx, in)
}
