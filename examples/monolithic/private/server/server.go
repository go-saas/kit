package server

import (
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/dal"
	kitdi "github.com/go-saas/kit/pkg/di"
	ksaas "github.com/go-saas/kit/pkg/saas"
	sserver "github.com/go-saas/kit/saas/private/server"
	sysserver "github.com/go-saas/kit/sys/private/server"
	userver "github.com/go-saas/kit/user/private/server"
	uservice "github.com/go-saas/kit/user/private/service"
	"github.com/go-saas/saas"
	"github.com/go-saas/saas/seed"
	"github.com/goava/di"
)

// ProviderSet is server providers.
var ProviderSet = kitdi.NewSet(
	kitdi.NewProvider(NewHTTPServer, di.As(new(transport.Server))),
	kitdi.NewProvider(NewGRPCServer, di.As(new(transport.Server))),
	kitdi.NewProvider(NewJobServer, di.As(new(transport.Server))),
	kitdi.NewProvider(NewEventServer, di.As(new(transport.Server))),

	NewSeeder,
	kitdi.Value(ConnName),
	NewAuthorizationOption,
	userver.NewSeeding,
	sserver.NewSeeding,
	sysserver.NewSeeding,
)

const ConnName = dal.ConnName("default")

func NewSeeder(ts saas.TenantStore, seeds []seed.Contrib) seed.Seeder {
	res := seed.NewDefaultSeeder(ksaas.NewTraceContrib(ksaas.SeedChangeTenant(ts, seeds...)))
	return res
}

func NewAuthorizationOption(userRole *uservice.UserRoleContrib) *authz.Option {
	return authz.NewAuthorizationOption(userRole)
}
