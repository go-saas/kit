package server

import (
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	api2 "github.com/goxiaoy/go-saas-kit/sys/api"
	"github.com/goxiaoy/go-saas-kit/sys/private/biz"
	"github.com/goxiaoy/go-saas-kit/sys/private/data"
	"github.com/goxiaoy/go-saas/seed"
	"github.com/goxiaoy/uow"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer, NewSeeder, wire.Value(ClientName))

var ClientName api.ClientName = api2.ServiceName

func NewSeeder(uow uow.Manager, migrate *data.Migrate, menu *biz.MenuSeed) seed.Seeder {
	return seed.NewDefaultSeeder(seed.NewUowContributor(uow, seed.Chain(migrate, menu)))
}
