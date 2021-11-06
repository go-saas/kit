package data

import (
	"context"
	"errors"
	"github.com/a8m/rql"
	gorm2 "github.com/goxiaoy/go-saas-kit/pkg/gorm"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/user/v1"
	"github.com/goxiaoy/go-saas-kit/user/internal_/biz"
	"gorm.io/gorm"
)

type RoleRepo struct {
	Repo
}

func NewRoleRepo(data *Data) biz.RoleRepo {
	return &RoleRepo{
		Repo{
			DbProvider: data.DbProvider,
		},
	}
}

func buildRoleScope(filter *v1.RoleFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		ret := db
		if filter.IdIn != nil {
			ret = ret.Where("id IN ?", filter.IdIn)
		}
		if filter.NameIn != nil {
			ret = ret.Where("name IN ?", filter.NameIn)
		}
		return ret
	}
}

func (r *RoleRepo) List(ctx context.Context, query *v1.ListRolesRequest) ([]*biz.Role, error) {
	db := r.GetDb(ctx).Model(&biz.Role{})
	db = db.Scopes(buildRoleScope(query.Filter), gorm2.SortScope(query, []string{"-created_at"}), gorm2.PageScope(query))
	var items []*biz.Role
	res := db.Find(&items)
	return items, res.Error
}

func (r *RoleRepo) First(ctx context.Context, query *v1.RoleFilter) (*biz.Role, error) {
	db := r.GetDb(ctx).Model(&biz.Role{})
	db = db.Scopes(buildRoleScope(query))
	var item = biz.Role{}
	if err := db.First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *RoleRepo) FindByName(ctx context.Context, name string) (*biz.Role, error) {
	db := r.GetDb(ctx)
	var item = &biz.Role{}
	if err := db.Where("normalized_name = ?", name).First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (r *RoleRepo) Count(ctx context.Context, query *v1.RoleFilter) (total int64, filtered int64, err error) {
	db := r.GetDb(ctx).Model(&biz.Role{})
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	db = db.Scopes(buildRoleScope(query))
	if err != nil {
		return
	}
	err = db.Count(&filtered).Error
	return
}

func (r *RoleRepo) Get(ctx context.Context, id string) (*biz.Role, error) {
	db := r.GetDb(ctx)
	var item = &biz.Role{}
	if err := db.First(item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (r *RoleRepo) Create(ctx context.Context, role *biz.Role) error {
	db := r.GetDb(ctx)
	return db.Create(role).Error
}

func (r *RoleRepo) Update(ctx context.Context, id string, role *biz.Role, p rql.Select) error {
	db := r.GetDb(ctx)
	return db.Where("id=?", id).Updates(role).Error
}

func (r *RoleRepo) Delete(ctx context.Context, id string) error {
	panic("implement me")
}
