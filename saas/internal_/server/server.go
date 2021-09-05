package server

import (
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/auth/jwt"
	"github.com/goxiaoy/go-saas-kit/saas/internal_/conf"
	"github.com/goxiaoy/go-saas-kit/saas/internal_/data"
	"github.com/goxiaoy/go-saas/seed"
	"github.com/goxiaoy/uow"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer, jwt.NewTokenizer, NewSeeder)

func NewSeeder(c *conf.Data, uow uow.Manager, migrate *data.Migrate) seed.Seeder {
	var opt = seed.NewSeedOption(migrate)
	// seed host
	opt.TenantIds = []string{""}

	return seed.NewDefaultSeeder(opt, uow, map[string]interface{}{})
}
