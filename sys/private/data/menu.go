package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/goxiaoy/go-eventbus"
	kitgorm "github.com/goxiaoy/go-saas-kit/pkg/gorm"
	"github.com/goxiaoy/go-saas-kit/pkg/query"
	v1 "github.com/goxiaoy/go-saas-kit/sys/api/menu/v1"
	"github.com/goxiaoy/go-saas-kit/sys/private/biz"
	sgorm "github.com/goxiaoy/go-saas/gorm"
	"gorm.io/gorm"
)

type MenuRepo struct {
	*kitgorm.Repo[biz.Menu, string, v1.ListMenuRequest]
}

var _ biz.MenuRepo = (*MenuRepo)(nil)

func NewMenuRepo(dbProvider sgorm.DbProvider, eventbus *eventbus.EventBus) biz.MenuRepo {
	res := &MenuRepo{}
	res.Repo = kitgorm.NewRepo[biz.Menu, string, v1.ListMenuRequest](dbProvider, eventbus, res)
	return res
}

func (c *MenuRepo) GetDb(ctx context.Context) *gorm.DB {
	return GetDb(ctx, c.Repo.DbProvider)
}

//BuildDetailScope preload relations
func (c *MenuRepo) BuildDetailScope(withDetail bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("Requirement")
	}
}

//BuildFilterScope filter
func (c *MenuRepo) BuildFilterScope(q *v1.ListMenuRequest) func(db *gorm.DB) *gorm.DB {
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

		if filter.IdIn != nil {
			ret = ret.Where("id IN ?", filter.IdIn)
		}
		if filter.NameIn != nil {
			ret = ret.Where("name IN ?", filter.NameIn)
		}
		if filter.ParentIn != nil {
			ret = ret.Where("parent IN ?", filter.ParentIn)
		}
		return ret
	}
}

func (c *MenuRepo) DefaultSorting() []string {
	return []string{"created_at"}
}

func (c *MenuRepo) FindByName(ctx context.Context, name string) (*biz.Menu, error) {
	db := c.GetDb(ctx).Model(&biz.Menu{})
	db = db.Scopes(c.BuildDetailScope(true)).Where("name = ?", name)
	var item = biz.Menu{}
	err := db.First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (c *MenuRepo) Update(ctx context.Context, id string, entity *biz.Menu, p query.Select) error {
	if entity.Requirement != nil {
		if err := c.GetDb(ctx).Model(entity).Association("Requirement").Replace(entity.Requirement); err != nil {
			return err
		}
	}
	return c.Repo.Update(ctx, id, entity, p)
}

func (c *MenuRepo) Delete(ctx context.Context, id string) error {
	if err := c.GetDb(ctx).Delete(&biz.MenuPermissionRequirement{}, "menu_id = ?", id).Error; err != nil {
		return err
	}
	return c.Repo.Delete(ctx, id)
}
