package data

import (
	"context"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/kit/product/private/biz"
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
	//make sure database exists
	ctx = kitgorm.NewDbGuardianContext(ctx)
	db := GetDb(ctx, m.data.DbProvider)
	return migrateDb(db)
}

func migrateDb(db *gorm.DB) error {
	return db.AutoMigrate(
		&biz.Brand{}, &biz.BrandTrans{}, &biz.ProductCategory{},
		&biz.Product{}, &biz.ProductMedia{}, &biz.Badge{}, &biz.Keyword{}, &biz.CampaignRule{},
		&biz.Price{}, &biz.PriceCurrencyOption{}, &biz.PriceCurrencyOptionTier{}, &biz.PriceRecurring{}, &biz.PriceTier{},
		&biz.ProductAttribute{}, &biz.Stock{},
		&biz.ProductSku{},
		&biz.ProductSyncLink{})
}
