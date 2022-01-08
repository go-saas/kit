package extract_claim

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwt2 "github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"strings"
)

func ServerExtract(tokenizer jwt2.Tokenizer, logger log.Logger) middleware.Middleware {
	log := log.NewHelper(log.With(logger, "module", "auth.extract_claim"))
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			t := ""
			if info, ok := transport.FromServerContext(ctx); ok {
				auth := info.RequestHeader().Get(jwt2.AuthorizationHeader)
				if auth != "" {
					splitToken := strings.Split(auth, jwt2.BearerTokenType)
					if len(splitToken) == 2 {
						t = strings.TrimSpace(splitToken[1])
					}
				}
				if t == "" {
					if ht, ok := info.(*http.Transport); ok {
						t = ht.Request().URL.Query().Get("access_token")
					}
				}
			}
			if t != "" {
				//
				if claims, err := tokenizer.Parse(t); err != nil {
					//errors
					log.Error(err)
					return handler(ctx, req)
				} else {
					if err := claims.Valid(); err != nil {
						log.Error(err)
						return handler(ctx, req)
					}
					return handler(jwt2.NewClaimsContext(jwt2.NewJWTContext(ctx, t), claims), req)
				}
			}
			return handler(ctx, req)
		}
	}
}
