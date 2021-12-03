package api

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/authorization/authorization"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/user/v1"
)

type RemoteRoleContributor struct {
	client v1.UserServiceClient
}
var _ authorization.SubjectContributor = (*RemoteRoleContributor)(nil)

func NewRemoteRoleContributor(	client v1.UserServiceClient) *RemoteRoleContributor  {
	return &RemoteRoleContributor{
		client: client,
	}
}

func (r *RemoteRoleContributor) Process(ctx context.Context, subject authorization.Subject) ([]authorization.Subject, error) {
	if us, ok := subject.(*authorization.UserSubject); ok {
		if us.GetUserId() != "" {
			roles,err:=r.client.GetUserRoles(ctx,&v1.GetUserRoleRequest{Id: us.GetUserId()})
			if err != nil {
				return nil, err
			}
			roleSubjects := make([]authorization.Subject, len(roles.Roles))
			for i := range roles.Roles {
				roleSubjects[i] = authorization.NewRoleSubject(roles.Roles[i].Id)
			}
			return roleSubjects, nil
		}
	}
	return nil, nil
}
