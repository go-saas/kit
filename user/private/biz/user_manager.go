package biz

import (
	"context"
	"fmt"
	"github.com/eko/gocache/v3/cache"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-saas/kit/event"
	cache2 "github.com/go-saas/kit/pkg/cache"
	"github.com/go-saas/kit/pkg/query"
	kithttp "github.com/go-saas/kit/pkg/server/http"
	v12 "github.com/go-saas/kit/user/api/auth/v1"
	v1 "github.com/go-saas/kit/user/api/user/v1"
	v13 "github.com/go-saas/kit/user/event/v1"
	"github.com/go-saas/kit/user/private/conf"
	"github.com/go-saas/saas"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"time"
)

type UserManager struct {
	cfg              *conf.UserConf
	userRepo         UserRepo
	pwdHasher        PasswordHasher
	userValidator    UserValidator
	pwdValidator     PasswordValidator
	lookupNormalizer LookupNormalizer
	userTokenRepo    UserTokenRepo
	refreshTokenRepo RefreshTokenRepo
	userTenantRepo   UserTenantRepo
	emailToken       *EmailTokenProvider
	phoneToken       *PhoneTokenProvider
	forgetPwdTwoStep *TwoStepTokenProvider[*ForgetPasswordTwoStepTokenPayload]
	otp              OtpTokenProvider
	userRoleCache    *cache2.Helper[*UserRoleCacheItem]
	log              log.Logger
}

func NewUserManager(
	cfg *conf.UserConf,
	userRepo UserRepo,
	pwdHasher PasswordHasher,
	userValidator UserValidator,
	pwdValidator PasswordValidator,
	lookupNormalizer LookupNormalizer,
	userTokenRepo UserTokenRepo,
	refreshTokenRepo RefreshTokenRepo,
	userTenantRepo UserTenantRepo,
	emailToken *EmailTokenProvider,
	phoneToken *PhoneTokenProvider,
	otp OtpTokenProvider,
	strCache cache.CacheInterface[string],
	logger log.Logger) *UserManager {
	return &UserManager{
		cfg:              cfg,
		userRepo:         userRepo,
		pwdHasher:        pwdHasher,
		userValidator:    userValidator,
		pwdValidator:     pwdValidator,
		lookupNormalizer: lookupNormalizer,
		userTokenRepo:    userTokenRepo,
		refreshTokenRepo: refreshTokenRepo,
		userTenantRepo:   userTenantRepo,
		emailToken:       emailToken,
		phoneToken:       phoneToken,
		forgetPwdTwoStep: NewTwoStepTokenProvider(func() *ForgetPasswordTwoStepTokenPayload {
			return &ForgetPasswordTwoStepTokenPayload{}
		}, strCache),
		otp:           otp,
		userRoleCache: cache2.NewHelper[*UserRoleCacheItem](cache2.NewProtoCache(func() *UserRoleCacheItem { return &UserRoleCacheItem{} }, strCache)),
		log:           log.With(logger, "module", "user/biz/user_manager")}
}

func (um *UserManager) List(ctx context.Context, query *v1.ListUsersRequest) ([]*User, error) {
	return um.userRepo.List(ctx, query)
}

func (um *UserManager) Count(ctx context.Context, query *v1.ListUsersRequest) (total int64, filtered int64, err error) {
	return um.userRepo.Count(ctx, query)
}

func (um *UserManager) ListAdmin(ctx context.Context, query *v1.AdminListUsersRequest) ([]*User, error) {
	return um.userRepo.ListAdmin(ctx, query)
}

func (um *UserManager) CountAdmin(ctx context.Context, query *v1.AdminListUsersRequest) (total int64, filtered int64, err error) {
	return um.userRepo.CountAdmin(ctx, query)
}

func (um *UserManager) Create(ctx context.Context, u *User) (err error) {
	err = um.normalize(ctx, u)
	if err != nil {
		return err
	}
	if err = um.validateUser(ctx, u); err != nil {
		return
	}
	if err := um.userRepo.Create(ctx, u); err != nil {
		return err
	}
	return nil
}

func (um *UserManager) CreateWithPassword(ctx context.Context, u *User, pwd string, validate bool) (err error) {
	var password *string
	if len(pwd) > 0 {
		password = &pwd
	}
	if err = um.updatePassword(ctx, u, password, validate); err != nil {
		return
	}
	return um.Create(ctx, u)
}

func (um *UserManager) FindByID(ctx context.Context, id string) (user *User, err error) {
	return um.userRepo.FindByID(ctx, id)
}

func (um *UserManager) FindByName(ctx context.Context, name string) (user *User, err error) {
	name, err = um.lookupNormalizer.Name(ctx, name)
	if err != nil {
		return nil, err
	}
	return um.userRepo.FindByName(ctx, name)
}

func (um *UserManager) FindByPhone(ctx context.Context, phone string) (user *User, err error) {
	phone, err = um.lookupNormalizer.Phone(ctx, phone)
	if err != nil {
		return nil, err
	}
	return um.userRepo.FindByPhone(ctx, phone)
}

func (um *UserManager) FindByEmail(ctx context.Context, email string) (user *User, err error) {
	email, err = um.lookupNormalizer.Email(ctx, email)
	if err != nil {
		return nil, err
	}
	return um.userRepo.FindByEmail(ctx, email)
}

func (um *UserManager) FindByIdentity(ctx context.Context, identity string) (user *User, err error) {
	//try to find by id
	if uid, err := uuid.Parse(identity); err == nil {
		return um.FindByID(ctx, uid.String())
	}
	if _, err := um.lookupNormalizer.Email(ctx, identity); err == nil {
		return um.FindByEmail(ctx, identity)
	}
	if phone, err := um.lookupNormalizer.Phone(ctx, identity); err == nil {
		return um.FindByPhone(ctx, phone)
	}
	return um.FindByName(ctx, identity)
}

func (um *UserManager) FindByLogin(ctx context.Context, loginProvider string, providerKey string) (*User, error) {
	return um.userRepo.FindByLogin(ctx, loginProvider, providerKey)
}

func (um *UserManager) ListLogin(ctx context.Context, user *User) ([]*UserLogin, error) {
	return um.userRepo.ListLogin(ctx, user)
}

func (um *UserManager) AddLogin(ctx context.Context, user *User, logins []UserLogin) error {
	//find logins
	existing, err := um.userRepo.ListLogin(ctx, user)
	if err != nil {
		return err
	}
	for _, login := range logins {
		if lo.ContainsBy(existing, func(l *UserLogin) bool {
			return l.LoginProvider == login.LoginProvider && l.ProviderKey == login.ProviderKey
		}) {
			continue
		}
		err = um.userRepo.AddLogin(ctx, user, &UserLogin{
			UserId:        user.ID,
			LoginProvider: login.LoginProvider,
			ProviderKey:   login.ProviderKey,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (um *UserManager) Update(ctx context.Context, user *User, q query.Select) (err error) {
	err = um.normalize(ctx, user)
	if err != nil {
		return err
	}
	if err = um.validateUser(ctx, user); err != nil {
		return
	}
	return um.userRepo.Update(ctx, user.ID.String(), user, q)
}

func (um *UserManager) Delete(ctx context.Context, user *User) error {
	return um.userRepo.Delete(ctx, user.ID.String())
}

func (um *UserManager) CheckPassword(ctx context.Context, user *User, password string) error {
	v := um.checkPassword(ctx, user, password)
	if v == PasswordVerificationSuccess {
		return nil
	}
	if v == PasswordVerificationSuccessRehashNeeded {
		if err := um.updatePassword(ctx, user, &password, false); err != nil {
			return err
		}
		err := um.userRepo.Update(ctx, user.ID.String(), user, query.NewField(&fieldmaskpb.FieldMask{Paths: []string{"password"}}))
		return err
	}
	//fail
	return v12.ErrorInvalidCredentialsLocalized(ctx, nil, nil)
}

func (um *UserManager) ChangePassword(ctx context.Context, user *User, current string, newPwd string) error {
	if v := um.checkPassword(ctx, user, current); v == PasswordVerificationFail {
		return v12.ErrorInvalidCredentialsLocalized(ctx, nil, nil)
	}
	if err := um.updatePassword(ctx, user, &newPwd, true); err != nil {
		return err
	}
	return um.Update(ctx, user, query.NewField(&fieldmaskpb.FieldMask{Paths: []string{"password"}}))
}

func (um *UserManager) UpdatePassword(ctx context.Context, user *User, newPwd string) error {
	if err := um.updatePassword(ctx, user, &newPwd, true); err != nil {
		return err
	}
	return um.Update(ctx, user, query.NewField(&fieldmaskpb.FieldMask{Paths: []string{"password"}}))
}

func (um *UserManager) GenerateEmailForgetPasswordToken(ctx context.Context, user *User) (string, error) {
	duration := 5 * time.Minute
	if um.cfg.EmailRecoverExpiry != nil {
		duration = um.cfg.EmailRecoverExpiry.AsDuration()
	}
	return um.emailToken.Generate(ctx, RecoverPurpose, user, duration)
}

func (um *UserManager) VerifyEmailForgetPasswordToken(ctx context.Context, email, token string) error {
	user, err := um.FindByPhone(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return v12.ErrorEmailRecoverFailedLocalized(ctx, nil, nil)
	}
	ok, err := um.emailToken.Validate(ctx, RecoverPurpose, token, user)
	if err != nil {
		return err
	}
	if !ok {
		return v12.ErrorEmailRecoverFailedLocalized(ctx, nil, nil)
	}
	return nil
}

func (um *UserManager) GeneratePhoneForgetPasswordToken(ctx context.Context, user *User) (string, error) {
	duration := 5 * time.Minute
	if um.cfg.PhoneRecoverExpiry != nil {
		duration = um.cfg.PhoneRecoverExpiry.AsDuration()
	}
	return um.phoneToken.Generate(ctx, RecoverPurpose, user, duration)
}

func (um *UserManager) VerifyPhoneForgetPasswordToken(ctx context.Context, phone, token string) error {
	user, err := um.FindByPhone(ctx, phone)
	if err != nil {
		return err
	}
	if user == nil {
		return v12.ErrorPhoneRecoverFailedLocalized(ctx, nil, nil)
	}
	ok, err := um.phoneToken.Validate(ctx, RecoverPurpose, token, user)
	if err != nil {
		return err
	}
	if !ok {
		return v12.ErrorPhoneRecoverFailedLocalized(ctx, nil, nil)
	}
	return nil
}

func (um *UserManager) GenerateEmailLoginPasswordlessToken(ctx context.Context, email string) (string, error) {
	email, err := um.lookupNormalizer.Email(ctx, email)
	if err != nil {
		return "", err
	}
	duration := 5 * time.Minute
	if um.cfg.LoginPasswordlessExpiry != nil {
		duration = um.cfg.LoginPasswordlessExpiry.AsDuration()
	}
	return um.otp.GenerateOtp(ctx, EmailLoginPurpose, email, duration)
}

func (um *UserManager) VerifyEmailLoginPasswordlessToken(ctx context.Context, email, token string) (*User, error) {
	email, err := um.lookupNormalizer.Email(ctx, email)
	if err != nil {
		return nil, err
	}
	ok, err := um.otp.VerifyOtp(ctx, EmailLoginPurpose, email, token)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, v12.ErrorEmailPasswordlessFailedLocalized(ctx, nil, nil)
	}
	//find or create user
	user, err := um.FindByEmail(ctx, email)
	if err != nil {
		return user, err
	}
	if user == nil {
		user = &User{}
		user.SetEmail(email, true)
		if err := um.Create(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (um *UserManager) GeneratePhoneLoginPasswordlessToken(ctx context.Context, phone string) (string, error) {
	phone, err := um.lookupNormalizer.Phone(ctx, phone)
	if err != nil {
		return "", err
	}
	duration := 5 * time.Minute
	if um.cfg.LoginPasswordlessExpiry != nil {
		duration = um.cfg.LoginPasswordlessExpiry.AsDuration()
	}
	return um.otp.GenerateOtp(ctx, PhoneLoginPurpose, phone, duration)
}

func (um *UserManager) VerifyPhoneLoginPasswordlessToken(ctx context.Context, phone, token string) (*User, error) {
	phone, err := um.lookupNormalizer.Phone(ctx, phone)
	if err != nil {
		return nil, err
	}
	ok, err := um.otp.VerifyOtp(ctx, PhoneLoginPurpose, phone, token)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, v12.ErrorPhonePasswordlessFailedLocalized(ctx, nil, nil)
	}
	//find or create user
	user, err := um.FindByPhone(ctx, phone)
	if err != nil {
		return user, err
	}
	if user == nil {
		user = &User{}
		user.SetPhone(phone, true)
		if err := um.Create(ctx, user); err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (um *UserManager) GenerateForgetPasswordTwoStepToken(ctx context.Context, user *User) (string, error) {
	duration := 5 * time.Minute
	return um.forgetPwdTwoStep.Generate(ctx, RecoverChangePasswordPurpose, &ForgetPasswordTwoStepTokenPayload{UserId: user.ID.String()}, duration)
}

func (um *UserManager) ChangePasswordByToken(ctx context.Context, token, newPwd string) error {
	user, err := um.retrieveTwoStepForgetPasswordToken(ctx, token)
	if err != nil {
		return err
	}
	if user == nil {
		return v1.ErrorUserNotFoundLocalized(ctx, nil, nil)
	}
	return um.UpdatePassword(ctx, user, newPwd)
}

func (um *UserManager) retrieveTwoStepForgetPasswordToken(ctx context.Context, token string) (*User, error) {
	payload, err := um.forgetPwdTwoStep.Retrieve(ctx, RecoverChangePasswordPurpose, token)
	if err != nil {
		return nil, err
	}
	if payload == nil {
		return nil, v12.ErrorTwoStepFailedLocalized(ctx, nil, nil)
	}
	user, err := um.FindByID(ctx, payload.UserId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (um *UserManager) GetRoles(ctx context.Context, userId string) ([]Role, error) {
	return um.userRepo.GetRoles(ctx, userId)
}

func getUserRoleCacheKey(userId string) string {
	return fmt.Sprintf("userrole:%s", userId)
}

func (um *UserManager) RemoveUserRoleCache(ctx context.Context, userId string) error {
	return um.userRoleCache.Delete(ctx, getUserRoleCacheKey(userId))
}

func (um *UserManager) GetUserRoleIds(ctx context.Context, userId string, currentTenantOnly bool) ([]*UserRoleCacheItem_UserRole, error) {
	item, err, _ := um.userRoleCache.GetOrSet(ctx, getUserRoleCacheKey(userId), func(ctx context.Context) (*UserRoleCacheItem, error) {
		roles, err := um.GetRoles(ctx, userId)
		if err != nil {
			return nil, err
		}
		return &UserRoleCacheItem{Role: lo.Map(roles, func(t Role, _ int) *UserRoleCacheItem_UserRole {
			return &UserRoleCacheItem_UserRole{
				RoleId:   t.ID.String(),
				TenantId: t.TenantId.String,
			}
		})}, nil
	})
	if err != nil {
		return nil, err
	}
	roles := item.Role
	if currentTenantOnly {
		ct, _ := saas.FromCurrentTenant(ctx)
		roles = lo.Filter(roles, func(t *UserRoleCacheItem_UserRole, _ int) bool { return t.TenantId == ct.GetId() })
	}
	return roles, nil
}

func (um *UserManager) UpdateRoles(ctx context.Context, user *User, roles []Role) error {
	e, _ := event.NewMessageFromProto(&v13.UserRoleChangeEvent{UserId: user.ID.String()})
	user.AppendEvent(e)
	if err := um.userRepo.UpdateRoles(ctx, user, roles); err != nil {
		return err
	}
	return nil
}

func (um *UserManager) AddToRole(ctx context.Context, user *User, role *Role) error {
	e, _ := event.NewMessageFromProto(&v13.UserRoleChangeEvent{UserId: user.ID.String()})
	user.AppendEvent(e)
	if err := um.userRepo.AddToRole(ctx, user, role); err != nil {
		return err
	}
	return nil
}
func (um *UserManager) RemoveFromRole(ctx context.Context, user *User, role *Role) error {
	e, _ := event.NewMessageFromProto(&v13.UserRoleChangeEvent{UserId: user.ID.String()})
	user.AppendEvent(e)
	if err := um.userRepo.RemoveFromRole(ctx, user, role); err != nil {
		return err
	}
	return nil
}

func (um *UserManager) CheckDeleted(ctx context.Context, u *User) (bool, error) {
	if u.DeletedAt.Valid && u.DeletedAt.Time.Before(time.Now()) {
		return true, nil
	}
	return false, nil
}

func (um *UserManager) CheckLocked(ctx context.Context, u *User) (bool, error) {
	if u.LockoutEndDateUtc.After(time.Now()) {
		return true, nil
	}
	return false, nil
}

func (um *UserManager) GenerateRememberToken(ctx context.Context, duration time.Duration, uid uuid.UUID) (string, error) {
	token := NewRefreshToken(uid, duration, kithttp.ClientUserAgent(ctx), kithttp.ClientIP(ctx))
	if err := um.refreshTokenRepo.Create(ctx, token); err != nil {
		return "", err
	} else {
		return token.Token, nil
	}
}

func (um *UserManager) RefreshRememberToken(ctx context.Context, token string, duration time.Duration) (*User, string, error) {
	//find token
	currTime := time.Now()
	if t, err := um.refreshTokenRepo.Find(ctx, token, false); err != nil {
		return nil, "", err
	} else {
		if t == nil {
			return nil, "", v12.ErrorRememberTokenNotFoundLocalized(ctx, nil, nil)
		}
		if t.Used && t.Expires.After(currTime.Add(-time.Minute*3)) {
			//for some concurrency refreshing
			return nil, "", v12.ErrorRememberTokenUsedLocalized(ctx, nil, nil)
		}
		if !t.Valid() {
			return nil, "", v12.ErrorRememberTokenNotFoundLocalized(ctx, nil, nil)
		}
		//find user
		user, err := um.FindByID(ctx, t.UserId.String())
		if err != nil {
			return nil, "", err
		}
		if user == nil {
			//user not found
			return nil, "", v12.ErrorRememberTokenNotFoundLocalized(ctx, nil, nil)
		}
		//TODO check locked?
		//refresh token
		newToken, err := um.GenerateRememberToken(ctx, duration, t.UserId)
		if err != nil {
			return nil, "", err
		}
		err = um.refreshTokenRepo.Revoke(ctx, token, true)
		if err != nil {
			return user, "", err
		}
		return user, newToken, nil
	}
}

func (um *UserManager) IsInTenant(ctx context.Context, uid, tenantId string) (bool, error) {
	return um.userTenantRepo.IsIn(ctx, uid, tenantId)
}

// JoinTenant add user into tenant. safe to call when user already in
func (um *UserManager) JoinTenant(ctx context.Context, uid, tenantId string) error {
	if in, err := um.userTenantRepo.IsIn(ctx, uid, tenantId); err != nil {
		return err
	} else if in {
		return nil
	}
	_, err := um.userTenantRepo.JoinTenant(ctx, uid, tenantId, Active)
	return err
}

func (um *UserManager) RemoveFromTenant(ctx context.Context, uid, tenantId string) error {
	err := um.userTenantRepo.RemoveFromTenant(ctx, uid, tenantId)
	return err
}

func (um *UserManager) validateUser(ctx context.Context, u *User) (err error) {
	err = um.userValidator.Validate(ctx, um, u)
	return
}

func (um *UserManager) normalize(ctx context.Context, u *User) error {
	//normalize
	if u.Username != nil {
		n, err := um.lookupNormalizer.Name(ctx, *u.Username)
		if err != nil {
			return err
		}
		u.NormalizedUsername = &n
	}
	if u.Email != nil {
		e, err := um.lookupNormalizer.Email(ctx, *u.Email)
		if err != nil {
			return err
		}
		u.NormalizedEmail = &e
	}
	if u.Phone != nil {
		phone, err := um.lookupNormalizer.Phone(ctx, *u.Phone)
		if err != nil {
			return err
		}
		u.Phone = &phone
	}
	t, _ := saas.FromCurrentTenant(ctx)
	if len(t.GetId()) > 0 {
		ti := t.GetId()
		u.CreatedTenant = &ti
	}
	return nil
}
func (um *UserManager) updatePassword(ctx context.Context, u *User, password *string, validate bool) error {
	if password != nil && validate {
		if err := um.pwdValidator.Validate(ctx, *password); err != nil {
			return err
		}
	}
	if password == nil {
		u.Password = nil
		return nil
	}
	// hash password
	h, err := um.pwdHasher.HashPassword(ctx, u, *password)
	if err != nil {
		return err
	}
	u.Password = &h
	return nil
}

func (um *UserManager) checkPassword(ctx context.Context, u *User, password string) PasswordVerificationResult {
	if u.Password == nil || password == "" {
		return PasswordVerificationFail
	}
	return um.pwdHasher.VerifyHashedPassword(ctx, u, *u.Password, password)
}
