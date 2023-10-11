package biz

import (
	"context"
	"fmt"
	"github.com/eko/gocache/v3/cache"
	"github.com/eko/gocache/v3/store"
	cache2 "github.com/go-saas/kit/pkg/cache"
	v1 "github.com/go-saas/kit/user/api/auth/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"time"
)

type UserTokenProvider interface {
	// Name of this provider
	Name() string
	Generate(ctx context.Context, purpose TokenPurpose, user *User, duration time.Duration) (token string, err error)
	Validate(ctx context.Context, purpose TokenPurpose, token string, user *User) (bool, error)
	CanGenerate(ctx context.Context, user *User) error
}

type TokenPurpose string

const (
	EmailName   = "email"
	PhoneName   = "phone"
	TwoStepName = "twostep"
	OtpName     = "otp"

	EmailLoginPurpose = "login_email"
	PhoneLoginPurpose = "login_phone"

	ConfirmPurpose               TokenPurpose = "confirm"
	RecoverPurpose               TokenPurpose = "recover"
	RecoverChangePasswordPurpose TokenPurpose = "recover_change_password"
)

var (
	_ UserTokenProvider = (*PhoneTokenProvider)(nil)
	_ UserTokenProvider = (*EmailTokenProvider)(nil)
)

type PhoneTokenProvider struct {
	r cache.CacheInterface[string]
}

func NewPhoneTokenProvider(r cache.CacheInterface[string]) *PhoneTokenProvider {
	return &PhoneTokenProvider{r: r}
}

func (p *PhoneTokenProvider) Name() string {
	return PhoneName
}

func (p *PhoneTokenProvider) Generate(ctx context.Context, purpose TokenPurpose, user *User, duration time.Duration) (string, error) {
	if err := p.CanGenerate(ctx, user); err != nil {
		return "", err
	}
	key := fmt.Sprintf("usertoken:%s:%s:%s", user.ID.String(), purpose, *user.Phone)
	token, err := GenerateOtp()
	if err != nil {
		return "", err
	}
	err = p.r.Set(ctx, key, token, store.WithExpiration(duration))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (p *PhoneTokenProvider) Validate(ctx context.Context, purpose TokenPurpose, token string, user *User) (bool, error) {
	key := fmt.Sprintf("usertoken:%s:%s:%s", user.ID.String(), purpose, *user.Phone)
	val, err := p.r.Get(ctx, key)
	if err != nil {
		if (store.NotFound{}).Is(err) {
			return false, nil
		}
		return false, err
	}
	return val == token, nil

}

func (p *PhoneTokenProvider) CanGenerate(ctx context.Context, user *User) error {
	if user.Phone != nil && user.PhoneConfirmed {
		return nil
	}
	return v1.ErrorPhoneNotConfirmedLocalized(ctx, nil, nil)
}

type EmailTokenProvider struct {
	r cache.CacheInterface[string]
}

func NewEmailTokenProvider(r cache.CacheInterface[string]) *EmailTokenProvider {
	return &EmailTokenProvider{r: r}
}

func (e *EmailTokenProvider) Name() string {
	return EmailName
}

func (e *EmailTokenProvider) Generate(ctx context.Context, purpose TokenPurpose, user *User, duration time.Duration) (string, error) {
	if err := e.CanGenerate(ctx, user); err != nil {
		return "", err
	}
	key := fmt.Sprintf("usertoken:%s:%s:%s", user.ID.String(), purpose, *user.NormalizedEmail)
	var token string
	var err error
	if purpose == ConfirmPurpose {
		token = uuid.New().String()
	} else {
		token, err = GenerateOtp()
	}
	if err != nil {
		return "", err
	}
	err = e.r.Set(ctx, key, token, store.WithExpiration(duration))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (e *EmailTokenProvider) Validate(ctx context.Context, purpose TokenPurpose, token string, user *User) (bool, error) {
	key := fmt.Sprintf("usertoken:%s:%s:%s", user.ID.String(), purpose, *user.NormalizedEmail)
	val, err := e.r.Get(ctx, key)
	if err != nil {
		if (store.NotFound{}).Is(err) {
			return false, nil
		}

		return false, err
	}
	return val == token, nil
}

func (e *EmailTokenProvider) CanGenerate(ctx context.Context, user *User) error {
	if user.Email != nil && user.EmailConfirmed {
		return nil
	}
	return v1.ErrorEmailNotConfirmedLocalized(ctx, nil, nil)
}

type TwoStepTokenProvider[T proto.Message] struct {
	c     *cache2.ProtoCache[T]
	proxy cache.CacheInterface[string]
}

func NewTwoStepTokenProvider[T proto.Message](creator func() T, proxy cache.CacheInterface[string]) *TwoStepTokenProvider[T] {
	return &TwoStepTokenProvider[T]{c: cache2.NewProtoCache[T](creator, proxy), proxy: proxy}
}

func (p *TwoStepTokenProvider[T]) Name() string {
	return TwoStepName
}

func (p *TwoStepTokenProvider[T]) Generate(ctx context.Context, purpose TokenPurpose, payload T, duration time.Duration) (string, error) {

	token := uuid.New().String()
	key := fmt.Sprintf("%s:%s:%s", TwoStepName, purpose, token)
	err := p.c.Set(ctx, key, payload, store.WithExpiration(duration))
	if err != nil {
		return "", err
	}
	return token, nil
}
func (p *TwoStepTokenProvider[T]) Retrieve(ctx context.Context, purpose TokenPurpose, token string) (T, error) {
	key := fmt.Sprintf("%s:%s:%s", TwoStepName, purpose, token)
	t, err := p.c.Get(ctx, key)
	if err != nil {
		var n T
		if (store.NotFound{}).Is(err) {
			return n, nil
		}
		return n, err
	}
	return t, nil
}

type OtpTokenProvider interface {
	GenerateOtp(ctx context.Context, purpose TokenPurpose, extraKey string, duration time.Duration) (string, error)
	VerifyOtp(ctx context.Context, purpose TokenPurpose, extraKey string, token string) (bool, error)
}

type DefaultOtpTokenProvider struct {
	c cache.CacheInterface[string]
}

func NewOtpTokenProvider(c cache.CacheInterface[string]) *DefaultOtpTokenProvider {
	return &DefaultOtpTokenProvider{c: c}
}

func (p *DefaultOtpTokenProvider) GenerateOtp(ctx context.Context, purpose TokenPurpose, extraKey string, duration time.Duration) (string, error) {
	token, err := GenerateOtp()
	if err != nil {
		return "", err
	}
	key := fmt.Sprintf("%s:%s:%s", OtpName, purpose, extraKey)
	err = p.c.Set(ctx, key, token, store.WithExpiration(duration))
	if err != nil {
		return "", err
	}
	return token, nil
}
func (p *DefaultOtpTokenProvider) VerifyOtp(ctx context.Context, purpose TokenPurpose, extraKey string, token string) (bool, error) {
	key := fmt.Sprintf("%s:%s:%s", OtpName, purpose, extraKey)
	otp, err := p.c.Get(ctx, key)
	if err != nil {
		if (store.NotFound{}).Is(err) {
			return false, nil
		}
		return false, err
	}
	return otp == token, nil
}
