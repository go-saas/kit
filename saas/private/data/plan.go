package data

import (
	"context"
	"fmt"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	v1 "github.com/go-saas/kit/saas/api/plan/v1"
	"github.com/go-saas/kit/saas/private/biz"
	"github.com/goxiaoy/go-eventbus"
	"gorm.io/gorm"
)

type PlanRepo struct {
	*kitgorm.Repo[biz.Plan, string, *v1.ListPlanRequest]
}

func NewPlanRepo(eventbus *eventbus.EventBus, data *Data) biz.PlanRepo {
	res := &PlanRepo{}
	res.Repo = kitgorm.NewRepo[biz.Plan, string, *v1.ListPlanRequest](data.DbProvider, eventbus, res)
	return res
}

func (g *PlanRepo) GetDb(ctx context.Context) *gorm.DB {
	ret := GetDb(ctx, g.DbProvider)
	return ret
}

func (g *PlanRepo) BuildDetailScope(withDetail bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("Plan")
	}
}
func (g *PlanRepo) DefaultSorting() []string {
	return []string{
		"-created_at",
	}
}

func (g *PlanRepo) BuildFilterScope(q *v1.ListPlanRequest) func(db *gorm.DB) *gorm.DB {
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
		ret = ret.Scopes(kitgorm.BuildStringFilter("`key`", filter.Key))
		ret = ret.Scopes(kitgorm.BuildStringFilter("`display_name`", filter.DisplayName))
		ret = ret.Scopes(kitgorm.BuildBooleanFilter("`active`", filter.Active))
		return ret
	}
}

func (g *PlanRepo) BuildPrimaryField() string {
	return "`key`"
}
