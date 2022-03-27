package api

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
)

type TrustedContextValidator interface {
	Trusted(ctx context.Context) (bool, error)
}

type ClientTrustedContextValidator struct {
}

func NewClientTrustedContextValidator() TrustedContextValidator {
	return &ClientTrustedContextValidator{}
}

func (c *ClientTrustedContextValidator) Trusted(ctx context.Context) (bool, error) {
	if claims, ok := jwt.FromClaimsContext(ctx); ok {
		//TODO trusted server to server communication
		if claims.ClientId != "" {
			return true, nil
		}
	}
	return false, nil
}

var _ TrustedContextValidator = (*ClientTrustedContextValidator)(nil)
