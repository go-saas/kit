// Code generated by protoc-gen-go-grpc-proxy. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc-proxy v1.2.0
// - protoc             (unknown)
// source: sys/api/menu/v1/menu.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

var _ MenuServiceServer = (*menuServiceClientProxy)(nil)

// menuServiceClientProxy is the proxy to turn MenuService client to server interface.
type menuServiceClientProxy struct {
	cc MenuServiceClient
}

func NewMenuServiceClientProxy(cc MenuServiceClient) MenuServiceServer {
	return &menuServiceClientProxy{cc}
}

func (c *menuServiceClientProxy) ListMenu(ctx context.Context, in *ListMenuRequest) (*ListMenuReply, error) {
	return c.cc.ListMenu(ctx, in)
}
func (c *menuServiceClientProxy) GetMenu(ctx context.Context, in *GetMenuRequest) (*Menu, error) {
	return c.cc.GetMenu(ctx, in)
}
func (c *menuServiceClientProxy) CreateMenu(ctx context.Context, in *CreateMenuRequest) (*Menu, error) {
	return c.cc.CreateMenu(ctx, in)
}
func (c *menuServiceClientProxy) UpdateMenu(ctx context.Context, in *UpdateMenuRequest) (*Menu, error) {
	return c.cc.UpdateMenu(ctx, in)
}
func (c *menuServiceClientProxy) DeleteMenu(ctx context.Context, in *DeleteMenuRequest) (*DeleteMenuReply, error) {
	return c.cc.DeleteMenu(ctx, in)
}
func (c *menuServiceClientProxy) GetAvailableMenus(ctx context.Context, in *GetAvailableMenusRequest) (*GetAvailableMenusReply, error) {
	return c.cc.GetAvailableMenus(ctx, in)
}
