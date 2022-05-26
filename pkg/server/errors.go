package server

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	v1 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	"github.com/goxiaoy/go-saas/common"
)

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
