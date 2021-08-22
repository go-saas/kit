package biz

import (
	concurrency "github.com/goxiaoy/gorm-concurrency"
	"github.com/goxiaoy/go-saas-kit/pkg/gorm"
)

type UserInfo struct {
	gorm.UIDBase
	concurrency.Version

	gorm.AuditedModel

	UserId string `json:"user_id" gorm:"index"`
}
