package data

import (
	"context"
	"github.com/go-saas/kit/pkg/authn"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	v1 "github.com/go-saas/kit/realtime/api/notification/v1"
	"github.com/go-saas/kit/realtime/private/biz"
	sgorm "github.com/go-saas/saas/gorm"
	"github.com/goxiaoy/go-eventbus"
	"gorm.io/gorm"
)

type NotificationRepo struct {
	*kitgorm.Repo[biz.Notification, string, *v1.ListNotificationRequest]
}

func NewNotificationRepo(dbProvider sgorm.DbProvider, eventbus *eventbus.EventBus) biz.NotificationRepo {
	res := &NotificationRepo{}
	res.Repo = kitgorm.NewRepo[biz.Notification, string, *v1.ListNotificationRequest](dbProvider, eventbus, res)
	return res
}

func (c *NotificationRepo) GetDb(ctx context.Context) *gorm.DB {
	return GetDb(ctx, c.Repo.DbProvider)
}

// BuildDetailScope preload relations
func (c *NotificationRepo) BuildDetailScope(withDetail bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

// BuildFilterScope filter
func (c *NotificationRepo) BuildFilterScope(q *v1.ListNotificationRequest) func(db *gorm.DB) *gorm.DB {
	filter := q.Filter
	return func(db *gorm.DB) *gorm.DB {
		ret := db

		u, _ := authn.FromUserContext(db.Statement.Context)
		ret = ret.Where("`user_id` = ?", u.GetId())
		if filter == nil {
			return ret
		}

		if filter.HasRead != nil {
			ret = ret.Scopes(kitgorm.BuildBooleanFilter("`has_read`", filter.HasRead))
		}
		return ret
	}
}

func (c *NotificationRepo) DefaultSorting() []string {
	return []string{"created_at"}
}

func (c *NotificationRepo) MyUnreadCount(ctx context.Context) (int32, error) {
	db := c.GetDb(ctx).Model(&biz.Notification{})
	u, err := authn.ErrIfUnauthenticated(db.Statement.Context)
	if err != nil {
		return 0, err
	}
	var count int64
	err = db.Where("`user_id` = ? AND has_read = ?", u.GetId(), false).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}
func (c *NotificationRepo) SetMyRead(ctx context.Context, idFilter string) error {
	db := c.GetDb(ctx).Model(&biz.Notification{})
	u, err := authn.ErrIfUnauthenticated(db.Statement.Context)
	if err != nil {
		return err
	}
	db = db.Where("`user_id` = ?", u.GetId())
	if len(idFilter) > 0 {
		db = db.Where("id = ?", idFilter)
	}
	return db.Update("has_read", true).Error
}

func (c *NotificationRepo) DeleteMy(ctx context.Context, id string) error {
	db := c.GetDb(ctx).Model(&biz.Notification{})
	u, err := authn.ErrIfUnauthenticated(db.Statement.Context)
	if err != nil {
		return err
	}
	return db.Delete(&biz.Notification{}, "`user_id` = ? AND id = ?", u.GetId(), id).Error
}
