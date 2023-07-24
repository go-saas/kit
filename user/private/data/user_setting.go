package data

import (
	"context"
	v1 "github.com/go-saas/kit/user/api/account/v1"
	"github.com/goxiaoy/go-eventbus"
	"github.com/samber/lo"

	kitgorm "github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/kit/user/private/biz"
	"gorm.io/gorm"
)

type UserSettingRepo struct {
	*kitgorm.Repo[biz.UserSetting, string, *v1.GetSettingsRequest]
}

func NewUserSettingRepo(data *Data, eventbus *eventbus.EventBus) biz.UserSettingRepo {
	res := &UserSettingRepo{}
	res.Repo = kitgorm.NewRepo[biz.UserSetting, string, *v1.GetSettingsRequest](data.DbProvider, eventbus, res)
	return res
}

func (r *UserSettingRepo) GetDb(ctx context.Context) *gorm.DB {
	return GetDb(ctx, r.DbProvider)
}

// BuildFilterScope filter
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

func (r *UserSettingRepo) UpdateByUser(ctx context.Context, userId string, updateBatch []biz.UpdateUserSetting) error {
	var e biz.UserSetting
	dels := lo.Map(updateBatch, func(t biz.UpdateUserSetting, _ int) string {
		return t.Key
	})
	db := r.GetDb(ctx).Model(&e)
	if err := db.Delete(&e, "user_id = ? AND `key` in (?)", userId, dels).Error; err != nil {
		return err
	}
	updates := lo.Map(lo.Filter(updateBatch, func(v biz.UpdateUserSetting, _ int) bool { return !v.Delete }), func(t biz.UpdateUserSetting, _ int) *biz.UserSetting {
		return &biz.UserSetting{
			UserId: userId,
			Key:    t.Key,
			Value:  *t.Value,
		}
	})
	if err := db.CreateInBatches(updates, 100).Error; err != nil {
		return err
	}
	return nil
}
