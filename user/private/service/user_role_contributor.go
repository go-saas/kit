package service

import (
	"context"
	authorization2 "github.com/goxiaoy/go-saas-kit/pkg/authz/authorization"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
)

type UserRoleContributor struct {
	userRepo biz.UserRepo
}

func NewUserRoleContributor(userRepo biz.UserRepo) *UserRoleContributor {
	return &UserRoleContributor{userRepo: userRepo}
}

func (u *UserRoleContributor) Process(ctx context.Context, subject authorization2.Subject) ([]authorization2.Subject, error) {
	if us, ok := subject.(*authorization2.UserSubject); ok {
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
			roleSubjects := make([]authorization2.Subject, len(roles))
			for i := range roles {
				roleSubjects[i] = authorization2.NewRoleSubject(roles[i].ID.String())
			}
			return roleSubjects, nil
		}
	}
	return nil, nil
}

var _ authorization2.SubjectContributor = (*UserRoleContributor)(nil)
