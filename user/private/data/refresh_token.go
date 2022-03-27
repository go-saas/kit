package data

import (
	"context"
	"errors"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	"gorm.io/gorm"
	"time"
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

func (r *RefreshTokenRepo) Find(ctx context.Context, token string, validOnly bool) (*biz.RefreshToken, error) {
	var t biz.RefreshToken
	q := r.GetDb(ctx).Model(&biz.RefreshToken{}).Where("token = ?", token)

	if validOnly {
		nowTime := time.Now()
		q = q.Where("expires is NULL or expires > ?", nowTime)
	}
	err := q.First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *RefreshTokenRepo) Revoke(ctx context.Context, token string) (err error) {
	err = r.GetDb(ctx).Delete(&biz.RefreshToken{}, "token = ? ", token).Error
	return
}
