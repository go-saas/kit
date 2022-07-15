package biz

import (
	"context"
	"github.com/go-saas/kit/pkg/data"
	gorm2 "github.com/go-saas/kit/pkg/gorm"
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
	gorm2.AuditedModel
	UserId   string           `gorm:"type:char(36);primary_key" json:"user_id"`
	TenantId string           `gorm:"type:char(36);primary_key" json:"tenant_id" `
	JoinTime time.Time        `json:"join_time"`
	Status   UserTenantStatus `json:"status;index"`
	Extra    data.JSONMap
}

func (u *UserTenant) SetTenantId(id string) *UserTenant {
	u.TenantId = id
	return u
}

func (u *UserTenant) GetTenantId() string {
	return u.TenantId
}

type UserTenantRepo interface {
	JoinTenant(ctx context.Context, userId string, tenantId string, status UserTenantStatus) (*UserTenant, error)
	RemoveFromTenant(ctx context.Context, userId string, tenantId string) error
	Get(ctx context.Context, userId string, tenantId string) (*UserTenant, error)
	IsIn(ctx context.Context, userId string, tenantId string) (bool, error)
	Update(ctx context.Context, userTenant *UserTenant) error
}
