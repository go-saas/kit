package server

import (
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/pkg/dal"
	ksaas "github.com/goxiaoy/go-saas-kit/pkg/saas"
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	sserver "github.com/goxiaoy/go-saas-kit/saas/private/server"
	sservice "github.com/goxiaoy/go-saas-kit/saas/private/service"
	sysserver "github.com/goxiaoy/go-saas-kit/sys/private/server"
	sysservice "github.com/goxiaoy/go-saas-kit/sys/private/service"
	userver "github.com/goxiaoy/go-saas-kit/user/private/server"
	uservice "github.com/goxiaoy/go-saas-kit/user/private/service"
	"github.com/goxiaoy/go-saas/seed"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
	NewHTTPServer,
	NewGRPCServer,
	NewJobServer,
	NewEventServer,
	NewHttpServiceRegister,
	NewGrpcServiceRegister,
	NewSeeder,
	wire.Value(dal.ConnName("default")),
	NewAuthorizationOption,
	userver.NewSeeding,
	sserver.NewSeeding,
	sysserver.NewSeeding,
)

type HttpServerRegister server.HttpServiceRegister
type GrpcServerRegister server.GrpcServiceRegister

func NewSeeder(ts saas.TenantStore, user userver.Seeding, sys sysserver.Seeding, saas sserver.Seeding) seed.Seeder {
	return seed.NewDefaultSeeder(ksaas.SeedChangeTenant(ts, user, sys, saas))
}

func NewHttpServiceRegister(user uservice.HttpServerRegister, sys sysservice.HttpServerRegister, saas sservice.HttpServerRegister) HttpServerRegister {
	return server.ChainHttpServiceRegister(user, sys, saas)
}

func NewGrpcServiceRegister(user uservice.GrpcServerRegister, sys sysservice.GrpcServerRegister, saas sservice.GrpcServerRegister) GrpcServerRegister {
	return server.ChainGrpcServiceRegister(user, sys, saas)
}

func NewAuthorizationOption(userRole *uservice.UserRoleContrib) *authz.Option {
	return authz.NewAuthorizationOption(userRole)
}