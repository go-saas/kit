package service

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/authorization/authorization"
	"github.com/goxiaoy/go-saas-kit/user/internal_/biz"
)

type UserRoleContributor struct {
	userRepo biz.UserRepo
}

func NewUserRoleContributor(userRepo biz.UserRepo) *UserRoleContributor {
	return &UserRoleContributor{userRepo: userRepo}
}

func (u *UserRoleContributor) Process(ctx context.Context, subject authorization.Subject) ([]authorization.Subject, error) {
	if us, ok := subject.(*authorization.UserSubject); ok {
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
			roleSubjects := make([]authorization.Subject, len(roles))
			for i := range roles {
				roleSubjects[i] = authorization.NewRoleSubject(roles[i].ID.String())
			}
			return roleSubjects, nil
		}
	}
	return nil, nil
}

var _ authorization.SubjectContributor = (*UserRoleContributor)(nil)
