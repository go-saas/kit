package server

import (
	dtmdata "github.com/go-saas/kit/dtm/data"
	dtmservice "github.com/go-saas/kit/dtm/service"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/dal"
	ksaas "github.com/go-saas/kit/pkg/saas"
	"github.com/go-saas/kit/pkg/server"
	sserver "github.com/go-saas/kit/saas/private/server"
	sservice "github.com/go-saas/kit/saas/private/service"
	sysserver "github.com/go-saas/kit/sys/private/server"
	sysservice "github.com/go-saas/kit/sys/private/service"
	userver "github.com/go-saas/kit/user/private/server"
	uservice "github.com/go-saas/kit/user/private/service"
	"github.com/go-saas/saas"
	"github.com/go-saas/saas/seed"
	"github.com/google/wire"
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

func NewSeeder(dtmMigrator *dtmdata.Migrator, ts saas.TenantStore, user userver.Seeding, sys sysserver.Seeding, saas sserver.Seeding) seed.Seeder {
	return seed.NewDefaultSeeder(ksaas.NewTraceContrib(ksaas.SeedChangeTenant(ts, dtmMigrator, user, sys, saas)))
}

func NewHttpServiceRegister(dtmRegister dtmservice.HttpServerRegister, user uservice.HttpServerRegister, sys sysservice.HttpServerRegister, saas sservice.HttpServerRegister) HttpServerRegister {
	return server.ChainHttpServiceRegister(dtmRegister, user, sys, saas)
}

func NewGrpcServiceRegister(dtmRegister dtmservice.GrpcServerRegister, user uservice.GrpcServerRegister, sys sysservice.GrpcServerRegister, saas sservice.GrpcServerRegister) GrpcServerRegister {
	return server.ChainGrpcServiceRegister(dtmRegister, user, sys, saas)
}

func NewAuthorizationOption(userRole *uservice.UserRoleContrib) *authz.Option {
	return authz.NewAuthorizationOption(userRole)
}
