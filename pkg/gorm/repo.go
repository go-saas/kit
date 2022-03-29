package gorm

import (
	"context"
	"errors"

	"github.com/goxiaoy/go-saas-kit/pkg/data"
	"github.com/goxiaoy/go-saas-kit/pkg/query"
	sgorm "github.com/goxiaoy/go-saas/gorm"
	"gorm.io/gorm"
)

type Repo[TEntity any, TKey any, TQuery any] struct {
	DbProvider sgorm.DbProvider
	override   interface{}
}

//var _ data.Repo[interface{}, interface{}, interface{}] = (*Repo[interface{}, interface{}, interface{}])(nil)

func NewRepo[TEntity any, TKey any, TQuery any](dbProvider sgorm.DbProvider, override interface{}) *Repo[TEntity, TKey, TQuery] {
	return &Repo[TEntity, TKey, TQuery]{DbProvider: dbProvider, override: override}
}

func (r *Repo[TEntity, TKey, TQuery]) GetDb(ctx context.Context) *gorm.DB {
	return r.DbProvider.Get(ctx, "")
}

type BuildDetailScope interface {
	BuildDetailScope(withDetail bool) func(db *gorm.DB) *gorm.DB
}

//BuildDetailScope preload relations
func (r *Repo[TEntity, TKey, TQuery]) buildDetailScope(withDetail bool) func(db *gorm.DB) *gorm.DB {
	if override, ok := r.override.(BuildDetailScope); ok {
		return override.BuildDetailScope(withDetail)
	}
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

type BuildFilterScope[TQuery any] interface {
	BuildFilterScope(q *TQuery) func(db *gorm.DB) *gorm.DB
}

//BuildFilterScope filter
func (r *Repo[TEntity, TKey, TQuery]) buildFilterScope(q *TQuery) func(db *gorm.DB) *gorm.DB {
	if override, ok := r.override.(BuildFilterScope[TQuery]); ok {
		return override.BuildFilterScope(q)
	}
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

type DefaultSorting interface {
	DefaultSorting() []string
}

//DefaultSorting get default sorting
func (r *Repo[TEntity, TKey, TQuery]) defaultSorting() []string {
	if override, ok := r.override.(DefaultSorting); ok {
		return override.DefaultSorting()
	}
	return nil
}

type BuildSortScope[TQuery any] interface {
	BuildSortScope(q *TQuery) func(db *gorm.DB) *gorm.DB
}

//BuildSortScope build sorting query
func (r *Repo[TEntity, TKey, TQuery]) buildSortScope(q *TQuery) func(db *gorm.DB) *gorm.DB {
	if override, ok := r.override.(BuildSortScope[TQuery]); ok {
		return override.BuildSortScope(q)
	}
	f, ok := (interface{})(q).(query.Sort)
	if ok {
		return SortScope(f, r.defaultSorting())
	}
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

type BuildPageScope[TQuery any] interface {
	BuildPageScope(q *TQuery) func(db *gorm.DB) *gorm.DB
}

//BuildPageScope page query
func (r *Repo[TEntity, TKey, TQuery]) buildPageScope(q *TQuery) func(db *gorm.DB) *gorm.DB {
	if override, ok := r.override.(BuildPageScope[TQuery]); ok {
		return override.BuildPageScope(q)
	}

	f, ok := (interface{})(q).(query.Page)
	if ok {
		return PageScope(f)
	}
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

func (r *Repo[TEntity, TKey, TQuery]) List(ctx context.Context, query *TQuery) ([]*TEntity, error) {
	var e TEntity
	db := r.GetDb(ctx).Model(&e)
	db = db.Scopes(r.buildFilterScope(query), r.buildDetailScope(false), r.buildSortScope(query), r.buildPageScope(query))
	var items []*TEntity
	res := db.Find(&items)
	return items, res.Error
}

func (r *Repo[TEntity, TKey, TQuery]) First(ctx context.Context, query *TQuery) (*TEntity, error) {
	var e TEntity
	db := r.GetDb(ctx).Model(&e)
	db = db.Scopes(r.buildFilterScope(query), r.buildDetailScope(true))
	var item TEntity
	err := db.First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *Repo[TEntity, TKey, TQuery]) Count(ctx context.Context, query *TQuery) (total int64, filtered int64, err error) {
	var e TEntity
	db := r.GetDb(ctx).Model(&e)
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	db = db.Scopes(r.buildFilterScope(query))
	err = db.Count(&filtered).Error
	return
}
func (r *Repo[TEntity, TKey, TQuery]) Get(ctx context.Context, id TKey) (*TEntity, error) {
	var entity TEntity
	err := r.GetDb(ctx).Model(&entity).Scopes(r.buildDetailScope(true)).First(&entity, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}
func (r *Repo[TEntity, TKey, TQuery]) Create(ctx context.Context, entity *TEntity) error {
	return r.GetDb(ctx).Create(entity).Error
}

func (r *Repo[TEntity, TKey, TQuery]) Update(ctx context.Context, id TKey, entity *TEntity, p query.Select) error {
	var e TEntity
	db := r.GetDb(ctx).Model(&e)
	if p == nil {
		db = db.Select("*")
	}
	return db.Where("id = ?", id).Updates(entity).Error
}
func (r *Repo[TEntity, TKey, TQuery]) Delete(ctx context.Context, id TKey) error {
	var e TEntity
	return r.GetDb(ctx).Delete(&e, "id = ?", id).Error
}

func PageScope(page query.Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == nil {
			return db
		}
		ret := db
		if page.GetPageOffset() > 0 {
			ret = db.Offset(int(page.GetPageOffset()))
		}
		if page.GetPageSize() > 0 {
			ret = db.Limit(int(page.GetPageSize()))
		}
		return ret
	}
}

// SortScope build sorting by sort and default d
func SortScope(sort query.Sort, d []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var s []string
		if sort != nil {
			s = sort.GetSort()
		}
		if len(s) == 0 {
			s = d
		}
		parsed := data.ParseSort(s)
		ret := db
		if parsed != "" {
			ret = ret.Order(parsed)
		}
		return ret
	}
}
