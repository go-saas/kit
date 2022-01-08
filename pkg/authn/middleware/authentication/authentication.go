package authentication

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/goxiaoy/go-saas-kit/pkg/authn"
	jwt2 "github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/middleware/extract_claim"
)

func ServerExtractAndAuth(tokenizer jwt2.Tokenizer, logger log.Logger) middleware.Middleware {
	return middleware.Chain(extract_claim.ServerExtract(tokenizer, logger), ServerAuth())
}

func ServerAuth() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if claims, ok := jwt2.FromClaimsContext(ctx); !ok {
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
				uc := authn.NewUserContext(ctx, authn.NewUserInfo(uid))
				// set client id context
				clientId := claims.ClientId
				if clientId != "" {
					uc = authn.NewClientContext(uc, clientId)
				}
				return handler(uc, req)
			}
		}
	}
}
