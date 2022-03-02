package service

import (
	"context"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/uuid"
	"github.com/goxiaoy/go-saas-kit/pkg/authn"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	v13 "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	v12 "github.com/goxiaoy/go-saas-kit/user/api/role/v1"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/user/v1"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"os"
	"path/filepath"

	pb "github.com/goxiaoy/go-saas-kit/user/api/account/v1"
)

type AccountService struct {
	pb.UnimplementedAccountServer
	um            *biz.UserManager
	tenantService v13.TenantServiceClient
	blob          blob.Factory
}

func NewAccountService(um *biz.UserManager, blob blob.Factory, tenantService v13.TenantServiceClient) *AccountService {
	return &AccountService{
		um:            um,
		blob:          blob,
		tenantService: tenantService,
	}
}

func (s *AccountService) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	userInfo, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	u, err := s.um.FindByID(ctx, userInfo.GetId())
	if err != nil {
		return nil, errors.Forbidden("", "")
	}
	res := &pb.GetProfileResponse{
		Id:       u.ID.String(),
		Username: &wrapperspb.StringValue{Value: *u.Username},
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
		v, ok := v1.Gender_value[*u.Gender]
		if ok {
			res.Gender = *v1.Gender(v).Enum()
		}
	}
	for _, role := range u.Roles {
		res.Roles = append(res.Roles, &v12.Role{
			Id:   role.ID.String(),
			Name: role.Name,
		})
	}

	tenantIds := make([]string, len(u.Tenants))
	linq.From(u.Tenants).SelectT(func(t biz.UserTenant) string { return t.TenantId.String }).ToSlice(&tenantIds)
	if len(tenantIds) > 0 {
		tenants, err := s.tenantService.ListTenant(ctx, &v13.ListTenantRequest{Filter: &v13.TenantFilter{IdIn: tenantIds}})
		if err != nil {
			return nil, err
		}

		reTenants := make([]*pb.UserTenant, len(u.Tenants))
		linq.From(u.Tenants).SelectT(func(ut biz.UserTenant) *pb.UserTenant {
			//get tenant info
			if ut.TenantId.Valid {
				//host
				return &pb.UserTenant{UserId: ut.UserId, TenantId: ut.TenantId.String, IsHost: true}
			}

			t := linq.From(tenants.Items).FirstWithT(func(t *v13.Tenant) bool { return t.Id == ut.TenantId.String })
			if t == nil {
				return nil
			}
			tt := t.(*v13.Tenant)
			return &pb.UserTenant{UserId: ut.UserId, TenantId: ut.TenantId.String, Tenant: &pb.UserTenant_Tenant{
				Id:          tt.Id,
				Name:        tt.Name,
				DisplayName: tt.DisplayName,
				Region:      tt.Region,
			}}
		}).ToSlice(&reTenants)
		res.Tenants = reTenants
	}
	//avatar
	res.Avatar = mapAvatar(ctx, s.blob, u)

	return res, nil
}
func (s *AccountService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	_, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateProfileResponse{}, nil
}
func (s *AccountService) GetSettings(ctx context.Context, req *pb.GetSettingsRequest) (*pb.GetSettingsResponse, error) {
	_, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetSettingsResponse{}, nil
}
func (s *AccountService) UpdateSettings(ctx context.Context, req *pb.UpdateSettingsRequest) (*pb.UpdateSettingsResponse, error) {
	_, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateSettingsResponse{}, nil
}
func (s *AccountService) GetAddresses(ctx context.Context, req *pb.GetAddressesRequest) (*pb.GetAddressesReply, error) {
	_, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.GetAddressesReply{}, nil
}

func (s *AccountService) UpdateAddresses(ctx context.Context, req *pb.UpdateAddressesRequest) (*pb.UpdateAddressesReply, error) {
	_, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateAddressesReply{}, nil
}

func (s *AccountService) UpdateAvatar(ctx http.Context) error {
	req := ctx.Request()
	//TODO do not know why should read form file first ...
	if _, _, err := req.FormFile("file"); err != nil {
		return err
	}
	h := ctx.Middleware(func(ctx context.Context, _ interface{}) (interface{}, error) {
		user, err := authn.ErrIfUnauthenticated(ctx)
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
		//update avatar field
		u, err := s.um.FindByID(ctx, user.GetId())
		if err != nil {
			return nil, err
		}
		u.Avatar = &normalizedName
		err = s.um.Update(ctx, u, &fieldmaskpb.FieldMask{Paths: []string{"avatar"}})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	_, err := h(ctx, nil)
	if err != nil {
		return err
	}
	return ctx.Returns(201, nil)
}

func mapAvatar(ctx context.Context, factory blob.Factory, user *biz.User) *blob.BlobFile {
	if user.Avatar == nil {
		return nil
	}
	profile := biz.ProfileBlob(ctx, factory)

	url, _ := profile.GeneratePublicUrl(*user.Avatar)
	return &blob.BlobFile{
		Id:   *user.Avatar,
		Name: "",
		Url:  url,
	}
}
