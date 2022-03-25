package api

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"strings"
)

const defaultPrefix Prefix = "internal."

type Option struct {
	HeaderPrefix Prefix
	Propagators  []Propagator
	BypassToken  bool
}

type Prefix string

func PrefixOrDefault(prefix Prefix) Prefix {
	if prefix == "" {
		return defaultPrefix
	}
	return prefix
}

func NewOption(prefix Prefix, bypassToken bool, propagators ...Propagator) *Option {
	prefix = PrefixOrDefault(prefix)
	return &Option{HeaderPrefix: prefix, BypassToken: bypassToken, Propagators: propagators}
}

func NewDefaultOption(saas *SaasPropagator, logger log.Logger) *Option {
	return NewOption("", false, saas, NewUserContributor(logger), NewClientContributor(true, logger))
}

type Header interface {
	Get(key string) string
	Set(key, value string)
	HasKey(key string) bool
}

type HeaderCarrier map[string]string

func (h HeaderCarrier) Get(key string) string {
	if r, ok := h[key]; ok {
		return r
	} else {
		return ""
	}
}

func (h HeaderCarrier) Set(key, value string) {
	h[key] = value
}

func (h HeaderCarrier) HasKey(key string) bool {
	_, ok := h[key]
	return ok
}

//Propagator propagates cross-cutting concerns as key-value text
//pairs within a carrier that travels in-band across process boundaries to keep same state across services
type Propagator interface {
	Extract(ctx context.Context, carrier Header) (context.Context, error)
	Inject(ctx context.Context, carrier Header) error
	Fields() []string
}

func ServerPropagation(opt *Option, logger log.Logger) middleware.Middleware {
	l := log.NewHelper(log.With(logger, "module", "api.ServerPropagation"))
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				//find client claims
				if claims, ok := jwt.FromClaimsContext(ctx); ok {
					//TODO trusted server to server communication
					if claims.ClientId != "" {
						l.Debugf("find trusted context with client: %s", claims.ClientId)
						//preserve all request header
						//recover context
						newCtx := ctx
						var err error
						cleanHeaders := HeaderCarrier(map[string]string{})
						for _, key := range tr.RequestHeader().Keys() {
							key = strings.ToLower(key)
							headerPrefix := strings.ToLower(string(opt.HeaderPrefix))
							if strings.HasPrefix(key, headerPrefix) {
								k := strings.TrimPrefix(key, headerPrefix)
								v := tr.RequestHeader().Get(key)
								l.Debugf("set clean header key: %s,v: %s", k, v)
								cleanHeaders.Set(k, v)
							}
						}

						for i := range opt.Propagators {
							newCtx, err = opt.Propagators[i].Extract(newCtx, cleanHeaders)
							if err != nil {
								return nil, err
							}
						}
						return handler(newCtx, req)
					}
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
					h := fmt.Sprintf("%s%s", opt.HeaderPrefix, k)
					l.Debugf("set header: %s,value: %s", h, v)
					tr.RequestHeader().Set(h, v)
				}
			}
			return handler(ctx, req)
		}
	}
}
