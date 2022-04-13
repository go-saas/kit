package remote

import (
	"context"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/session"
	errors2 "github.com/goxiaoy/go-saas-kit/pkg/errors"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/auth/v1"
)

func NewRefreshProvider(client v1.AuthClient, logger klog.Logger) session.RefreshTokenProvider {
	l := klog.NewHelper(klog.With(logger, "module", "remote.RefreshTokenProvider"))
	return session.RefreshTokenProviderFunc(func(ctx context.Context, token, userId string) (err error) {
		if writer, ok := session.FromClientStateWriterContext(ctx); ok {
			handlerError := func(err error) error {
				err = kerrors.FromError(err)
				if errors2.NotBizError(err) {
					return err
				} else {
					if !v1.IsRememberTokenUsed(err) {
						//clean remember token
						if err := writer.SignOutRememberToken(ctx); err != nil {
							return err
						}
						if err := writer.Save(ctx); err != nil {
							return err
						}
					}
					return err
				}
			}

			rep, err := client.RefreshRememberToken(ctx, &v1.RefreshRememberTokenRequest{RmToken: token})
			if err != nil {
				return handlerError(err)
			}
			l.Infof("refresh user %s remember token successfully", rep.UserId)
			if err := writer.SetUid(ctx, rep.UserId); err != nil {
				return err
			}
			if err := writer.SetRememberToken(ctx, rep.NewRmToken, userId); err != nil {
				return err
			}
			if err := writer.Save(ctx); err != nil {
				return err
			}
			return nil

		} else {
			panic("writer not found")
		}
	})
}
