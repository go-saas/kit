package biz

import (
	"github.com/goxiaoy/go-saas-kit/pkg/data"
	"github.com/goxiaoy/go-saas-kit/pkg/gorm"
)

// UserSetting contains key/value pair of user settings
type UserSetting struct {
	gorm.UIDBase
	UserId string     `json:"user_id"`
	Key    string     `json:"key" gorm:"index"`
	Value  data.Value `gorm:"embedded"`
}
