package api

import (
	"context"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authn/session"
	errors2 "github.com/go-saas/kit/pkg/errors"
	v1 "github.com/go-saas/kit/user/api/auth/v1"
)

// NewRefreshProvider return session.RefreshTokenProvider
//
// Read session -> Call v1.AuthServer to exchange token -> invalid, sign out.
func NewRefreshProvider(srv v1.AuthServer, logger klog.Logger) session.RefreshTokenProvider {
	l := klog.NewHelper(klog.With(logger, "module", "remote.RefreshTokenProvider"))
	return session.RefreshTokenProviderFunc(func(ctx context.Context, token string) (*session.RefreshNewToken, error) {
		//replace withe trusted environment to skip trusted check if in same process
		ctx = api.NewTrustedContext(ctx)
		if writer, ok := session.FromClientStateWriterContext(ctx); ok {
			handlerError := func(err error) error {
				err = kerrors.FromError(err)
				if errors2.UnRecoverableError(err) {
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
			// call remote or local service to exchange token
			rep, err := srv.RefreshRememberToken(ctx, &v1.RefreshRememberTokenRequest{RmToken: token})
			if err != nil {
				return nil, handlerError(err)
			}
			t := &session.RefreshNewToken{
				UserId:   rep.UserId,
				NewToken: rep.NewRmToken,
			}
			l.Infof("refresh user %s remember token successfully", rep.UserId)
			if err := writer.SetUid(ctx, rep.UserId); err != nil {
				return t, err
			}
			if err := writer.SetRememberToken(ctx, rep.NewRmToken, t.UserId); err != nil {
				return t, err
			}
			if err := writer.Save(ctx); err != nil {
				return t, err
			}
			return t, nil

		} else {
			panic("writer not found")
		}
	})
}
