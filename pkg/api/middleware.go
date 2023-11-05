package api

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-saas/kit/pkg/authn/jwt"
	"github.com/go-saas/kit/pkg/conf"
	"strings"
)

type Option struct {
	Propagators []Propagator
	BypassToken bool
	Insecure    bool
	Subset      *int
}

func NewOption(bypassToken bool, propagators ...Propagator) *Option {
	return &Option{BypassToken: bypassToken, Propagators: propagators}
}

func (o *Option) WithInsecure() *Option {
	o.Insecure = true
	return o
}

// WithSubset with client disocvery subset size.
// zero value means subset filter disabled
func (o *Option) WithSubset(size int) *Option {
	o.Subset = &size
	return o
}

func NewDefaultOption(logger log.Logger) *Option {
	return NewOption(
		false,
		NewSaasPropagator(logger),
		NewUserPropagator(logger),
		NewClientPropagator(true, logger),
	).WithInsecure()
}

type Header interface {
	Get(key string) string
	Set(key, value string)
	HasKey(key string) bool
	Keys() []string
}

type HeaderCarrier map[string]string

func (h HeaderCarrier) Get(key string) string {
	key = strings.ToLower(key)
	if r, ok := h[key]; ok {
		return r
	} else {
		return ""
	}
}

func (h HeaderCarrier) Set(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func (h HeaderCarrier) HasKey(key string) bool {
	key = strings.ToLower(key)
	_, ok := h[key]
	return ok
}

func (h HeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	return keys
}

// Propagator propagates cross-cutting concerns as key-value text
// pairs within a carrier that travels in-band across process boundaries to keep same state across services
type Propagator interface {
	Extract(ctx context.Context, carrier Header) (context.Context, error)
	Inject(ctx context.Context, carrier Header) error
	Fields() []string
}

func ServerPropagation(opt *Option, validator TrustedContextValidator, logger log.Logger) middleware.Middleware {
	l := log.NewHelper(log.With(logger, "module", "api.ServerPropagation"))
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				if ok, err := validator.Trusted(ctx); err != nil {
					return nil, err
				} else if ok {
					l.Debugf("find trusted context")
					//preserve all request header
					//recover context
					newCtx := ctx
					var err error
					headers := HeaderCarrier(map[string]string{})
					for _, key := range tr.RequestHeader().Keys() {
						headers.Set(key, tr.RequestHeader().Get(key))
					}
					for i := range opt.Propagators {
						newCtx, err = opt.Propagators[i].Extract(newCtx, headers)
						if err != nil {
							return nil, err
						}
					}
					return handler(newCtx, req)
				}

			}
			return handler(ctx, req)
		}
	}
}

func ClientPropagation(client *conf.Client, opt *Option, tokenMgr TokenManager, logger log.Logger) middleware.Middleware {
	l := log.NewHelper(log.With(logger, "module", "api.ClientPropagation"))
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromClientContext(ctx); ok {
				if opt.BypassToken {
					l.Debugf("bypass token")
					if rawToken, ok := jwt.FromJWTContext(ctx); ok {
						//bypass raw token
						tr.RequestHeader().Set(jwt.AuthorizationHeader, fmt.Sprintf("%s %s", jwt.BearerTokenType, rawToken))
					}
				} else if client != nil && client.ClientId != "" {
					//use token mgr
					token, err := tokenMgr.GetOrGenerateToken(ctx, client)
					l.Debugf("replace with client %s token", client.ClientId)
					if err != nil {
						return nil, err
					}
					tr.RequestHeader().Set(jwt.AuthorizationHeader, fmt.Sprintf("%s %s", jwt.BearerTokenType, token))
				}
				headers := HeaderCarrier(map[string]string{})
				//contributor create header
				for _, contributor := range opt.Propagators {
					err := contributor.Inject(ctx, headers)
					if err != nil {
						return nil, err
					}
				}
				for k, v := range headers {
					l.Debugf("set header: %s,value: %s", k, v)
					tr.RequestHeader().Set(k, v)
				}
			}
			return handler(ctx, req)
		}
	}
}
