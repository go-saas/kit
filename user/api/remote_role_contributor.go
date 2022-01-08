package api

import (
	"context"
	authorization2 "github.com/goxiaoy/go-saas-kit/pkg/authz/authorization"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/user/v1"
)

type RemoteRoleContributor struct {
	client v1.UserServiceClient
}

var _ authorization2.SubjectContributor = (*RemoteRoleContributor)(nil)

func NewRemoteRoleContributor(client v1.UserServiceClient) *RemoteRoleContributor {
	return &RemoteRoleContributor{
		client: client,
	}
}

func (r *RemoteRoleContributor) Process(ctx context.Context, subject authorization2.Subject) ([]authorization2.Subject, error) {
	if us, ok := subject.(*authorization2.UserSubject); ok {
		if us.GetUserId() != "" {
			roles, err := r.client.GetUserRoles(ctx, &v1.GetUserRoleRequest{Id: us.GetUserId()})
			if err != nil {
				return nil, err
			}
			roleSubjects := make([]authorization2.Subject, len(roles.Roles))
			for i := range roles.Roles {
				roleSubjects[i] = authorization2.NewRoleSubject(roles.Roles[i].Id)
			}
			return roleSubjects, nil
		}
	}
	return nil, nil
}
