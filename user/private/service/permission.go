package service

import (
	"context"
	api2 "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authz/authz"
	pb "github.com/go-saas/kit/user/api/permission/v1"
	"github.com/go-saas/kit/user/util"
	"github.com/samber/lo"
)

type PermissionService struct {
	pb.UnimplementedPermissionServiceServer
	auth          authz.Service
	permissionMgr authz.PermissionManagementService
	sr            authz.SubjectResolver
	trust         api2.TrustedContextValidator
}

var _ pb.PermissionServiceServer = (*PermissionService)(nil)

func NewPermissionService(auth authz.Service, permissionMgr authz.PermissionManagementService, sr authz.SubjectResolver, trust api2.TrustedContextValidator) *PermissionService {
	return &PermissionService{auth: auth, permissionMgr: permissionMgr, sr: sr, trust: trust}
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

	acl := lo.Map(beans, func(bean authz.PermissionBean, _ int) *pb.Permission {
		t := &pb.Permission{}
		util.MapPermissionBeanToPb(bean, t)
		return t
	})
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
