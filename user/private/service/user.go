package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	errors2 "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/uuid"
	"github.com/goxiaoy/go-saas-kit/pkg/authn"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	v12 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/role/v1"
	"github.com/goxiaoy/go-saas/common"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pb "github.com/goxiaoy/go-saas-kit/user/api/user/v1"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	"github.com/mennanov/fmutils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	um   *biz.UserManager
	rm   *biz.RoleManager
	auth authz.Service
	blob blob.Factory
}

func NewUserService(um *biz.UserManager, rm *biz.RoleManager, auth authz.Service, blob blob.Factory) *UserService {
	return &UserService{
		um:   um,
		rm:   rm,
		auth: auth,
		blob: blob,
	}
}

func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource("user.user", "*"), authz.ListAction); err != nil {
		return nil, err
	}
	ret := &pb.ListUsersResponse{}
	totalCount, filterCount, err := s.um.Count(ctx, req.Filter)
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
	if authResult, err := s.auth.Check(ctx, authz.NewEntityResource("user.user", req.Id), authz.GetAction); err != nil {
		return nil, err
	} else if !authResult.Allowed {
		return nil, errors2.Forbidden("", "")
	}
	u, err := s.um.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors2.NotFound("", "")
	}
	res := MapBizUserToApi(ctx, u, s.blob)
	return res, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource("user.user", "*"), authz.CreateAction); err != nil {
		return nil, err
	}
	// check confirm password
	if req.Password != "" {
		if req.ConfirmPassword != req.Password {
			return nil, pb.ErrorConfirmPasswordMismatch("")
		}
	}
	u := biz.User{}
	if req.Name != nil {
		u.Name = &req.Name.Value
	}
	if req.Username != nil {
		u.Username = &req.Username.Value
	}
	if req.Email != nil {
		u.Email = &req.Email.Value
	}
	if req.Phone != nil {
		u.Phone = &req.Phone.Value
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
		err = s.um.CreateWithPassword(ctx, &u, req.Password)
	} else {
		err = s.um.Create(ctx, &u)
	}
	if err != nil {
		return nil, ConvertError(err)
	}
	//set roles
	var roles []*biz.Role
	for _, r := range req.RolesId {
		if rr, err := s.rm.Get(ctx, r); err != nil {
			return nil, ConvertError(err)
		} else if rr == nil {
			return nil, errors2.NotFound("", "role not found")
		} else {
			roles = append(roles, rr)
		}
	}

	if err := s.um.UpdateRoles(ctx, &u, roles); err != nil {
		return nil, err
	}
	res := MapBizUserToApi(ctx, &u, s.blob)
	return res, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {

	if _, err := s.auth.Check(ctx, authz.NewEntityResource("user.user", req.User.Id), authz.UpdateAction); err != nil {
		return nil, err
	}

	// check confirm password
	if req.User.Password != "" {
		if req.User.ConfirmPassword != req.User.Password {
			return nil, pb.ErrorConfirmPasswordMismatch("")
		}
	}

	if req.UpdateMask != nil {
		fmutils.Filter(req, req.UpdateMask.Paths)
	}
	u, err := s.um.FindByID(ctx, req.User.Id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors2.NotFound("", "")
	}

	if req.User.Password != "" {
		//reset password
		if err := s.um.UpdatePassword(ctx, u, req.User.Password); err != nil {
			return nil, err
		}
	}
	if req.GetUser().GetUsername() != nil {
		v := req.GetUser().GetUsername().Value
		u.Username = &v
	}

	if req.GetUser().GetName() != nil {
		v := req.GetUser().GetName().Value
		u.Name = &v
	}
	if req.GetUser().GetPhone() != nil {
		v := req.GetUser().GetPhone().Value
		u.Phone = &v
	}
	if req.GetUser().GetEmail() != nil {
		v := req.GetUser().GetEmail().Value
		u.Email = &v
	}
	if req.GetUser().GetBirthday() != nil {
		v := req.GetUser().GetBirthday().AsTime()
		u.Birthday = &v
	}
	if len(req.User.Avatar) > 0 {
		u.Avatar = &req.User.Avatar
	}
	g := req.GetUser().Gender.Enum().String()
	u.Gender = &g
	if err := s.um.Update(ctx, u, nil); err != nil {
		return nil, ConvertError(err)
	}

	//set roles
	var roles []*biz.Role
	for _, r := range req.User.RolesId {
		if rr, err := s.rm.Get(ctx, r); err != nil {
			return nil, ConvertError(err)
		} else if rr == nil {
			return nil, errors2.NotFound("", "role not found")
		} else {
			roles = append(roles, rr)
		}
	}

	if err := s.um.UpdateRoles(ctx, u, roles); err != nil {
		return nil, err
	}

	res := MapBizUserToApi(ctx, u, s.blob)
	return res, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource("user.user", req.Id), authz.DeleteAction); err != nil {
		return nil, err
	}
	u, err := s.um.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors2.NotFound("", "")
	}
	ti, _ := common.FromCurrentTenant(ctx)
	if req.Force {

		creatorId := ""
		if u.CreatedTenant != nil {
			creatorId = *u.CreatedTenant
		}
		if ti.GetId() != creatorId {
			//TODO ??other matching?
			return nil, errors2.Forbidden("", "")
		}
		if err := s.um.Delete(ctx, u); err != nil {
			return nil, err
		}
		return &pb.DeleteUserResponse{}, nil
	} else {
		//just remove from tenant
		if err := s.um.RemoveFromTenant(ctx, u.ID.String(), ti.GetId()); err != nil {
			return nil, err
		}
	}
	return &pb.DeleteUserResponse{}, nil
}

func (s *UserService) GetUserRoles(ctx context.Context, req *pb.GetUserRoleRequest) (*pb.GetUserRoleReply, error) {
	//TODO frequency call. use cache
	if _, err := s.auth.Check(ctx, authz.NewEntityResource("user.user", req.Id), authz.GetAction); err != nil {
		return nil, err
	}
	u, err := s.um.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors2.NotFound("", "")
	}
	roles, err := s.um.GetRoles(ctx, u)
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
	if _, err := s.auth.Check(ctx, authz.NewEntityResource("user.user", "*"), authz.CreateAction); err != nil {
		return nil, err
	}

	ctx = biz.NewIgnoreUserTenantsContext(ctx, true)
	//find user
	u, err := s.um.FindByIdentity(ctx, req.Identify)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors2.NotFound("", "")
	}
	ti, _ := common.FromCurrentTenant(ctx)
	//TODO confirm??
	err = s.um.JoinTenant(ctx, u.ID.String(), ti.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.InviteUserReply{RequiredConfirm: false}, nil
}

//CheckUserTenant internal api for check user tenant
func (s *UserService) CheckUserTenant(ctx context.Context, req *pb.CheckUserTenantRequest) (*pb.CheckUserTenantReply, error) {
	//check permission
	if _, err := s.auth.CheckInTenant(ctx, authz.NewEntityResource("user.user", req.UserId), authz.ActionStr("check_user_tenant"), "*"); err != nil {
		return nil, err
	}
	ok, err := s.CheckUserTenantInternal(ctx, req.UserId, req.TenantId)
	if err != nil {
		return nil, err
	}
	return &pb.CheckUserTenantReply{Ok: ok}, nil
}

func (s *UserService) CheckUserTenantInternal(ctx context.Context, userId, tenantId string) (bool, error) {
	ok, err := s.um.IsInTenant(ctx, userId, tenantId)
	if err != nil {
		return false, err
	}
	if ok {
		//user in this tenant
		return true, nil
	}
	tenantCtx := common.NewCurrentTenant(ctx, tenantId, "")
	if _, err := s.auth.Check(tenantCtx, authz.NewEntityResource("*", "*"), authz.ActionStr("*")); err != nil {
		//no permission
		return false, v12.ErrorTenantForbidden("")
	}
	//super permission
	return true, nil
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
			if _, err := s.auth.Check(ctx, authz.NewEntityResource("user.user", userId), authz.UpdateAction); err != nil {
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
		normalizedName := fmt.Sprintf("avatar/%s%s", uuid.New().String(), ext)
		profileBlob := biz.ProfileBlob(ctx, s.blob)
		a := profileBlob.GetAfero()
		err = a.MkdirAll("avatar", 0755)
		if err != nil {
			return nil, err
		}
		f, err := a.OpenFile(normalizedName, os.O_WRONLY|os.O_CREATE, 0o666)
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
			err = s.um.Update(ctx, u, &fieldmaskpb.FieldMask{Paths: []string{"avatar"}})
			if err != nil {
				return nil, err
			}
		}

		profile := biz.ProfileBlob(ctx, s.blob)
		url, _ := profile.GeneratePublicUrl(normalizedName)
		return &blob.BlobFile{
			Id:   normalizedName,
			Name: "",
			Url:  url,
		}, nil
	})
	out, err := h(ctx, nil)
	if err != nil {
		return err
	}
	return ctx.Result(201, out)
}

func MapBizUserToApi(ctx context.Context, u *biz.User, b blob.Factory) *pb.User {
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
			return &v1.Role{
				Id:   i.ID.String(),
				Name: i.Name,
			}
		})
		res.Roles = returnRoles
	}
	res.Avatar = mapAvatar(ctx, b, u)
	return res
}
