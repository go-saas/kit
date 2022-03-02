package data

import (
	"context"
	"errors"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	gorm2 "github.com/goxiaoy/go-saas/gorm"
	"gorm.io/gorm"
	"time"
)

type UserTenantRepo struct {
	Repo
}

func NewUserTenantRepo(data *Data) biz.UserTenantRepo {
	return &UserTenantRepo{
		Repo{
			DbProvider: data.DbProvider,
		},
	}
}

func (u *UserTenantRepo) JoinTenant(ctx context.Context, userId string, tenantId string) (*biz.UserTenant, error) {
	if ut, err := u.Get(ctx, userId, tenantId); err != nil {
		return nil, err
	} else if ut != nil {
		//already in
		return ut, nil
	}
	//not present
	t := &biz.UserTenant{
		UserId:   userId,
		TenantId: gorm2.NewTenantId(tenantId),
		JoinTime: time.Now(),
		Status:   biz.Active,
		Extra:    nil,
	}
	if err := u.GetDb(ctx).Save(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (u *UserTenantRepo) IsIn(ctx context.Context, userId string, tenantId string) (bool, error) {
	ut, err := u.Get(ctx, userId, tenantId)
	if err != nil {
		return false, err
	}
	return ut != nil, nil
}

func (u *UserTenantRepo) RemoveFromTenant(ctx context.Context, userId string, tenantId string) error {
	var err error
	if len(tenantId) > 0 {
		err = u.GetDb(ctx).Delete(&biz.UserTenant{}, "user_id = ? and tenant_id = ?", userId, tenantId).Error
	} else {
		err = u.GetDb(ctx).Delete(&biz.UserTenant{}, "user_id = ? and tenant_id is NULL ?", userId).Error
	}
	return err
}

func (u *UserTenantRepo) Get(ctx context.Context, userId string, tenantId string) (*biz.UserTenant, error) {
	t := &biz.UserTenant{}
	var err error
	if len(tenantId) > 0 {
		err = u.GetDb(ctx).Where("user_id = ? and tenant_id = ?", userId, tenantId).First(t).Error
	} else {
		err = u.GetDb(ctx).Where("user_id = ? and tenant_id is NULL", userId).First(t).Error
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return t, nil
}

func (u *UserTenantRepo) Update(ctx context.Context, userTenant *biz.UserTenant) error {
	return u.GetDb(ctx).Updates(userTenant).Error
}
