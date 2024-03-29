// Code generated by protoc-gen-go-grpc-proxy. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc-proxy v1.2.0
// - protoc             (unknown)
// source: user/api/permission/v1/permission_internal.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

var _ PermissionInternalServiceServer = (*permissionInternalServiceClientProxy)(nil)

// permissionInternalServiceClientProxy is the proxy to turn PermissionInternalService client to server interface.
type permissionInternalServiceClientProxy struct {
	cc PermissionInternalServiceClient
}

func NewPermissionInternalServiceClientProxy(cc PermissionInternalServiceClient) PermissionInternalServiceServer {
	return &permissionInternalServiceClientProxy{cc}
}

func (c *permissionInternalServiceClientProxy) CheckForSubjects(ctx context.Context, in *CheckSubjectsPermissionRequest) (*CheckSubjectsPermissionReply, error) {
	return c.cc.CheckForSubjects(ctx, in)
}
func (c *permissionInternalServiceClientProxy) AddSubjectPermission(ctx context.Context, in *AddSubjectPermissionRequest) (*AddSubjectPermissionResponse, error) {
	return c.cc.AddSubjectPermission(ctx, in)
}
func (c *permissionInternalServiceClientProxy) ListSubjectPermission(ctx context.Context, in *ListSubjectPermissionRequest) (*ListSubjectPermissionResponse, error) {
	return c.cc.ListSubjectPermission(ctx, in)
}
func (c *permissionInternalServiceClientProxy) UpdateSubjectPermission(ctx context.Context, in *UpdateSubjectPermissionRequest) (*UpdateSubjectPermissionResponse, error) {
	return c.cc.UpdateSubjectPermission(ctx, in)
}
func (c *permissionInternalServiceClientProxy) RemoveSubjectPermission(ctx context.Context, in *RemoveSubjectPermissionRequest) (*RemoveSubjectPermissionReply, error) {
	return c.cc.RemoveSubjectPermission(ctx, in)
}
