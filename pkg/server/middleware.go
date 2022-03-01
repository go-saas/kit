package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/goxiaoy/go-saas/common"
	shttp "github.com/goxiaoy/go-saas/common/http"
	"github.com/goxiaoy/go-saas/kratos/saas"
)

func Saas(hmtOpt *shttp.WebMultiTenancyOption, ts common.TenantStore) middleware.Middleware {
	return selector.Server(saas.Server(hmtOpt, ts)).Match(func(ctx context.Context, operation string) bool {
		_, ok := common.FromCurrentTenant(ctx)
		return !ok
	}).Build()
}
