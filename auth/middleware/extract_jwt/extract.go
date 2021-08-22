package extract_jwt

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/goxiaoy/go-saas-kit/auth/jwt"
	"strings"
)

func ServerExtract() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			t := ""
			if info, ok := transport.FromServerContext(ctx); ok {
				auth := info.RequestHeader().Get("Authorization")
				if auth != "" {
					splitToken := strings.Split(auth, "Bearer")
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
				return handler(jwt.NewJWTContext(ctx, t), req)
			}
			return handler(ctx, req)
		}
	}
}
