package api

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	grpcmd "google.golang.org/grpc/metadata"
	"time"
)

var _ transport.Transporter = &Transport{}

// Transport is a gRPC transport.
type Transport struct {
	endpoint    string
	operation   string
	reqHeader   headerCarrier
	replyHeader headerCarrier
	filters     []selector.Filter
}

// Kind returns the transport kind.
func (tr *Transport) Kind() transport.Kind {
	return transport.KindGRPC
}

// Endpoint returns the transport endpoint.
func (tr *Transport) Endpoint() string {
	return tr.endpoint
}

// Operation returns the transport operation.
func (tr *Transport) Operation() string {
	return tr.operation
}

// RequestHeader returns the request header.
func (tr *Transport) RequestHeader() transport.Header {
	return tr.reqHeader
}

// ReplyHeader returns the reply header.
func (tr *Transport) ReplyHeader() transport.Header {
	return tr.replyHeader
}

// SelectFilters returns the client select filters.
func (tr *Transport) SelectFilters() []selector.Filter {
	return tr.filters
}

type headerCarrier metadata.MD

// Get returns the value associated with the passed key.
func (mc headerCarrier) Get(key string) string {
	vals := metadata.MD(mc).Get(key)
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// Set stores the key-value pair.
func (mc headerCarrier) Set(key string, value string) {
	metadata.MD(mc).Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (mc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(mc))
	for k := range metadata.MD(mc) {
		keys = append(keys, k)
	}
	return keys
}

// UnaryClientInterceptor
//
// https://github.com/go-kratos/kratos/blob/4ca25e46799ca1579d9e877e5be23f454ca23299/transport/grpc/client.go#L159-L192
func UnaryClientInterceptor(ms []middleware.Middleware, timeout time.Duration, filters []selector.Filter) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = transport.NewClientContext(ctx, &Transport{
			endpoint:  cc.Target(),
			operation: method,
			reqHeader: headerCarrier{},
			filters:   filters,
		})
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromClientContext(ctx); ok {
				header := tr.RequestHeader()
				keys := header.Keys()
				keyvals := make([]string, 0, len(keys))
				for _, k := range keys {
					keyvals = append(keyvals, k, header.Get(k))
				}
				ctx = grpcmd.AppendToOutgoingContext(ctx, keyvals...)
			}
			return reply, invoker(ctx, method, req, reply, cc, opts...)
		}
		if len(ms) > 0 {
			h = middleware.Chain(ms...)(h)
		}
		_, err := h(ctx, req)
		return err
	}
}
