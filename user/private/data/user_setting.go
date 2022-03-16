package data

import (
	"context"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/account/v1"

	kitgorm "github.com/goxiaoy/go-saas-kit/pkg/gorm"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	"gorm.io/gorm"
)

type UserSettingRepo struct {
	*kitgorm.Repo[biz.UserSetting, string, v1.GetSettingsRequest]
}

func NewUserSettingRepo(data *Data) biz.UserSettingRepo {
	return &UserSettingRepo{
		Repo: kitgorm.NewRepo[biz.UserSetting, string, v1.GetSettingsRequest](data.DbProvider),
	}
}

func (r *UserSettingRepo) GetDb(ctx context.Context) *gorm.DB {
	return GetDb(ctx, r.DbProvider)
}

//BuildFilterScope filter
func (r *UserSettingRepo) BuildFilterScope(q *v1.GetSettingsRequest) func(db *gorm.DB) *gorm.DB {
	filter := q.Filter
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		ret := db
		if filter.Key != nil {
			ret = ret.Scopes(kitgorm.BuildStringFilter("`key`", filter.Key))
		}
		return ret
	}
}

func (r *UserSettingRepo) FindByUser(ctx context.Context, userId string, query *v1.GetSettingsRequest) ([]*biz.UserSetting, error) {
	var e biz.UserSetting
	db := r.GetDb(ctx).Model(&e)
	db = db.Scopes(kitgorm.WhereUserId(userId), r.BuildFilterScope(query))
	var items []*biz.UserSetting
	res := db.Find(&items)
	return items, res.Error
}
