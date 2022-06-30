//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	dtmserver "github.com/go-saas/kit/dtm/server"
	"github.com/go-saas/kit/examples/monolithic/private/server"
	kapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authn/jwt"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/authz/casbin"
	kitconf "github.com/go-saas/kit/pkg/conf"
	kdal "github.com/go-saas/kit/pkg/dal"
	"github.com/go-saas/kit/pkg/job"
	kserver "github.com/go-saas/kit/pkg/server"
	"github.com/google/wire"

	sbiz "github.com/go-saas/kit/saas/private/biz"
	sconf "github.com/go-saas/kit/saas/private/conf"
	sdata "github.com/go-saas/kit/saas/private/data"
	sservice "github.com/go-saas/kit/saas/private/service"

	sysbiz "github.com/go-saas/kit/sys/private/biz"
	sysdata "github.com/go-saas/kit/sys/private/data"
	sysservice "github.com/go-saas/kit/sys/private/service"

	ubiz "github.com/go-saas/kit/user/private/biz"
	uconf "github.com/go-saas/kit/user/private/conf"
	udata "github.com/go-saas/kit/user/private/data"
	uservice "github.com/go-saas/kit/user/private/service"
)

// initApp init kratos application.
func initApp(*kitconf.Services, *kitconf.Security, *kitconf.Data, *sconf.SaasConf, *uconf.UserConf, log.Logger, *kitconf.AppConfig, ...grpc.ClientOption) (*kratos.App, func(), error) {
	panic(wire.Build(authz.ProviderSet, jwt.ProviderSet, kserver.DefaultProviderSet, kserver.NewWebMultiTenancyOption, kapi.DefaultProviderSet, kdal.DefaultProviderSet, job.DefaultProviderSet, dtmserver.DtmProviderSet,
		sdata.ProviderSet, sbiz.ProviderSet, sservice.ProviderSet,
		sysdata.ProviderSet, sysbiz.ProviderSet, sysservice.ProviderSet,
		udata.ProviderSet, ubiz.ProviderSet, uservice.ProviderSet,
		casbin.PermissionProviderSet, server.ProviderSet,
		newApp))
}
