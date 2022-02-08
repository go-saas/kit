package server

import (
	"context"
	"net/http"
)

type writerContextKey struct {
}

func Writer() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, writerContextKey{}, w)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func FromWriterContext(ctx context.Context) (http.ResponseWriter, bool) {
	w, ok := ctx.Value(writerContextKey{}).(http.ResponseWriter)
	return w, ok
}
