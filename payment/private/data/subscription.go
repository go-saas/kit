package data

import (
	"context"
	v1 "github.com/go-saas/kit/payment/api/subscription/v1"
	"github.com/go-saas/kit/payment/private/biz"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	sgorm "github.com/go-saas/saas/gorm"
	"github.com/goxiaoy/go-eventbus"
	"gorm.io/gorm"
)

type SubscriptionRepo struct {
	*kitgorm.Repo[biz.Subscription, string, biz.SubscriptionListPrams]
}

var _ biz.SubscriptionRepo = (*SubscriptionRepo)(nil)

func NewSubscriptionRepo(dbProvider sgorm.DbProvider, eventbus *eventbus.EventBus) biz.SubscriptionRepo {
	res := &SubscriptionRepo{}
	res.Repo = kitgorm.NewRepo[biz.Subscription, string, biz.SubscriptionListPrams](dbProvider, eventbus, res)
	return res
}

func (c *SubscriptionRepo) GetDb(ctx context.Context) *gorm.DB {
	return GetDb(ctx, c.DbProvider)
}

// BuildDetailScope preload relations
func (c *SubscriptionRepo) BuildDetailScope(withDetail bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Preload("Items")
		return db
	}
}

func buildSubsScope(filter *v1.SubscriptionFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		ret := db
		if filter == nil {
			return ret
		}
		if len(filter.And) > 0 {
			for _, filter := range filter.And {
				ret = ret.Where(buildSubsScope(filter)(db.Session(&gorm.Session{NewDB: true})))
			}
		}
		if len(filter.Or) > 0 {
			for _, filter := range filter.Or {
				ret = ret.Or(buildSubsScope(filter)(db.Session(&gorm.Session{NewDB: true})))
			}
		}
		if filter.Id != nil {
			ret = ret.Scopes(kitgorm.BuildStringFilter("`id`", filter.Id))
		}
		if filter.UserId != nil {
			ret = ret.Scopes(kitgorm.BuildStringFilter("`user_id`", filter.UserId))
		}
		if filter.Provider != nil {
			ret = ret.Scopes(kitgorm.BuildStringFilter("`provider`", filter.Provider))
		}
		if filter.ProviderKey != nil {
			ret = ret.Scopes(kitgorm.BuildStringFilter("`provider_key`", filter.ProviderKey))
		}

		return ret
	}
}

// BuildFilterScope filter
func (c *SubscriptionRepo) BuildFilterScope(q biz.SubscriptionListPrams) func(db *gorm.DB) *gorm.DB {
	filter := q.GetFilter()
	return buildSubsScope(filter)
}

func (c *SubscriptionRepo) DefaultSorting() []string {
	return []string{"-created_at"}
}

func (c *SubscriptionRepo) FindByProvider(ctx context.Context, provider, providerKey string) (*biz.Subscription, error) {
	g := &biz.Subscription{}
	err := c.GetDb(ctx).Model(&biz.Subscription{}).Scopes(c.BuildDetailScope(true)).First(g, "provider = ? AND provider_key = ?", provider, providerKey).Error
	if err != nil {
		return nil, err
	}
	return g, nil
}
