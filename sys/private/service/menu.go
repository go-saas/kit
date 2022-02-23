package service

import (
	"context"
	"github.com/ahmetb/go-linq/v3"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/sys/private/biz"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/permission/v1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"

	pb "github.com/goxiaoy/go-saas-kit/sys/api/menu/v1"
)

type MenuService struct {
	pb.UnimplementedMenuServiceServer
	auth authz.Service
	repo biz.MenuRepo
}

func NewMenuService(auth authz.Service, repo biz.MenuRepo) *MenuService {
	return &MenuService{auth: auth, repo: repo}
}

func (s *MenuService) ListMenu(ctx context.Context, req *pb.ListMenuRequest) (*pb.ListMenuReply, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource("sys.menu", "*"), authz.ListAction); err != nil {
		return nil, err
	}
	ret := &pb.ListMenuReply{}

	totalCount, filterCount, err := s.repo.Count(ctx, req.Search, req.Filter)
	ret.TotalSize = int32(totalCount)
	ret.FilterSize = int32(filterCount)

	if err != nil {
		return ret, err
	}
	items, err := s.repo.List(ctx, req)
	if err != nil {
		return ret, err
	}
	rItems := make([]*pb.Menu, len(items))

	linq.From(items).SelectT(func(g *biz.Menu) *pb.Menu {
		b := &pb.Menu{}
		MapBizMenu2Pb(g, b)
		return b
	}).ToSlice(&rItems)

	ret.Items = rItems
	return ret, nil
}
func (s *MenuService) GetMenu(ctx context.Context, req *pb.GetMenuRequest) (*pb.Menu, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource("sys.menu", req.Id), authz.GetAction); err != nil {
		return nil, err
	}
	g, err := s.repo.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.NotFound("", "")
	}
	res := &pb.Menu{}
	MapBizMenu2Pb(g, res)
	return res, nil
}
func (s *MenuService) CreateMenu(ctx context.Context, req *pb.CreateMenuRequest) (*pb.Menu, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource("sys.menu", "*"), authz.CreateAction); err != nil {
		return nil, err
	}

	//check duplicate name
	if dbP, err := s.repo.FindByName(ctx, normalizeName(req.Name)); err != nil {
		return nil, err
	} else if dbP != nil {
		return nil, pb.ErrorMenuNameDuplicate("", "")
	}
	e := &biz.Menu{}
	MapCreatePbMenu2Biz(req, e)
	err := s.repo.Create(ctx, e)
	if err != nil {
		return nil, err
	}
	res := &pb.Menu{}
	MapBizMenu2Pb(e, res)
	return res, nil
}
func (s *MenuService) UpdateMenu(ctx context.Context, req *pb.UpdateMenuRequest) (*pb.Menu, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource("sys.menu", req.Menu.Id), authz.UpdateAction); err != nil {
		return nil, err
	}

	//check duplicate name
	if dbP, err := s.repo.FindByName(ctx, normalizeName(req.Menu.Name)); err != nil {
		return nil, err
	} else if dbP != nil && dbP.ID.String() != req.Menu.Id {
		return nil, pb.ErrorMenuNameDuplicate("", "")
	}

	g, err := s.repo.Get(ctx, req.Menu.Id)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.NotFound("", "")
	}
	MapUpdatePbMenu2Biz(req.Menu, g)
	if err := s.repo.Update(ctx, g, nil); err != nil {
		return nil, err
	}
	res := &pb.Menu{}
	MapBizMenu2Pb(g, res)
	return res, nil
}

func (s *MenuService) DeleteMenu(ctx context.Context, req *pb.DeleteMenuRequest) (*pb.DeleteMenuReply, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource("sys.menu", req.Id), authz.DeleteAction); err != nil {
		return nil, err
	}
	g, err := s.repo.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.NotFound("", "")
	}
	err = s.repo.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteMenuReply{Id: g.ID.String(), Name: g.Name}, nil
}

func MapBizMenu2Pb(a *biz.Menu, b *pb.Menu) {
	b.Id = a.ID.String()
	b.Name = a.Name
	b.Desc = a.Desc
	b.CreatedAt = timestamppb.New(a.CreatedAt)
	b.Component = a.Component
	var requirement []*v1.PermissionRequirement
	linq.From(a.Requirement).SelectT(func(a biz.MenuPermissionRequirement) *v1.PermissionRequirement {
		ret := &v1.PermissionRequirement{
			Namespace: a.Namespace,
			Resource:  a.Resource,
			Action:    a.Action,
		}
		return ret
	}).ToSlice(&requirement)
	b.Requirement = requirement
	b.Parent = a.Parent
	if a.Props != nil {
		b.Props, _ = structpb.NewStruct(a.Props)
	}
	b.FullPath = a.FullPath
	b.Priority = a.Priority
	b.IgnoreAuth = a.IgnoreAuth
	b.Icon = a.Icon
	b.Iframe = a.Iframe
	b.MicroApp = a.MicroApp
	if a.Meta != nil {
		b.Meta, _ = structpb.NewStruct(a.Meta)
	}
	b.Title = a.Title
}

func MapListBizMenu2Pb(a []biz.Menu, b *[]*pb.Menu) {
	linq.From(a).SelectT(func(c biz.Menu) *pb.Menu {
		ca := pb.Menu{}
		MapBizMenu2Pb(&c, &ca)
		return &ca
	}).ToSlice(b)
}

func MapUpdatePbMenu2Biz(a *pb.UpdateMenu, b *biz.Menu) {

	b.Name = normalizeName(a.Name)
	b.Desc = a.Desc

	b.Component = a.Component
	var requirement []biz.MenuPermissionRequirement
	linq.From(a.Requirement).SelectT(func(a *v1.PermissionRequirement) biz.MenuPermissionRequirement {
		ret := biz.MenuPermissionRequirement{
			Namespace: a.Namespace,
			Resource:  a.Resource,
			Action:    a.Action,
		}
		return ret
	}).ToSlice(&requirement)
	b.Requirement = requirement
	b.Parent = a.Parent
	if a.Props != nil {
		b.Props = a.Props.AsMap()
	}
	b.FullPath = a.FullPath
	b.Priority = a.Priority
	b.IgnoreAuth = a.IgnoreAuth
	b.Icon = a.Icon
	b.Iframe = a.Iframe
	b.MicroApp = a.MicroApp
	if a.Meta != nil {
		b.Meta = a.Meta.AsMap()
	}
	b.Title = a.Title
}

func MapCreatePbMenu2Biz(a *pb.CreateMenuRequest, b *biz.Menu) {

	b.Name = normalizeName(a.Name)
	b.Desc = a.Desc

	b.Component = a.Component
	var requirement []biz.MenuPermissionRequirement
	linq.From(a.Requirement).SelectT(func(a *v1.PermissionRequirement) biz.MenuPermissionRequirement {
		ret := biz.MenuPermissionRequirement{
			Namespace: a.Namespace,
			Resource:  a.Resource,
			Action:    a.Action,
		}
		return ret
	}).ToSlice(&requirement)
	b.Requirement = requirement
	b.Parent = a.Parent
	if a.Props != nil {
		b.Props = a.Props.AsMap()
	}
	b.FullPath = a.FullPath
	b.Priority = a.Priority
	b.IgnoreAuth = a.IgnoreAuth
	b.Icon = a.Icon
	b.Iframe = a.Iframe
	b.MicroApp = a.MicroApp
	if a.Meta != nil {
		b.Meta = a.Meta.AsMap()
	}
	b.Title = a.Title
}

func normalizeName(name string) string {
	return strings.ToLower(name)
}
