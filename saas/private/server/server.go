package server

import (
	dtmdata "github.com/go-saas/kit/dtm/data"
	kapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authz/authz"
	ksaas "github.com/go-saas/kit/pkg/saas"
	uow2 "github.com/go-saas/kit/pkg/uow"
	"github.com/go-saas/kit/saas/api"
	"github.com/go-saas/kit/saas/private/biz"
	"github.com/go-saas/kit/saas/private/data"
	"github.com/go-saas/saas"
	"github.com/google/wire"

	"github.com/go-saas/saas/seed"
	"github.com/go-saas/uow"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer, NewJobServer, NewEventServer, NewSeeder, wire.Value(ClientName), wire.Value(biz.ConnName), NewSeeding, NewAuthorizationOption)

var ClientName kapi.ClientName = api.ServiceName

// Seeding workaround for https://github.com/google/wire/issues/207
type Seeding seed.Contrib

func NewSeeding(uow uow.Manager, migrate *data.Migrate) Seeding {
	return uow2.NewUowContrib(uow, seed.Chain(migrate))
}

func NewSeeder(ts saas.TenantStore, dtmMigrator *dtmdata.Migrator, ss Seeding) seed.Seeder {
	return seed.NewDefaultSeeder(ksaas.NewTraceContrib(ksaas.SeedChangeTenant(ts, dtmMigrator, ss)))
}

func NewAuthorizationOption() *authz.Option {
	return authz.NewAuthorizationOption()
}
