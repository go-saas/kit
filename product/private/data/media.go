package data

import (
	"context"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/kit/product/private/biz"
	sgorm "github.com/go-saas/saas/gorm"
	"github.com/goxiaoy/go-eventbus"
	"gorm.io/gorm"
)

type ProductMediaRepo struct {
	*kitgorm.Repo[biz.ProductMedia, string, interface{}]
}

func NewProductMediaRepo(dbProvider sgorm.DbProvider, eventbus *eventbus.EventBus) biz.ProductMediaRepo {
	res := &ProductMediaRepo{}
	res.Repo = kitgorm.NewRepo[biz.ProductMedia, string, interface{}](dbProvider, eventbus, res)
	return res
}

func (c *ProductMediaRepo) GetDb(ctx context.Context) *gorm.DB {
	return GetDb(ctx, c.Repo.DbProvider)
}
