package biz

import (
	"github.com/goxiaoy/go-saas-kit/pkg/gorm"
	concurrency "github.com/goxiaoy/gorm-concurrency"
)

type UserInfo struct {
	gorm.UIDBase
	concurrency.Version

	gorm.AuditedModel

	UserId string `json:"user_id" gorm:"index"`
}
