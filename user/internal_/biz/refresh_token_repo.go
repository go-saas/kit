package biz

import "context"

type RefreshTokenRepo interface {
	Create(ctx context.Context, t *RefreshToken) (err error)
	FindUser(ctx context.Context, token string) (uId string, err error)
	// Refresh delete previous token and generate a new token
	Refresh(ctx context.Context, pre *RefreshToken, now *RefreshToken) (err error)
	// Revoke refresh token
	Revoke(ctx context.Context, token string) (err error)
}
