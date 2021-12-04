package auth

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/goxiaoy/go-saas-kit/auth/current"
)

func ErrIfUnauthenticated(ctx context.Context) (current.UserInfo, error) {
	user, ok := current.FromUserContext(ctx)
	if !ok || user.GetId() == "" {
		return user, errors.Unauthorized("", "")
	}
	return user, nil
}
