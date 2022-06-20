package server

import (
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	uow2 "github.com/goxiaoy/go-saas-kit/pkg/uow"
	api2 "github.com/goxiaoy/go-saas-kit/sys/api"
	"github.com/goxiaoy/go-saas-kit/sys/private/biz"
	"github.com/goxiaoy/go-saas-kit/sys/private/data"
	"github.com/goxiaoy/go-saas/seed"
	"github.com/goxiaoy/uow"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer, NewSeeder, wire.Value(ClientName), wire.Value(data.ConnName), NewSeeding, NewAuthorizationOption)

var ClientName api.ClientName = api2.ServiceName

type Seeding seed.Contributor

func NewSeeding(uow uow.Manager, migrate *data.Migrate, menu *biz.MenuSeed) Seeding {
	return uow2.NewUowContributor(uow, seed.Chain(migrate, menu))
}

func NewSeeder(ss Seeding) seed.Seeder {
	return seed.NewDefaultSeeder(ss)
}

func NewAuthorizationOption() *authz.Option {
	return authz.NewAuthorizationOption()
}
