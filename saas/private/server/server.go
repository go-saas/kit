package server

import (
	"github.com/go-kratos/kratos/v2/transport"
	kapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authz/authz"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/go-saas/kit/pkg/server"
	"github.com/go-saas/kit/saas/api"
	"github.com/go-saas/kit/saas/private/biz"
	"github.com/go-saas/kit/saas/private/data"
	"github.com/go-saas/saas"
	"github.com/go-saas/saas/seed"
	"github.com/go-saas/uow"
	"github.com/goava/di"
)

// ProviderSet is server providers.
var ProviderSet = kitdi.NewSet(
	kitdi.NewProvider(NewHTTPServer, di.As(new(transport.Server))),
	kitdi.NewProvider(NewGRPCServer, di.As(new(transport.Server))),
	kitdi.NewProvider(NewJobServer, di.As(new(transport.Server))),
	kitdi.NewProvider(NewEventServer, di.As(new(transport.Server))),
	kitdi.Value(ClientName),
	kitdi.Value(biz.ConnName),
	NewSeeder,
	NewSeeding,
	NewAuthorizationOption,
)

var ClientName kapi.ClientName = api.ServiceName

func NewSeeding(uow uow.Manager, migrate *data.Migrate) seed.Contrib {
	return server.NewUowContrib(uow, seed.Chain(migrate))
}

func NewSeeder(ts saas.TenantStore, seeds []seed.Contrib) seed.Seeder {
	res := seed.NewDefaultSeeder(server.NewTraceContrib(server.SeedChangeTenant(ts, seeds...)))
	return res
}

func NewAuthorizationOption() *authz.Option {
	return authz.NewAuthorizationOption()
}
