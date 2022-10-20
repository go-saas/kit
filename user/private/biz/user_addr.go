package biz

import (
	"context"
	"github.com/go-saas/kit/pkg/data"
)

type UserAddress struct {
	data.UIDBase
	data.AuditedModel
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
