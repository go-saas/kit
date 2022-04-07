package biz

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	Token     string    `gorm:"type:char(36);primaryKey"`
	UserId    uuid.UUID `gorm:"type:char(36);index" json:"user_id"`
	Expires   *time.Time
	Ip        string
	UserAgent string
	Used      bool
}

func NewRefreshToken(userId uuid.UUID, duration time.Duration, userAgent string, ip string) *RefreshToken {
	var t *time.Time
	if duration > 0 {
		e := time.Now().Add(duration)
		t = &e
	}
	return &RefreshToken{
		Token:     uuid.New().String(),
		UserId:    userId,
		Expires:   t,
		Ip:        ip,
		UserAgent: userAgent,
	}
}

func (r *RefreshToken) Valid() bool {
	if r.Expires == nil {
		return true
	}
	t := *(r.Expires)
	return t.After(time.Now())
}

type RefreshTokenRepo interface {
	Create(ctx context.Context, t *RefreshToken) (err error)
	Find(ctx context.Context, token string, validOnly bool) (*RefreshToken, error)
	// Revoke refresh token
	Revoke(ctx context.Context, token string, used bool) (err error)
}
