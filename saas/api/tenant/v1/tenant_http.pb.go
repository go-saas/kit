// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// protoc-gen-go-http v2.1.2

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

type TenantServiceHTTPServer interface {
	CreateTenant(context.Context, *CreateTenantRequest) (*Tenant, error)
	DeleteTenant(context.Context, *DeleteTenantRequest) (*DeleteTenantReply, error)
	GetTenant(context.Context, *GetTenantRequest) (*Tenant, error)
	ListTenant(context.Context, *ListTenantRequest) (*ListTenantReply, error)
	UpdateTenant(context.Context, *UpdateTenantRequest) (*Tenant, error)
}

func RegisterTenantServiceHTTPServer(s *http.Server, srv TenantServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/v1/saas/tenant", _TenantService_CreateTenant0_HTTP_Handler(srv))
	r.PATCH("/v1/saas/tenant/{tenant.id}", _TenantService_UpdateTenant0_HTTP_Handler(srv))
	r.PUT("/v1/saas/tenant/{tenant.id}", _TenantService_UpdateTenant1_HTTP_Handler(srv))
	r.DELETE("/v1/saas/tenant/{id}", _TenantService_DeleteTenant0_HTTP_Handler(srv))
	r.GET("/v1/saas/tenant/{id_or_name}", _TenantService_GetTenant0_HTTP_Handler(srv))
	r.POST("/v1/saas/tenant/list", _TenantService_ListTenant0_HTTP_Handler(srv))
	r.GET("/v1/saas/tenants", _TenantService_ListTenant1_HTTP_Handler(srv))
}

func _TenantService_CreateTenant0_HTTP_Handler(srv TenantServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateTenantRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.tenant.v1.TenantService/CreateTenant")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateTenant(ctx, req.(*CreateTenantRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Tenant)
		return ctx.Result(200, reply)
	}
}

func _TenantService_UpdateTenant0_HTTP_Handler(srv TenantServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateTenantRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.tenant.v1.TenantService/UpdateTenant")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateTenant(ctx, req.(*UpdateTenantRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Tenant)
		return ctx.Result(200, reply)
	}
}

func _TenantService_UpdateTenant1_HTTP_Handler(srv TenantServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateTenantRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.tenant.v1.TenantService/UpdateTenant")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateTenant(ctx, req.(*UpdateTenantRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Tenant)
		return ctx.Result(200, reply)
	}
}

func _TenantService_DeleteTenant0_HTTP_Handler(srv TenantServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteTenantRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.tenant.v1.TenantService/DeleteTenant")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteTenant(ctx, req.(*DeleteTenantRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteTenantReply)
		return ctx.Result(200, reply)
	}
}

func _TenantService_GetTenant0_HTTP_Handler(srv TenantServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetTenantRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.tenant.v1.TenantService/GetTenant")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetTenant(ctx, req.(*GetTenantRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Tenant)
		return ctx.Result(200, reply)
	}
}

func _TenantService_ListTenant0_HTTP_Handler(srv TenantServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListTenantRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.tenant.v1.TenantService/ListTenant")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListTenant(ctx, req.(*ListTenantRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListTenantReply)
		return ctx.Result(200, reply)
	}
}

func _TenantService_ListTenant1_HTTP_Handler(srv TenantServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListTenantRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.tenant.v1.TenantService/ListTenant")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListTenant(ctx, req.(*ListTenantRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListTenantReply)
		return ctx.Result(200, reply)
	}
}

type TenantServiceHTTPClient interface {
	CreateTenant(ctx context.Context, req *CreateTenantRequest, opts ...http.CallOption) (rsp *Tenant, err error)
	DeleteTenant(ctx context.Context, req *DeleteTenantRequest, opts ...http.CallOption) (rsp *DeleteTenantReply, err error)
	GetTenant(ctx context.Context, req *GetTenantRequest, opts ...http.CallOption) (rsp *Tenant, err error)
	ListTenant(ctx context.Context, req *ListTenantRequest, opts ...http.CallOption) (rsp *ListTenantReply, err error)
	UpdateTenant(ctx context.Context, req *UpdateTenantRequest, opts ...http.CallOption) (rsp *Tenant, err error)
}

type TenantServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewTenantServiceHTTPClient(client *http.Client) TenantServiceHTTPClient {
	return &TenantServiceHTTPClientImpl{client}
}

func (c *TenantServiceHTTPClientImpl) CreateTenant(ctx context.Context, in *CreateTenantRequest, opts ...http.CallOption) (*Tenant, error) {
	var out Tenant
	pattern := "/v1/saas/tenant"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation("/api.tenant.v1.TenantService/CreateTenant"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *TenantServiceHTTPClientImpl) DeleteTenant(ctx context.Context, in *DeleteTenantRequest, opts ...http.CallOption) (*DeleteTenantReply, error) {
	var out DeleteTenantReply
	pattern := "/v1/saas/tenant/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.tenant.v1.TenantService/DeleteTenant"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *TenantServiceHTTPClientImpl) GetTenant(ctx context.Context, in *GetTenantRequest, opts ...http.CallOption) (*Tenant, error) {
	var out Tenant
	pattern := "/v1/saas/tenant/{id_or_name}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.tenant.v1.TenantService/GetTenant"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *TenantServiceHTTPClientImpl) ListTenant(ctx context.Context, in *ListTenantRequest, opts ...http.CallOption) (*ListTenantReply, error) {
	var out ListTenantReply
	pattern := "/v1/saas/tenants"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.tenant.v1.TenantService/ListTenant"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *TenantServiceHTTPClientImpl) UpdateTenant(ctx context.Context, in *UpdateTenantRequest, opts ...http.CallOption) (*Tenant, error) {
	var out Tenant
	pattern := "/v1/saas/tenant/{tenant.id}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation("/api.tenant.v1.TenantService/UpdateTenant"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
