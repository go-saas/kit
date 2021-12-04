package api

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goxiaoy/go-saas-kit/auth/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"strings"
)

type Option struct {
	HeaderPrefix string
	Contributor  []Contributor
	BypassToken  bool
}

func NewOption(prefix string, bypassToken bool, contributor ...Contributor) *Option {
	if prefix == "" {
		prefix = "internal."
	}
	return &Option{HeaderPrefix: prefix, BypassToken: bypassToken, Contributor: contributor}
}

func NewDefaultOption(saas *SaasContributor, user *UserContributor) *Option {
	return NewOption("", false, saas, user)
}

type Header interface {
	Get(key string) string
	Set(key string, value string)
	Keys() []string
	HasKey(key string) bool
}

type headerCarrier map[string]string

func (h headerCarrier) Get(key string) string {
	return h[key]
}

func (h headerCarrier) Set(key string, value string) {
	h[key] = value
}

func (h headerCarrier) Keys() []string {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	return keys
}

func (h headerCarrier) HasKey(key string) bool {
	_, ok := h[key]
	return ok
}

type Contributor interface {
	RecoverContext(ctx context.Context, headers Header) (context.Context, error)
	CreateHeader(ctx context.Context) map[string]string
}

func ServerMiddleware(opt *Option) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				//find client claims
				if claims, ok := jwt.FromClaimsContext(ctx); ok {
					//TODO trusted server to server communication
					if claims.ClientId != "" {
						//preserve all request header
						//recover context
						newCtx := ctx
						var err error
						cleanHeaders := headerCarrier(map[string]string{})
						for _, key := range tr.RequestHeader().Keys() {
							if strings.HasPrefix(key, opt.HeaderPrefix) {
								cleanHeaders.Set(strings.TrimPrefix(key, opt.HeaderPrefix), tr.RequestHeader().Get(key))
							}
						}
						for i := range opt.Contributor {
							newCtx, err = opt.Contributor[i].RecoverContext(newCtx, cleanHeaders)
							if err != nil {
								return nil, err
							}
						}
						return handler(newCtx, req)
					}
				}
				//clean internal headers
				for _, key := range tr.RequestHeader().Keys() {
					if strings.HasPrefix(key, opt.HeaderPrefix) {
						tr.RequestHeader().Set(key, "")
					}
				}
			}
			return handler(ctx, req)
		}
	}
}

func ClientMiddleware(client *conf.Client, opt *Option, tokenMgr TokenManager) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromClientContext(ctx); ok {
				if opt.BypassToken {
					if rawToken, ok := jwt.FromJWTContext(ctx); ok {
						//bypass raw token
						tr.RequestHeader().Set(jwt.AuthorizationHeader, fmt.Sprintf("%s %s", jwt.BearerTokenType, rawToken))
					}
				} else if client != nil && client.ClientId != "" {
					//use token mgr
					token, err := tokenMgr.GetOrGenerateToken(ctx, client)
					if err != nil {
						return nil, err
					}
					tr.RequestHeader().Set(jwt.AuthorizationHeader, fmt.Sprintf("%s %s", jwt.BearerTokenType, token))
				}
				//contributor create header
				for _, contributor := range opt.Contributor {
					headers := contributor.CreateHeader(ctx)
					if headers != nil {
						for k, v := range headers {
							tr.RequestHeader().Set(fmt.Sprintf("%s%s", opt.HeaderPrefix, k), v)
						}
					}
				}
			}
			return handler(ctx, req)
		}
	}
}
