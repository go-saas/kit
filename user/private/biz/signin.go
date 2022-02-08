package biz

import (
	"context"
	"errors"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/session"
)

var (
	ErrUserDeleted    = errors.New("user deleted")
	ErrUserLocked     = errors.New("user locked")
	ErrWriterNotFound = errors.New("writer not found")
)

type SignInManager struct {
	um *UserManager
}

func NewSignInManager(um *UserManager) *SignInManager {
	return &SignInManager{um: um}
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
		return ErrUserDeleted
	}
	if locked, err := s.um.CheckLocked(ctx, u); err != nil {
		return err
	} else if locked {
		return ErrUserLocked
	}
	return nil
}

//RefreshSignIn refresh sign in
func (s *SignInManager) RefreshSignIn(ctx context.Context, u *User) error {
	panic("")
}

func (s *SignInManager) SignIn(ctx context.Context, u *User, isPersistent bool) {
	panic("")
}

func (s *SignInManager) SignOut(ctx context.Context) {
	panic("")
}

func (s *SignInManager) ValidateSecurityStamp(ctx context.Context, u *User, securityStamp string) {
	panic("")
}

func (s *SignInManager) PasswordSignIn(ctx context.Context, u *User, pwd string, isPersistent bool, tryLockoutOnFailure bool) error {
	if u == nil {
		return ErrInvalidCredential
	}
	err := s.um.CheckPassword(ctx, u, pwd)
	if err != nil {
		if err == ErrInvalidCredential {
			if tryLockoutOnFailure {
				//TODO lock out
			}
		}
		return err
	}
	//password correct
	if writer, ok := session.FromClientStateWriterContext(ctx); ok {
		if err := writer.SetUid(ctx, u.ID.String()); err != nil {
			return err
		}
		if isPersistent {
			//TODO generate remember me token and set in writer
		}
	}
	return ErrWriterNotFound
}
func (s *SignInManager) PasswordSignInWithUsername(ctx context.Context, username, pwd string, isPersistent bool, tryLockoutOnFailure bool) error {
	u, err := s.um.FindByName(ctx, username)
	if err != nil {
		return err
	}
	return s.PasswordSignIn(ctx, u, pwd, isPersistent, tryLockoutOnFailure)
}
func (s *SignInManager) PasswordSignInWithPhone(ctx context.Context, phone, pwd string, isPersistent bool, tryLockoutOnFailure bool) error {
	u, err := s.um.FindByPhone(ctx, phone)
	if err != nil {
		return err
	}
	return s.PasswordSignIn(ctx, u, pwd, isPersistent, tryLockoutOnFailure)
}
func (s *SignInManager) PasswordSignInWithEmail(ctx context.Context, email, pwd string, isPersistent bool, tryLockoutOnFailure bool) error {
	u, err := s.um.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	return s.PasswordSignIn(ctx, u, pwd, isPersistent, tryLockoutOnFailure)
}

func (s *SignInManager) IsTwoFactorClientRemembered(ctx context.Context, u *User) (bool, error) {
	if state, ok := session.FromClientStateContext(ctx); ok {
		return state.GetTwoFactorClientRemembered(), nil
	}
	return false, nil
}

func (s *SignInManager) RememberTwoFactorClient(ctx context.Context, u *User) error {
	if writer, ok := session.FromClientStateWriterContext(ctx); ok {
		return writer.SetTwoFactorClientRemembered(ctx)
	}
	return ErrWriterNotFound
}

func (s *SignInManager) ForgetTwoFactorClient(ctx context.Context) error {
	if writer, ok := session.FromClientStateWriterContext(ctx); ok {
		return writer.SignOutTwoFactorClientRemembered(ctx)
	}
	return ErrWriterNotFound
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
