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

type HandlerFunc[TRet any] func(w http.ResponseWriter, r *http.Request) (TRet, error)

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

//HttpResponseAndErrorEncoder must before any filter may use http.ResponseWriter
func HttpResponseAndErrorEncoder(resEncoder khttp.EncodeResponseFunc, errorHandler ErrorHandler) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return errorHandler.Wrap(func(w http.ResponseWriter, r *http.Request) error {
			handler.ServeHTTP(w, r)
			res := r.Context().Value(respKey{})
			err, _ := r.Context().Value(errorKey{}).(error)
			if err != nil {
				return err
			}
			if res != nil {
				//after encoder. w is frozen
				if err := resEncoder(w, r, res); err != nil {
					return err
				}
			}
			return nil
		})
	}
}

func HandlerWrap[TRet any](handler HandlerFunc[TRet]) http.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) {
		res, err := handler(w, r)
		//put into context
		*r = *r.WithContext(context.WithValue(r.Context(), respKey{}, res))
		*r = *r.WithContext(context.WithValue(r.Context(), errorKey{}, err))
	}
	return f
}
