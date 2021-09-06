package biz

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"time"
)

type UserTenantStatus int32

const (
	Active      UserTenantStatus = 0
	Inactive      UserTenantStatus = 1
)

func (p UserTenantStatus) String() string {
	switch p {
	case Active: return "ACTIVE"
	case Inactive: return "INACTIVE"
	default:         return "UNKNOWN"
	}
}

type UserTenant struct {
	UserId    uuid.UUID `gorm:"type:char(36);primaryKey" json:"user_id"`
	TenantId  string    `json:"tenant_id" gorm:"type:char(36);primaryKey"`
	JoinTime time.Time `json:"joint_time"`
	Status    UserTenantStatus    `json:"status"`
	Extra datatypes.JSONMap
}

type UserTenantRepo interface {
	JoinTenant(userId string,tenantId string)(*UserTenant,error)
	RemoveFromTenant(userId string,tenantId string)error
	Get(userId string,tenantId string)(*UserTenant,error)
	Update(userTenant *UserTenant)error
}