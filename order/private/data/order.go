package data

import (
	"context"
	"fmt"
	v1 "github.com/go-saas/kit/order/api/order/v1"
	"github.com/go-saas/kit/order/private/biz"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	sgorm "github.com/go-saas/saas/gorm"
	"github.com/goxiaoy/go-eventbus"
	"gorm.io/gorm"
)

type OrderRepo struct {
	*kitgorm.Repo[biz.Order, string, *v1.ListOrderRequest]
}

var _ biz.OrderRepo = (*OrderRepo)(nil)

func NewOrderRepo(dbProvider sgorm.DbProvider, eventbus *eventbus.EventBus) biz.OrderRepo {
	res := &OrderRepo{}
	res.Repo = kitgorm.NewRepo[biz.Order, string, *v1.ListOrderRequest](dbProvider, eventbus, res)
	return res
}

func (c *OrderRepo) GetDb(ctx context.Context) *gorm.DB {
	return GetDb(ctx, c.Repo.DbProvider)
}

// BuildDetailScope preload relations
func (c *OrderRepo) BuildDetailScope(withDetail bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Preload("Items")
		if withDetail {
			db = db.Preload("PaymentProviders")
		}
		return db
	}
}

// BuildFilterScope filter
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
		if filter.CustomerId != nil {
			ret = ret.Scopes(kitgorm.BuildStringFilter("`customer_id`", filter.Name))
		}
		return ret
	}
}

func (c *OrderRepo) DefaultSorting() []string {
	return []string{"-created_at"}
}

func (c *OrderRepo) FindByPaymentProvider(ctx context.Context, provider, providerKey string) (*biz.Order, error) {
	g := &biz.Order{}
	err := c.GetDb(ctx).Model(&biz.Order{}).Scopes(c.BuildDetailScope(true)).
		Joins("left join order_payment_providers on order_payment_providers.order_id = orders.id").
		Where("order_payment_providers.provider=? and order_payment_providers.provider_key=?", provider, providerKey).First(g).Error
	return g, err
}

func (c *OrderRepo) UpsertPaymentProvider(ctx context.Context, order *biz.Order, provider *biz.OrderPaymentProvider) error {
	if len(provider.Provider) == 0 {
		return fmt.Errorf("provider is required")
	}
	err := c.GetDb(ctx).Model(&biz.OrderPaymentProvider{}).Delete(&biz.OrderPaymentProvider{}, "order_id = ? AND provider = ?", order.ID, provider.Provider).Error
	if err != nil {
		return err
	}
	provider.OrderID = order.ID
	return c.GetDb(ctx).Model(&biz.OrderPaymentProvider{}).Create(provider).Error
}

func (c *OrderRepo) ListPaymentProviders(ctx context.Context, order *biz.Order) ([]biz.OrderPaymentProvider, error) {
	var providers []biz.OrderPaymentProvider
	err := c.GetDb(ctx).Model(&biz.OrderPaymentProvider{}).Where("order_id = ?", order.ID).Find(providers).Error
	return providers, err
}
