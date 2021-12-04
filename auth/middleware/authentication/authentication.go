package authentication

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/goxiaoy/go-saas-kit/auth/current"
	"github.com/goxiaoy/go-saas-kit/auth/jwt"
	"github.com/goxiaoy/go-saas-kit/auth/middleware/extract_claim"
)

func ServerExtractAndAuth(l log.Logger, tokenizer jwt.Tokenizer) middleware.Middleware {
	return middleware.Chain(extract_claim.ServerExtract(tokenizer), ServerAuth(l))
}

func ServerAuth(l log.Logger) middleware.Middleware {
	logger := log.NewHelper(l)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if claims, ok := jwt.FromClaimsContext(ctx); !ok {
				//no jwt
				return handler(ctx, req)
			} else {
				//extract user and set user id context
				uid := ""
				if claims.Subject != "" {
					uid = claims.Subject
				} else {
					uid = claims.Uid
				}
				logger.Debugf("Current User: %v", uid)
				uc := current.NewUserContext(ctx, current.NewUserInfo(uid))
				// set client id context
				clientId := claims.ClientId
				if clientId != "" {
					uc = current.NewClientContext(uc, clientId)
				}
				return handler(uc, req)
			}
		}
	}
}
