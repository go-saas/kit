package remote

import (
	"context"
	"errors"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/session"
	errors2 "github.com/goxiaoy/go-saas-kit/pkg/errors"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/auth/v1"
)

func NewRefreshProvider(client v1.AuthClient, logger klog.Logger) session.RefreshTokenProvider {
	l := klog.NewHelper(klog.With(logger, "module", "remote.RefreshTokenProvider"))
	return func(ctx context.Context, token string) (err error) {
		if writer, ok := session.FromClientStateWriterContext(ctx); ok {
			handlerError := func(err error) error {
				if errors2.NotBizError(err) {
					return err
				} else {
					l.Errorf("fail to refresh with error %v", err)
					//just clean remember token
					err = writer.SignOutRememberToken(ctx)
					if err != nil {
						return err
					}
					err = writer.Save(ctx)
					if err != nil {
						return err
					}
					return nil
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
			err = writer.SetRememberToken(ctx, rep.NewRmToken)
			if err != nil {
				return err
			}
			err = writer.Save(ctx)
			if err != nil {
				return err
			}
			return nil

		} else {
			return errors.New("writer not found")
		}
	}
}
