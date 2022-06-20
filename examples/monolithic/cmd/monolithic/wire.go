//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/examples/monolithic/private/server"
	kapi "github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/casbin"
	kitconf "github.com/goxiaoy/go-saas-kit/pkg/conf"
	kdal "github.com/goxiaoy/go-saas-kit/pkg/dal"
	"github.com/goxiaoy/go-saas-kit/pkg/job"
	kserver "github.com/goxiaoy/go-saas-kit/pkg/server"

	sbiz "github.com/goxiaoy/go-saas-kit/saas/private/biz"
	sconf "github.com/goxiaoy/go-saas-kit/saas/private/conf"
	sdata "github.com/goxiaoy/go-saas-kit/saas/private/data"
	sservice "github.com/goxiaoy/go-saas-kit/saas/private/service"

	sysbiz "github.com/goxiaoy/go-saas-kit/sys/private/biz"
	sysdata "github.com/goxiaoy/go-saas-kit/sys/private/data"
	sysservice "github.com/goxiaoy/go-saas-kit/sys/private/service"

	ubiz "github.com/goxiaoy/go-saas-kit/user/private/biz"
	uconf "github.com/goxiaoy/go-saas-kit/user/private/conf"
	udata "github.com/goxiaoy/go-saas-kit/user/private/data"
	uservice "github.com/goxiaoy/go-saas-kit/user/private/service"
)

// initApp init kratos application.
func initApp(*kitconf.Services, *kitconf.Security, *kitconf.Data, *sconf.SaasConf, *uconf.UserConf, log.Logger, *kitconf.AppConfig, ...grpc.ClientOption) (*kratos.App, func(), error) {
	panic(wire.Build(authz.ProviderSet, jwt.ProviderSet, kserver.DefaultCodecProviderSet, kserver.NewWebMultiTenancyOption, kapi.DefaultProviderSet, kdal.DefaultProviderSet, job.DefaultProviderSet,
		sdata.ProviderSet, sbiz.ProviderSet, sservice.ProviderSet,
		sysdata.ProviderSet, sysbiz.ProviderSet, sysservice.ProviderSet,
		udata.ProviderSet, ubiz.ProviderSet, uservice.ProviderSet,
		casbin.PermissionProviderSet, server.ProviderSet,
		newApp))
}
