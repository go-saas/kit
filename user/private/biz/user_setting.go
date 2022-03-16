package biz

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/data"
	"github.com/goxiaoy/go-saas-kit/pkg/gorm"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/account/v1"
)

// UserSetting contains key/value pair of user settings
type UserSetting struct {
	gorm.UIDBase
	UserId string     `json:"user_id"`
	Key    string     `json:"key" gorm:"index"`
	Value  data.Value `gorm:"embedded"`
}

type UserSettingRepo interface {
	data.Repo[UserSetting, string, v1.GetSettingsRequest]
	FindByUser(ctx context.Context, userId string, query *v1.GetSettingsRequest) ([]*UserSetting, error)
}
