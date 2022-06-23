//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	conf2 "github.com/goxiaoy/go-saas-kit/pkg/conf"
	sapi "github.com/goxiaoy/go-saas-kit/saas/api"
	uapi "github.com/goxiaoy/go-saas-kit/user/api"
	shttp "github.com/goxiaoy/go-saas/http"
)

func initApp(*conf2.Services, *conf2.Security, *shttp.WebMultiTenancyOption, api.ClientName, klog.Logger, ...grpc.ClientOption) (*App, func(), error) {
	panic(wire.Build(ProviderSet, authz.ProviderSet, sapi.GrpcProviderSet, uapi.GrpcProviderSet, jwt.ProviderSet, newApp))
}
