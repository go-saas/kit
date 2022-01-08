package service

import (
	"context"
	errors2 "github.com/go-kratos/kratos/v2/errors"
	authorization2 "github.com/goxiaoy/go-saas-kit/pkg/authz/authorization"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	"github.com/mennanov/fmutils"

	pb "github.com/goxiaoy/go-saas-kit/user/api/role/v1"
)

type RoleService struct {
	repo biz.RoleRepo
	auth authorization2.Service
	pb.UnimplementedRoleServiceServer
}

func NewRoleServiceService(repo biz.RoleRepo, auth authorization2.Service) *RoleService {
	return &RoleService{repo: repo, auth: auth}
}

func (s *RoleService) ListRoles(ctx context.Context, req *pb.ListRolesRequest) (*pb.ListRolesResponse, error) {
	if authResult, err := s.auth.Check(ctx, authorization2.NewEntityResource("user.role", "*"), authorization2.ListAction); err != nil {
		return nil, err
	} else if !authResult.Allowed {
		return nil, errors2.Forbidden("", "")
	}
	ret := &pb.ListRolesResponse{}
	totalCount, filterCount, err := s.repo.Count(ctx, req.Filter)
	if err != nil {
		return nil, err
	}
	ret.TotalSize = int32(totalCount)
	ret.FilterSize = int32(filterCount)

	items, err := s.repo.List(ctx, req)
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
		return nil, errors2.BadRequest("", "id or name can not be empty")
	}
	var u *biz.Role
	var err error
	if req.Id != "" {
		u, err = s.repo.FindById(ctx, req.Id)
		if err != nil {
			return nil, err
		}
	}
	if req.Name != "" {
		u, err = s.repo.FindByName(ctx, req.Name)
		if err != nil {
			return nil, err
		}
	}
	if u == nil {
		return nil, errors2.Forbidden("", "")
	}
	if authResult, err := s.auth.Check(ctx, authorization2.NewEntityResource("user.role", u.ID.String()), authorization2.GetAction); err != nil {
		return nil, err
	} else if !authResult.Allowed {
		return nil, errors2.Forbidden("", "")
	}
	res := &pb.Role{}
	MapBizRoleToApi(u, res)
	return res, nil
}

func (s *RoleService) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.Role, error) {
	return &pb.Role{}, nil
}
func (s *RoleService) UpdateRole(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.Role, error) {
	return &pb.Role{}, nil
}
func (s *RoleService) DeleteRole(ctx context.Context, req *pb.DeleteRoleRequest) (*pb.DeleteRoleResponse, error) {
	return &pb.DeleteRoleResponse{}, nil
}

func MapBizRoleToApi(u *biz.Role, b *pb.Role) {
	b.Id = u.ID.String()
	b.Name = u.Name
}
