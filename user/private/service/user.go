package service

import (
	"context"
	"fmt"
	klog "github.com/go-kratos/kratos/v2/log"
	api2 "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/query"
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

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pb "github.com/go-saas/kit/user/api/user/v1"
	"github.com/go-saas/kit/user/private/biz"
	"github.com/mennanov/fmutils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	um     *biz.UserManager
	rm     *biz.RoleManager
	auth   authz.Service
	blob   vfs.Blob
	trust  api2.TrustedContextValidator
	logger *klog.Helper
}

func NewUserService(
	um *biz.UserManager,
	rm *biz.RoleManager,
	auth authz.Service,
	blob vfs.Blob,
	trust api2.TrustedContextValidator,
	l klog.Logger,
) *UserService {
	return &UserService{
		um:     um,
		rm:     rm,
		auth:   auth,
		blob:   blob,
		trust:  trust,
		logger: klog.NewHelper(klog.With(l, "module", "user.UserService")),
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
		if req.Name != nil {
			u.Name = &req.Name.Value
		}
		if req.Username != nil {
			u.Username = &req.Username.Value
		}
		if req.Email != nil {
			u.SetEmail(req.Email.Value, false)
		}
		if req.Phone != nil {
			u.SetPhone(req.Phone.Value, false)
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

// SearchUser is for inviting user or creating user
func (s *UserService) SearchUser(ctx context.Context, req *pb.SearchUserRequest) (*pb.SearchUserResponse, error) {
	if _, err := authn.ErrIfUnauthenticated(ctx); err != nil {
		return nil, err
	}
	var user *biz.User
	var err error
	if len(req.Identity) > 0 {
		user, err = s.um.FindByIdentity(ctx, req.Identity)
	} else if req.Email != nil {
		user, err = s.um.FindByEmail(ctx, req.Email.Value)
	} else if len(req.Phone) > 0 {
		user, err = s.um.FindByPhone(ctx, req.Phone)
	} else if len(req.Username) > 0 {
		user, err = s.um.FindByName(ctx, req.Username)
	}
	if err != nil {
		return nil, err
	}
	ret := &pb.SearchUserResponse{}
	if user == nil {
		return ret, nil
	}

	ret.User = &pb.SearchUserResponse_SearchUser{
		Id: user.ID.String(),
	}
	if user.Username != nil {
		ret.User.Username = &wrapperspb.StringValue{Value: *user.Username}
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
	userId := req.FormValue("id")
	h := ctx.Middleware(func(ctx context.Context, _ interface{}) (interface{}, error) {
		if len(userId) > 0 {
			if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceUser, userId), authz.UpdateAction); err != nil {
				return nil, err
			}
		} else {
			_, err := authn.ErrIfUnauthenticated(ctx)
			if err != nil {
				return nil, err
			}
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
		if len(userId) > 0 {
			//update avatar field
			u, err := s.um.FindByID(ctx, userId)
			if err != nil {
				return nil, err
			}
			u.Avatar = &normalizedName
			err = s.um.Update(ctx, u, query.NewField(&fieldmaskpb.FieldMask{Paths: []string{"avatar"}}))
			if err != nil {
				return nil, err
			}
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
	if u.Username != nil {
		res.Username = &wrapperspb.StringValue{Value: *u.Username}
	}
	if u.Name != nil {
		res.Name = &wrapperspb.StringValue{Value: *u.Name}
	}
	if u.Phone != nil {
		res.Phone = &wrapperspb.StringValue{Value: *u.Phone}
	}
	if u.Email != nil {
		res.Email = &wrapperspb.StringValue{Value: *u.Email}
	}
	if u.Birthday != nil {
		res.Birthday = timestamppb.New(*u.Birthday)
	}
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
