// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.0
// - protoc             (unknown)
// source: dtm/api/dtm/v1/dtm.proto

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

const OperationMsgServiceQueryPrepared = "/dtm.api.dtm.v1.MsgService/QueryPrepared"

type MsgServiceHTTPServer interface {
	QueryPrepared(context.Context, *QueryPreparedRequest) (*emptypb.Empty, error)
}

func RegisterMsgServiceHTTPServer(s *http.Server, srv MsgServiceHTTPServer) {
	r := s.Route("/")
	r.GET("/v1/{service}/dtm/query-prepared", _MsgService_QueryPrepared0_HTTP_Handler(srv))
}

func _MsgService_QueryPrepared0_HTTP_Handler(srv MsgServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in QueryPreparedRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationMsgServiceQueryPrepared)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.QueryPrepared(ctx, req.(*QueryPreparedRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

type MsgServiceHTTPClient interface {
	QueryPrepared(ctx context.Context, req *QueryPreparedRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
}

type MsgServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewMsgServiceHTTPClient(client *http.Client) MsgServiceHTTPClient {
	return &MsgServiceHTTPClientImpl{client}
}

func (c *MsgServiceHTTPClientImpl) QueryPrepared(ctx context.Context, in *QueryPreparedRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/v1/{service}/dtm/query-prepared"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationMsgServiceQueryPrepared))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
