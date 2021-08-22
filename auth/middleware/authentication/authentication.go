package authentication

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/goxiaoy/go-saas-kit/auth/current"
	"github.com/goxiaoy/go-saas-kit/auth/jwt"
	"github.com/goxiaoy/go-saas-kit/auth/middleware/extract_jwt"
)

func ServerExtractAndAuth(l log.Logger, tokenizer jwt.Tokenizer) middleware.Middleware {
	return middleware.Chain(extract_jwt.ServerExtract(), ServerAuth(l, tokenizer))
}

func ServerAuth(l log.Logger, tokenizer jwt.Tokenizer) middleware.Middleware {
	logger := log.NewHelper(l)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if jwt, ok := jwt.FromJWTContext(ctx); !ok {
				//no jwt
				return handler(ctx, req)
			} else {
				//validate
				if claims, err := tokenizer.Parse(jwt); err != nil {
					return handler(ctx, req)
				} else {
					if err := claims.Valid(); err != nil {
						return handler(ctx, req)
					}
					//extract user
					uid := ""
					if claims.Subject != "" {
						uid = claims.Subject
					} else {
						uid = claims.Uid
					}
					logger.Debugf("Current User: %v", uid)
					uc := current.NewUserContext(ctx, current.UserInfo{Id: uid})
					return handler(uc, req)
				}
			}
		}
	}
}
