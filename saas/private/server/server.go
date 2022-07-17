package server

import (
	kapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/dal"
	kitdi "github.com/go-saas/kit/pkg/di"
	ksaas "github.com/go-saas/kit/pkg/saas"
	uow2 "github.com/go-saas/kit/pkg/uow"
	"github.com/go-saas/kit/saas/api"
	"github.com/go-saas/kit/saas/private/biz"
	"github.com/go-saas/kit/saas/private/data"
	"github.com/go-saas/saas"
	"github.com/go-saas/saas/seed"
	"github.com/go-saas/uow"
)

// ProviderSet is server providers.
var ProviderSet = kitdi.NewSet(
	NewHTTPServer,
	NewGRPCServer,
	NewJobServer,
	NewEventServer,
	NewSeeder,
	func() kapi.ClientName { return ClientName },
	func() dal.ConnName { return biz.ConnName },
	NewSeeding,
	NewAuthorizationOption,
)

var ClientName kapi.ClientName = api.ServiceName

// Seeding workaround for https://github.com/google/wire/issues/207
type Seeding seed.Contrib

func NewSeeding(uow uow.Manager, migrate *data.Migrate) Seeding {
	return uow2.NewUowContrib(uow, seed.Chain(migrate))
}

func NewSeeder(ts saas.TenantStore, ss Seeding) seed.Seeder {
	return seed.NewDefaultSeeder(ksaas.NewTraceContrib(ksaas.SeedChangeTenant(ts, ss)))
}

func NewAuthorizationOption() *authz.Option {
	return authz.NewAuthorizationOption()
}
