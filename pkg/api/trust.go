package api

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-saas/kit/pkg/authn/jwt"
)

type (
	trustKey struct{}
)

// TrustedContextValidator validate whether the communication is behind authed gateway or server to server communication
type TrustedContextValidator interface {
	Trusted(ctx context.Context) (bool, error)
}

type ClientTrustedContextValidator struct {
}

func NewClientTrustedContextValidator() TrustedContextValidator {
	return &ClientTrustedContextValidator{}
}

// NewTrustedContext create a trusted (or not) context without propaganda to other services
func NewTrustedContext(ctx context.Context, trust ...bool) context.Context {
	t := true
	if len(trust) > 0 {
		t = trust[0]
	}
	return context.WithValue(ctx, trustKey{}, t)
}

func FromTrustedContext(ctx context.Context) (bool, bool) {
	v, ok := ctx.Value(trustKey{}).(bool)
	if ok {
		return v, ok
	}
	return false, false
}

func (c *ClientTrustedContextValidator) Trusted(ctx context.Context) (bool, error) {
	if v, ok := FromTrustedContext(ctx); ok {
		return v, nil
	}
	if claims, ok := jwt.FromClaimsContext(ctx); ok {
		//TODO trusted server to server communication
		if claims.ClientId != "" {
			return true, nil
		}
	}
	return false, nil
}

var _ TrustedContextValidator = (*ClientTrustedContextValidator)(nil)

func ErrIfUntrusted(ctx context.Context, t TrustedContextValidator) error {
	if ok, err := t.Trusted(ctx); err != nil {
		return err
	} else if !ok {
		return errors.Forbidden("", "")
	}
	return nil
}
