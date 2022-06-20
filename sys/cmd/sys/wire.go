//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	kapi "github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	sconf "github.com/goxiaoy/go-saas-kit/pkg/conf"
	kdal "github.com/goxiaoy/go-saas-kit/pkg/dal"
	"github.com/goxiaoy/go-saas-kit/pkg/job"
	kserver "github.com/goxiaoy/go-saas-kit/pkg/server"
	"github.com/goxiaoy/go-saas-kit/sys/private/biz"
	"github.com/goxiaoy/go-saas-kit/sys/private/data"
	"github.com/goxiaoy/go-saas-kit/sys/private/server"
	"github.com/goxiaoy/go-saas-kit/sys/private/service"
	uapi "github.com/goxiaoy/go-saas-kit/user/api"
	shttp "github.com/goxiaoy/go-saas/common/http"

	sapi "github.com/goxiaoy/go-saas-kit/saas/api"
)

// initApp init kratos application.
func initApp(*sconf.Services, *sconf.Security, *shttp.WebMultiTenancyOption, *sconf.Data, log.Logger, ...grpc.ClientOption) (*kratos.App, func(), error) {
	panic(wire.Build(authz.ProviderSet, jwt.ProviderSet, kserver.DefaultCodecProviderSet, kapi.DefaultProviderSet, kdal.DefaultProviderSet, job.DefaultProviderSet,
		uapi.GrpcProviderSet,
		sapi.GrpcProviderSet,
		server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
