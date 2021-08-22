package server

import (
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas/seed"
	"github.com/goxiaoy/go-saas-kit/auth/jwt"
	"github.com/goxiaoy/go-saas-kit/user/internal_/biz"
	"github.com/goxiaoy/go-saas-kit/user/internal_/conf"
	"github.com/goxiaoy/go-saas-kit/user/internal_/data"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(jwt.NewTokenizer,NewHTTPServer, NewGRPCServer,NewSeeder)

func NewSeeder(c *conf.Data,migrate *data.Migrate,roleSeed *biz.RoleSeed,userSeed *biz.UserSeed) seed.Seeder  {
	var opt =seed.NewSeedOption(migrate,roleSeed,userSeed)
	// seed host
	opt.TenantIds = []string{""}

	return seed.NewDefaultSeeder(opt, map[string]interface{}{
		biz.AdminUserNameKey:c.Admin.GetUsername(),
		biz.AdminPasswordKey: c.Admin.GetPassword(),
	})
}