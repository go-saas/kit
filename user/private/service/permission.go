package service

import (
	"context"
	"github.com/ahmetb/go-linq/v3"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"

	pb "github.com/goxiaoy/go-saas-kit/user/api/permission/v1"
)

type PermissionService struct {
	pb.UnimplementedPermissionServiceServer
	auth          authz.Service
	permissionMgr authz.PermissionManagementService
	sr            authz.SubjectResolver
}

func NewPermissionService(auth authz.Service, permissionMgr authz.PermissionManagementService, sr authz.SubjectResolver) *PermissionService {
	return &PermissionService{auth: auth, permissionMgr: permissionMgr, sr: sr}
}

func (s *PermissionService) GetCurrent(ctx context.Context, req *pb.GetCurrentPermissionRequest) (*pb.GetCurrentPermissionReply, error) {
	subjects, err := s.sr.ResolveFromContext(ctx)
	if err != nil {
		return nil, err
	}
	newSubjects, err := s.sr.ResolveProcessed(ctx, subjects...)
	if err != nil {
		return nil, err
	}
	beans, err := s.permissionMgr.ListAcl(ctx, newSubjects...)
	if err != nil {
		return nil, err
	}
	var acl []*pb.Permission
	linq.From(beans).SelectT(func(bean authz.PermissionBean) *pb.Permission {
		t := &pb.Permission{}
		mapPermissionBeanToPb(bean, t)
		return t
	}).ToSlice(&acl)
	return &pb.GetCurrentPermissionReply{Acl: acl}, nil
}

func (s *PermissionService) CheckCurrent(ctx context.Context, req *pb.CheckPermissionRequest) (*pb.CheckPermissionReply, error) {
	grant, err := s.auth.Check(ctx, authz.NewEntityResource(req.Namespace, req.Resource), authz.ActionStr(req.Action))
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
	if grant, err := s.auth.Check(ctx, authz.NewEntityResource("permission", "*"), authz.GetAction); err != nil {
		return nil, err
	} else if !grant.Allowed {
		return nil, errors.Forbidden("", "")
	}
	subjects := make([]authz.Subject, len(req.Subjects))
	for i, subject := range req.Subjects {
		subjects[i] = authz.SubjectStr(subject)
	}
	grant, err := s.auth.CheckForSubjects(ctx, authz.NewEntityResource(req.Namespace, req.Resource), authz.ActionStr(req.Action), subjects...)
	if err != nil {
		return nil, err
	}
	effect := pb.Effect_FORBIDDEN
	if grant.Allowed {
		effect = pb.Effect_GRANT
	}
	return &pb.CheckSubjectsPermissionReply{Effect: effect}, nil
}
func (s *PermissionService) UpdateSubjectPermission(ctx context.Context, req *pb.UpdateSubjectPermissionRequest) (*pb.UpdateSubjectPermissionResponse, error) {
	//check update permission
	if grant, err := s.auth.Check(ctx, authz.NewEntityResource("permission", req.Subject), authz.UpdateAction); err != nil {
		return nil, err
	} else if !grant.Allowed {
		return nil, errors.Forbidden("", "")
	}
	var acl []authz.UpdateSubjectPermission
	linq.From(req.Acl).SelectT(func(a *pb.UpdateSubjectPermissionAcl) authz.UpdateSubjectPermission {
		effect := authz.EffectUnknown
		switch a.Effect {
		case pb.Effect_GRANT:
			effect = authz.EffectGrant
			break
		case pb.Effect_FORBIDDEN:
			effect = authz.EffectForbidden
			break
		}
		return authz.UpdateSubjectPermission{
			Resource: authz.NewEntityResource("", a.Resource),
			Action:   authz.ActionStr(a.Action),
			Effect:   effect,
		}
	}).ToSlice(&acl)
	if err := s.permissionMgr.UpdateGrant(ctx, authz.SubjectStr(req.Subject), acl); err != nil {
		return nil, err
	}
	return &pb.UpdateSubjectPermissionResponse{}, nil
}

func mapPermissionBeanToPb(bean authz.PermissionBean, t *pb.Permission) {
	t.Subject = bean.Subject
	t.Namespace = bean.Namespace
	t.Resource = bean.Resource
	t.Action = bean.Action

	switch bean.Effect {
	case authz.EffectUnknown:
		t.Effect = pb.Effect_UNKNOWN
		break
	case authz.EffectGrant:
		t.Effect = pb.Effect_GRANT
		break
	case authz.EffectForbidden:
		t.Effect = pb.Effect_FORBIDDEN
		break
	}
}
