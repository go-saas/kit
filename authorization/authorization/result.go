package authorization

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/goxiaoy/go-saas-kit/pkg/auth/jwt"
)

type Result struct {
	Allowed      bool
	Requirements []Requirement
}

func NewAllowAuthorizationResult() Result {
	return Result{Allowed: true}
}

func NewDisallowAuthorizationResult(requirements []Requirement) Result {
	return Result{Allowed: false, Requirements: requirements}
}

func FormatError(ctx context.Context, result Result) error {
	if result.Allowed {
		return nil
	}
	if _, ok := jwt.FromClaimsContext(ctx); ok {
		//TODO format error
		return errors.Forbidden("", "")
	}
	//no claims
	return errors.Unauthorized("", "")
}
