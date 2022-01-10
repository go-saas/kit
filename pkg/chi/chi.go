package chi

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"net/http"
)

func MiddlewareConvert(m ...middleware.Middleware) func(http.Handler) http.Handler {
	chain := middleware.Chain(m...)
	return func(handler http.Handler) http.Handler {
		var req interface{}
		handleFunc := func(w http.ResponseWriter, r *http.Request) {
			//impl: Bind(r, &req)
			next := func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, nil
			}
			a := chain(next)
			_, _ = a(r.Context(), &req)
			return
		}
		return http.HandlerFunc(handleFunc)
	}
}
