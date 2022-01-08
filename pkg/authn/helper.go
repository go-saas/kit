package authn

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

func ErrIfUnauthenticated(ctx context.Context) (UserInfo, error) {
	user, ok := FromUserContext(ctx)
	if !ok || user.GetId() == "" {
		return user, errors.Unauthorized("", "")
	}
	return user, nil
}
