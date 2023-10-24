package data

import (
	"context"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	v1 "github.com/go-saas/kit/product/api/price/v1"
	"github.com/go-saas/kit/product/private/biz"
	sgorm "github.com/go-saas/saas/gorm"
	"github.com/goxiaoy/go-eventbus"
	"gorm.io/gorm"
)

type PriceRepo struct {
	*kitgorm.Repo[biz.Price, string, *v1.ListPriceRequest]
}

var _ biz.PriceRepo = (*PriceRepo)(nil)

func NewPriceRepo(dbProvider sgorm.DbProvider, eventbus *eventbus.EventBus) biz.PriceRepo {
	res := &PriceRepo{}
	res.Repo = kitgorm.NewRepo[biz.Price, string, *v1.ListPriceRequest](dbProvider, eventbus, res)
	return res
}

func (c *PriceRepo) GetDb(ctx context.Context) *gorm.DB {
	return GetDb(ctx, c.DbProvider)
}

// BuildDetailScope preload relations
func (c *PriceRepo) BuildDetailScope(withDetail bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Preload("CurrencyOptions").Preload("CurrencyOptions.Tiers").Preload("Recurring").Preload("Tiers")
		return db
	}
}

// BuildFilterScope filter
func (c *PriceRepo) BuildFilterScope(q *v1.ListPriceRequest) func(db *gorm.DB) *gorm.DB {
	filter := q.Filter
	return func(db *gorm.DB) *gorm.DB {
		ret := db
		if filter == nil {
			return ret
		}

		if filter.OwnerType != nil {
			ret = ret.Scopes(kitgorm.BuildStringFilter("`owner_type`", filter.OwnerType))
		}
		if filter.OwnerId != nil {
			ret = ret.Scopes(kitgorm.BuildStringFilter("`owner_id`", filter.OwnerId))
		}
		return ret
	}
}

func (c *PriceRepo) DefaultSorting() []string {
	return []string{"created_at"}
}
