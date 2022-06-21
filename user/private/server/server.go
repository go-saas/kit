package server

import (
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/pkg/event/event"
	uow2 "github.com/goxiaoy/go-saas-kit/pkg/uow"
	api2 "github.com/goxiaoy/go-saas-kit/user/api"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	"github.com/goxiaoy/go-saas-kit/user/private/data"
	"github.com/goxiaoy/go-saas-kit/user/private/service"
	"github.com/goxiaoy/go-saas/seed"
	"github.com/goxiaoy/uow"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer, NewJobServer, wire.Value(ClientName), wire.Value(biz.ConnName), NewSeeding, NewSeeder, NewAuthorizationOption, NewEventHandler)

var ClientName api.ClientName = api2.ServiceName

type Seeding seed.Contributor

func NewSeeding(uow uow.Manager,
	migrate *data.Migrate,
	roleSeed *biz.RoleSeed,
	userSeed *biz.UserSeed,
	p *biz.PermissionSeeder) Seeding {
	return seed.Chain(migrate, uow2.NewUowContributor(uow, seed.Chain(roleSeed, userSeed, p)))
}

func NewSeeder(us Seeding) seed.Seeder {
	res := seed.NewDefaultSeeder(us)
	return res
}

func NewAuthorizationOption(userRole *service.UserRoleContributor) *authz.Option {
	return authz.NewAuthorizationOption(userRole)
}

func NewEventHandler(e biz.UserEventHandler) event.Handler {
	return event.Handler(e)
}
