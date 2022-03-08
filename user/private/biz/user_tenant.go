package biz

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/data"
	gorm2 "github.com/goxiaoy/go-saas-kit/pkg/gorm"
	gg "gorm.io/gorm"
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
	UserId    string           `gorm:"type:char(36)" json:"user_id"`
	TenantId  *string          `json:"tenant_id" gorm:"type:char(36)"`
	JoinTime  time.Time        `json:"join_time"`
	Status    UserTenantStatus `json:"status"`
	DeletedAt gg.DeletedAt     `gorm:"index"`
	Extra     data.JSONMap
}

func (u *UserTenant) SetTenantId(id string) *UserTenant {
	if len(id) == 0 {
		u.TenantId = nil
	} else {
		u.TenantId = &id
	}
	return u
}

func (u *UserTenant) GetTenantId() string {
	if u.TenantId == nil {
		return ""
	}
	return *u.TenantId
}

type UserTenantRepo interface {
	JoinTenant(ctx context.Context, userId string, tenantId string, status UserTenantStatus) (*UserTenant, error)
	RemoveFromTenant(ctx context.Context, userId string, tenantId string) error
	Get(ctx context.Context, userId string, tenantId string) (*UserTenant, error)
	IsIn(ctx context.Context, userId string, tenantId string) (bool, error)
	Update(ctx context.Context, userTenant *UserTenant) error
}
