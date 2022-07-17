package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-saas/kit/pkg/authz/casbin"
	_ "github.com/go-saas/kit/pkg/blob/memory"
	_ "github.com/go-saas/kit/pkg/blob/os"
	_ "github.com/go-saas/kit/pkg/blob/s3"
	conf2 "github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/kit/pkg/dal"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/go-saas/kit/user/private/biz"
	"github.com/go-saas/saas/gorm"
	g "gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = kitdi.NewSet(
	NewData,
	NewEnforcerProvider,
	NewUserRepo,
	NewRefreshTokenRepo,
	NewRoleRepo,
	NewUserTenantRepo,
	NewMigrate,
	NewUserTokenRepo,
	NewUserSettingRepo,
	NewUserAddrRepo,
)

// Data .
type Data struct {
	DbProvider gorm.DbProvider
}

func GetDb(ctx context.Context, provider gorm.DbProvider) *g.DB {
	db := provider.Get(ctx, string(biz.ConnName))
	if err := db.SetupJoinTable(&biz.User{}, "Roles", &biz.UserRole{}); err != nil {
		panic(err)
	}
	return db
}

// NewData .
func NewData(c *conf2.Data, dbProvider dal.ConstDbProvider, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		logger.Log(log.LevelInfo, log.DefaultMessageKey, "closing the data resources")
	}
	return &Data{
		DbProvider: dbProvider,
	}, cleanup, nil
}

func NewEnforcerProvider(logger log.Logger, dbProvider gorm.DbProvider) (*casbin.EnforcerProvider, error) {
	return casbin.NewEnforcerProvider(logger, dbProvider, string(biz.ConnName))
}
