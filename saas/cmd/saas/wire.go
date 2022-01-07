//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/authorization/authorization"
	api2 "github.com/goxiaoy/go-saas-kit/pkg/api"
	jwt2 "github.com/goxiaoy/go-saas-kit/pkg/auth/jwt"
	conf2 "github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/saas/private/biz"
	"github.com/goxiaoy/go-saas-kit/saas/private/conf"
	"github.com/goxiaoy/go-saas-kit/saas/private/data"
	"github.com/goxiaoy/go-saas-kit/saas/private/server"
	"github.com/goxiaoy/go-saas-kit/saas/private/service"
	"github.com/goxiaoy/go-saas-kit/user/api"
	"github.com/goxiaoy/go-saas/common/http"
	"github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/uow"
)

// initApp init kratos application.
func initApp(*conf2.Services, *conf2.Security, *conf.Data, log.Logger, *uow.Config, *gorm.Config, *http.WebMultiTenancyOption, ...grpc.ClientOption) (*kratos.App, func(), error) {
	panic(wire.Build(authorization.ProviderSet, jwt2.ProviderSet, server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, api.GrpcProviderSet, api2.DefaultProviderSet, newApp))
}
