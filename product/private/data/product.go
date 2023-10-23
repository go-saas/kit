package data

import (
	"context"
	"fmt"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	v1 "github.com/go-saas/kit/product/api/product/v1"
	"github.com/go-saas/kit/product/private/biz"
	sgorm "github.com/go-saas/saas/gorm"
	"github.com/goxiaoy/go-eventbus"
	"gorm.io/gorm"
)

type ProductRepo struct {
	*kitgorm.Repo[biz.Product, string, v1.ListProductRequest]
}

func NewProductRepo(dbProvider sgorm.DbProvider, eventbus *eventbus.EventBus) biz.ProductRepo {
	res := &ProductRepo{}
	res.Repo = kitgorm.NewRepo[biz.Product, string, v1.ListProductRequest](dbProvider, eventbus, res)
	return res
}

func (c *ProductRepo) GetDb(ctx context.Context) *gorm.DB {
	return GetDb(ctx, c.Repo.DbProvider)
}

// BuildDetailScope preload relations
func (c *ProductRepo) BuildDetailScope(withDetail bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

// BuildFilterScope filter
func (c *ProductRepo) BuildFilterScope(q *v1.ListProductRequest) func(db *gorm.DB) *gorm.DB {
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

func (c *ProductRepo) DefaultSorting() []string {
	return []string{"created_at"}
}
