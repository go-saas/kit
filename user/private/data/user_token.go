package data

import (
	"context"
	"errors"
	"github.com/go-saas/kit/user/private/biz"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserTokenRepo struct {
	Repo
}

func NewUserTokenRepo(data *Data) biz.UserTokenRepo {
	return &UserTokenRepo{
		Repo{
			DbProvider: data.DbProvider,
		},
	}
}

func (u *UserTokenRepo) FindByUserIdAndLoginProvider(ctx context.Context, userId, loginProvider string) ([]*biz.UserToken, error) {
	var entity []*biz.UserToken
	err := u.GetDb(ctx).Model(&biz.UserToken{}).Find(entity, "user_id = ? and login_provider = ? ", userId, loginProvider).Error
	return entity, err
}

func (u *UserTokenRepo) FindByUserIdAndLoginProviderAndName(ctx context.Context, userId, loginProvider, name string) (*biz.UserToken, error) {
	entity := &biz.UserToken{}
	err := u.GetDb(ctx).Model(&biz.UserToken{}).First(entity, "user_id = ? and login_provider = ? and name = ?", userId, loginProvider, name).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return entity, nil
}

func (u *UserTokenRepo) DeleteByUserIdAndLoginProvider(ctx context.Context, userId, loginProvider string) error {
	err := u.GetDb(ctx).Model(&biz.UserToken{}).Delete(&biz.UserToken{}, "user_id = ? and login_provider = ? ", userId, loginProvider).Error
	return err
}

func (u *UserTokenRepo) DeleteByUserIdAndLoginProviderAndName(ctx context.Context, userId, loginProvider, name string) error {
	err := u.GetDb(ctx).Model(&biz.UserToken{}).Delete(&biz.UserToken{}, "user_id = ? and login_provider = ? and name = ?", userId, loginProvider, name).Error
	return err
}

func (u *UserTokenRepo) Create(ctx context.Context, userId, loginProvider, name, value string) (*biz.UserToken, error) {
	if err := u.DeleteByUserIdAndLoginProviderAndName(ctx, userId, loginProvider, name); err != nil {
		return nil, err
	}
	entity := &biz.UserToken{UserId: uuid.MustParse(userId), LoginProvider: loginProvider, Name: name, Value: value}
	err := u.GetDb(ctx).Model(&biz.UserToken{}).Create(entity).Error
	return entity, err
}
