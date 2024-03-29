// Code generated by protoc-gen-go-grpc-proxy. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc-proxy v1.2.0
// - protoc             (unknown)
// source: user/api/permission/v1/permission.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

var _ PermissionServiceServer = (*permissionServiceClientProxy)(nil)

// permissionServiceClientProxy is the proxy to turn PermissionService client to server interface.
type permissionServiceClientProxy struct {
	cc PermissionServiceClient
}

func NewPermissionServiceClientProxy(cc PermissionServiceClient) PermissionServiceServer {
	return &permissionServiceClientProxy{cc}
}

func (c *permissionServiceClientProxy) GetCurrent(ctx context.Context, in *GetCurrentPermissionRequest) (*GetCurrentPermissionReply, error) {
	return c.cc.GetCurrent(ctx, in)
}
func (c *permissionServiceClientProxy) CheckCurrent(ctx context.Context, in *CheckPermissionRequest) (*CheckPermissionReply, error) {
	return c.cc.CheckCurrent(ctx, in)
}
