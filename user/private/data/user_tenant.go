package data

import (
	"context"
	"errors"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
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
	t := &biz.UserTenant{
		UserId:   userId,
		TenantId: tenantId,
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
	return u.GetDb(ctx).Delete(&biz.UserTenant{}, "user_id = ? and tenant_id = ?", userId, tenantId).Error
}

func (u *UserTenantRepo) Get(ctx context.Context, userId string, tenantId string) (*biz.UserTenant, error) {
	t := &biz.UserTenant{}
	err := u.GetDb(ctx).Where("user_id = ? and tenant_id = ?", userId, tenantId).First(t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return t, nil
}

func (u *UserTenantRepo) Update(ctx context.Context, userTenant *biz.UserTenant) error {
	err := u.GetDb(ctx).Where("user_id = ? and tenant_id = ?", userTenant.UserId, userTenant.TenantId).Updates(userTenant).Error
	return err
}
