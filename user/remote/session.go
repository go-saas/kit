package remote

import (
	"context"
	"errors"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/session"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/auth/v1"
)

func NewRefreshProvider(client v1.AuthClient) session.RefreshTokenProvider {
	return func(ctx context.Context, token string) (err error) {
		if writer, ok := session.FromClientStateWriterContext(ctx); ok {
			rep, err := client.RefreshRememberToken(ctx, &v1.RefreshRememberTokenRequest{RmToken: token})
			if err != nil {
				return err
			}
			if err := writer.SetUid(ctx, rep.UserId); err != nil {
				return err
			}
			err = writer.SetRememberToken(ctx, rep.NewRmToken)
			if err != nil {
				return err
			}
			return nil

		} else {
			return errors.New("writer not found")
		}
	}
}
