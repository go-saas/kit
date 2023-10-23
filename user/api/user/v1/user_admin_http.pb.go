// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.0
// - protoc             (unknown)
// source: user/api/user/v1/user_admin.proto

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

const OperationUserAdminServiceCreateUserAdmin = "/user.api.user.v1.UserAdminService/CreateUserAdmin"
const OperationUserAdminServiceDeleteUserAdmin = "/user.api.user.v1.UserAdminService/DeleteUserAdmin"
const OperationUserAdminServiceGetUserAdmin = "/user.api.user.v1.UserAdminService/GetUserAdmin"
const OperationUserAdminServiceListUsersAdmin = "/user.api.user.v1.UserAdminService/ListUsersAdmin"
const OperationUserAdminServiceUpdateUserAdmin = "/user.api.user.v1.UserAdminService/UpdateUserAdmin"

type UserAdminServiceHTTPServer interface {
	// CreateUserAdmin CreateUser
	// authz: user.admin.user,*,create
	CreateUserAdmin(context.Context, *AdminCreateUserRequest) (*User, error)
	// DeleteUserAdminDeleteUser
	// authz: user.admin.user,id,delete
	DeleteUserAdmin(context.Context, *AdminDeleteUserRequest) (*AdminDeleteUserResponse, error)
	// GetUserAdminGetUser
	// authz: user.admin.user,id,get
	GetUserAdmin(context.Context, *AdminGetUserRequest) (*User, error)
	// ListUsersAdminListUsers
	// authz: user.admin.user,*,list
	ListUsersAdmin(context.Context, *AdminListUsersRequest) (*AdminListUsersResponse, error)
	// UpdateUserAdminUpdateUser
	// authz: user.admin.user,id,update
	UpdateUserAdmin(context.Context, *AdminUpdateUserRequest) (*User, error)
}

func RegisterUserAdminServiceHTTPServer(s *http.Server, srv UserAdminServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/v1/user/admin/user/list", _UserAdminService_ListUsersAdmin0_HTTP_Handler(srv))
	r.GET("/v1/user/admin/users", _UserAdminService_ListUsersAdmin1_HTTP_Handler(srv))
	r.GET("/v1/user/admin/user/{id}", _UserAdminService_GetUserAdmin0_HTTP_Handler(srv))
	r.POST("/v1/user/admin/user", _UserAdminService_CreateUserAdmin0_HTTP_Handler(srv))
	r.PATCH("/v1/user/admin/user/{user.id}", _UserAdminService_UpdateUserAdmin0_HTTP_Handler(srv))
	r.PUT("/v1/user/admin/user/{user.id}", _UserAdminService_UpdateUserAdmin1_HTTP_Handler(srv))
	r.DELETE("/v1/user/admin/user/{id}", _UserAdminService_DeleteUserAdmin0_HTTP_Handler(srv))
}

func _UserAdminService_ListUsersAdmin0_HTTP_Handler(srv UserAdminServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in AdminListUsersRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserAdminServiceListUsersAdmin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListUsersAdmin(ctx, req.(*AdminListUsersRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*AdminListUsersResponse)
		return ctx.Result(200, reply)
	}
}

func _UserAdminService_ListUsersAdmin1_HTTP_Handler(srv UserAdminServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in AdminListUsersRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserAdminServiceListUsersAdmin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListUsersAdmin(ctx, req.(*AdminListUsersRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*AdminListUsersResponse)
		return ctx.Result(200, reply)
	}
}

func _UserAdminService_GetUserAdmin0_HTTP_Handler(srv UserAdminServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in AdminGetUserRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserAdminServiceGetUserAdmin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetUserAdmin(ctx, req.(*AdminGetUserRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*User)
		return ctx.Result(200, reply)
	}
}

func _UserAdminService_CreateUserAdmin0_HTTP_Handler(srv UserAdminServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in AdminCreateUserRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserAdminServiceCreateUserAdmin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateUserAdmin(ctx, req.(*AdminCreateUserRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*User)
		return ctx.Result(200, reply)
	}
}

func _UserAdminService_UpdateUserAdmin0_HTTP_Handler(srv UserAdminServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in AdminUpdateUserRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserAdminServiceUpdateUserAdmin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateUserAdmin(ctx, req.(*AdminUpdateUserRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*User)
		return ctx.Result(200, reply)
	}
}

func _UserAdminService_UpdateUserAdmin1_HTTP_Handler(srv UserAdminServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in AdminUpdateUserRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserAdminServiceUpdateUserAdmin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateUserAdmin(ctx, req.(*AdminUpdateUserRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*User)
		return ctx.Result(200, reply)
	}
}

func _UserAdminService_DeleteUserAdmin0_HTTP_Handler(srv UserAdminServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in AdminDeleteUserRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserAdminServiceDeleteUserAdmin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteUserAdmin(ctx, req.(*AdminDeleteUserRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*AdminDeleteUserResponse)
		return ctx.Result(200, reply)
	}
}

type UserAdminServiceHTTPClient interface {
	CreateUserAdmin(ctx context.Context, req *AdminCreateUserRequest, opts ...http.CallOption) (rsp *User, err error)
	DeleteUserAdmin(ctx context.Context, req *AdminDeleteUserRequest, opts ...http.CallOption) (rsp *AdminDeleteUserResponse, err error)
	GetUserAdmin(ctx context.Context, req *AdminGetUserRequest, opts ...http.CallOption) (rsp *User, err error)
	ListUsersAdmin(ctx context.Context, req *AdminListUsersRequest, opts ...http.CallOption) (rsp *AdminListUsersResponse, err error)
	UpdateUserAdmin(ctx context.Context, req *AdminUpdateUserRequest, opts ...http.CallOption) (rsp *User, err error)
}

type UserAdminServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewUserAdminServiceHTTPClient(client *http.Client) UserAdminServiceHTTPClient {
	return &UserAdminServiceHTTPClientImpl{client}
}

func (c *UserAdminServiceHTTPClientImpl) CreateUserAdmin(ctx context.Context, in *AdminCreateUserRequest, opts ...http.CallOption) (*User, error) {
	var out User
	pattern := "/v1/user/admin/user"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserAdminServiceCreateUserAdmin))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *UserAdminServiceHTTPClientImpl) DeleteUserAdmin(ctx context.Context, in *AdminDeleteUserRequest, opts ...http.CallOption) (*AdminDeleteUserResponse, error) {
	var out AdminDeleteUserResponse
	pattern := "/v1/user/admin/user/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationUserAdminServiceDeleteUserAdmin))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *UserAdminServiceHTTPClientImpl) GetUserAdmin(ctx context.Context, in *AdminGetUserRequest, opts ...http.CallOption) (*User, error) {
	var out User
	pattern := "/v1/user/admin/user/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationUserAdminServiceGetUserAdmin))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *UserAdminServiceHTTPClientImpl) ListUsersAdmin(ctx context.Context, in *AdminListUsersRequest, opts ...http.CallOption) (*AdminListUsersResponse, error) {
	var out AdminListUsersResponse
	pattern := "/v1/user/admin/users"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationUserAdminServiceListUsersAdmin))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *UserAdminServiceHTTPClientImpl) UpdateUserAdmin(ctx context.Context, in *AdminUpdateUserRequest, opts ...http.CallOption) (*User, error) {
	var out User
	pattern := "/v1/user/admin/user/{user.id}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserAdminServiceUpdateUserAdmin))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
