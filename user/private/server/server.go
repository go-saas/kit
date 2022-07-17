package server

import (
	dtmdata "github.com/go-saas/kit/dtm/data"
	"github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/dal"
	kitdi "github.com/go-saas/kit/pkg/di"
	ksaas "github.com/go-saas/kit/pkg/saas"
	uow2 "github.com/go-saas/kit/pkg/uow"
	api2 "github.com/go-saas/kit/user/api"
	"github.com/go-saas/kit/user/private/biz"
	"github.com/go-saas/kit/user/private/data"
	"github.com/go-saas/kit/user/private/service"
	"github.com/go-saas/saas"
	"github.com/go-saas/saas/seed"
	"github.com/go-saas/uow"
)

// ProviderSet is server providers.
var ProviderSet = kitdi.NewSet(
	NewHTTPServer,
	NewGRPCServer,
	NewJobServer,
	NewEventServer,
	func() api.ClientName { return ClientName },
	func() dal.ConnName { return biz.ConnName },
	NewSeeding,
	NewSeeder,
	NewAuthorizationOption,
)

var ClientName api.ClientName = api2.ServiceName

// Seeding workaround for https://github.com/google/wire/issues/207
type Seeding seed.Contrib

func NewSeeding(uow uow.Manager,
	migrate *data.Migrate,
	dtmMigrator *dtmdata.BarrierMigrator, //barrier only
	roleSeed *biz.RoleSeed,
	userSeed *biz.UserSeed,
	p *biz.PermissionSeeder) Seeding {
	return seed.Chain(migrate, dtmMigrator, uow2.NewUowContrib(uow, roleSeed, userSeed, p))
}

func NewSeeder(ts saas.TenantStore, us Seeding) seed.Seeder {
	res := seed.NewDefaultSeeder(ksaas.NewTraceContrib(ksaas.SeedChangeTenant(ts, us)))
	return res
}

func NewAuthorizationOption(userRole *service.UserRoleContrib) *authz.Option {
	return authz.NewAuthorizationOption(userRole)
}
