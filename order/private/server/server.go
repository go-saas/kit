package server

import (
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas/seed"
	api2 "order/api"
	"order/private/biz"
	"order/private/data"
	"github.com/goxiaoy/uow"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer, NewSeeder, wire.Value(ClientName))

var ClientName api.ClientName = api2.ServiceName

func NewSeeder(uow uow.Manager, migrate *data.Migrate, post *biz.PostSeeder) seed.Seeder {
	return seed.NewDefaultSeeder(seed.NewUowContributor(uow, seed.Chain(migrate, post)))
}
