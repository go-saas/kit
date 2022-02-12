package service

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
)

type UserRoleContributor struct {
	userRepo biz.UserRepo
}

func NewUserRoleContributor(userRepo biz.UserRepo) *UserRoleContributor {
	return &UserRoleContributor{userRepo: userRepo}
}

func (u *UserRoleContributor) Process(ctx context.Context, subject authz.Subject) ([]authz.Subject, error) {
	if us, ok := subject.(*authz.UserSubject); ok {
		if us.GetUserId() != "" {
			user, err := u.userRepo.FindByID(ctx, us.GetUserId())
			if err != nil {
				return nil, err
			}
			if user == nil {
				return nil, nil
			}
			roles, err := u.userRepo.GetRoles(ctx, user)
			if err != nil {
				return nil, err
			}
			roleSubjects := make([]authz.Subject, len(roles))
			for i := range roles {
				roleSubjects[i] = authz.NewRoleSubject(roles[i].ID.String())
			}
			return roleSubjects, nil
		}
	}
	return nil, nil
}

var _ authz.SubjectContributor = (*UserRoleContributor)(nil)
