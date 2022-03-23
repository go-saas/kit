package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/user/api"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/permission/v1"
	pb "github.com/goxiaoy/go-saas-kit/user/api/role/v1"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	"github.com/goxiaoy/go-saas-kit/user/util"
	"github.com/goxiaoy/go-saas/common"
	"github.com/mennanov/fmutils"
	"github.com/samber/lo"
)

type RoleService struct {
	mgr           *biz.RoleManager
	auth          authz.Service
	permissionMgr authz.PermissionManagementService
	pb.UnimplementedRoleServiceServer
}

func NewRoleServiceService(repo *biz.RoleManager, auth authz.Service, permissionMgr authz.PermissionManagementService) *RoleService {
	return &RoleService{mgr: repo, auth: auth, permissionMgr: permissionMgr}
}

func (s *RoleService) ListRoles(ctx context.Context, req *pb.ListRolesRequest) (*pb.ListRolesResponse, error) {
	if authResult, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceRole, "*"), authz.ReadAction); err != nil {
		return nil, err
	} else if !authResult.Allowed {
		return nil, errors.Forbidden("", "")
	}
	ret := &pb.ListRolesResponse{}
	totalCount, filterCount, err := s.mgr.Count(ctx, req)
	if err != nil {
		return nil, err
	}
	ret.TotalSize = int32(totalCount)
	ret.FilterSize = int32(filterCount)

	items, err := s.mgr.List(ctx, req)
	if err != nil {
		return ret, err
	}
	rItems := make([]*pb.Role, len(items))
	for index, u := range items {
		res := &pb.Role{}
		MapBizRoleToApi(u, res)
		if req.Fields != nil {
			fmutils.Filter(res, req.Fields.Paths)
		}
		rItems[index] = res
	}
	ret.Items = rItems
	return ret, nil
}

func (s *RoleService) GetRole(ctx context.Context, req *pb.GetRoleRequest) (*pb.Role, error) {
	if req.Id == "" && req.Name == "" {
		return nil, errors.BadRequest("", "id or name can not be empty")
	}
	var u *biz.Role
	var err error
	if req.Id != "" {
		u, err = s.mgr.Get(ctx, req.Id)
		if err != nil {
			return nil, err
		}
	}
	if req.Name != "" {
		u, err = s.mgr.FindByName(ctx, req.Name)
		if err != nil {
			return nil, err
		}
	}
	if u == nil {
		return nil, errors.Forbidden("", "")
	}
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceRole, u.ID.String()), authz.ReadAction); err != nil {
		return nil, err
	}
	res := &pb.Role{}
	MapBizRoleToApi(u, res)
	return res, nil
}

func (s *RoleService) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.Role, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceRole, "*"), authz.CreateAction); err != nil {
		return nil, err
	}
	r := &biz.Role{
		Name:        req.Name,
		IsPreserved: false,
	}
	if err := s.mgr.Create(ctx, r); err != nil {
		return nil, err
	}
	ret := &pb.Role{}
	MapBizRoleToApi(r, ret)
	return ret, nil
}
func (s *RoleService) UpdateRole(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.Role, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceRole, req.Role.Id), authz.UpdateAction); err != nil {
		return nil, err
	}
	r, err := s.mgr.Get(ctx, req.Role.Id)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, errors.NotFound("", "")
	}
	r.Name = req.Role.Name
	if err := s.mgr.Update(ctx, r.ID.String(), r, nil); err != nil {
		return nil, err
	}
	ret := &pb.Role{}
	MapBizRoleToApi(r, ret)
	return ret, nil
}

func (s *RoleService) DeleteRole(ctx context.Context, req *pb.DeleteRoleRequest) (*pb.DeleteRoleResponse, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceRole, req.Id), authz.DeleteAction); err != nil {
		return nil, err
	}
	if err := s.mgr.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &pb.DeleteRoleResponse{}, nil
}
func (s *RoleService) GetRolePermission(ctx context.Context, req *pb.GetRolePermissionRequest) (*pb.GetRolePermissionResponse, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceRole, req.Id), authz.ReadAction); err != nil {
		return nil, err
	}
	r, err := s.mgr.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, errors.NotFound("", "")
	}

	acl, err := s.permissionMgr.ListAcl(ctx, authz.NewRoleSubject(req.Id))
	if err != nil {
		return nil, err
	}
	resItems := make([]*v1.Permission, len(acl))
	for i, bean := range acl {
		r := &v1.Permission{}
		util.MapPermissionBeanToPb(bean, r)
		resItems[i] = r
	}
	return &pb.GetRolePermissionResponse{
		Acl: resItems,
	}, nil
}

func (s *RoleService) UpdateRolePermission(ctx context.Context, req *pb.UpdateRolePermissionRequest) (*pb.UpdateRolePermissionResponse, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceRole, req.Id), authz.ReadAction); err != nil {
		return nil, err
	}
	r, err := s.mgr.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, errors.NotFound("", "")
	}
	if r.IsPreserved {
		return nil, pb.ErrorRolePreserved("")
	}
	ti, _ := common.FromCurrentTenant(ctx)
	var acl = lo.Map(req.Acl, func(a *pb.UpdateRolePermissionAcl, _ int) authz.UpdateSubjectPermission {
		effect := util.MapPbEffect2AuthEffect(a.Effect)
		return authz.UpdateSubjectPermission{
			Resource: authz.NewEntityResource(a.Namespace, a.Resource),
			Action:   authz.ActionStr(a.Action),
			TenantID: ti.GetId(),
			Effect:   effect,
		}
	})
	if err := s.permissionMgr.UpdateGrant(ctx, authz.NewRoleSubject(req.Id), acl); err != nil {
		return nil, err
	}
	return &pb.UpdateRolePermissionResponse{}, nil
}

func MapBizRoleToApi(u *biz.Role, b *pb.Role) {
	b.Id = u.ID.String()
	b.Name = u.Name
	b.IsPreserved = u.IsPreserved
}
