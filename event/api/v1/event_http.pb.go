// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.5.3
// - protoc             (unknown)
// source: event/api/v1/event.proto

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

const OperationEventServiceEvent = "/event.api.v1.EventService/Event"

type EventServiceHTTPServer interface {
	Event(context.Context, *EventRequest) (*emptypb.Empty, error)
}

func RegisterEventServiceHTTPServer(s *http.Server, srv EventServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/v1/{service}/proxy/event", _EventService_Event0_HTTP_Handler(srv))
}

func _EventService_Event0_HTTP_Handler(srv EventServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in EventRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationEventServiceEvent)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Event(ctx, req.(*EventRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

type EventServiceHTTPClient interface {
	Event(ctx context.Context, req *EventRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
}

type EventServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewEventServiceHTTPClient(client *http.Client) EventServiceHTTPClient {
	return &EventServiceHTTPClientImpl{client}
}

func (c *EventServiceHTTPClientImpl) Event(ctx context.Context, in *EventRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/v1/{service}/proxy/event"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationEventServiceEvent))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
