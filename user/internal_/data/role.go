package data

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"github.com/goxiaoy/go-saas-kit/pkg/rql"
	"github.com/goxiaoy/go-saas-kit/user/internal_/biz"
)

type RoleRepo struct {
	Repo
}
func NewRoleRepo(data *Data) biz.RoleRepo  {
	return &RoleRepo{
		Repo{
			DbProvider: data.DbProvider,
		},
	}
}

func (r RoleRepo) List(ctx context.Context, query interface{}) ([]*biz.Role, error) {
	db := r.GetDb(ctx)
	db, err := r.BuildQuery(db, &biz.Role{}, query)
	if err != nil {
		return nil, err
	}
	var items []*biz.Role
	res := db.Find(&items)
	return items, res.Error
}

func (r RoleRepo) First(ctx context.Context, query interface{}) (*biz.Role, error) {
	db := r.GetDb(ctx)
	db, err := r.BuildFilter(db, &biz.Role{}, query)
	if err != nil {
		return nil,err
	}
	var item = biz.Role{}
	if err = db.First(&item).Error;err!=nil {
		if errors.Is(err,gorm.ErrRecordNotFound){
			return nil,nil
		}
		return nil,err
	}
	return &item,nil
}

func (r RoleRepo) FindByName(ctx context.Context, name string) (*biz.Role, error) {
	db := r.GetDb(ctx)
	var item = &biz.Role{}
	if err :=db.Where("normalized_name = ?",name).First(item).Error;err!=nil {
		if errors.Is(err,gorm.ErrRecordNotFound){
			return nil,nil
		}
		return nil,err
	}
	return item,nil
}

func (r RoleRepo) Count(ctx context.Context, query interface{}) (total int64, filtered int64, err error) {
	db := r.GetDb(ctx)
	err = db.Model(&biz.Role{}).Count(&total).Error
	if err != nil {
		return
	}
	db, err = r.BuildFilter(db, &biz.Role{}, query)
	if err != nil {
		return
	}
	err = db.Count(&filtered).Error
	return
}

func (r RoleRepo) Get(ctx context.Context, id string) (*biz.Role, error) {
	db := r.GetDb(ctx)
	var item = &biz.Role{}
	if err:=db.First(item,id).Error;err!=nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
			return nil,nil
		}
		return nil,err
	}
	return item,nil
}

func (r RoleRepo) Create(ctx context.Context, role *biz.Role) error {
	db := r.GetDb(ctx)
	return db.Create(role).Error
}

func (r RoleRepo) Update(ctx context.Context, id string, role *biz.Role, p rql.Select) error {
	db := r.GetDb(ctx)
	return db.Where("id=?",id).Updates(role).Error
}

func (r RoleRepo) Delete(ctx context.Context, id string) error {
	panic("implement me")
}
