// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.0
// - protoc             (unknown)
// source: product/api/category/v1/category.proto

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

const OperationProductCategoryServiceCreateCategory = "/product.api.category.v1.ProductCategoryService/CreateCategory"
const OperationProductCategoryServiceDeleteCategory = "/product.api.category.v1.ProductCategoryService/DeleteCategory"
const OperationProductCategoryServiceGetCategory = "/product.api.category.v1.ProductCategoryService/GetCategory"
const OperationProductCategoryServiceListCategory = "/product.api.category.v1.ProductCategoryService/ListCategory"
const OperationProductCategoryServiceUpdateCategory = "/product.api.category.v1.ProductCategoryService/UpdateCategory"

type ProductCategoryServiceHTTPServer interface {
	CreateCategory(context.Context, *CreateCategoryRequest) (*Category, error)
	DeleteCategory(context.Context, *DeleteCategoryRequest) (*DeleteCategoryReply, error)
	GetCategory(context.Context, *GetCategoryRequest) (*Category, error)
	ListCategory(context.Context, *ListCategoryRequest) (*ListCategoryReply, error)
	UpdateCategory(context.Context, *UpdateCategoryRequest) (*Category, error)
}

func RegisterProductCategoryServiceHTTPServer(s *http.Server, srv ProductCategoryServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/v1/product/category/list", _ProductCategoryService_ListCategory0_HTTP_Handler(srv))
	r.GET("/v1/product/category", _ProductCategoryService_ListCategory1_HTTP_Handler(srv))
	r.GET("/v1/product/category/{key}", _ProductCategoryService_GetCategory0_HTTP_Handler(srv))
	r.POST("/v1/product/category", _ProductCategoryService_CreateCategory0_HTTP_Handler(srv))
	r.PATCH("/v1/product/category/{category.key}", _ProductCategoryService_UpdateCategory0_HTTP_Handler(srv))
	r.PUT("/v1/product/category/{category.key}", _ProductCategoryService_UpdateCategory1_HTTP_Handler(srv))
	r.DELETE("/v1/product/category/{key}", _ProductCategoryService_DeleteCategory0_HTTP_Handler(srv))
}

func _ProductCategoryService_ListCategory0_HTTP_Handler(srv ProductCategoryServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListCategoryRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProductCategoryServiceListCategory)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListCategory(ctx, req.(*ListCategoryRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListCategoryReply)
		return ctx.Result(200, reply)
	}
}

func _ProductCategoryService_ListCategory1_HTTP_Handler(srv ProductCategoryServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListCategoryRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProductCategoryServiceListCategory)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListCategory(ctx, req.(*ListCategoryRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListCategoryReply)
		return ctx.Result(200, reply)
	}
}

func _ProductCategoryService_GetCategory0_HTTP_Handler(srv ProductCategoryServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetCategoryRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProductCategoryServiceGetCategory)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetCategory(ctx, req.(*GetCategoryRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Category)
		return ctx.Result(200, reply)
	}
}

func _ProductCategoryService_CreateCategory0_HTTP_Handler(srv ProductCategoryServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateCategoryRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProductCategoryServiceCreateCategory)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateCategory(ctx, req.(*CreateCategoryRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Category)
		return ctx.Result(200, reply)
	}
}

func _ProductCategoryService_UpdateCategory0_HTTP_Handler(srv ProductCategoryServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateCategoryRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProductCategoryServiceUpdateCategory)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateCategory(ctx, req.(*UpdateCategoryRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Category)
		return ctx.Result(200, reply)
	}
}

func _ProductCategoryService_UpdateCategory1_HTTP_Handler(srv ProductCategoryServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateCategoryRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProductCategoryServiceUpdateCategory)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateCategory(ctx, req.(*UpdateCategoryRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Category)
		return ctx.Result(200, reply)
	}
}

func _ProductCategoryService_DeleteCategory0_HTTP_Handler(srv ProductCategoryServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteCategoryRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProductCategoryServiceDeleteCategory)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteCategory(ctx, req.(*DeleteCategoryRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteCategoryReply)
		return ctx.Result(200, reply)
	}
}

type ProductCategoryServiceHTTPClient interface {
	CreateCategory(ctx context.Context, req *CreateCategoryRequest, opts ...http.CallOption) (rsp *Category, err error)
	DeleteCategory(ctx context.Context, req *DeleteCategoryRequest, opts ...http.CallOption) (rsp *DeleteCategoryReply, err error)
	GetCategory(ctx context.Context, req *GetCategoryRequest, opts ...http.CallOption) (rsp *Category, err error)
	ListCategory(ctx context.Context, req *ListCategoryRequest, opts ...http.CallOption) (rsp *ListCategoryReply, err error)
	UpdateCategory(ctx context.Context, req *UpdateCategoryRequest, opts ...http.CallOption) (rsp *Category, err error)
}

type ProductCategoryServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewProductCategoryServiceHTTPClient(client *http.Client) ProductCategoryServiceHTTPClient {
	return &ProductCategoryServiceHTTPClientImpl{client}
}

func (c *ProductCategoryServiceHTTPClientImpl) CreateCategory(ctx context.Context, in *CreateCategoryRequest, opts ...http.CallOption) (*Category, error) {
	var out Category
	pattern := "/v1/product/category"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationProductCategoryServiceCreateCategory))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *ProductCategoryServiceHTTPClientImpl) DeleteCategory(ctx context.Context, in *DeleteCategoryRequest, opts ...http.CallOption) (*DeleteCategoryReply, error) {
	var out DeleteCategoryReply
	pattern := "/v1/product/category/{key}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationProductCategoryServiceDeleteCategory))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *ProductCategoryServiceHTTPClientImpl) GetCategory(ctx context.Context, in *GetCategoryRequest, opts ...http.CallOption) (*Category, error) {
	var out Category
	pattern := "/v1/product/category/{key}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationProductCategoryServiceGetCategory))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *ProductCategoryServiceHTTPClientImpl) ListCategory(ctx context.Context, in *ListCategoryRequest, opts ...http.CallOption) (*ListCategoryReply, error) {
	var out ListCategoryReply
	pattern := "/v1/product/category"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationProductCategoryServiceListCategory))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *ProductCategoryServiceHTTPClientImpl) UpdateCategory(ctx context.Context, in *UpdateCategoryRequest, opts ...http.CallOption) (*Category, error) {
	var out Category
	pattern := "/v1/product/category/{category.key}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationProductCategoryServiceUpdateCategory))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
