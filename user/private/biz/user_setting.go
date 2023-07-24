package biz

import (
	"context"
	"github.com/go-saas/kit/pkg/data"
	"github.com/go-saas/kit/pkg/gorm"
	v1 "github.com/go-saas/kit/user/api/account/v1"
)

// UserSetting contains key/value pairs of user settings
type UserSetting struct {
	gorm.UIDBase
	UserId string     `json:"user_id" gorm:"index"`
	Key    string     `json:"key" gorm:"index"`
	Value  data.Value `gorm:"embedded"`
}

type UpdateUserSetting struct {
	Key string
	//Value new value
	Value  *data.Value
	Delete bool
}

type UserSettingRepo interface {
	data.Repo[UserSetting, string, *v1.GetSettingsRequest]
	FindByUser(ctx context.Context, userId string, query *v1.GetSettingsRequest) ([]*UserSetting, error)
	UpdateByUser(ctx context.Context, userId string, updateBatch []UpdateUserSetting) error
}
