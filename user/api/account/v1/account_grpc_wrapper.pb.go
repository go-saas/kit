// Code generated by protoc-gen-go-grpc-wrapper. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc-wrapper v1.2.0
// - protoc             (unknown)
// source: user/api/account/v1/account.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

var _ AccountServer = (*accountClientWrapper)(nil)

// accountClientWrapper is the wrapper to turn Account client to server interface.
//
type accountClientWrapper struct {
	cc AccountClient
}

func NewAccountClientWrapper(cc AccountClient) AccountServer {
	return &accountClientWrapper{cc}
}

func (c *accountClientWrapper) GetProfile(ctx context.Context, in *GetProfileRequest) (*GetProfileResponse, error) {
	return c.cc.GetProfile(ctx, in)
}
func (c *accountClientWrapper) UpdateProfile(ctx context.Context, in *UpdateProfileRequest) (*UpdateProfileResponse, error) {
	return c.cc.UpdateProfile(ctx, in)
}
func (c *accountClientWrapper) GetSettings(ctx context.Context, in *GetSettingsRequest) (*GetSettingsResponse, error) {
	return c.cc.GetSettings(ctx, in)
}
func (c *accountClientWrapper) UpdateSettings(ctx context.Context, in *UpdateSettingsRequest) (*UpdateSettingsResponse, error) {
	return c.cc.UpdateSettings(ctx, in)
}
func (c *accountClientWrapper) GetAddresses(ctx context.Context, in *GetAddressesRequest) (*GetAddressesReply, error) {
	return c.cc.GetAddresses(ctx, in)
}
func (c *accountClientWrapper) CreateAddresses(ctx context.Context, in *CreateAddressesRequest) (*CreateAddressReply, error) {
	return c.cc.CreateAddresses(ctx, in)
}
func (c *accountClientWrapper) UpdateAddresses(ctx context.Context, in *UpdateAddressesRequest) (*UpdateAddressesReply, error) {
	return c.cc.UpdateAddresses(ctx, in)
}
func (c *accountClientWrapper) DeleteAddresses(ctx context.Context, in *DeleteAddressRequest) (*DeleteAddressesReply, error) {
	return c.cc.DeleteAddresses(ctx, in)
}
