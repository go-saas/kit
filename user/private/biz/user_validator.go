package biz

import (
	"context"
	v1 "github.com/go-saas/kit/user/api/user/v1"
)

// UserValidator validate user before create and update
type UserValidator interface {
	Validate(ctx context.Context, um *UserManager, user *User) (err error)
}

type userValidator struct {
}

var _ UserValidator = (*userValidator)(nil)

func NewUserValidator() UserValidator {
	return &userValidator{}
}

func (u *userValidator) Validate(ctx context.Context, um *UserManager, user *User) (err error) {
	//check duplicate email/phone/username
	if user.Username != nil {
		u, err := um.FindByName(ctx, *user.Username)
		if err != nil {
			return err
		}
		if u != nil && u.ID != user.ID {
			return v1.ErrorDuplicateUsernameLocalized(ctx, nil, nil)
		}
	}
	if user.Email != nil {
		u, err := um.FindByEmail(ctx, *user.Email)
		if err != nil {
			return err
		}
		if u != nil && u.ID != user.ID {
			return v1.ErrorDuplicateEmailLocalized(ctx, nil, nil)
		}
	}
	if user.Phone != nil {
		u, err := um.FindByPhone(ctx, *user.Phone)
		if err != nil {
			return err
		}
		if u != nil && u.ID != user.ID {
			return v1.ErrorDuplicatePhoneLocalized(ctx, nil, nil)
		}
	}
	return nil
}
