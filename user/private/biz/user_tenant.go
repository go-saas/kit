package biz

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/data"
	gorm2 "github.com/goxiaoy/go-saas-kit/pkg/gorm"
	"github.com/goxiaoy/go-saas/gorm"
	"time"
)

type UserTenantStatus int32

const (
	Active   UserTenantStatus = 0
	Inactive UserTenantStatus = 1
)

func (p UserTenantStatus) String() string {
	switch p {
	case Active:
		return "ACTIVE"
	case Inactive:
		return "INACTIVE"
	default:
		return "UNKNOWN"
	}
}

type UserTenant struct {
	gorm2.UIDBase
	UserId   string           `gorm:"type:char(36)" json:"user_id"`
	TenantId gorm.HasTenant   `json:"tenant_id" gorm:"type:char(36)"`
	JoinTime time.Time        `json:"join_time"`
	Status   UserTenantStatus `json:"status"`
	Extra    data.JSONMap
}

type UserTenantRepo interface {
	JoinTenant(ctx context.Context, userId string, tenantId string) (*UserTenant, error)
	RemoveFromTenant(ctx context.Context, userId string, tenantId string) error
	Get(ctx context.Context, userId string, tenantId string) (*UserTenant, error)
	IsIn(ctx context.Context, userId string, tenantId string) (bool, error)
	Update(ctx context.Context, userTenant *UserTenant) error
}
