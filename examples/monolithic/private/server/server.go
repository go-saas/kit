package server

import (
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/pkg/dal"
	"github.com/goxiaoy/go-saas-kit/pkg/event/event"
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	sbiz "github.com/goxiaoy/go-saas-kit/saas/private/biz"
	sserver "github.com/goxiaoy/go-saas-kit/saas/private/server"
	sservice "github.com/goxiaoy/go-saas-kit/saas/private/service"
	sysserver "github.com/goxiaoy/go-saas-kit/sys/private/server"
	sysservice "github.com/goxiaoy/go-saas-kit/sys/private/service"
	ubiz "github.com/goxiaoy/go-saas-kit/user/private/biz"
	userver "github.com/goxiaoy/go-saas-kit/user/private/server"
	uservice "github.com/goxiaoy/go-saas-kit/user/private/service"
	"github.com/goxiaoy/go-saas/seed"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
	NewHTTPServer,
	NewGRPCServer,
	NewJobServer,
	NewHttpServiceRegister,
	NewGrpcServiceRegister,
	NewSeeder,
	NewEventHandler,
	wire.Value(dal.ConnName("default")),
	NewAuthorizationOption,
	userver.NewHydra,
	userver.NewSeeding,
	sserver.NewSeeding,
	sysserver.NewSeeding,
)

type HttpServerRegister server.HttpServiceRegister
type GrpcServerRegister server.GrpcServiceRegister

func NewSeeder(user userver.Seeding, sys sysserver.Seeding, saas sserver.Seeding) seed.Seeder {
	return seed.NewDefaultSeeder(user, sys, saas)
}

func NewHttpServiceRegister(user uservice.HttpServerRegister, sys sysservice.HttpServerRegister, saas sservice.HttpServerRegister) HttpServerRegister {
	return server.ChainHttpServiceRegister(user, sys, saas)
}

func NewGrpcServiceRegister(user uservice.GrpcServerRegister, sys sysservice.GrpcServerRegister, saas sservice.GrpcServerRegister) GrpcServerRegister {
	return server.ChainGrpcServiceRegister(user, sys, saas)
}

func NewEventHandler(user ubiz.UserEventHandler, saas sbiz.SaasEventHandler) event.Handler {
	return event.ChainHandler(event.Handler(user), event.Handler(saas))
}

func NewAuthorizationOption(userRole *uservice.UserRoleContributor) *authz.Option {
	return authz.NewAuthorizationOption(userRole)
}
