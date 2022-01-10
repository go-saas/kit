package biz

import (
	"context"
	"github.com/google/uuid"
	"github.com/goxiaoy/go-saas-kit/pkg/gorm"
	"github.com/volatiletech/authboss/v3"
	"github.com/volatiletech/authboss/v3/otp/twofactor/sms2fa"
	"github.com/volatiletech/authboss/v3/otp/twofactor/totp2fa"
	"time"
)

// This pattern is useful in real code to ensure that
// we've got the right interfaces implemented.
var (
	assertUser   = &User{}
	assertStorer = &AuthbossStoreWrapper{}

	_ authboss.User            = assertUser
	_ authboss.AuthableUser    = assertUser
	_ authboss.ConfirmableUser = assertUser
	_ authboss.LockableUser    = assertUser
	_ authboss.RecoverableUser = assertUser
	_ authboss.ArbitraryUser   = assertUser

	_ totp2fa.User             = assertUser
	_ sms2fa.User              = assertUser
	_ sms2fa.SMSNumberProvider = assertUser

	_ authboss.CreatingServerStorer    = assertStorer
	_ authboss.ConfirmingServerStorer  = assertStorer
	_ authboss.RecoveringServerStorer  = assertStorer
	_ authboss.RememberingServerStorer = assertStorer
)

func (u *User) GetPID() (pid string) {
	return u.ID.String()
}

func (u *User) PutPID(pid string) {
	u.UIDBase = gorm.UIDBase{ID: uuid.MustParse(pid)}
}

func (u *User) GetPassword() (password string) {
	if u.Password == nil {
		return ""
	} else {
		return *u.Password
	}
}

func (u *User) PutPassword(password string) {
	panic("use PasswordHasher")
}
func (u *User) GetEmail() (email string) {
	if u.Email == nil {
		return ""
	} else {
		return *u.Email
	}
}

func (u *User) GetConfirmed() (confirmed bool) {
	//do not need a confirmation
	if u.Email != nil {
		return u.EmailConfirmed
	}
	return true
}

func (u *User) GetConfirmSelector() (selector string) {
	return u.ConfirmSelector
}

func (u *User) GetConfirmVerifier() (verifier string) {
	return u.ConfirmVerifier
}

func (u *User) PutEmail(email string) {
	if email == "" {
		u.Email = nil
	} else {
		u.Email = &email
	}
}

func (u *User) PutConfirmed(confirmed bool) {
	u.EmailConfirmed = confirmed
}

func (u *User) PutConfirmSelector(selector string) {
	u.ConfirmSelector = selector
}

func (u *User) PutConfirmVerifier(verifier string) {
	u.ConfirmVerifier = verifier
}

func (u *User) GetAttemptCount() (attempts int) {
	return u.AccessFailedCount
}

func (u *User) GetLastAttempt() (last time.Time) {
	return u.LastLoginAttempt
}

func (u *User) GetLocked() (locked time.Time) {
	return u.LockoutEndDateUtc
}

func (u *User) PutAttemptCount(attempts int) {
	u.AccessFailedCount = attempts
}

func (u *User) PutLastAttempt(last time.Time) {
	u.LastLoginAttempt = last
}

func (u *User) PutLocked(locked time.Time) {
	u.LockoutEndDateUtc = locked
}

func (u *User) GetRecoverSelector() (selector string) {
	return u.RecoverSelector
}

func (u *User) GetRecoverVerifier() (verifier string) {
	return u.RecoverVerifier
}

func (u *User) GetRecoverExpiry() (expiry time.Time) {
	return u.RecoverTokenExpiry
}

func (u *User) PutRecoverSelector(selector string) {
	u.RecoverSelector = selector
}

func (u *User) PutRecoverVerifier(verifier string) {
	u.RecoverVerifier = verifier
}

func (u *User) PutRecoverExpiry(expiry time.Time) {
	u.RecoverTokenExpiry = expiry
}

func (u *User) GetArbitrary() (arbitrary map[string]string) {
	name := ""
	if u.Name != nil {
		name = *u.Name
	}
	return map[string]string{
		"name": name,
	}
}

func (u *User) PutArbitrary(arbitrary map[string]string) {
	if n, ok := arbitrary["name"]; ok && n != "" {
		u.Name = &n
	}
}

func (u *User) GetRecoveryCodes() string {
	return u.RecoveryCodes
}

func (u *User) PutRecoveryCodes(codes string) {
	u.RecoveryCodes = codes
}

func (u *User) GetTOTPSecretKey() string {
	return u.TOTPSecretKey
}

func (u *User) PutTOTPSecretKey(s string) {
	u.TOTPSecretKey = s
}

func (u *User) GetSMSPhoneNumber() string {
	if u.Phone == nil {
		return ""
	} else {
		return *u.Phone
	}
}

func (u *User) PutSMSPhoneNumber(s string) {
	if s == "" {
		u.Phone = nil
	} else {
		u.Phone = &s
	}
}

func (u *User) GetSMSPhoneNumberSeed() string {
	return ""
}

type AuthbossStoreWrapper struct {
	userManager   *UserManager
	userTokenRepo UserTokenRepo
}

func NewAuthbossStoreWrapper(userManager *UserManager, userTokenRepo UserTokenRepo) *AuthbossStoreWrapper {
	return &AuthbossStoreWrapper{userManager: userManager, userTokenRepo: userTokenRepo}
}

func (a *AuthbossStoreWrapper) AddRememberToken(ctx context.Context, pid, token string) error {
	_, err := a.userTokenRepo.Create(ctx, pid, InternalLoginProvider, InternalRememberTokenName, token)
	return err
}

func (a *AuthbossStoreWrapper) DelRememberTokens(ctx context.Context, pid string) error {
	return a.userTokenRepo.DeleteByUserIdAndLoginProviderAndName(ctx, pid, InternalLoginProvider, InternalRememberTokenName)
}

func (a *AuthbossStoreWrapper) UseRememberToken(ctx context.Context, pid, token string) error {
	//find token
	t, err := a.userTokenRepo.FindByUserIdAndLoginProviderAndName(ctx, pid, InternalLoginProvider, InternalRememberTokenName)
	if err != nil {
		return err
	}
	if t == nil || t.Value != token {
		return authboss.ErrTokenNotFound
	} else {
		//delete token
		if err := a.userTokenRepo.DeleteByUserIdAndLoginProviderAndName(ctx, pid, InternalLoginProvider, InternalRememberTokenName); err != nil {
			return err
		}
	}
	return nil
}

func (a *AuthbossStoreWrapper) LoadByRecoverSelector(ctx context.Context, selector string) (authboss.RecoverableUser, error) {
	if user, err := a.userManager.FindByRecoverSelector(ctx, selector); err != nil {
		return nil, err
	} else {
		if user == nil {
			return nil, authboss.ErrUserNotFound
		}
		return user, nil
	}
}

func (a *AuthbossStoreWrapper) LoadByConfirmSelector(ctx context.Context, selector string) (authboss.ConfirmableUser, error) {
	if user, err := a.userManager.FindByConfirmSelector(ctx, selector); err != nil {
		return nil, err
	} else {
		if user == nil {
			return nil, authboss.ErrUserNotFound
		}
		return user, nil
	}
}

func (a *AuthbossStoreWrapper) Load(ctx context.Context, key string) (authboss.User, error) {
	if user, err := a.userManager.FindByID(ctx, key); err != nil {
		return nil, err
	} else {
		if user == nil {
			return nil, authboss.ErrUserNotFound
		}
		return user, nil
	}
}

func (a *AuthbossStoreWrapper) Save(ctx context.Context, user authboss.User) error {
	u := user.(*User)
	return a.userManager.Update(ctx, u)
}

func (a *AuthbossStoreWrapper) New(ctx context.Context) authboss.User {
	return &User{}
}

func (a *AuthbossStoreWrapper) Create(ctx context.Context, user authboss.User) error {
	u := user.(*User)
	return a.userManager.Create(ctx, u)
}
