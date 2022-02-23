package data

import (
	"context"
	"errors"
	"fmt"
	gorm2 "github.com/goxiaoy/go-saas-kit/pkg/gorm"
	"github.com/goxiaoy/go-saas-kit/pkg/query"
	v1 "github.com/goxiaoy/go-saas-kit/sys/api/menu/v1"
	"github.com/goxiaoy/go-saas-kit/sys/private/biz"
	"github.com/goxiaoy/go-saas/gorm"
	g "gorm.io/gorm"
)

type MenuRepo struct {
	Repo
}

func NewMenuRepo(dbProvider gorm.DbProvider) biz.MenuRepo {
	return &MenuRepo{Repo: Repo{DbProvider: dbProvider}}
}

func buildMenuScope(search string, filter *v1.MenuFilter) func(db *g.DB) *g.DB {
	return func(db *g.DB) *g.DB {
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

func preloadMenuScope() func(db *g.DB) *g.DB {
	return func(db *g.DB) *g.DB {
		return db.Preload("Requirement")
	}
}

func (c *MenuRepo) List(ctx context.Context, query *v1.ListMenuRequest) ([]*biz.Menu, error) {
	db := c.GetDb(ctx).Model(&biz.Menu{})
	db = db.Scopes(buildMenuScope(query.Search, query.Filter), preloadMenuScope(), gorm2.SortScope(query, []string{"-created_at"}), gorm2.PageScope(query))
	var items []*biz.Menu
	res := db.Find(&items)
	return items, res.Error
}

func (c *MenuRepo) First(ctx context.Context, search string, query *v1.MenuFilter) (*biz.Menu, error) {
	db := c.GetDb(ctx).Model(&biz.Menu{})
	db = db.Scopes(buildMenuScope(search, query), preloadMenuScope())
	var item = biz.Menu{}
	err := db.First(&item).Error
	if err != nil {
		if errors.Is(err, g.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (c *MenuRepo) FindByName(ctx context.Context, name string) (*biz.Menu, error) {
	db := c.GetDb(ctx).Model(&biz.Menu{})
	db = db.Scopes(preloadMenuScope()).Where("name = ?", name)
	var item = biz.Menu{}
	err := db.First(&item).Error
	if err != nil {
		if errors.Is(err, g.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (c *MenuRepo) Count(ctx context.Context, search string, query *v1.MenuFilter) (total int64, filtered int64, err error) {
	db := c.GetDb(ctx).Model(&biz.Menu{})
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	db = db.Scopes(buildMenuScope(search, query))
	err = db.Count(&filtered).Error
	return
}

func (c *MenuRepo) Get(ctx context.Context, id string) (*biz.Menu, error) {
	var entity = &biz.Menu{}
	err := c.GetDb(ctx).Model(&biz.Menu{}).Scopes(preloadMenuScope()).First(entity, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, g.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return entity, nil
}

func (c *MenuRepo) Create(ctx context.Context, entity *biz.Menu) error {
	return c.GetDb(ctx).Create(entity).Error
}

func (c *MenuRepo) Update(ctx context.Context, entity *biz.Menu, p query.Select) error {
	return c.GetDb(ctx).Updates(entity).Error
}

func (c *MenuRepo) Delete(ctx context.Context, id string) error {
	return c.GetDb(ctx).Delete(&biz.Menu{}, "id = ?", id).Error
}
