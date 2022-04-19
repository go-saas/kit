package biz

import (
	"context"
	"fmt"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/auth/v1"
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

	ConfirmPurpose               TokenPurpose = "confirm"
	RecoverPurpose               TokenPurpose = "recover"
	RecoverChangePasswordPurpose TokenPurpose = "recover_change_password"
)

var (
	_ UserTokenProvider = (*PhoneTokenProvider)(nil)
	_ UserTokenProvider = (*EmailTokenProvider)(nil)
)

type PhoneTokenProvider struct {
	r *redis.Client
}

func NewPhoneTokenProvider(r *redis.Client) *PhoneTokenProvider {
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
	err = p.r.Set(ctx, key, token, duration).Err()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (p *PhoneTokenProvider) Validate(ctx context.Context, purpose TokenPurpose, token string, user *User) (bool, error) {
	key := fmt.Sprintf("usertoken:%s:%s:%s", user.ID.String(), purpose, *user.Phone)
	val, err := p.r.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
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
	return v1.ErrorPhoneNotConfirmed("")
}

type EmailTokenProvider struct {
	r *redis.Client
}

func NewEmailTokenProvider(r *redis.Client) *EmailTokenProvider {
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
	err = e.r.Set(ctx, key, token, duration).Err()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (e *EmailTokenProvider) Validate(ctx context.Context, purpose TokenPurpose, token string, user *User) (bool, error) {
	key := fmt.Sprintf("usertoken:%s:%s:%s", user.ID.String(), purpose, *user.NormalizedEmail)
	val, err := e.r.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
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
	return v1.ErrorEmailNotConfirmed("")
}

type TwoStepTokenProvider struct {
	c *cache.Cache
}

func NewTwoStepTokenProvider(c *cache.Cache) *TwoStepTokenProvider {
	return &TwoStepTokenProvider{c: c}
}

func (p *TwoStepTokenProvider) Name() string {
	return TwoStepName
}

func (p *TwoStepTokenProvider) Generate(ctx context.Context, purpose TokenPurpose, payload interface{}, duration time.Duration) (string, error) {

	token := uuid.New().String()
	key := fmt.Sprintf("%s:%s:%s", TwoStepName, purpose, token)
	err := p.c.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: payload,
		TTL:   &duration,
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func (p *TwoStepTokenProvider) Retrieve(ctx context.Context, purpose TokenPurpose, token string, dest interface{}) (bool, error) {
	key := fmt.Sprintf("%s:%s:%s", TwoStepName, purpose, token)
	err := p.c.Get(ctx, key, dest)
	if err != nil {
		if err == cache.ErrCacheMiss {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
