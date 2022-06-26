//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	kapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authn/jwt"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/authz/casbin"
	kconf "github.com/go-saas/kit/pkg/conf"
	kdal "github.com/go-saas/kit/pkg/dal"
	"github.com/go-saas/kit/pkg/job"
	kserver "github.com/go-saas/kit/pkg/server"
	sapi "github.com/go-saas/kit/saas/api"
	"github.com/go-saas/kit/user/private/biz"
	"github.com/go-saas/kit/user/private/conf"
	"github.com/go-saas/kit/user/private/data"
	"github.com/go-saas/kit/user/private/server"
	"github.com/go-saas/kit/user/private/service"
	"github.com/go-saas/saas/http"
)

// initApp init kratos application.
func initApp(*kconf.Services, *kconf.Security, *conf.UserConf, *kconf.Data, log.Logger, *http.WebMultiTenancyOption, ...grpc.ClientOption) (*kratos.App, func(), error) {
	panic(wire.Build(authz.ProviderSet, kserver.DefaultProviderSet, jwt.ProviderSet, kapi.DefaultProviderSet, kdal.DefaultProviderSet, job.DefaultProviderSet,
		sapi.GrpcProviderSet,
		casbin.PermissionProviderSet, server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
