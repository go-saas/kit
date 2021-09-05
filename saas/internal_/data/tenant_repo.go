package data

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/goxiaoy/go-saas-kit/saas/internal_/biz"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/gorm"
	gg "gorm.io/gorm"
)

type TenantRepo struct {
	Repo
}

func NewTenantRepo() biz.TenantRepo {
	return &TenantRepo{Repo{
		DbProvider: GlobalData.DbProvider,
	}}
}

func (g *TenantRepo) Db(ctx context.Context, preload bool) *gg.DB {
	ret := GetDb(ctx, g.DbProvider)
	if preload {
		ret = ret.Preload("Conn").Preload("Features")
	}
	return ret
}

func (g *TenantRepo) FindByIdOrName(ctx context.Context, idOrName string) (*biz.Tenant, error) {
	var t = new(biz.Tenant)
	var tDb Tenant
	//parse
	if idOrName == "" {
		return t, nil
	}
	//parse uuid
	id, err := uuid.Parse(idOrName)
	if err == nil {
		//id
		err = g.Db(ctx, true).Where("id = ?", id.String()).First(&tDb).Error
	} else {
		err = g.Db(ctx, true).Where("name = ?", idOrName).First(&tDb).Error
	}
	if err != nil {
		if errors.Is(err, gg.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	err = common.Copy(&tDb, t)
	return t, err
}

func (g *TenantRepo) GetCount(ctx context.Context) (int64, error) {
	var count int64
	//check count
	tx := g.Db(ctx, false).Model(&Tenant{}).Count(&count)
	return count, tx.Error
}

func (g *TenantRepo) GetPaged(ctx context.Context, p common.Pagination) (c int64, t []*biz.Tenant, err error) {
	err = g.Db(ctx, false).Model(&Tenant{}).Count(&c).Error
	var tDb Tenants
	if err != nil {
		return c, nil, err
	}
	err = gorm.BuildPage(g.Db(ctx, false), p).Find(&tDb).Error
	//copy
	common.Copy(tDb, &t)
	return
}

func (g *TenantRepo) Create(ctx context.Context, t biz.Tenant) error {
	var tDb = new(Tenant)
	common.Copy(&t, tDb)
	d := g.Db(ctx, true)
	ret := d.Create(tDb)
	return ret.Error

}

func (g *TenantRepo) Update(ctx context.Context, id string, t biz.Tenant) error {
	var tDb = new(Tenant)
	common.Copy(&t, tDb)
	d := g.Db(ctx, true)
	return d.Model(&Tenant{}).Where("id = ?", id).Updates(tDb).Error
}
