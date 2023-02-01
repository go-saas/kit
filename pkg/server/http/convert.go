package http

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"net/http"
)

type Handler[TRet any] interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request) (TRet, error)
}

type HandlerFunc[TRet any] func(w http.ResponseWriter, r *http.Request) (TRet, error)

type (
	respKey  struct{}
	errorKey struct{}
)

// MiddlewareConvert convert kratos middleware into standard http middleware
func MiddlewareConvert(errEncoder khttp.EncodeErrorFunc, m ...middleware.Middleware) func(handler http.Handler) http.Handler {
	chain := middleware.Chain(m...)
	return func(handler http.Handler) http.Handler {
		var req interface{}
		handleFunc := func(w http.ResponseWriter, r *http.Request) {
			next := func(ctx context.Context, req interface{}) (interface{}, error) {
				//replace context and request
				*r = *r.WithContext(ctx)
				handler.ServeHTTP(w, r)
				//compatible with other middlewares
				res := r.Context().Value(respKey{})
				err, _ := r.Context().Value(errorKey{}).(error)
				return res, err
			}
			a := chain(next)
			_, err := a(r.Context(), &req)
			if err != nil {
				//encode error
				errEncoder(w, r, err)
			}
		}
		return http.HandlerFunc(handleFunc)
	}
}

func HandlerWrap[TRet any](resEncoder khttp.EncodeResponseFunc, handler HandlerFunc[TRet]) http.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) {

		res, err := handler(w, r)
		//put into context
		*r = *r.WithContext(context.WithValue(r.Context(), respKey{}, res))

		*r = *r.WithContext(context.WithValue(r.Context(), errorKey{}, err))

		//after encoder. w is frozen
		if err == nil {
			if err := resEncoder(w, r, res); err != nil {
				*r = *r.WithContext(context.WithValue(r.Context(), errorKey{}, err))
			}
		}
	}
	return f
}
