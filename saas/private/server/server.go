package server

import (
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas"
	kapi "github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/pkg/event/event"
	ksaas "github.com/goxiaoy/go-saas-kit/pkg/saas"
	uow2 "github.com/goxiaoy/go-saas-kit/pkg/uow"
	"github.com/goxiaoy/go-saas-kit/saas/api"
	"github.com/goxiaoy/go-saas-kit/saas/private/biz"
	"github.com/goxiaoy/go-saas-kit/saas/private/data"

	"github.com/goxiaoy/go-saas/seed"
	"github.com/goxiaoy/uow"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer, NewJobServer, NewSeeder, wire.Value(ClientName), wire.Value(biz.ConnName), NewSeeding, NewEventHandler, NewAuthorizationOption)

var ClientName kapi.ClientName = api.ServiceName

// Seeding workaround for https://github.com/google/wire/issues/207
type Seeding seed.Contrib

func NewSeeding(uow uow.Manager, migrate *data.Migrate) Seeding {
	return uow2.NewUowContrib(uow, seed.Chain(migrate))
}

func NewSeeder(ts saas.TenantStore, ss Seeding) seed.Seeder {
	return seed.NewDefaultSeeder(ksaas.SeedChangeTenant(ts, ss))
}

func NewAuthorizationOption() *authz.Option {
	return authz.NewAuthorizationOption()
}

func NewEventHandler(saas biz.SaasEventHandler) event.Handler {
	return event.Handler(saas)
}
