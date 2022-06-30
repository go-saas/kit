package server

import (
	"github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authz/authz"
	ksaas "github.com/go-saas/kit/pkg/saas"
	uow2 "github.com/go-saas/kit/pkg/uow"
	api2 "github.com/go-saas/kit/sys/api"
	"github.com/go-saas/kit/sys/private/biz"
	"github.com/go-saas/kit/sys/private/data"
	"github.com/go-saas/saas"
	"github.com/google/wire"

	"github.com/go-saas/saas/seed"
	"github.com/go-saas/uow"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer, NewJobServer, NewSeeder, wire.Value(ClientName), wire.Value(biz.ConnName), NewSeeding, NewAuthorizationOption)

var ClientName api.ClientName = api2.ServiceName

// Seeding workaround for https://github.com/google/wire/issues/207
type Seeding seed.Contrib

func NewSeeding(uow uow.Manager, migrate *data.Migrate, menu *biz.MenuSeed) Seeding {
	return uow2.NewUowContrib(uow, seed.Chain(migrate, menu))
}

func NewSeeder(ts saas.TenantStore, ss Seeding) seed.Seeder {
	return seed.NewDefaultSeeder(ksaas.NewTraceContrib(ksaas.SeedChangeTenant(ts, ss)))
}

func NewAuthorizationOption() *authz.Option {
	return authz.NewAuthorizationOption()
}
