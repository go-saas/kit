package data

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/user/internal_/biz"
	"github.com/goxiaoy/go-saas/seed"
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
		&biz.UserTenant{}); err != nil {
		return err
	}

	return nil
}
