package auth

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/goxiaoy/go-saas-kit/auth/current"
)

func ErrIfUnauthorized(ctx context.Context) error {
	user, ok := current.FromUserContext(ctx)
	if !ok || user.Id == "" {
		return errors.Unauthorized("", "")
	}
	return nil
}
