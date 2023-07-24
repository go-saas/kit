package data

import (
	"context"
	"errors"
	"github.com/goxiaoy/go-eventbus"

	kitgorm "github.com/go-saas/kit/pkg/gorm"
	v12 "github.com/go-saas/kit/user/api/role/v1"
	"github.com/go-saas/kit/user/private/biz"
	"gorm.io/gorm"
)

type RoleRepo struct {
	*kitgorm.Repo[biz.Role, string, *v12.ListRolesRequest]
}

func NewRoleRepo(data *Data, eventbus *eventbus.EventBus) biz.RoleRepo {
	res := &RoleRepo{}
	res.Repo = kitgorm.NewRepo[biz.Role, string, *v12.ListRolesRequest](data.DbProvider, eventbus, res)
	return res
}

func (r *RoleRepo) GetDb(ctx context.Context) *gorm.DB {
	return GetDb(ctx, r.DbProvider)
}

// BuildFilterScope filter
func (r *RoleRepo) BuildFilterScope(q *v12.ListRolesRequest) func(db *gorm.DB) *gorm.DB {
	filter := q.Filter
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		ret := db
		ret = ret.Scopes(kitgorm.BuildStringFilter("`id`", filter.Id))
		ret = ret.Scopes(kitgorm.BuildStringFilter("`name`", filter.Name))
		return ret
	}
}

// DefaultSorting get default sorting
func (r *RoleRepo) DefaultSorting() []string {
	return []string{"-created_at"}
}

func (r *RoleRepo) FindByName(ctx context.Context, name string) (*biz.Role, error) {
	db := r.GetDb(ctx).Model(&biz.Role{})
	var item = &biz.Role{}
	if err := db.Where("normalized_name = ?", name).First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}
