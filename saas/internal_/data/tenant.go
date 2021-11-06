package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/google/uuid"
	gorm2 "github.com/goxiaoy/go-saas-kit/pkg/gorm"
	v1 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	"github.com/goxiaoy/go-saas-kit/saas/internal_/biz"
	"github.com/goxiaoy/go-saas/gorm"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	gg "gorm.io/gorm"
	"time"
)

type Tenant struct {
	gorm2.UIDBase
	//unique name. usually for domain name
	Name string `gorm:"column:name;index;size:255;"`
	//localed display name
	DisplayName string `gorm:"column:display_name;index;size:255;"`
	//region of this tenant
	Region    string     `gorm:"column:region;index;size:255;"`
	CreatedAt time.Time  `gorm:"column:created_at;index;"`
	UpdatedAt time.Time  `gorm:"column:updated_at;index;"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index;"`

	//connection
	Conn []TenantConn `gorm:"foreignKey:TenantId"`
	//edition
	Features []TenantFeature `gorm:"foreignKey:TenantId"`
}

type TenantRepo struct {
	DbProviderGetter
}

type DbProviderGetter func() gorm.DbProvider

func NewTenantRepo() biz.TenantRepo {
	return &TenantRepo{DbProviderGetter: func() gorm.DbProvider { return GlobalData.DbProvider }}
}

func (g *TenantRepo) GetDb(ctx context.Context) *gg.DB {
	ret := GetDb(ctx, g.DbProviderGetter())
	return ret
}

func preloadUserScope(withDetail bool) func(db *gg.DB) *gg.DB {
	return func(db *gg.DB) *gg.DB {
		if withDetail {
			return db.Preload("Conn").Preload("Features")
		}
		return db
	}
}

func buildTenantScope(search string, filter *v1.TenantFilter) func(db *gg.DB) *gg.DB {
	return func(db *gg.DB) *gg.DB {
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
		if filter.NameLike != "" {
			ret = ret.Where("name like ?", fmt.Sprintf("%%%v%%", filter.NameLike))
		}

		if filter.RegionIn != nil {
			ret = ret.Where("region IN ?", filter.RegionIn)
		}

		return ret
	}
}

func (g *TenantRepo) FindByIdOrName(ctx context.Context, idOrName string) (*biz.Tenant, error) {
	//parse
	if idOrName == "" {
		return nil, nil
	}
	//parse uuid
	id, err := uuid.Parse(idOrName)
	var tDb Tenant
	if err == nil {
		//id
		err = g.GetDb(ctx).Model(&Tenant{}).Scopes(preloadUserScope(true)).Where("id = ?", id.String()).First(&tDb).Error
	} else {
		err = g.GetDb(ctx).Model(&Tenant{}).Scopes(preloadUserScope(true)).Where("name = ?", idOrName).First(&tDb).Error
	}
	if err != nil {
		if errors.Is(err, gg.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	var ret = &biz.Tenant{}
	mapDataTenantToBizTenant(&tDb, ret)
	return ret, err
}

func (g *TenantRepo) List(ctx context.Context, query *v1.ListTenantRequest) ([]*biz.Tenant, error) {
	db := g.GetDb(ctx).Model(&Tenant{})
	db = db.Scopes(buildTenantScope(query.Search, query.Filter), gorm2.SortScope(query, []string{"-created_at"}), gorm2.PageScope(query))
	var items []*Tenant
	res := db.Find(&items)
	var rItems []*biz.Tenant
	linq.From(items).SelectT(func(t *Tenant) *biz.Tenant {
		res := &biz.Tenant{}
		mapDataTenantToBizTenant(t, res)
		return res
	}).ToSlice(&rItems)
	return rItems, res.Error
}

func (g *TenantRepo) First(ctx context.Context, search string, query *v1.TenantFilter) (*biz.Tenant, error) {
	db := g.GetDb(ctx).Model(&Tenant{})
	db = db.Scopes(buildTenantScope(search, query))
	var item = Tenant{}
	if err := db.First(&item).Error; err != nil {
		if errors.Is(err, gg.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	var ret = &biz.Tenant{}
	mapDataTenantToBizTenant(&item, ret)
	return ret, nil
}

func (g *TenantRepo) Count(ctx context.Context, search string, query *v1.TenantFilter) (total int64, filtered int64, err error) {
	db := g.GetDb(ctx).Model(&Tenant{})
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	db = db.Scopes(buildTenantScope(search, query))
	if err != nil {
		return
	}
	err = db.Count(&filtered).Error
	return
}

func (g *TenantRepo) Get(ctx context.Context, id string) (*biz.Tenant, error) {
	db := g.GetDb(ctx)
	var item = &Tenant{}
	if err := db.First(item, "id = ?", id).Error; err != nil {
		if errors.Is(err, gg.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	var ret = &biz.Tenant{}
	mapDataTenantToBizTenant(item, ret)

	return ret, nil
}

func (g *TenantRepo) Create(ctx context.Context, entity *biz.Tenant) error {
	var tDb = &Tenant{}
	mapBizTenantToDataTenant(entity, tDb)
	d := g.GetDb(ctx)
	if err := d.Create(tDb).Error; err != nil {
		return err
	}
	mapDataTenantToBizTenant(tDb, entity)
	return nil
}

func (g *TenantRepo) Update(ctx context.Context, entity *biz.Tenant, p *fieldmaskpb.FieldMask) error {
	var tDb = &Tenant{}
	mapBizTenantToDataTenant(entity, tDb)
	d := g.GetDb(ctx)

	if tDb.Conn != nil {
		d.Model(tDb).Association("Conn").Replace(tDb.Conn)
	}
	if tDb.Features != nil {
		d.Model(tDb).Association("Features").Replace(tDb.Features)
	}
	err := d.Model(&Tenant{}).Where("id = ?", entity.ID).Updates(tDb).Error
	if err != nil {
		return err
	}
	mapDataTenantToBizTenant(tDb, entity)
	return nil
}

func (g *TenantRepo) Delete(ctx context.Context, id string) error {
	return g.GetDb(ctx).Delete(&Tenant{}, "id = ?", id).Error
}

func mapBizTenantToDataTenant(a *biz.Tenant, b *Tenant) {
	var conn []TenantConn
	linq.From(a.Conn).SelectT(func(c biz.TenantConn) TenantConn {
		return TenantConn{
			TenantId:  c.TenantId,
			Key:       c.Key,
			Value:     c.Value,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}
	}).ToSlice(&conn)

	var features []TenantFeature

	linq.From(a.Features).SelectT(func(c biz.TenantFeature) TenantFeature {
		return TenantFeature{
			TenantId:  c.TenantId,
			Key:       c.Key,
			Value:     c.Value,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}
	}).ToSlice(&features)
	var id uuid.UUID
	if a.ID != "" {
		id = uuid.MustParse(a.ID)
	}
	b.UIDBase = gorm2.UIDBase{ID: id}
	b.Name = a.Name
	b.DisplayName = a.DisplayName
	b.Region = a.Region
	b.CreatedAt = a.CreatedAt
	b.UpdatedAt = a.UpdatedAt
	b.DeletedAt = a.DeletedAt
	b.Conn = conn
	b.Features = features
}

func mapDataTenantToBizTenant(a *Tenant, b *biz.Tenant) {
	var conn []biz.TenantConn
	linq.From(a.Conn).SelectT(func(c TenantConn) biz.TenantConn {
		return biz.TenantConn{
			TenantId:  c.TenantId,
			Key:       c.Key,
			Value:     c.Value,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}
	}).ToSlice(&conn)

	var features []biz.TenantFeature

	linq.From(a.Features).SelectT(func(c TenantFeature) biz.TenantFeature {
		return biz.TenantFeature{
			TenantId:  c.TenantId,
			Key:       c.Key,
			Value:     c.Value,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}
	}).ToSlice(&features)
	b.ID = a.ID.String()
	b.Name = a.Name
	b.DisplayName = a.DisplayName
	b.Region = a.Region
	b.CreatedAt = a.CreatedAt
	b.UpdatedAt = a.UpdatedAt
	b.DeletedAt = a.DeletedAt
	b.Conn = conn
	b.Features = features
}
