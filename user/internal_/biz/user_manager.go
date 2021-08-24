package biz

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
)

type UserManager struct {
	cfg              *Config
	userRepo         UserRepo
	pwdHasher        PasswordHasher
	userValidator    UserValidator
	pwdValidator     PasswordValidator
	lookupNormalizer LookupNormalizer
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
	//tokenFactory UserTwoFactorTokenProviderFactory,
	logger log.Logger) *UserManager {
	return &UserManager{
		//cfg:       cfg,
		userRepo:  userRepo,
		pwdHasher: pwdHasher,
		userValidator: userValidator,
		pwdValidator: pwdValidator,
		lookupNormalizer: lookupNormalizer,
		//tokenFactory: tokenFactory,
		log: log.With(logger, "module", "/biz/user_manager")}
}

type Config struct {
}

func (um *UserManager) List(ctx context.Context, query interface{}) ([]*User, error) {
	return um.userRepo.List(ctx,query)
}

func (um *UserManager)  Count(ctx context.Context, query interface{}) (total int64, filtered int64, err error) {
	return um.userRepo.Count(ctx,query)
}

func (um *UserManager) Create(ctx context.Context, u *User) (err error) {
	um.normalize(ctx, u)
	if err = um.validateUser(ctx, u); err != nil {
		return
	}
	return um.userRepo.Create(ctx, u)
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
	name = um.lookupNormalizer.Name(name)
	return um.userRepo.FindByName(ctx, name)
}
func (um *UserManager) FindByPhone(ctx context.Context, phone string) (user *User, err error) {
	return um.userRepo.FindByPhone(ctx, phone)
}

func (um *UserManager) FindByEmail(ctx context.Context, email string) (user *User, err error) {
	email = um.lookupNormalizer.Email(email)
	return um.userRepo.FindByEmail(ctx, email)
}

func (um *UserManager) Update(ctx context.Context, user *User) (err error) {
	um.normalize(ctx, user)
	if err = um.validateUser(ctx, user); err != nil {
		return
	}
	return um.userRepo.Update(ctx, user)
}

func (um *UserManager) Delete(ctx context.Context, user *User) error {
	return um.userRepo.Delete(ctx, user)
}

func (um *UserManager) CheckPassword(ctx context.Context, user *User, password string) (ok bool, err error) {

	v := um.checkPassword(ctx, user, password)
	if v == PasswordVerificationSuccess {
		return true, nil
	}
	if v == PasswordVerificationSuccessRehashNeeded {
		if err = um.updatePassword(ctx, user, &password, false); err != nil {
			return ok, err
		}
		err = um.userRepo.Update(ctx, user)
		return true, err
	}
	//fail
	return false, ErrInvalidPassword
}

func (um *UserManager) ChangePassword(ctx context.Context, user *User, current string, newPwd string) error {
	if v := um.checkPassword(ctx, user, current); v == PasswordVerificationFail {
		return ErrInvalidPassword
	}
	if err := um.updatePassword(ctx, user, &newPwd, true); err != nil {
		return err
	}
	return um.Update(ctx, user)
}

func (um *UserManager) UpdatePassword(ctx context.Context, user *User, newPwd string) error {
	if err := um.updatePassword(ctx, user, &newPwd, true); err != nil {
		return err
	}
	return um.Update(ctx, user)
}

func (um *UserManager) GetRoles(ctx context.Context, user *User) ([]*Role,error){
	return um.userRepo.GetRoles(ctx,user)
}

func (um *UserManager) UpdateRoles(ctx context.Context, user *User,roles []*Role)error{
	return um.userRepo.UpdateRoles(ctx,user,roles)
}
func (um *UserManager) AddToRole(ctx context.Context, user *User,role *Role) error {
	return um.userRepo.AddToRole(ctx,user,role)
}
func (um *UserManager) RemoveFromRole(ctx context.Context, user *User,role *Role) error {
	return um.userRepo.RemoveFromRole(ctx,user,role)
}

func (um *UserManager) validateUser(ctx context.Context, u *User) (err error) {
	err = um.userValidator.Validate(ctx, um, u)
	return
}

func (um *UserManager) normalize(_ context.Context, u *User) {
	//normalize
	if u.Username != nil {
		n := um.lookupNormalizer.Name(*u.Username)
		u.NormalizedUsername = &n
	}
	if u.Email != nil {
		e := um.lookupNormalizer.Name(*u.Email)
		u.NormalizedEmail = &e
	}
}
func (um *UserManager) updatePassword(ctx context.Context, u *User, password *string, validate bool) error {
	if password != nil && validate {
		if err := um.pwdValidator.Validate(ctx, u, *password); err != nil {
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
