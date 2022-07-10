package data

import (
	"context"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/kit/user/private/biz"
	"github.com/go-saas/saas/seed"
	"gorm.io/gorm"
)

type Migrate struct {
	data *Data
}

func NewMigrate(data *Data) *Migrate {
	return &Migrate{
		data: data,
	}
}
func (m *Migrate) Seed(ctx context.Context, sCtx *seed.Context) error {
	if len(sCtx.TenantId) > 0 {
		//host only migrate
		return nil
	}
	//make sure database exists
	ctx = kitgorm.NewDbGuardianContext(ctx)
	db := GetDb(ctx, m.data.DbProvider)
	return migrateDb(db)
}

func migrateDb(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&biz.User{},
		&biz.Role{},
		&biz.UserRole{},
		&biz.UserLogin{},
		&biz.UserSetting{},
		&biz.UserToken{},
		&biz.RefreshToken{},
		&biz.UserTenant{},
		&biz.UserAddress{}); err != nil {
		return err
	}
	//migrate casbin
	if _, err := gormadapter.NewAdapterByDB(db); err != nil {
		return err
	}

	return nil
}
