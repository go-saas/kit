package service

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	"github.com/goxiaoy/go-saas/data"
)

type UserRoleContributor struct {
	userRepo biz.UserRepo
}

func NewUserRoleContributor(userRepo biz.UserRepo) *UserRoleContributor {
	return &UserRoleContributor{userRepo: userRepo}
}

func (u *UserRoleContributor) Process(ctx context.Context, subject authz.Subject) ([]authz.Subject, error) {
	if us, ok := authz.ParseUserSubject(subject); ok {
		if us.GetUserId() != "" {

			ctx = data.NewDisableMultiTenancyDataFilter(ctx)
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

var _ authz.SubjectContributor = (*UserRoleContributor)(nil)
