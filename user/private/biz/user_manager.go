package biz

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/user/v1"
	"github.com/goxiaoy/go-saas/common"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"time"
)

var (
	ErrInvalidCredential     = errors.New("invalid credential")
	ErrRememberTokenNotFound = errors.New("remember token not found")
)

type UserManager struct {
	cfg              *Config
	userRepo         UserRepo
	pwdHasher        PasswordHasher
	userValidator    UserValidator
	pwdValidator     PasswordValidator
	lookupNormalizer LookupNormalizer
	userTokenRepo    UserTokenRepo
	refreshTokenRepo RefreshTokenRepo
	userTenantRepo   UserTenantRepo
	tokenFactory     UserTwoFactorTokenProviderFactory
	log              log.Logger
}

func NewUserManager(
	//cfg *Config,
	userRepo UserRepo,
	pwdHasher PasswordHasher,
	userValidator UserValidator,
	pwdValidator PasswordValidator,
	lookupNormalizer LookupNormalizer,
	userTokenRepo UserTokenRepo,
	refreshTokenRepo RefreshTokenRepo,
	userTenantRepo UserTenantRepo,
	//tokenFactory UserTwoFactorTokenProviderFactory,
	logger log.Logger) *UserManager {
	return &UserManager{
		//cfg:       cfg,
		userRepo:         userRepo,
		pwdHasher:        pwdHasher,
		userValidator:    userValidator,
		pwdValidator:     pwdValidator,
		lookupNormalizer: lookupNormalizer,
		userTokenRepo:    userTokenRepo,
		refreshTokenRepo: refreshTokenRepo,
		userTenantRepo:   userTenantRepo,
		//tokenFactory: tokenFactory,
		log: log.With(logger, "module", "/biz/user_manager")}
}

type Config struct {
}

func (um *UserManager) List(ctx context.Context, query *v1.ListUsersRequest) ([]*User, error) {
	return um.userRepo.List(ctx, query)
}

func (um *UserManager) Count(ctx context.Context, query *v1.UserFilter) (total int64, filtered int64, err error) {
	return um.userRepo.Count(ctx, query)
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
	ct, _ := common.FromCurrentTenant(ctx)
	if err := um.JoinTenant(ctx, u.UIDBase.ID.String(), ct.GetId()); err != nil {
		return err
	}
	return nil
}

func (um *UserManager) CreateWithPassword(ctx context.Context, u *User, pwd string) (err error) {
	if err = um.updatePassword(ctx, u, &pwd, true); err != nil {
		return
	}
	return um.Create(ctx, u)
}

func (um *UserManager) FindByID(ctx context.Context, id string) (user *User, err error) {
	return um.userRepo.FindByID(ctx, id)
}
func (um *UserManager) FindByName(ctx context.Context, name string) (user *User, err error) {
	name, err = um.lookupNormalizer.Name(name)
	if err != nil {
		return nil, err
	}
	return um.userRepo.FindByName(ctx, name)
}

func (um *UserManager) FindByPhone(ctx context.Context, phone string) (user *User, err error) {
	phone, err = um.lookupNormalizer.Phone(phone)
	if err != nil {
		return nil, err
	}
	return um.userRepo.FindByPhone(ctx, phone)
}

func (um *UserManager) FindByEmail(ctx context.Context, email string) (user *User, err error) {
	email, err = um.lookupNormalizer.Email(email)
	if err != nil {
		return nil, err
	}
	return um.userRepo.FindByEmail(ctx, email)
}

func (um *UserManager) FindByIdentity(ctx context.Context, identity string) (user *User, err error) {
	//try to find by id
	if uid, err := uuid.Parse(identity); err == nil {
		//
		return um.FindByID(ctx, uid.String())
	}
	if _, err := um.lookupNormalizer.Email(identity); err == nil {
		return um.FindByEmail(ctx, identity)
	}
	if phone, err := um.lookupNormalizer.Phone(identity); err == nil {
		return um.FindByPhone(ctx, phone)
	}
	return um.FindByName(ctx, identity)
}

func (um *UserManager) FindByRecoverSelector(ctx context.Context, r string) (user *User, err error) {
	return um.userRepo.FindByRecoverSelector(ctx, r)
}
func (um *UserManager) FindByConfirmSelector(ctx context.Context, c string) (user *User, err error) {
	return um.userRepo.FindByConfirmSelector(ctx, c)
}

func (um *UserManager) Update(ctx context.Context, user *User, p *fieldmaskpb.FieldMask) (err error) {
	err = um.normalize(ctx, user)
	if err != nil {
		return err
	}
	if err = um.validateUser(ctx, user); err != nil {
		return
	}
	return um.userRepo.Update(ctx, user, p)
}

func (um *UserManager) Delete(ctx context.Context, user *User) error {
	return um.userRepo.Delete(ctx, user)
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
		err := um.userRepo.Update(ctx, user, &fieldmaskpb.FieldMask{Paths: []string{"password"}})
		return err
	}
	//fail
	return ErrInvalidCredential
}

func (um *UserManager) ChangePassword(ctx context.Context, user *User, current string, newPwd string) error {
	if v := um.checkPassword(ctx, user, current); v == PasswordVerificationFail {
		return ErrInvalidCredential
	}
	if err := um.updatePassword(ctx, user, &newPwd, true); err != nil {
		return err
	}
	return um.Update(ctx, user, &fieldmaskpb.FieldMask{Paths: []string{"password"}})
}

func (um *UserManager) UpdatePassword(ctx context.Context, user *User, newPwd string) error {
	if err := um.updatePassword(ctx, user, &newPwd, true); err != nil {
		return err
	}
	return um.Update(ctx, user, &fieldmaskpb.FieldMask{Paths: []string{"password"}})
}

func (um *UserManager) GetRoles(ctx context.Context, user *User) ([]*Role, error) {
	return um.userRepo.GetRoles(ctx, user)
}

func (um *UserManager) UpdateRoles(ctx context.Context, user *User, roles []*Role) error {
	return um.userRepo.UpdateRoles(ctx, user, roles)
}
func (um *UserManager) AddToRole(ctx context.Context, user *User, role *Role) error {
	return um.userRepo.AddToRole(ctx, user, role)
}
func (um *UserManager) RemoveFromRole(ctx context.Context, user *User, role *Role) error {
	return um.userRepo.RemoveFromRole(ctx, user, role)
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

func (um *UserManager) GenerateRememberToken(ctx context.Context, uid uuid.UUID) (string, error) {
	//TODO? use refresh token?
	token := NewRefreshToken(uid, 0, server.ClientUserAgent(ctx), server.ClientIP(ctx))
	if err := um.refreshTokenRepo.Create(ctx, token); err != nil {
		return "", err
	} else {
		return token.Token, nil
	}
}

func (um *UserManager) RefreshRememberToken(ctx context.Context, uid uuid.UUID, token string) (string, error) {
	//find token
	if t, err := um.refreshTokenRepo.Find(ctx, token, true); err != nil {
		return "", err
	} else {
		if t == nil || t.UserId != uid {
			return "", ErrRememberTokenNotFound
		}
		//refresh token
		newToken, err := um.GenerateRememberToken(ctx, uid)
		if err != nil {
			return "", err
		}
		err = um.refreshTokenRepo.Revoke(ctx, token)
		if err != nil {
			return "", err
		}
		return newToken, nil
	}
}

func (um *UserManager) IsInTenant(ctx context.Context, uid, tenantId string) (bool, error) {
	return um.userTenantRepo.IsIn(ctx, uid, tenantId)
}

//JoinTenant add user into tenant. safe to call when user already in
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
		n, err := um.lookupNormalizer.Name(*u.Username)
		if err != nil {
			return err
		}
		u.NormalizedUsername = &n
	}
	if u.Email != nil {
		e, err := um.lookupNormalizer.Name(*u.Email)
		if err != nil {
			return err
		}
		u.NormalizedEmail = &e
	}
	if u.Phone != nil {
		phone, err := um.lookupNormalizer.Phone(*u.Phone)
		if err != nil {
			return err
		}
		u.Phone = &phone
	}
	t, _ := common.FromCurrentTenant(ctx)
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
