//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	jwt2 "github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/casbin"
	conf2 "github.com/goxiaoy/go-saas-kit/pkg/conf"
	server2 "github.com/goxiaoy/go-saas-kit/pkg/server"
	api2 "github.com/goxiaoy/go-saas-kit/saas/api"
	"github.com/goxiaoy/go-saas-kit/saas/remote"
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
func initApp(*conf2.Services, *conf2.Security, *conf.UserConf, *conf.Data, log.Logger, *biz.PasswordValidatorConfig, *uow.Config, *gorm.Config, *http.WebMultiTenancyOption, ...grpc.ClientOption) (*kratos.App, func(), error) {
	panic(wire.Build(authz.ProviderSet, casbin.PermissionProviderSet, server2.DefaultCodecProviderSet, jwt2.ProviderSet, api.DefaultProviderSet,
		api2.GrpcProviderSet, remote.GrpcProviderSet,
		server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
