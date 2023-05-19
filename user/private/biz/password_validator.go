package biz

import (
	"context"
	v1 "github.com/go-saas/kit/user/api/user/v1"
	"github.com/go-saas/kit/user/private/conf"
	"github.com/nbutton23/zxcvbn-go"
)

type PasswordValidator interface {
	// Validate password
	Validate(ctx context.Context, password string) error
}

type passwordValidator struct {
	config *conf.UserConf
}

func NewPasswordValidator(c *conf.UserConf) PasswordValidator {
	return &passwordValidator{
		config: c,
	}
}

func (p *passwordValidator) Validate(ctx context.Context, password string) (err error) {
	if len(password) > 100 {
		password = password[:100]
	}

	strength := zxcvbn.PasswordStrength(password, []string{})
	ok := strength.Score >= int(p.config.PasswordScoreMin)
	if !ok {
		return v1.ErrorPasswordInsufficientStrengthLocalized(ctx, nil, nil)
	}
	return nil
}
