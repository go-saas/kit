package biz

import (
	"context"
	"errors"
	"github.com/go-saas/kit/pkg/authn/session"
	"github.com/go-saas/kit/pkg/conf"
	v12 "github.com/go-saas/kit/user/api/auth/v1"
	"time"
)

var (
	ErrWriterNotFound = errors.New("writer not found")
)

type SignInManager struct {
	um          *UserManager
	securityCfg *conf.Security
}

func NewSignInManager(um *UserManager, securityCfg *conf.Security) *SignInManager {
	return &SignInManager{um: um, securityCfg: securityCfg}
}

func (s *SignInManager) IsSignedIn(ctx context.Context) (bool, error) {
	if state, ok := session.FromClientStateContext(ctx); ok {
		return len(state.GetUid()) > 0, nil
	}
	return false, nil
}

func (s *SignInManager) CheckCanSignIn(ctx context.Context, u *User) error {
	if d, err := s.um.CheckDeleted(ctx, u); err != nil {
		return err
	} else if d {
		return v12.ErrorUserDeletedLocalized(ctx, nil, nil)
	}
	if locked, err := s.um.CheckLocked(ctx, u); err != nil {
		return err
	} else if locked {
		return v12.ErrorUserLockedLocalized(ctx, nil, nil)
	}
	return nil
}

func (s *SignInManager) SignIn(ctx context.Context, u *User, isPersistent bool) error {
	if writer, ok := session.FromClientStateWriterContext(ctx); ok {
		if err := writer.SetUid(ctx, u.ID.String()); err != nil {
			return err
		}
		if isPersistent {
			duration := 0
			if s.securityCfg != nil && s.securityCfg.RememberCookie != nil && s.securityCfg.RememberCookie.MaxAge != nil {
				duration = int(s.securityCfg.RememberCookie.MaxAge.Value)
			}
			duration = session.RememberMeExpireSecondsOrDefault(duration)
			rememberToken, err := s.um.GenerateRememberToken(ctx, time.Duration(duration)*time.Second, u.ID)
			if err != nil {
				return err
			}
			err = writer.SetRememberToken(ctx, rememberToken, u.ID.String())
			if err != nil {
				return err
			}
		}
		//save session
		return writer.Save(ctx)
	} else {
		panic(ErrWriterNotFound)
	}
}

func (s *SignInManager) SignOut(ctx context.Context) error {
	if writer, ok := session.FromClientStateWriterContext(ctx); ok {
		if err := writer.Clear(ctx); err != nil {
			return err
		}
		return writer.Save(ctx)
	} else {
		panic(ErrWriterNotFound)
	}

}

func (s *SignInManager) ValidateSecurityStamp(ctx context.Context, u *User, securityStamp string) {
	panic("")
}

func (s *SignInManager) PasswordSignIn(ctx context.Context, u *User, pwd string, isPersistent bool, tryLockoutOnFailure bool) error {
	if u == nil {
		return v12.ErrorInvalidCredentialsLocalized(ctx, nil, nil)
	}
	err := s.um.CheckPassword(ctx, u, pwd)
	if err != nil {
		if v12.IsInvalidCredentials(err) {
			if tryLockoutOnFailure {
				//TODO lock out
			}
		}
		return err
	}
	//password correct
	return s.SignIn(ctx, u, isPersistent)
}
func (s *SignInManager) PasswordSignInWithUsername(ctx context.Context, username, pwd string, isPersistent bool, tryLockoutOnFailure bool) (error, string) {
	u, err := s.um.FindByName(ctx, username)
	if err != nil {
		return err, ""
	}
	id := ""
	if u != nil {
		id = u.ID.String()
	}
	return s.PasswordSignIn(ctx, u, pwd, isPersistent, tryLockoutOnFailure), id
}
func (s *SignInManager) PasswordSignInWithPhone(ctx context.Context, phone, pwd string, isPersistent bool, tryLockoutOnFailure bool) (error, string) {
	u, err := s.um.FindByPhone(ctx, phone)
	if err != nil {
		return err, ""
	}
	id := ""
	if u != nil {
		id = u.ID.String()
	}
	return s.PasswordSignIn(ctx, u, pwd, isPersistent, tryLockoutOnFailure), id
}
func (s *SignInManager) PasswordSignInWithEmail(ctx context.Context, email, pwd string, isPersistent bool, tryLockoutOnFailure bool) (error, string) {
	u, err := s.um.FindByEmail(ctx, email)
	if err != nil {
		return err, ""
	}
	id := ""
	if u != nil {
		id = u.ID.String()
	}
	return s.PasswordSignIn(ctx, u, pwd, isPersistent, tryLockoutOnFailure), id
}

func (s *SignInManager) IsTwoFactorClientRemembered(ctx context.Context, u *User) (bool, error) {
	if state, ok := session.FromClientStateContext(ctx); ok {
		return state.GetTwoFactorClientRemembered(), nil
	}
	return false, nil
}

func (s *SignInManager) RememberTwoFactorClient(ctx context.Context, u *User) error {
	if writer, ok := session.FromClientStateWriterContext(ctx); ok {
		err := writer.SetTwoFactorClientRemembered(ctx)
		if err != nil {
			return err
		}
		return writer.Save(ctx)
	}
	panic(ErrWriterNotFound)
}

func (s *SignInManager) ForgetTwoFactorClient(ctx context.Context) error {
	if writer, ok := session.FromClientStateWriterContext(ctx); ok {
		return writer.SignOutTwoFactorClientRemembered(ctx)
	}
	panic(ErrWriterNotFound)
}

func (s *SignInManager) TwoFactorAuthenticatorSignIn(ctx context.Context, code string, isPersistent, rememberClient bool) error {
	panic("")
	//TODO find two factor info
	//Compare code
	//Sign out

}
func (s *SignInManager) TwoFactorSignIn(ctx context.Context, provider, code string, isPersistent, rememberClient bool) error {
	panic("")
	//TODO find two factor info
	//Compare code
	//Sign out
}

func (s *SignInManager) GetTwoFactorAuthenticationUser(ctx context.Context) (*User, error) {
	panic("")
}
func (s *SignInManager) ExternalLoginSignInAsync(ctx context.Context, loginProvider, providerKey string, isPersistent, bypassTwoFactor bool) error {
	panic("")
}

func (s *SignInManager) isLockedOut(ctx context.Context, u *User) (bool, error) {
	panic("")
}
func (s *SignInManager) preSignInCheck(ctx context.Context, u *User) error {
	if err := s.CheckCanSignIn(ctx, u); err != nil {
		return err
	}
	panic("")
}
func (s *SignInManager) isTfaEnabled(ctx context.Context, u *User) (bool, error) {
	panic("")
}

func (s *SignInManager) resetLockedOut(ctx context.Context, u *User) error {
	panic("")
}
