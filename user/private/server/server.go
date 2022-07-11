package server

import (
	dtmdata "github.com/go-saas/kit/dtm/data"
	"github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authz/authz"
	ksaas "github.com/go-saas/kit/pkg/saas"
	uow2 "github.com/go-saas/kit/pkg/uow"
	api2 "github.com/go-saas/kit/user/api"
	"github.com/go-saas/kit/user/private/biz"
	"github.com/go-saas/kit/user/private/data"
	"github.com/go-saas/kit/user/private/service"
	"github.com/go-saas/saas"
	"github.com/go-saas/saas/seed"
	"github.com/go-saas/uow"
	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer, NewJobServer, NewEventServer, wire.Value(ClientName), wire.Value(biz.ConnName), NewSeeding, NewSeeder, NewAuthorizationOption)

var ClientName api.ClientName = api2.ServiceName

// Seeding workaround for https://github.com/google/wire/issues/207
type Seeding seed.Contrib

func NewSeeding(uow uow.Manager,
	migrate *data.Migrate,
	roleSeed *biz.RoleSeed,
	userSeed *biz.UserSeed,
	p *biz.PermissionSeeder) Seeding {
	return seed.Chain(migrate, uow2.NewUowContrib(uow, roleSeed, userSeed, p))
}

func NewSeeder(ts saas.TenantStore, dtmMigrator *dtmdata.Migrator, us Seeding) seed.Seeder {
	res := seed.NewDefaultSeeder(ksaas.NewTraceContrib(dtmMigrator, ksaas.SeedChangeTenant(ts, us)))
	return res
}

func NewAuthorizationOption(userRole *service.UserRoleContrib) *authz.Option {
	return authz.NewAuthorizationOption(userRole)
}
