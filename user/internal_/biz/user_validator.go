package biz

import (
	"context"
	"errors"
)

var (
	ErrDuplicateUserName = errors.New("duplicate userName")
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicatePhone    = errors.New("duplicate phone")
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
	if user.UserName != nil {
		u, err := um.FindByName(ctx, *user.UserName)
		if err != nil {
			return err
		}
		if u != nil && u.ID != user.ID {
			return ErrDuplicateUserName
		}
	}
	if user.Email != nil {
		u, err := um.FindByEmail(ctx, *user.Email)
		if err != nil {
			return err
		}
		if u != nil && u.ID != user.ID {
			return ErrDuplicateEmail
		}
	}
	if user.Phone != nil {
		u, err := um.FindByPhone(ctx, *user.Phone)
		if err != nil {
			return err
		}
		if u != nil && u.ID != user.ID {
			return ErrDuplicatePhone
		}
	}
	return nil
}
