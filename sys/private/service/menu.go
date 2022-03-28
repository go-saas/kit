package service

import (
	"context"
	"fmt"
	klog "github.com/go-kratos/kratos/v2/log"
	"strings"

	"github.com/goxiaoy/go-saas-kit/sys/api"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/goxiaoy/go-saas-kit/pkg/authn"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/sys/private/biz"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/permission/v1"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/goxiaoy/go-saas-kit/sys/api/menu/v1"
)

type MenuService struct {
	pb.UnimplementedMenuServiceServer
	auth   authz.Service
	repo   biz.MenuRepo
	logger *klog.Helper
}

func NewMenuService(auth authz.Service, repo biz.MenuRepo, logger klog.Logger) *MenuService {
	return &MenuService{auth: auth, repo: repo, logger: klog.NewHelper(klog.With(logger, "module", "MenuService"))}
}

func (s *MenuService) ListMenu(ctx context.Context, req *pb.ListMenuRequest) (*pb.ListMenuReply, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceMenu, "*"), authz.ReadAction); err != nil {
		return nil, err
	}
	ret := &pb.ListMenuReply{}

	totalCount, filterCount, err := s.repo.Count(ctx, req)
	ret.TotalSize = int32(totalCount)
	ret.FilterSize = int32(filterCount)

	if err != nil {
		return ret, err
	}
	items, err := s.repo.List(ctx, req)
	if err != nil {
		return ret, err
	}
	rItems := lo.Map(items, func(g *biz.Menu, _ int) *pb.Menu {
		b := &pb.Menu{}
		MapBizMenu2Pb(g, b)
		return b
	})

	ret.Items = rItems
	return ret, nil
}
func (s *MenuService) GetMenu(ctx context.Context, req *pb.GetMenuRequest) (*pb.Menu, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceMenu, req.Id), authz.ReadAction); err != nil {
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
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceMenu, "*"), authz.CreateAction); err != nil {
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
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceMenu, req.Menu.Id), authz.UpdateAction); err != nil {
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
	if g.IsPreserved {
		return nil, pb.ErrorMenuPreserved("", "")
	}
	MapUpdatePbMenu2Biz(req.Menu, g)
	if err := s.repo.Update(ctx, g.ID.String(), g, nil); err != nil {
		return nil, err
	}
	res := &pb.Menu{}
	MapBizMenu2Pb(g, res)
	return res, nil
}

func (s *MenuService) DeleteMenu(ctx context.Context, req *pb.DeleteMenuRequest) (*pb.DeleteMenuReply, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceMenu, req.Id), authz.DeleteAction); err != nil {
		return nil, err
	}
	g, err := s.repo.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.NotFound("", "")
	}
	if g.IsPreserved {
		return nil, pb.ErrorMenuPreserved("", "")
	}
	err = s.repo.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteMenuReply{Id: g.ID.String(), Name: g.Name}, nil
}

func (s *MenuService) GetAvailableMenus(ctx context.Context, req *pb.GetAvailableMenusRequest) (*pb.GetAvailableMenusReply, error) {
	//allow public call
	items, err := s.repo.List(ctx, &pb.ListMenuRequest{
		PageOffset: 0,
		PageSize:   -1,
	})
	if err != nil {
		return nil, err
	}
	//filter by permission
	var filter []*biz.Menu
	var disAllowMenuId []string

	var waitForCheckerRequirements = []lo.Tuple2[string, []biz.MenuPermissionRequirement]{}

	for _, item := range items {
		if item.IgnoreAuth {
			filter = append(filter, item)
			continue
		}
		if len(item.Requirement) > 0 {
			waitForCheckerRequirements = append(waitForCheckerRequirements, lo.Tuple2[string, []biz.MenuPermissionRequirement]{A: item.ID.String(), B: item.Requirement})

		} else {
			//just check if login
			if ui, ok := authn.FromUserContext(ctx); ok && len(ui.GetId()) > 0 {
				filter = append(filter, item)
			} else {
				disAllowMenuId = append(disAllowMenuId, item.ID.String())
			}
		}
	}
	requirementConv := func(t biz.MenuPermissionRequirement, _ int) *authz.Requirement {
		return authz.NewRequirement(authz.NewEntityResource(t.Namespace, t.Resource), authz.ActionStr(t.Action))
	}
	requirementKeyFunc := func(r *authz.Requirement) string {
		return fmt.Sprintf("%s/%s@%s", r.Resource.GetNamespace(), r.Resource.GetIdentity(), r.Action.GetIdentity())
	}
	if len(waitForCheckerRequirements) > 0 {
		//check
		allReqEffectMap := map[string]*authz.Result{}
		rl := lo.UniqBy(lo.Map(lo.FlatMap(waitForCheckerRequirements, func(t lo.Tuple2[string, []biz.MenuPermissionRequirement], _ int) []biz.MenuPermissionRequirement {
			return t.B
		}), requirementConv), requirementKeyFunc)
		grantList, err := s.auth.BatchCheck(ctx, rl)
		if err != nil {
			return nil, err
		}
		for i, r := range rl {
			allReqEffectMap[requirementKeyFunc(r)] = grantList[i]
		}
		for _, menuRequirements := range waitForCheckerRequirements {
			for _, mr := range menuRequirements.B {
				if !allReqEffectMap[requirementKeyFunc(requirementConv(mr, 0))].Allowed {
					disAllowMenuId = append(disAllowMenuId, menuRequirements.A)
					continue
				}
			}
		}
	}

	var filterChildren []*biz.Menu
	for _, f := range filter {
		allow := true
		for _, dis := range disAllowMenuId {
			if f.Parent == dis {
				allow = false
				break
			}
		}
		if allow {
			filterChildren = append(filterChildren, f)
		}
	}

	var retItems = lo.Map(filterChildren, func(a *biz.Menu, _ int) *pb.Menu {
		ret := &pb.Menu{}
		MapBizMenu2Pb(a, ret)
		return ret
	})
	return &pb.GetAvailableMenusReply{Items: retItems}, nil
}
func MapBizMenu2Pb(a *biz.Menu, b *pb.Menu) {
	b.Id = a.ID.String()
	b.Name = a.Name
	b.Desc = a.Desc
	b.CreatedAt = timestamppb.New(a.CreatedAt)
	b.Component = a.Component

	requirement := lo.Map(a.Requirement, func(a biz.MenuPermissionRequirement, _ int) *v1.PermissionRequirement {
		ret := &v1.PermissionRequirement{
			Namespace: a.Namespace,
			Resource:  a.Resource,
			Action:    a.Action,
		}
		return ret
	})

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
	b.Path = a.Path
	b.Redirect = a.Redirect
}

func MapUpdatePbMenu2Biz(a *pb.UpdateMenu, b *biz.Menu) {

	b.Name = normalizeName(a.Name)
	b.Desc = a.Desc

	b.Component = a.Component
	requirement := lo.Map(a.Requirement, func(a *v1.PermissionRequirement, _ int) biz.MenuPermissionRequirement {
		ret := biz.MenuPermissionRequirement{
			Namespace: a.Namespace,
			Resource:  a.Resource,
			Action:    a.Action,
		}
		return ret
	})

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
	b.Title = a.Title
	b.Path = a.Path
	b.Redirect = a.Redirect
}

func MapCreatePbMenu2Biz(a *pb.CreateMenuRequest, b *biz.Menu) {

	b.Name = normalizeName(a.Name)
	b.Desc = a.Desc

	b.Component = a.Component
	requirement := lo.Map(a.Requirement, func(a *v1.PermissionRequirement, _ int) biz.MenuPermissionRequirement {
		ret := biz.MenuPermissionRequirement{
			Namespace: a.Namespace,
			Resource:  a.Resource,
			Action:    a.Action,
		}
		return ret
	})
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
	b.Title = a.Title
	b.Path = a.Path
	b.Redirect = a.Redirect
}

func normalizeName(name string) string {
	return strings.ToLower(name)
}
