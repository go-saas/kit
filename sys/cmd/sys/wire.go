//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	api2 "github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	sconf "github.com/goxiaoy/go-saas-kit/pkg/conf"
	server2 "github.com/goxiaoy/go-saas-kit/pkg/server"
	sapi "github.com/goxiaoy/go-saas-kit/saas/api"
	sremote "github.com/goxiaoy/go-saas-kit/saas/remote"
	"github.com/goxiaoy/go-saas-kit/sys/private/biz"
	"github.com/goxiaoy/go-saas-kit/sys/private/conf"
	"github.com/goxiaoy/go-saas-kit/sys/private/data"
	"github.com/goxiaoy/go-saas-kit/sys/private/server"
	"github.com/goxiaoy/go-saas-kit/sys/private/service"
	uapi "github.com/goxiaoy/go-saas-kit/user/api"
	uremote "github.com/goxiaoy/go-saas-kit/user/remote"
	shttp "github.com/goxiaoy/go-saas/common/http"
	"github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/uow"
)

// initApp init kratos application.
func initApp(*sconf.Services, *sconf.Security, *uow.Config, *gorm.Config, *shttp.WebMultiTenancyOption, *conf.Data, log.Logger, ...grpc.ClientOption) (*kratos.App, func(), error) {
	panic(wire.Build(authz.ProviderSet, jwt.ProviderSet, server2.DefaultCodecProviderSet, api2.DefaultProviderSet,
		uapi.GrpcProviderSet, uremote.GrpcProviderSet,
		sapi.GrpcProviderSet, sremote.GrpcProviderSet,
		server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
