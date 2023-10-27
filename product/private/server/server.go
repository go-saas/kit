package server

import (
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authz/authz"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/go-saas/kit/pkg/server"
	api2 "github.com/go-saas/kit/product/api"
	"github.com/go-saas/kit/product/private/biz"
	"github.com/go-saas/kit/product/private/data"
	"github.com/go-saas/saas"
	"github.com/go-saas/saas/seed"
	"github.com/go-saas/uow"
	"github.com/goava/di"
)

// ProviderSet is server providers.
var ProviderSet = kitdi.NewSet(
	NewAuthorizationOption,
	kitdi.NewProvider(NewHTTPServer, di.As(new(transport.Server))),
	kitdi.NewProvider(NewGRPCServer, di.As(new(transport.Server))),
	kitdi.NewProvider(NewJobServer, di.As(new(transport.Server))),
	NewSeeder,
	NewSeeding,
	kitdi.Value(ClientName),
	kitdi.Value(biz.ConnName))

var ClientName api.ClientName = api2.ServiceName

func NewSeeding(uow uow.Manager, migrate *data.Migrate, post *biz.ProductSeeder) seed.Contrib {
	return seed.Chain(server.NewUowContrib(uow, seed.Chain(migrate, post)))
}

func NewSeeder(ts saas.TenantStore, seeds []seed.Contrib) seed.Seeder {
	res := seed.NewDefaultSeeder(server.NewTraceContrib(server.SeedChangeTenant(ts, seeds...)))
	return res
}

func NewAuthorizationOption() *authz.Option {
	return authz.NewAuthorizationOption()
}
