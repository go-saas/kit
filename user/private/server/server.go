package server

import (
	"github.com/go-kratos/kratos/v2/transport"
	dtmdata "github.com/go-saas/kit/dtm/data"
	"github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authz/authz"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/go-saas/kit/pkg/server"
	api2 "github.com/go-saas/kit/user/api"
	"github.com/go-saas/kit/user/private/biz"
	"github.com/go-saas/kit/user/private/data"
	"github.com/go-saas/kit/user/private/service"
	"github.com/go-saas/saas"
	"github.com/go-saas/saas/seed"
	"github.com/go-saas/uow"
	"github.com/goava/di"
)

// ProviderSet is server providers.
var ProviderSet = kitdi.NewSet(
	kitdi.NewProvider(NewHTTPServer, di.As(new(transport.Server))),
	kitdi.NewProvider(NewGRPCServer, di.As(new(transport.Server))),
	//kitdi.NewProvider(NewJobServer, di.As(new(transport.Server))),
	kitdi.NewProvider(NewEventServer, di.As(new(transport.Server))),
	kitdi.Value(ClientName),
	kitdi.Value(biz.ConnName),
	NewSeeding,
	NewSeeder,
	NewAuthorizationOption,
)

var ClientName api.ClientName = api2.ServiceName

// NewSeeding wrap all service migrator into one seed.Contrib, which grants the running sequence of those contribs
func NewSeeding(uow uow.Manager,
	migrate *data.Migrate,
	dtmMigrator *dtmdata.BarrierMigrator, //barrier only
	roleSeed *biz.RoleSeed,
	userSeed *biz.UserSeed,
	p *biz.PermissionSeeder) seed.Contrib {
	return seed.Chain(migrate, dtmMigrator, server.NewUowContrib(uow, roleSeed, userSeed, p))
}

func NewSeeder(ts saas.TenantStore, seeds []seed.Contrib) seed.Seeder {
	res := seed.NewDefaultSeeder(server.NewTraceContrib(server.SeedChangeTenant(ts, seeds...)))
	return res
}

func NewAuthorizationOption(userRole *service.UserRoleContrib) *authz.Option {
	return authz.NewAuthorizationOption(userRole)
}
