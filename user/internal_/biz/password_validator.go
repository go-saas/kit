package biz

import (
	"context"
	"errors"
	"github.com/nbutton23/zxcvbn-go"
)

var (
	ErrInsufficientStrength = errors.New("insufficient strength")
)

type PasswordValidator interface {
	// Validate password
	Validate(ctx context.Context, password string) error
}

type PasswordValidatorConfig struct {
	MinScore int
}
type passwordValidator struct {
	config *PasswordValidatorConfig
}

func NewPasswordValidator(c *PasswordValidatorConfig) PasswordValidator {
	return &passwordValidator{
		config: c,
	}
}

func (p *passwordValidator) Validate(ctx context.Context, password string) (err error) {
	if len(password) > 100 {
		password = password[:100]
	}

	strength := zxcvbn.PasswordStrength(password, []string{})
	ok := strength.Score >= p.config.MinScore
	if !ok {
		return ErrInsufficientStrength
	}
	return nil
}
