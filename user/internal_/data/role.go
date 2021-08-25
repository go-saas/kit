package data

import (
	"context"
	"errors"
	"github.com/a8m/rql"
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

func (r *RoleRepo) buildFilter(db *gorm.DB,query *v1.RoleFilter)*gorm.DB{
	ret := db.Model(&biz.Role{})
	if query==nil{
		return ret
	}
	if len(query.IdIn)>0{
		ret = ret.Where("id IN ?",query.IdIn)
	}
	if len(query.NameIn)>0{
		ret = ret.Where("name IN ?",query.IdIn)
	}
	return ret
}

func (r *RoleRepo) List(ctx context.Context, query v1.ListRolesRequest) ([]*biz.Role, error) {
	db := r.GetDb(ctx)
	db =  r.buildFilter(db,query.Filter)
	db = r.BuildSort(db,&query)
	db = r.BuildPage(db,&query)
	var items []*biz.Role
	res := db.Find(&items)
	return items, res.Error
}

func (r *RoleRepo) First(ctx context.Context, query v1.RoleFilter) (*biz.Role, error) {
	db := r.GetDb(ctx)
	db =  r.buildFilter(db,&query)
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

func (r *RoleRepo) Count(ctx context.Context, query v1.RoleFilter) (total int64, filtered int64, err error) {
	db := r.GetDb(ctx)
	err = db.Model(&biz.Role{}).Count(&total).Error
	if err != nil {
		return
	}
	db = r.buildFilter(db,&query)
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
