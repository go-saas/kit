package data

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/user/internal_/biz"
)

type RefreshTokenRepo struct {
	Repo
}

var _ biz.RefreshTokenRepo = (*RefreshTokenRepo)(nil)

func NewRefreshTokenRepo(data *Data) biz.RefreshTokenRepo {
	return &RefreshTokenRepo{
		Repo{DbProvider: data.DbProvider},
	}
}

func (r *RefreshTokenRepo) Create(ctx context.Context, t *biz.RefreshToken) (err error) {
	err = r.GetDb(ctx).Create(t).Error
	return
}

func (r *RefreshTokenRepo) FindUser(ctx context.Context, token string) (uId string, err error) {
	var t biz.RefreshToken
	err = r.GetDb(ctx).Model(&biz.RefreshToken{}).First(&t, "token=?", token).Error
	uId = t.UserId.String()
	return
}

func (r *RefreshTokenRepo) Refresh(ctx context.Context, pre *biz.RefreshToken, now *biz.RefreshToken) (err error) {
	err = r.GetDb(ctx).Model(&biz.RefreshToken{}).Delete(pre).Error
	if err != nil {
		return
	}
	err = r.Create(ctx, now)
	return
}

func (r *RefreshTokenRepo) Revoke(ctx context.Context, token string) (err error) {
	err = r.GetDb(ctx).Delete(&biz.RefreshToken{}, token).Error
	return
}
