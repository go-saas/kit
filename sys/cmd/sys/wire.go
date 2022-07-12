//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	kapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authn/jwt"
	"github.com/go-saas/kit/pkg/authz/authz"
	sconf "github.com/go-saas/kit/pkg/conf"
	kdal "github.com/go-saas/kit/pkg/dal"
	"github.com/go-saas/kit/pkg/job"
	kserver "github.com/go-saas/kit/pkg/server"
	"github.com/go-saas/kit/sys/private/biz"
	"github.com/go-saas/kit/sys/private/conf"
	"github.com/go-saas/kit/sys/private/data"
	"github.com/go-saas/kit/sys/private/server"
	"github.com/go-saas/kit/sys/private/service"
	uapi "github.com/go-saas/kit/user/api"
	shttp "github.com/go-saas/saas/http"
	"github.com/google/wire"

	sapi "github.com/go-saas/kit/saas/api"
)

// initApp init kratos application.
func initApp(*sconf.Services, *sconf.Security, *conf.SysConf, *shttp.WebMultiTenancyOption, *sconf.Data, log.Logger, ...grpc.ClientOption) (*kratos.App, func(), error) {
	panic(wire.Build(authz.ProviderSet, jwt.ProviderSet, kserver.DefaultProviderSet, kapi.DefaultProviderSet, kdal.DefaultProviderSet, job.DefaultProviderSet,
		uapi.GrpcProviderSet,
		sapi.GrpcProviderSet,
		server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
