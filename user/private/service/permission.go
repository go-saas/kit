package service

import (
	"context"
	"github.com/ahmetb/go-linq/v3"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	pb "github.com/goxiaoy/go-saas-kit/user/api/permission/v1"
	"github.com/goxiaoy/go-saas-kit/user/util"
	"github.com/goxiaoy/go-saas/common"
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
		util.MapPermissionBeanToPb(bean, t)
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
	if _, err := s.auth.Check(ctx, authz.NewEntityResource("permission", "*"), authz.GetAction); err != nil {
		return nil, err
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

func (s *PermissionService) AddSubjectPermission(ctx context.Context, req *pb.AddSubjectPermissionRequest) (*pb.AddSubjectPermissionResponse, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource("permission", req.Subject), authz.CreateAction); err != nil {
		return nil, err
	}
	if err := s.permissionMgr.AddGrant(ctx, authz.NewEntityResource(req.Namespace, req.Resource),
		authz.ActionStr(req.Action), authz.SubjectStr(req.Subject), req.TenantId, util.MapPbEffect2AuthEffect(req.Effect)); err != nil {
		return nil, err
	}
	return &pb.AddSubjectPermissionResponse{}, nil
}
func (s *PermissionService) ListSubjectPermission(ctx context.Context, req *pb.ListSubjectPermissionRequest) (*pb.ListSubjectPermissionResponse, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource("permission", "*"), authz.ListAction); err != nil {
		return nil, err
	}
	subs := make([]authz.Subject, len(req.Subjects))
	for i, subject := range req.Subjects {
		subs[i] = authz.SubjectStr(subject)
	}
	acl, err := s.permissionMgr.ListAcl(ctx, subs...)
	if err != nil {
		return nil, err
	}
	resItems := make([]*pb.Permission, len(acl))
	for i, bean := range acl {
		r := &pb.Permission{}
		util.MapPermissionBeanToPb(bean, r)
		resItems[i] = r
	}
	return &pb.ListSubjectPermissionResponse{
		Acl: resItems,
	}, nil
}

func (s *PermissionService) UpdateSubjectPermission(ctx context.Context, req *pb.UpdateSubjectPermissionRequest) (*pb.UpdateSubjectPermissionResponse, error) {
	//check update permission
	if _, err := s.auth.Check(ctx, authz.NewEntityResource("permission", req.Subject), authz.UpdateAction); err != nil {
		return nil, err
	}
	var acl []authz.UpdateSubjectPermission
	linq.From(req.Acl).SelectT(func(a *pb.UpdateSubjectPermissionAcl) authz.UpdateSubjectPermission {
		effect := util.MapPbEffect2AuthEffect(a.Effect)
		return authz.UpdateSubjectPermission{
			Resource: authz.NewEntityResource(a.Namespace, a.Resource),
			Action:   authz.ActionStr(a.Action),
			TenantID: normalizeTenantId(ctx, a.TenantId),
			Effect:   effect,
		}
	}).ToSlice(&acl)
	if err := s.permissionMgr.UpdateGrant(ctx, authz.SubjectStr(req.Subject), acl); err != nil {
		return nil, err
	}
	return &pb.UpdateSubjectPermissionResponse{}, nil
}

func (s *PermissionService) RemoveSubjectPermission(ctx context.Context, req *pb.RemoveSubjectPermissionRequest) (*pb.RemoveSubjectPermissionReply, error) {
	//check delete permission
	if _, err := s.auth.Check(ctx, authz.NewEntityResource("permission", req.Subject), authz.DeleteAction); err != nil {
		return nil, err
	}
	effList := make([]authz.Effect, len(req.Effects))
	for i, effect := range req.Effects {
		effList[i] = util.MapPbEffect2AuthEffect(effect)
	}
	if err := s.permissionMgr.RemoveGrant(ctx, authz.NewEntityResource(req.Namespace, req.Resource), authz.ActionStr(req.Action), authz.SubjectStr(req.Subject), normalizeTenantId(ctx, req.TenantId), effList); err != nil {
		return nil, err
	}
	return &pb.RemoveSubjectPermissionReply{}, nil
}

func normalizeTenantId(ctx context.Context, tenantId string) string {
	ti := common.FromCurrentTenant(ctx)
	if ti.GetId() == "" {
		//host side
		return tenantId
	}
	return ti.GetId()
}
