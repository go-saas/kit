package biz

import (
	"context"
	"github.com/go-saas/kit/pkg/data"
	"github.com/go-saas/kit/pkg/gorm"
)

type UserAddress struct {
	gorm.UIDBase
	gorm.AuditedModel
	UserId   string             `json:"user_id" gorm:"index:,type:char(36)"`
	Phone    string             `json:"phone"`
	Usage    string             `json:"usage"`
	Prefer   bool               `json:"prefer"`
	Address  data.AddressEntity `json:"address" gorm:"embedded"`
	Metadata data.JSONMap       `json:"metadata"`
}

type UserAddressRepo interface {
	data.Repo[UserAddress, string, interface{}]
	FindByUser(ctx context.Context, userId string) ([]*UserAddress, error)
	SetPrefer(ctx context.Context, addr *UserAddress) error
}
