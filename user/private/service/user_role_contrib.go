package service

import (
	"context"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/user/private/biz"
	"github.com/go-saas/saas/data"
)

type UserRoleContrib struct {
	userRepo biz.UserRepo
}

func NewUserRoleContrib(userRepo biz.UserRepo) *UserRoleContrib {
	return &UserRoleContrib{userRepo: userRepo}
}

func (u *UserRoleContrib) Process(ctx context.Context, subject authz.Subject) ([]authz.Subject, error) {
	if us, ok := authz.ParseUserSubject(subject); ok {
		if us.GetUserId() != "" {
			//TODO ?
			ctx = data.NewMultiTenancyDataFilter(ctx, false)
			user, err := u.userRepo.FindByID(ctx, us.GetUserId())
			if err != nil {
				return nil, err
			}
			if user == nil {
				return nil, nil
			}
			roleSubjects := make([]authz.Subject, len(user.Roles))
			for i, r := range user.Roles {
				roleSubjects[i] = authz.NewRoleSubject(r.ID.String())
			}
			return roleSubjects, nil
		}
	}
	return nil, nil
}

var _ authz.SubjectContrib = (*UserRoleContrib)(nil)
