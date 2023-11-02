package data

import (
	"context"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/kit/user/private/biz"
	"github.com/goxiaoy/go-eventbus"
	"gorm.io/gorm"
)

type UserAddrRepo struct {
	*kitgorm.Repo[biz.UserAddress, string, interface{}]
}

func NewUserAddrRepo(data *Data, eventbus *eventbus.EventBus) biz.UserAddressRepo {
	res := &UserAddrRepo{}
	res.Repo = kitgorm.NewRepo[biz.UserAddress, string, interface{}](data.DbProvider, eventbus, res)
	return res
}

func (u *UserAddrRepo) GetDb(ctx context.Context) *gorm.DB {
	return GetDb(ctx, u.DbProvider)
}

func (u *UserAddrRepo) FindByUser(ctx context.Context, userId string) ([]*biz.UserAddress, error) {
	var e biz.UserAddress
	db := u.GetDb(ctx).Model(&e)
	db = db.Scopes(kitgorm.WhereUserId(userId)).Order("prefer desc,updated_at desc")
	var items []*biz.UserAddress
	res := db.Find(&items)
	return items, res.Error
}

func (u *UserAddrRepo) SetPrefer(ctx context.Context, addr *biz.UserAddress) error {
	db := u.GetDb(ctx)
	//set other prefer as no
	if err := db.Model(&biz.UserAddress{}).Scopes(kitgorm.WhereUserId(addr.UserId.String())).Update("prefer", false); err != nil {
		return nil
	}
	if err := db.Model(addr).Update("prefer", true); err != nil {
		return nil
	}
	return nil
}
