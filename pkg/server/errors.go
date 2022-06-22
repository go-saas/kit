package server

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kerrors "github.com/goxiaoy/go-saas-kit/pkg/errors"
)

// Recovery wrap kratos recovery with handler
func Recovery() middleware.Middleware {
	return recovery.Recovery()
}

func Stack() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			r, err := handler(ctx, req)
			if err == nil {
				return r, err
			}
			err = fmt.Errorf("%w\n,%s", err, kerrors.Stack(0))
			return r, err
		}
	}
}
