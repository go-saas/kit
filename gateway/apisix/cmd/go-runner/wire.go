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
	conf2 "github.com/goxiaoy/go-saas-kit/pkg/conf"
	sapi "github.com/goxiaoy/go-saas-kit/saas/api"
	sremote "github.com/goxiaoy/go-saas-kit/saas/remote"
)

func initApp(*conf2.Services, *conf2.Security, api.ClientName, klog.Logger, ...grpc.ClientOption) (*App, func(), error) {
	panic(wire.Build(ProviderSet, sapi.GrpcProviderSet, sremote.GrpcProviderSet, jwt.ProviderSet, newApp))
}
