package server

import (
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/auth/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	api2 "github.com/goxiaoy/go-saas-kit/saas/api"
	"github.com/goxiaoy/go-saas-kit/saas/internal_/biz"
	"github.com/goxiaoy/go-saas-kit/saas/internal_/conf"
	"github.com/goxiaoy/go-saas-kit/saas/internal_/data"
	"github.com/goxiaoy/go-saas/seed"
	"github.com/goxiaoy/uow"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer, jwt.NewTokenizer, NewSeeder, wire.Value(ClientName))

var ClientName api.ClientName = api2.ServiceName

func NewSeeder(c *conf.Data, uow uow.Manager, migrate *data.Migrate, permission *biz.PermissionSeeder) seed.Seeder {
	var opt = seed.NewSeedOption(migrate, permission)
	// seed host
	opt.TenantIds = []string{""}

	return seed.NewDefaultSeeder(opt, uow, map[string]interface{}{})
}
