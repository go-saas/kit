//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	jwt2 "github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authorization"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/casbin"
	conf2 "github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	"github.com/goxiaoy/go-saas-kit/user/private/conf"
	"github.com/goxiaoy/go-saas-kit/user/private/data"
	"github.com/goxiaoy/go-saas-kit/user/private/server"
	"github.com/goxiaoy/go-saas-kit/user/private/service"
	"github.com/goxiaoy/go-saas/common/http"
	"github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/uow"
)

// initApp init kratos application.
func initApp(*conf2.Services, *conf2.Security, *conf.UserConf, *conf.Data, log.Logger, *biz.PasswordValidatorConfig, *uow.Config, *gorm.Config, *http.WebMultiTenancyOption) (*kratos.App, func(), error) {
	panic(wire.Build(authorization.ProviderSet, casbin.PermissionProviderSet, jwt2.ProviderSet, api.DefaultProviderSet, server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
