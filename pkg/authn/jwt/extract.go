package jwt

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"strings"
)

func ServerExtract(tokenizer Tokenizer, logger log.Logger) middleware.Middleware {
	log := log.NewHelper(log.With(logger, "module", "auth.extract_claim"))
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			t := ""
			if info, ok := transport.FromServerContext(ctx); ok {
				auth := info.RequestHeader().Get(AuthorizationHeader)
				if auth != "" {
					t = ExtractHeaderToken(auth)
				}
				if t == "" {
					if ht, ok := info.(*http.Transport); ok {
						t = ht.Request().URL.Query().Get(AuthorizationQuery)
					}
				}
			}
			if t != "" {
				//
				if claims, err := ExtractAndValidate(tokenizer, t); err != nil {
					//errors
					log.Error(err)
					return handler(ctx, req)
				} else {
					return handler(NewClaimsContext(NewJWTContext(ctx, t), claims), req)
				}
			}
			return handler(ctx, req)
		}
	}
}

func ExtractAndValidate(tokenizer Tokenizer, t string) (*Claims, error) {
	if claims, err := tokenizer.Parse(t); err != nil {
		return nil, err
	} else {
		if err := claims.Valid(); err != nil {
			return nil, err
		}
		return claims, nil
	}
}

func ExtractHeaderToken(token string) string {
	splitToken := strings.Split(token, BearerTokenType)
	if len(splitToken) == 2 {
		return strings.TrimSpace(splitToken[1])
	} else {
		return ""
	}
}
