package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kerrors "github.com/goxiaoy/go-saas-kit/pkg/errors"
	v1 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	"github.com/goxiaoy/go-saas/common"
)

// Recovery wrap kratos recovery with handler
func Recovery() middleware.Middleware {
	return recovery.Recovery(recovery.WithHandler(func(ctx context.Context, req, err interface{}) error {
		if rerr, ok := err.(error); ok {
			if errors.Is(rerr, common.ErrTenantNotFound) {
				return v1.ErrorTenantNotFound("")
			}
		}
		return recovery.ErrUnknownRequest
	}))
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
