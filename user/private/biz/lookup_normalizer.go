package biz

import (
	"context"
	"github.com/go-saas/kit/pkg/conf"
	v1 "github.com/go-saas/kit/user/api/user/v1"
	"github.com/nyaruka/phonenumbers"
	"net/mail"
	"strings"
)

type LookupNormalizer interface {
	// Name normalizer
	Name(ctx context.Context, name string) (string, error)
	// Email normalizer
	Email(ctx context.Context, email string) (string, error)
	// Phone normalizer
	Phone(ctx context.Context, phone string) (string, error)
}

type lookupNormalizer struct {
	phoneDefaultRegion string
}

func NewLookupNormalizer(cfg *conf.AppConfig) LookupNormalizer {
	ret := &lookupNormalizer{phoneDefaultRegion: "CN"}
	if cfg != nil && cfg.PhoneDefaultRegion != nil {
		ret.phoneDefaultRegion = *cfg.PhoneDefaultRegion
	}
	return ret
}
func (l *lookupNormalizer) Name(ctx context.Context, name string) (string, error) {
	if name == "" {
		return "", v1.ErrorInvalidUsernameLocalized(ctx, nil, nil)
	}
	if _, err := l.Email(ctx, name); err == nil {
		return "", v1.ErrorInvalidUsernameLocalized(ctx, nil, nil)
	}
	if _, err := l.Phone(ctx, name); err == nil {
		return "", v1.ErrorInvalidUsernameLocalized(ctx, nil, nil)
	}
	return strings.ToUpper(name), nil
}

func (l *lookupNormalizer) Email(ctx context.Context, email string) (string, error) {
	if email == "" {
		return "", v1.ErrorInvalidEmailLocalized(ctx, nil, nil)
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return "", v1.ErrorInvalidEmailLocalized(ctx, nil, nil)
	}
	return strings.ToUpper(email), nil
}

func (l *lookupNormalizer) Phone(ctx context.Context, phone string) (string, error) {
	if phone == "" {
		return "", v1.ErrorInvalidPhoneLocalized(ctx, nil, nil)
	}
	num, err := phonenumbers.Parse(phone, l.phoneDefaultRegion)
	if err != nil {
		return "", v1.ErrorInvalidPhoneLocalized(ctx, nil, nil)
	}
	if ok := phonenumbers.IsValidNumber(num); !ok {
		return "", v1.ErrorInvalidPhoneLocalized(ctx, nil, nil)
	}
	formattedNum := phonenumbers.Format(num, phonenumbers.E164)
	return formattedNum, err
}
