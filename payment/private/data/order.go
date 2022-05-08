package data

import (
	"context"
	"fmt"
	"github.com/goxiaoy/go-eventbus"
	v1 "github.com/goxiaoy/go-saas-kit/payment/api/order/v1"
	"github.com/goxiaoy/go-saas-kit/payment/private/biz"
	kitgorm "github.com/goxiaoy/go-saas-kit/pkg/gorm"
	sgorm "github.com/goxiaoy/go-saas/gorm"
	"gorm.io/gorm"
)

type OrderRepo struct {
	*kitgorm.Repo[biz.Order, string, v1.ListOrderRequest]
}

func NewOrderRepo(dbProvider sgorm.DbProvider, eventbus *eventbus.EventBus) biz.OrderRepo {
	res := &OrderRepo{}
	res.Repo = kitgorm.NewRepo[biz.Order, string, v1.ListOrderRequest](dbProvider, eventbus, res)
	return res
}

func (c *OrderRepo) GetDb(ctx context.Context) *gorm.DB {
	return GetDb(ctx, c.Repo.DbProvider)
}

//BuildDetailScope preload relations
func (c *OrderRepo) BuildDetailScope(withDetail bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

//BuildFilterScope filter
func (c *OrderRepo) BuildFilterScope(q *v1.ListOrderRequest) func(db *gorm.DB) *gorm.DB {
	search := q.Search
	filter := q.Filter
	return func(db *gorm.DB) *gorm.DB {
		ret := db

		if search != "" {
			ret = ret.Where("name like ?", fmt.Sprintf("%%%v%%", search))
		}
		if filter == nil {
			return ret
		}
		if filter.Id != nil {
			ret = ret.Scopes(kitgorm.BuildStringFilter("`id`", filter.Id))
		}
		if filter.Name != nil {
			ret = ret.Scopes(kitgorm.BuildStringFilter("`name`", filter.Name))
		}
		return ret
	}
}

func (c *OrderRepo) DefaultSorting() []string {
	return []string{"created_at"}
}
