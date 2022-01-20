package service

import (
	"context"
	"github.com/ahmetb/go-linq/v3"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authorization"

	pb "github.com/goxiaoy/go-saas-kit/user/api/permission/v1"
)

type PermissionService struct {
	pb.UnimplementedPermissionServiceServer
	auth          authorization.Service
	permissionMgr authorization.PermissionManagementService
	sr            authorization.SubjectResolver
}

func NewPermissionService(auth authorization.Service, permissionMgr authorization.PermissionManagementService, sr authorization.SubjectResolver) *PermissionService {
	return &PermissionService{auth: auth, permissionMgr: permissionMgr, sr: sr}
}

func (s *PermissionService) GetCurrent(ctx context.Context, req *pb.GetCurrentPermissionRequest) (*pb.GetCurrentPermissionReply, error) {
	subjects, err := s.sr.Resolve(ctx)
	if err != nil {
		return nil, err
	}
	beans, err := s.permissionMgr.ListAcl(ctx, subjects...)
	if err != nil {
		return nil, err
	}
	var acl []*pb.Permission
	linq.From(beans).SelectT(func(bean authorization.PermissionBean) *pb.Permission {
		t := &pb.Permission{}
		mapPermissionBeanToPb(bean, t)
		return t
	}).ToSlice(&acl)
	return &pb.GetCurrentPermissionReply{Acl: acl}, nil
}

func (s *PermissionService) CheckCurrent(ctx context.Context, req *pb.CheckPermissionRequest) (*pb.CheckPermissionReply, error) {
	grant, err := s.auth.Check(ctx, authorization.NewEntityResource(req.Namespace, req.Resource), authorization.ActionStr(req.Action))
	if err != nil {
		return nil, err
	}
	effect := pb.Effect_FORBIDDEN
	if grant.Allowed {
		effect = pb.Effect_GRANT
	}
	return &pb.CheckPermissionReply{Effect: effect}, nil
}

func (s *PermissionService) CheckForSubjects(ctx context.Context, req *pb.CheckSubjectsPermissionRequest) (*pb.CheckSubjectsPermissionReply, error) {
	//TODO
	return &pb.CheckSubjectsPermissionReply{}, nil
}
func (s *PermissionService) UpdateSubjectPermission(ctx context.Context, req *pb.UpdateSubjectPermissionRequest) (*pb.UpdateSubjectPermissionResponse, error) {
	//check update permission
	if grant, err := s.auth.Check(ctx, authorization.NewEntityResource("permission", req.Subject), authorization.UpdateAction); err != nil {
		return nil, err
	} else if !grant.Allowed {
		return nil, errors.Forbidden("", "")
	}
	var acl []authorization.UpdateSubjectPermission
	linq.From(req.Acl).SelectT(func(a *pb.UpdateSubjectPermissionAcl) authorization.UpdateSubjectPermission {
		effect := authorization.EffectUnknown
		switch a.Effect {
		case pb.Effect_GRANT:
			effect = authorization.EffectGrant
			break
		case pb.Effect_FORBIDDEN:
			effect = authorization.EffectForbidden
			break
		}
		return authorization.UpdateSubjectPermission{
			Resource: authorization.NewEntityResource("", a.Resource),
			Action:   authorization.ActionStr(a.Action),
			Effect:   effect,
		}
	}).ToSlice(&acl)
	if err := s.permissionMgr.UpdateGrant(ctx, authorization.SubjectStr(req.Subject), acl); err != nil {
		return nil, err
	}
	return &pb.UpdateSubjectPermissionResponse{}, nil
}

func mapPermissionBeanToPb(bean authorization.PermissionBean, t *pb.Permission) {
	t.Subject = bean.Subject
	t.Namespace = bean.Namespace
	t.Resource = bean.Resource
	t.Action = bean.Action

	switch bean.Effect {
	case authorization.EffectUnknown:
		t.Effect = pb.Effect_UNKNOWN
		break
	case authorization.EffectGrant:
		t.Effect = pb.Effect_GRANT
		break
	case authorization.EffectForbidden:
		t.Effect = pb.Effect_FORBIDDEN
		break
	}
}
