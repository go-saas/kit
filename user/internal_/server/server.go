package server

import (
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/auth/jwt"
	"github.com/goxiaoy/go-saas-kit/user/internal_/biz"
	"github.com/goxiaoy/go-saas-kit/user/internal_/conf"
	"github.com/goxiaoy/go-saas-kit/user/internal_/data"
	seed2 "github.com/goxiaoy/go-saas-kit/user/internal_/seed"
	"github.com/goxiaoy/go-saas/seed"
	"github.com/goxiaoy/uow"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(jwt.NewTokenizer, NewHTTPServer, NewGRPCServer, seed2.NewFake, NewSeeder)

func NewSeeder(c *conf.Data, uow uow.Manager, migrate *data.Migrate, roleSeed *biz.RoleSeed, userSeed *biz.UserSeed, fake *seed2.Fake) seed.Seeder {
	var opt = seed.NewSeedOption(migrate, roleSeed, userSeed, fake)
	// seed host
	opt.TenantIds = []string{""}

	return seed.NewDefaultSeeder(opt, uow, map[string]interface{}{
		biz.AdminUsernameKey: c.Admin.GetUsername(),
		biz.AdminPasswordKey: c.Admin.GetPassword(),
		seed2.FakeSeedKey: true,
	})
}
