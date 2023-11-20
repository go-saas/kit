package service

import (
	"context"
	"fmt"
	klog "github.com/go-kratos/kratos/v2/log"
	api2 "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/utils"
	v12 "github.com/go-saas/kit/user/api/permission/v1"
	"github.com/go-saas/kit/user/util"
	"github.com/go-saas/saas"
	"github.com/goxiaoy/vfs"
	"io"
	"os"
	"path/filepath"

	"github.com/go-saas/kit/user/api"

	errors2 "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-saas/kit/pkg/authn"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/blob"
	v1 "github.com/go-saas/kit/user/api/role/v1"
	"github.com/google/uuid"

	pb "github.com/go-saas/kit/user/api/user/v1"
	"github.com/go-saas/kit/user/private/biz"
	"github.com/mennanov/fmutils"
	"github.com/samber/lo"
)

type UserService struct {
	um            *biz.UserManager
	rm            *biz.RoleManager
	auth          authz.Service
	blob          vfs.Blob
	trust         api2.TrustedContextValidator
	permissionMgr authz.PermissionManagementService
	logger        *klog.Helper
}

func NewUserService(
	um *biz.UserManager,
	rm *biz.RoleManager,
	auth authz.Service,
	blob vfs.Blob,
	trust api2.TrustedContextValidator,
	permissionMgr authz.PermissionManagementService,
	l klog.Logger,
) *UserService {
	return &UserService{
		um:            um,
		rm:            rm,
		auth:          auth,
		blob:          blob,
		trust:         trust,
		permissionMgr: permissionMgr,
		logger:        klog.NewHelper(klog.With(l, "module", "user.UserService")),
	}
}

var _ pb.UserServiceServer = (*UserService)(nil)
var _ pb.UserAdminServiceServer = (*UserService)(nil)

func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceUser, "*"), authz.ReadAction); err != nil {
		return nil, err
	}
	ret := &pb.ListUsersResponse{}
	totalCount, filterCount, err := s.um.Count(ctx, req)
	if err != nil {
		return nil, err
	}
	ret.TotalSize = int32(totalCount)
	ret.FilterSize = int32(filterCount)

	items, err := s.um.List(ctx, req)
	if err != nil {
		return ret, err
	}
	rItems := make([]*pb.User, len(items))
	for index, u := range items {
		res := MapBizUserToApi(ctx, u, s.blob)
		if req.Fields != nil {
			fmutils.Filter(res, req.Fields.Paths)
		}
		rItems[index] = res
	}
	ret.Items = rItems
	return ret, nil

}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceUser, req.Id), authz.ReadAction); err != nil {
		return nil, err
	}
	u, err := s.um.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, pb.ErrorUserNotFoundLocalized(ctx, nil, nil)
	}
	if err := u.CheckInCurrentTenant(ctx); err != nil {
		return nil, err
	}
	res := MapBizUserToApi(ctx, u, s.blob)
	return res, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceUser, "*"), authz.CreateAction); err != nil {
		return nil, err
	}
	// check confirm password
	if req.Password != "" {
		if req.ConfirmPassword != req.Password {
			return nil, pb.ErrorConfirmPasswordMismatchLocalized(ctx, nil, nil)
		}
	}
	ct, _ := saas.FromCurrentTenant(ctx)
	u := biz.User{}
	if len(req.Id) > 0 {
		dbUser, err := s.um.FindByID(ctx, req.Id)
		if err != nil {
			return nil, err
		}
		if dbUser == nil {
			return nil, errors2.NotFound("", "")
		}
		u = *dbUser
	} else {
		u.Name = req.Name
		u.Username = req.Username
		if req.Email != nil {
			u.SetEmail(*req.Email, false)
		}
		if req.Phone != nil {
			u.SetPhone(*req.Phone, false)
		}
		if req.Birthday != nil {
			b := req.Birthday.AsTime()
			u.Birthday = &b
		}
		gender := req.Gender.String()
		if len(req.Avatar) > 0 {
			u.Avatar = &req.Avatar
		}
		u.Gender = &gender
		var err error
		if req.Password != "" {
			err = s.um.CreateWithPassword(ctx, &u, req.Password, true)
		} else {
			err = s.um.Create(ctx, &u)
		}
		if err != nil {
			return nil, err
		}
	}
	if err := s.um.JoinTenant(ctx, u.ID.String(), ct.GetId()); err != nil {
		return nil, err
	}
	//set roles
	var roles []biz.Role
	for _, r := range req.RolesId {
		if rr, err := s.rm.Get(ctx, r); err != nil {
			return nil, err
		} else if rr == nil {
			return nil, errors2.NotFound("", "role not found")
		} else {
			roles = append(roles, *rr)
		}
	}

	if err := s.um.UpdateRoles(ctx, &u, roles); err != nil {
		return nil, err
	}
	res := MapBizUserToApi(ctx, &u, s.blob)
	return res, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceUser, req.User.Id), authz.UpdateAction); err != nil {
		return nil, err
	}
	u, err := s.um.FindByID(ctx, req.User.Id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, pb.ErrorUserNotFoundLocalized(ctx, nil, nil)
	}
	if err := u.CheckInCurrentTenant(ctx); err != nil {
		return nil, err
	}

	//set roles
	var roles []biz.Role
	for _, r := range req.User.RolesId {
		if rr, err := s.rm.Get(ctx, r); err != nil {
			return nil, err
		} else if rr == nil {
			return nil, errors2.NotFound("", "role not found")
		} else {
			roles = append(roles, *rr)
		}
	}

	if err := s.um.UpdateRoles(ctx, u, roles); err != nil {
		return nil, err
	}

	res := MapBizUserToApi(ctx, u, s.blob)
	return res, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {

	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceUser, req.Id), authz.DeleteAction); err != nil {
		return nil, err
	}

	u, err := s.um.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, pb.ErrorUserNotFoundLocalized(ctx, nil, nil)
	}
	if err := u.CheckInCurrentTenant(ctx); err != nil {
		return nil, err
	}
	ti, _ := saas.FromCurrentTenant(ctx)
	//just remove from tenant
	if err := s.um.RemoveFromTenant(ctx, u.ID.String(), ti.GetId()); err != nil {
		return nil, err
	}
	return &pb.DeleteUserResponse{}, nil
}

func (s *UserService) GetUserRoles(ctx context.Context, req *pb.GetUserRoleRequest) (*pb.GetUserRoleReply, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceUser, req.Id), authz.ReadAction); err != nil {
		return nil, err
	}

	u, err := s.um.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, pb.ErrorUserNotFoundLocalized(ctx, nil, nil)
	}
	if err := u.CheckInCurrentTenant(ctx); err != nil {
		return nil, err
	}
	roles, err := s.um.GetRoles(ctx, u.ID.String())
	if err != nil {
		return nil, err
	}
	resp := &pb.GetUserRoleReply{}
	resp.Roles = make([]*pb.UserRole, len(roles))
	for i := range roles {
		resp.Roles[i] = &pb.UserRole{
			Id:   roles[i].ID.String(),
			Name: roles[i].Name,
		}
	}
	return resp, nil
}

func (s *UserService) GetUserPermission(ctx context.Context, req *pb.GetUserPermissionRequest) (*pb.GetUserPermissionReply, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceUser, req.Id), authz.ReadAction); err != nil {
		return nil, err
	}
	u, err := s.um.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, pb.ErrorUserNotFoundLocalized(ctx, nil, nil)
	}
	if err := u.CheckInCurrentTenant(ctx); err != nil {
		return nil, err
	}
	userSubject := authz.NewUserSubject(u.ID.String())
	acl, err := s.permissionMgr.ListAcl(ctx, userSubject)
	if err != nil {
		return nil, err
	}
	resItems := make([]*v12.Permission, len(acl))
	for i, bean := range acl {
		r := &v12.Permission{}
		util.MapPermissionBeanToPb(bean, r)
		resItems[i] = r
	}
	res := &pb.GetUserPermissionReply{
		Acl: resItems,
	}
	ti, _ := saas.FromCurrentTenant(ctx)
	var groups []*v12.PermissionDefGroup
	authz.WalkGroups(len(ti.GetId()) == 0, true, func(group *authz.PermissionDefGroup) {
		g := &v12.PermissionDefGroup{}
		mapGroupDef2Pb(ctx, group, g)
		groups = append(groups, g)
		var defs []*v12.PermissionDef
		group.Walk(len(ti.GetId()) == 0, true, func(def *authz.PermissionDef) {
			d := &v12.PermissionDef{}
			mapDef2Pb(ctx, def, d)
			defs = append(defs, d)
		})
		g.Def = defs
	})
	requirements := lo.FlatMap(groups, func(group *v12.PermissionDefGroup, _ int) []*authz.Requirement {
		return lo.Map(group.Def, func(def *v12.PermissionDef, _ int) *authz.Requirement {
			return authz.NewRequirement(authz.NewEntityResource(def.Namespace, authz.AnyResource), authz.ActionStr(def.Action))
		})
	})
	checkResults, err := s.auth.BatchCheckForSubjects(ctx, requirements, userSubject)
	if err != nil {
		return nil, err
	}
	i := 0
	for _, group := range groups {
		for _, def := range group.Def {
			def.Granted = checkResults[i].Allowed
			i++
		}
	}
	res.DefGroups = groups
	return res, nil
}

func (s *UserService) UpdateUserPermission(ctx context.Context, req *pb.UpdateUserPermissionRequest) (*pb.UpdateUserPermissionReply, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceUser, req.Id), authz.UpdateAction); err != nil {
		return nil, err
	}
	u, err := s.um.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, pb.ErrorUserNotFoundLocalized(ctx, nil, nil)
	}
	if err := u.CheckInCurrentTenant(ctx); err != nil {
		return nil, err
	}
	ti, _ := saas.FromCurrentTenant(ctx)
	var acl = lo.Map(req.Acl, func(a *v12.UpdateSubjectPermissionAcl, _ int) authz.UpdateSubjectPermission {
		effect := util.MapPbEffect2AuthEffect(a.Effect)
		return authz.UpdateSubjectPermission{
			Resource: authz.NewEntityResource(a.Namespace, a.Resource),
			Action:   authz.ActionStr(a.Action),
			TenantID: ti.GetId(),
			Effect:   effect,
		}
	})
	if err := s.permissionMgr.UpdateGrant(ctx, authz.NewUserSubject(u.ID.String()), acl); err != nil {
		return nil, err
	}
	return &pb.UpdateUserPermissionReply{}, nil

}

func (s *UserService) InviteUser(ctx context.Context, req *pb.InviteUserRequest) (*pb.InviteUserReply, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceUser, "*"), authz.CreateAction); err != nil {
		return nil, err
	}
	//find user
	u, err := s.um.FindByIdentity(ctx, req.Identify)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors2.NotFound("", "")
	}
	ti, _ := saas.FromCurrentTenant(ctx)
	//TODO confirm??
	err = s.um.JoinTenant(ctx, u.ID.String(), ti.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.InviteUserReply{RequiredConfirm: false}, nil
}

// PublicSearchUser is for inviting user or creating user
func (s *UserService) PublicSearchUser(ctx context.Context, req *pb.SearchUserRequest) (*pb.SearchUserResponse, error) {
	if _, err := authn.ErrIfUnauthenticated(ctx); err != nil {
		return nil, err
	}
	var user *biz.User
	var err error
	if req.Identity != nil {
		user, err = s.um.FindByIdentity(ctx, *req.Identity)
	} else if req.Email != nil {
		user, err = s.um.FindByEmail(ctx, *req.Email)
	} else if req.Phone != nil {
		user, err = s.um.FindByPhone(ctx, *req.Phone)
	} else if req.Username != nil {
		user, err = s.um.FindByName(ctx, *req.Username)
	}
	if err != nil {
		return nil, err
	}
	ret := &pb.SearchUserResponse{}
	if user == nil {
		return ret, nil
	}

	ret.User = &pb.SearchUserResponse_SearchUser{
		Id:       user.ID.String(),
		Username: user.Username,
	}

	ret.User.Avatar = mapAvatar(ctx, s.blob, user)
	return ret, err
}

func (s *UserService) UpdateAvatar(ctx http.Context) error {
	req := ctx.Request()
	//TODO do not know why should read form file first ...
	if _, _, err := req.FormFile("file"); err != nil {
		return err
	}
	h := ctx.Middleware(func(ctx context.Context, _ interface{}) (interface{}, error) {
		_, err := authn.ErrIfUnauthenticated(ctx)
		if err != nil {
			return nil, err
		}
		file, handle, err := req.FormFile("file")
		if err != nil {
			return nil, err
		}
		defer file.Close()
		fileName := handle.Filename
		ext := filepath.Ext(fileName)
		normalizedName := fmt.Sprintf("%s/%s%s", biz.UserAvatarPath, uuid.New().String(), ext)

		err = s.blob.MkdirAll(biz.UserAvatarPath, 0755)
		if err != nil {
			return nil, err
		}
		f, err := s.blob.OpenFile(normalizedName, os.O_WRONLY|os.O_CREATE, 0o666)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		_, err = io.Copy(f, file)
		if err != nil {
			return nil, err
		}
		url, _ := s.blob.PublicUrl(ctx, normalizedName)
		return &blob.BlobFile{
			Id:   normalizedName,
			Name: "",
			Url:  url.URL,
		}, nil
	})
	out, err := h(ctx, nil)
	if err != nil {
		return err
	}
	return ctx.Result(201, out)
}

func MapBizUserToApi(ctx context.Context, u *biz.User, b vfs.Blob) *pb.User {
	res := &pb.User{
		Id:    u.ID.String(),
		Roles: nil,
	}
	res.Username = u.Username
	res.Name = u.Name
	res.Phone = u.Phone
	res.Email = u.Email
	res.Birthday = utils.Time2Timepb(u.Birthday)
	if u.Gender != nil {
		if v, ok := pb.Gender_value[*u.Gender]; ok {
			res.Gender = pb.Gender(v)
		}
	}
	if u.Roles != nil {
		var returnRoles = lo.Map(u.Roles, func(i biz.Role, _ int) *v1.Role {
			r := &v1.Role{}
			MapBizRoleToApi(&i, r)
			return r
		})
		res.Roles = returnRoles
	}
	res.Avatar = mapAvatar(ctx, b, u)
	return res
}
