package data

import (
	"context"
	"errors"
	"fmt"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/kit/pkg/query"
	v1 "github.com/go-saas/kit/saas/api/tenant/v1"
	"github.com/go-saas/kit/saas/private/biz"
	"github.com/google/uuid"
	"github.com/goxiaoy/go-eventbus"
	"gorm.io/gorm"
)

type TenantRepo struct {
	*kitgorm.Repo[biz.Tenant, string, *v1.ListTenantRequest]
}

func NewTenantRepo(eventbus *eventbus.EventBus, data *Data) biz.TenantRepo {
	res := &TenantRepo{}
	res.Repo = kitgorm.NewRepo[biz.Tenant, string, *v1.ListTenantRequest](data.DbProvider, eventbus, res)
	return res
}

func (g *TenantRepo) GetDb(ctx context.Context) *gorm.DB {
	ret := GetDb(ctx, g.DbProvider)
	return ret
}

func (g *TenantRepo) BuildDetailScope(withDetail bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Preload("Plan")
		if withDetail {
			return db.Preload("Conn").Preload("Features")
		}
		return db
	}
}
func (g *TenantRepo) DefaultSorting() []string {
	return []string{
		"-created_at",
	}
}

func (g *TenantRepo) BuildFilterScope(q *v1.ListTenantRequest) func(db *gorm.DB) *gorm.DB {
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
		ret = ret.Scopes(kitgorm.BuildStringFilter("`id`", filter.Id))
		ret = ret.Scopes(kitgorm.BuildStringFilter("`name`", filter.Name))
		ret = ret.Scopes(kitgorm.BuildStringFilter("`region`", filter.Region))
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
	var tDb biz.Tenant
	if err == nil {
		//id
		err = g.GetDb(ctx).Model(&biz.Tenant{}).Scopes(g.BuildDetailScope(true)).Where("id = ?", id.String()).First(&tDb).Error
	} else {
		err = g.GetDb(ctx).Model(&biz.Tenant{}).Scopes(g.BuildDetailScope(true)).Where("name = ?", idOrName).First(&tDb).Error
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tDb, err
}

func (g *TenantRepo) UpdateAssociation(ctx context.Context, entity *biz.Tenant, p query.Select) error {
	if query.SelectContains(p, "Conn") {
		if err := g.GetDb(ctx).Model(entity).Session(&gorm.Session{FullSaveAssociations: true}).Association("Conn").Replace(entity.Conn); err != nil {
			return err
		}
	}
	if query.SelectContains(p, "Features") {
		if err := g.GetDb(ctx).Model(entity).Session(&gorm.Session{FullSaveAssociations: true}).Association("Features").Replace(entity.Features); err != nil {
			return err
		}
	}
	return nil
}
