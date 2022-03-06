package server

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	v1 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	"github.com/goxiaoy/go-saas/common"
	shttp "github.com/goxiaoy/go-saas/common/http"
	"github.com/goxiaoy/go-saas/kratos/saas"
)

func Saas(hmtOpt *shttp.WebMultiTenancyOption, ts common.TenantStore) middleware.Middleware {
	return selector.Server(saas.Server(hmtOpt, ts, ErrorFormatter())).Match(func(ctx context.Context, operation string) bool {
		_, ok := common.FromCurrentTenant(ctx)
		return !ok
	}).Build()
}

func ErrorFormatter() func(err error) (interface{}, error) {
	return func(err error) (interface{}, error) {
		if errors.Is(err, common.ErrTenantNotFound) {
			return nil, v1.ErrorTenantNotFound("")
		} else {
			return nil, err
		}
	}
}
