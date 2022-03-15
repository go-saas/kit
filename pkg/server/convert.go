package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"net/http"
)

type Handler[TRet any] interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request) (TRet, error)
}

type HandlerFunc[TRet any] struct {
	f func(w http.ResponseWriter, r *http.Request) (TRet, error)
}

func NewHandlerFunc[TRet any](f func(w http.ResponseWriter, r *http.Request) (TRet, error)) *HandlerFunc[TRet] {
	return &HandlerFunc[TRet]{
		f: f,
	}
}

func (h *HandlerFunc[TRet]) ServeHTTP(w http.ResponseWriter, r *http.Request) (TRet, error) {
	return h.f(w, r)
}

type (
	respKey  struct{}
	errorKey struct{}
)

func MiddlewareConvert(m ...middleware.Middleware) func(handler http.Handler) http.Handler {
	chain := middleware.Chain(m...)
	return func(handler http.Handler) http.Handler {
		var req interface{}
		handleFunc := func(w http.ResponseWriter, r *http.Request) {
			next := func(ctx context.Context, req interface{}) (interface{}, error) {
				handler.ServeHTTP(w, r)
				res := r.Context().Value(respKey{})
				err, _ := r.Context().Value(errorKey{}).(error)
				return res, err
			}
			a := chain(next)
			a(r.Context(), &req)
		}
		return http.HandlerFunc(handleFunc)
	}
}

func HandlerWrap[TRet any](resEncoder khttp.EncodeResponseFunc, errorHandler ErrorHandler, handler *HandlerFunc[TRet]) http.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) error {
		res, err := handler.ServeHTTP(w, r)
		//put into context
		*r = *r.WithContext(context.WithValue(r.Context(), respKey{}, res))
		*r = *r.WithContext(context.WithValue(r.Context(), errorKey{}, err))
		if err != nil {
			return err
		}
		if err := resEncoder(w, r, res); err != nil {
			return err
		}
		return nil
	}
	return errorHandler.Wrap(f).ServeHTTP
}
