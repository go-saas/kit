package server

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	v1 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	"github.com/goxiaoy/go-saas/common"
	shttp "github.com/goxiaoy/go-saas/common/http"
	"github.com/goxiaoy/go-saas/kratos/saas"
)

func Saas(hmtOpt *shttp.WebMultiTenancyOption, ts common.TenantStore, trustedContextValidator api.TrustedContextValidator, options ...common.ResolveOption) middleware.Middleware {
	return selector.Server(saas.Server(ts, saas.WithMultiTenancyOption(hmtOpt), saas.WithErrorFormatter(ErrorFormatter()), saas.WithResolveOption(options...))).Match(func(ctx context.Context, operation string) bool {
		ok, _ := trustedContextValidator.Trusted(ctx)
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

func NewWebMultiTenancyOption(opt *conf.AppConfig) *shttp.WebMultiTenancyOption {
	ret := shttp.NewDefaultWebMultiTenancyOption()
	if opt == nil {
		return ret
	}
	if opt.TenantKey != nil {
		ret.TenantKey = opt.TenantKey.Value
	}
	if opt.DomainFormat != nil {
		ret.DomainFormat = opt.DomainFormat.Value
	}
	return ret
}
