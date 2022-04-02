package server

import (
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	conf2 "github.com/goxiaoy/go-saas-kit/pkg/conf"
	api2 "github.com/goxiaoy/go-saas-kit/user/api"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	"github.com/goxiaoy/go-saas-kit/user/private/conf"
	"github.com/goxiaoy/go-saas-kit/user/private/data"
	"github.com/goxiaoy/go-saas-kit/user/private/server/http"
	"github.com/goxiaoy/go-saas/seed"
	"github.com/goxiaoy/uow"
	client "github.com/ory/hydra-client-go"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer, NewRefreshTokenProvider, wire.Value(ClientName), NewSeeder, NewHydra, http.NewAuth)

var ClientName api.ClientName = api2.ServiceName

func NewSeeder(c *conf.UserConf,
	uow uow.Manager,
	migrate *data.Migrate,
	roleSeed *biz.RoleSeed,
	userSeed *biz.UserSeed,
	p *biz.PermissionSeeder) seed.Seeder {
	res := seed.NewDefaultSeeder(migrate, seed.NewUowContributor(uow, seed.Chain(roleSeed, userSeed, p)))
	return res
}

func NewHydra(c *conf2.Security) *client.APIClient {
	cfg := client.NewConfiguration()
	cfg.Servers = client.ServerConfigurations{
		{
			URL: c.Oidc.Hydra.AdminUrl,
		},
	}
	return client.NewAPIClient(cfg)
}
