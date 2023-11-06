package data

import (
	"context"
	"fmt"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/kit/pkg/query"
	v1 "github.com/go-saas/kit/product/api/product/v1"
	"github.com/go-saas/kit/product/private/biz"
	sgorm "github.com/go-saas/saas/gorm"
	"github.com/goxiaoy/go-eventbus"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"gorm.io/gorm"
)

type ProductRepo struct {
	*kitgorm.Repo[biz.Product, string, *v1.ListProductRequest]
}

func NewProductRepo(dbProvider sgorm.DbProvider, eventbus *eventbus.EventBus) biz.ProductRepo {
	res := &ProductRepo{}
	res.Repo = kitgorm.NewRepo[biz.Product, string, *v1.ListProductRequest](dbProvider, eventbus, res)
	return res
}

func (c *ProductRepo) GetDb(ctx context.Context) *gorm.DB {
	return GetDb(ctx, c.Repo.DbProvider)
}

// BuildDetailScope preload relations
func (c *ProductRepo) BuildDetailScope(withDetail bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Preload("MainPic").
			Preload("Badges").Preload("Categories").Preload("MainCategory").Preload("Keywords").Preload("Attributes")
		if withDetail {
			db = db.Preload("Medias").Preload("CampaignRules").Preload("Stocks").Preload("SyncLinks").
				Preload("Prices").Preload("Prices.CurrencyOptions").Preload("Prices.CurrencyOptions.Tiers").Preload("Prices.Recurring").Preload("Prices.Tiers")
		}
		if !withDetail {
			db = db.Omit("Content")
		}
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
		if filter.Internal != nil {
			ret = ret.Scopes(kitgorm.BuildBooleanFilter("`internal`", filter.Internal))
		}
		return ret
	}
}

func (c *ProductRepo) UpdateAssociation(ctx context.Context, entity *biz.Product, p query.Select) error {
	if query.SelectContains(p, "Medias") {
		if err := c.GetDb(ctx).Model(entity).
			Association("Medias").Replace(entity.Medias); err != nil {
			return err
		}
	}
	if query.SelectContains(p, "Badges") {
		if err := c.GetDb(ctx).Model(entity).Session(&gorm.Session{FullSaveAssociations: true}).
			Association("Badges").Replace(entity.Badges); err != nil {
			return err
		}
	}
	if query.SelectContains(p, "Categories") {
		if err := c.GetDb(ctx).Model(entity).
			Association("Badges").Replace(entity.Badges); err != nil {
			return err
		}
	}
	if query.SelectContains(p, "Keywords") {
		if err := c.GetDb(ctx).Model(entity).Session(&gorm.Session{FullSaveAssociations: true}).
			Association("Keywords").Replace(entity.Keywords); err != nil {
			return err
		}
	}
	if query.SelectContains(p, "Attributes") {
		if err := c.GetDb(ctx).Model(entity).Session(&gorm.Session{FullSaveAssociations: true}).
			Association("Attributes").Replace(entity.Attributes); err != nil {
			return err
		}
	}
	if query.SelectContains(p, "CampaignRules") {
		if err := c.GetDb(ctx).Model(entity).Session(&gorm.Session{FullSaveAssociations: true}).
			Association("CampaignRules").Replace(entity.CampaignRules); err != nil {
			return err
		}
	}
	if query.SelectContains(p, "Stocks") {
		if err := c.GetDb(ctx).Model(entity).Session(&gorm.Session{FullSaveAssociations: true}).
			Association("Stocks").Replace(entity.Stocks); err != nil {
			return err
		}
	}
	if query.SelectContains(p, "Prices") {
		if err := c.GetDb(ctx).Model(entity).Session(&gorm.Session{FullSaveAssociations: true}).
			Association("Prices").Replace(entity.Prices); err != nil {
			return err
		}
	}
	return nil
}

func (c *ProductRepo) DefaultSorting() []string {
	return []string{"-created_at"}
}

func (c *ProductRepo) GetSyncLinks(ctx context.Context, product *biz.Product) ([]biz.ProductSyncLink, error) {
	var links []biz.ProductSyncLink
	if err := c.GetDb(ctx).Model(&biz.ProductSyncLink{}).Where("product_id = ?", product.ID.String()).Find(&links).Error; err != nil {
		return links, err
	}
	return links, nil
}

func (c *ProductRepo) UpdateSyncLink(ctx context.Context, product *biz.Product, syncLink *biz.ProductSyncLink) error {
	syncLink.ProductId = product.ID.String()
	if err := c.GetDb(ctx).Save(syncLink).Error; err != nil {
		return err
	}
	if product.ManageInfo.Managed && product.ManageInfo.ManagedBy == syncLink.ProviderName {
		product.ManageInfo.LastSyncTime = syncLink.LastSyncTime
		return c.Update(ctx, product.ID.String(), product, query.NewField(&fieldmaskpb.FieldMask{Paths: []string{"last_sync_time"}}))
	}
	return nil
}
