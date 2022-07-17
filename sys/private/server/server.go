package server

import (
	dtmdata "github.com/go-saas/kit/dtm/data"
	"github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/dal"
	kitdi "github.com/go-saas/kit/pkg/di"
	ksaas "github.com/go-saas/kit/pkg/saas"
	uow2 "github.com/go-saas/kit/pkg/uow"
	api2 "github.com/go-saas/kit/sys/api"
	"github.com/go-saas/kit/sys/private/biz"
	"github.com/go-saas/kit/sys/private/data"
	"github.com/go-saas/saas"
	"github.com/go-saas/saas/seed"
	"github.com/go-saas/uow"
)

// ProviderSet is server providers.
var ProviderSet = kitdi.NewSet(
	NewHTTPServer,
	NewGRPCServer,
	NewJobServer,
	NewSeeder,
	func() api.ClientName { return ClientName },

	func() dal.ConnName {
		return biz.ConnName
	},
	NewSeeding,
	NewAuthorizationOption,
)

var ClientName api.ClientName = api2.ServiceName

// Seeding workaround for https://github.com/google/wire/issues/207
type Seeding seed.Contrib

// NewSeeding sys seeding should migrate dtmsrv and dmtcli
func NewSeeding(apisixSeeder *biz.ApisixSeed, uow uow.Manager, dtmMigrator *dtmdata.Migrator, migrate *data.Migrate, menu *biz.MenuSeed) Seeding {
	return seed.Chain(apisixSeeder, dtmMigrator, migrate, uow2.NewUowContrib(uow, seed.Chain(menu)))
}

func NewSeeder(ts saas.TenantStore, ss Seeding) seed.Seeder {
	return seed.NewDefaultSeeder(ksaas.NewTraceContrib(ksaas.SeedChangeTenant(ts, ss)))
}

func NewAuthorizationOption() *authz.Option {
	return authz.NewAuthorizationOption()
}
