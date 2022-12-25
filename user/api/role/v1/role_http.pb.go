// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.5.3
// - protoc             (unknown)
// source: user/api/role/v1/role.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationRoleServiceCreateRole = "/user.api.role.v1.RoleService/CreateRole"
const OperationRoleServiceDeleteRole = "/user.api.role.v1.RoleService/DeleteRole"
const OperationRoleServiceGetRole = "/user.api.role.v1.RoleService/GetRole"
const OperationRoleServiceGetRolePermission = "/user.api.role.v1.RoleService/GetRolePermission"
const OperationRoleServiceListRoles = "/user.api.role.v1.RoleService/ListRoles"
const OperationRoleServiceUpdateRole = "/user.api.role.v1.RoleService/UpdateRole"
const OperationRoleServiceUpdateRolePermission = "/user.api.role.v1.RoleService/UpdateRolePermission"

type RoleServiceHTTPServer interface {
	CreateRole(context.Context, *CreateRoleRequest) (*Role, error)
	DeleteRole(context.Context, *DeleteRoleRequest) (*DeleteRoleResponse, error)
	GetRole(context.Context, *GetRoleRequest) (*Role, error)
	GetRolePermission(context.Context, *GetRolePermissionRequest) (*GetRolePermissionResponse, error)
	ListRoles(context.Context, *ListRolesRequest) (*ListRolesResponse, error)
	UpdateRole(context.Context, *UpdateRoleRequest) (*Role, error)
	UpdateRolePermission(context.Context, *UpdateRolePermissionRequest) (*UpdateRolePermissionResponse, error)
}

func RegisterRoleServiceHTTPServer(s *http.Server, srv RoleServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/v1/role/list", _RoleService_ListRoles0_HTTP_Handler(srv))
	r.GET("/v1/roles", _RoleService_ListRoles1_HTTP_Handler(srv))
	r.GET("/v1/role/{id}", _RoleService_GetRole0_HTTP_Handler(srv))
	r.POST("/v1/role", _RoleService_CreateRole0_HTTP_Handler(srv))
	r.PATCH("/v1/role/{role.id}", _RoleService_UpdateRole0_HTTP_Handler(srv))
	r.PUT("/v1/role/{role.id}", _RoleService_UpdateRole1_HTTP_Handler(srv))
	r.DELETE("/v1/role/{id}", _RoleService_DeleteRole0_HTTP_Handler(srv))
	r.GET("/v1/role/{id}/permission", _RoleService_GetRolePermission0_HTTP_Handler(srv))
	r.PUT("/v1/role/{id}/permission", _RoleService_UpdateRolePermission0_HTTP_Handler(srv))
}

func _RoleService_ListRoles0_HTTP_Handler(srv RoleServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListRolesRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoleServiceListRoles)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListRoles(ctx, req.(*ListRolesRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListRolesResponse)
		return ctx.Result(200, reply)
	}
}

func _RoleService_ListRoles1_HTTP_Handler(srv RoleServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListRolesRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoleServiceListRoles)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListRoles(ctx, req.(*ListRolesRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListRolesResponse)
		return ctx.Result(200, reply)
	}
}

func _RoleService_GetRole0_HTTP_Handler(srv RoleServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetRoleRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoleServiceGetRole)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetRole(ctx, req.(*GetRoleRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Role)
		return ctx.Result(200, reply)
	}
}

func _RoleService_CreateRole0_HTTP_Handler(srv RoleServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateRoleRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoleServiceCreateRole)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateRole(ctx, req.(*CreateRoleRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Role)
		return ctx.Result(200, reply)
	}
}

func _RoleService_UpdateRole0_HTTP_Handler(srv RoleServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateRoleRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoleServiceUpdateRole)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateRole(ctx, req.(*UpdateRoleRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Role)
		return ctx.Result(200, reply)
	}
}

func _RoleService_UpdateRole1_HTTP_Handler(srv RoleServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateRoleRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoleServiceUpdateRole)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateRole(ctx, req.(*UpdateRoleRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Role)
		return ctx.Result(200, reply)
	}
}

func _RoleService_DeleteRole0_HTTP_Handler(srv RoleServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteRoleRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoleServiceDeleteRole)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteRole(ctx, req.(*DeleteRoleRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteRoleResponse)
		return ctx.Result(200, reply)
	}
}

func _RoleService_GetRolePermission0_HTTP_Handler(srv RoleServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetRolePermissionRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoleServiceGetRolePermission)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetRolePermission(ctx, req.(*GetRolePermissionRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetRolePermissionResponse)
		return ctx.Result(200, reply)
	}
}

func _RoleService_UpdateRolePermission0_HTTP_Handler(srv RoleServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateRolePermissionRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoleServiceUpdateRolePermission)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateRolePermission(ctx, req.(*UpdateRolePermissionRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateRolePermissionResponse)
		return ctx.Result(200, reply)
	}
}

type RoleServiceHTTPClient interface {
	CreateRole(ctx context.Context, req *CreateRoleRequest, opts ...http.CallOption) (rsp *Role, err error)
	DeleteRole(ctx context.Context, req *DeleteRoleRequest, opts ...http.CallOption) (rsp *DeleteRoleResponse, err error)
	GetRole(ctx context.Context, req *GetRoleRequest, opts ...http.CallOption) (rsp *Role, err error)
	GetRolePermission(ctx context.Context, req *GetRolePermissionRequest, opts ...http.CallOption) (rsp *GetRolePermissionResponse, err error)
	ListRoles(ctx context.Context, req *ListRolesRequest, opts ...http.CallOption) (rsp *ListRolesResponse, err error)
	UpdateRole(ctx context.Context, req *UpdateRoleRequest, opts ...http.CallOption) (rsp *Role, err error)
	UpdateRolePermission(ctx context.Context, req *UpdateRolePermissionRequest, opts ...http.CallOption) (rsp *UpdateRolePermissionResponse, err error)
}

type RoleServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewRoleServiceHTTPClient(client *http.Client) RoleServiceHTTPClient {
	return &RoleServiceHTTPClientImpl{client}
}

func (c *RoleServiceHTTPClientImpl) CreateRole(ctx context.Context, in *CreateRoleRequest, opts ...http.CallOption) (*Role, error) {
	var out Role
	pattern := "/v1/role"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoleServiceCreateRole))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *RoleServiceHTTPClientImpl) DeleteRole(ctx context.Context, in *DeleteRoleRequest, opts ...http.CallOption) (*DeleteRoleResponse, error) {
	var out DeleteRoleResponse
	pattern := "/v1/role/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationRoleServiceDeleteRole))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *RoleServiceHTTPClientImpl) GetRole(ctx context.Context, in *GetRoleRequest, opts ...http.CallOption) (*Role, error) {
	var out Role
	pattern := "/v1/role/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationRoleServiceGetRole))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *RoleServiceHTTPClientImpl) GetRolePermission(ctx context.Context, in *GetRolePermissionRequest, opts ...http.CallOption) (*GetRolePermissionResponse, error) {
	var out GetRolePermissionResponse
	pattern := "/v1/role/{id}/permission"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationRoleServiceGetRolePermission))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *RoleServiceHTTPClientImpl) ListRoles(ctx context.Context, in *ListRolesRequest, opts ...http.CallOption) (*ListRolesResponse, error) {
	var out ListRolesResponse
	pattern := "/v1/roles"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationRoleServiceListRoles))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *RoleServiceHTTPClientImpl) UpdateRole(ctx context.Context, in *UpdateRoleRequest, opts ...http.CallOption) (*Role, error) {
	var out Role
	pattern := "/v1/role/{role.id}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoleServiceUpdateRole))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *RoleServiceHTTPClientImpl) UpdateRolePermission(ctx context.Context, in *UpdateRolePermissionRequest, opts ...http.CallOption) (*UpdateRolePermissionResponse, error) {
	var out UpdateRolePermissionResponse
	pattern := "/v1/role/{id}/permission"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoleServiceUpdateRolePermission))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
