package service

import (
	"context"
	"fmt"
	"github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/kit/pkg/data"
	"github.com/go-saas/kit/pkg/query"
	"github.com/go-saas/kit/pkg/utils"
	"github.com/go-saas/lbs"
	"github.com/go-saas/saas"
	"github.com/goxiaoy/vfs"
	"google.golang.org/protobuf/types/known/structpb"
	"io"
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-saas/kit/pkg/authn"
	"github.com/go-saas/kit/pkg/blob"
	v13 "github.com/go-saas/kit/saas/api/tenant/v1"
	v12 "github.com/go-saas/kit/user/api/role/v1"
	v1 "github.com/go-saas/kit/user/api/user/v1"
	"github.com/go-saas/kit/user/private/biz"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pb "github.com/go-saas/kit/user/api/account/v1"
)

type AccountService struct {
	pb.UnimplementedAccountServer
	um            *biz.UserManager
	tenantService v13.TenantInternalServiceServer
	blob          vfs.Blob
	userSetting   biz.UserSettingRepo
	userAddr      biz.UserAddressRepo
	normalizer    biz.LookupNormalizer
	appCfg        *conf.AppConfig
}

func NewAccountService(
	um *biz.UserManager,
	blob vfs.Blob,
	tenantService v13.TenantInternalServiceServer,
	userSetting biz.UserSettingRepo,
	userAddr biz.UserAddressRepo,
	normalizer biz.LookupNormalizer,
	appCfg *conf.AppConfig,
) *AccountService {
	return &AccountService{
		um:            um,
		blob:          blob,
		tenantService: tenantService,
		userSetting:   userSetting,
		userAddr:      userAddr,
		normalizer:    normalizer,
		appCfg:        appCfg,
	}
}

func (s *AccountService) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	userInfo, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	u, err := s.um.FindByID(ctx, userInfo.GetId())
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.Unauthorized("", "")
	}

	res := &pb.GetProfileResponse{
		Id:       u.ID.String(),
		Username: u.Username,
		Name:     u.Name,
		Phone:    u.Phone,
		Email:    u.Email,
		Birthday: utils.Time2Timepb(u.Birthday),
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
	var tenantIds []string
	tenantIds = lo.Map(u.Tenants, func(t biz.UserTenant, _ int) string {
		return t.GetTenantId()
	})
	currentTenant, _ := saas.FromCurrentTenant(ctx)
	tenantIds = append(tenantIds, currentTenant.GetId())

	hostTenantInfo := &pb.UserTenant{UserId: userInfo.GetId(), TenantId: "", IsHost: true, Tenant: &v13.TenantInfo{
		DisplayName: s.appCfg.HostDisplayName,
		Logo:        &blob.BlobFile{Url: s.appCfg.HostLogo},
	}}
	if len(tenantIds) > 0 {
		tenants, err := s.tenantService.ListTenant(ctx,
			&v13.ListTenantRequest{Filter: &v13.TenantFilter{
				Id: &query.StringFilterOperation{In: lo.Map(tenantIds, func(t string, _ int) *wrapperspb.StringValue {
					return &wrapperspb.StringValue{Value: t}
				})}}})
		if err != nil {
			return nil, err
		}

		reTenants := lo.Map(u.Tenants, func(ut biz.UserTenant, _ int) *pb.UserTenant {
			//get tenant info
			if len(ut.TenantId) == 0 {
				//host
				return hostTenantInfo
			}
			t, ok := lo.Find(tenants.Items, func(t *v13.Tenant) bool { return t.Id == ut.GetTenantId() })
			if !ok {
				return nil
			}
			return &pb.UserTenant{UserId: ut.UserId, TenantId: ut.GetTenantId(), Tenant: mapTenant2TenantInfo(t)}
		})
		for _, tt := range tenants.Items {
			if currentTenant.GetId() == tt.GetId() {
				res.CurrentTenant = &pb.UserTenant{UserId: userInfo.GetId(), TenantId: currentTenant.GetId(), Tenant: mapTenant2TenantInfo(tt)}
				break
			}
		}
		res.Tenants = reTenants
	}
	if len(currentTenant.GetId()) == 0 {
		//host
		res.CurrentTenant = hostTenantInfo
	}
	//avatar
	res.Avatar = mapAvatar(ctx, s.blob, u)

	return res, nil
}

func mapTenant2TenantInfo(tt *v13.Tenant) *v13.TenantInfo {
	return &v13.TenantInfo{
		Id:          tt.Id,
		Name:        tt.Name,
		DisplayName: tt.DisplayName,
		Region:      tt.Region,
		Logo:        tt.Logo,
		PlanKey:     tt.PlanKey,
		Plan:        tt.Plan,
	}
}

func (s *AccountService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	userInfo, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	u, err := s.um.FindByID(ctx, userInfo.GetId())
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.Unauthorized("", "")
	}
	u.Username = req.Username
	u.Name = req.Name
	if req.Gender != v1.Gender_UNKNOWN {
		g := req.Gender.String()
		u.Gender = &g
	}
	if err := s.um.Update(ctx, u, nil); err != nil {
		return nil, err
	}
	return &pb.UpdateProfileResponse{}, nil
}
func (s *AccountService) GetSettings(ctx context.Context, req *pb.GetSettingsRequest) (*pb.GetSettingsResponse, error) {

	u, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	entities, err := s.userSetting.FindByUser(ctx, u.GetId(), req)
	if err != nil {
		return nil, err
	}
	set := lo.Map(entities, func(t *biz.UserSetting, _ int) *pb.Settings {
		return &pb.Settings{
			Key:   t.Key,
			Value: t.Value.ToDynamicValue(),
		}
	})
	return &pb.GetSettingsResponse{Settings: set}, nil
}
func (s *AccountService) UpdateSettings(ctx context.Context, req *pb.UpdateSettingsRequest) (*pb.UpdateSettingsResponse, error) {

	u, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.userSetting.UpdateByUser(ctx, u.GetId(), lo.Map(req.Settings, func(t *pb.UpdateSettings, _ int) biz.UpdateUserSetting {
		return biz.UpdateUserSetting{
			Key:    t.Key,
			Value:  data.NewFromDynamicValue(t.Value),
			Delete: t.Reset_,
		}
	})); err != nil {
		return nil, err
	}

	entities, err := s.userSetting.FindByUser(ctx, u.GetId(), new(pb.GetSettingsRequest))
	if err != nil {
		return nil, err
	}
	set := lo.Map(entities, func(t *biz.UserSetting, _ int) *pb.Settings {
		return &pb.Settings{
			Key:   t.Key,
			Value: t.Value.ToDynamicValue(),
		}
	})
	return &pb.UpdateSettingsResponse{Settings: set}, nil
}

func (s *AccountService) GetAddresses(ctx context.Context, req *pb.GetAddressesRequest) (*pb.GetAddressesReply, error) {
	u, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	addres, err := s.userAddr.FindByUser(ctx, u.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.GetAddressesReply{
		Addresses: lo.Map(addres, func(t *biz.UserAddress, _ int) *pb.UserAddress {
			res := &pb.UserAddress{}
			mapBizUserAddr2Pb(t, res)
			return res
		}),
	}, nil
}
func (s *AccountService) CreateAddresses(ctx context.Context, req *pb.CreateAddressesRequest) (*pb.CreateAddressReply, error) {
	u, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	addr := &biz.UserAddress{}
	mapCreateAddr2Biz(req, addr)
	addr.UserId = uuid.MustParse(u.GetId())

	if len(addr.Phone) > 0 {
		p, err := s.normalizer.Phone(ctx, addr.Phone)
		if err != nil {
			return nil, err
		}
		addr.Phone = p
	}
	if err := s.userAddr.Create(ctx, addr); err != nil {
		return nil, err
	}
	if addr.Prefer {
		if err := s.userAddr.SetPrefer(ctx, addr); err != nil {
			return nil, err
		}
	}
	return &pb.CreateAddressReply{}, nil
}

func (s *AccountService) UpdateAddresses(ctx context.Context, req *pb.UpdateAddressesRequest) (*pb.UpdateAddressesReply, error) {

	ui, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	addr, err := s.userAddr.Get(ctx, req.Address.Id)
	if err != nil {
		return nil, err
	}
	if addr == nil || addr.UserId.String() != ui.GetId() {
		return nil, errors.NotFound("", "")
	}
	mapUpdateAddr2Biz(req.Address, addr)

	if err := s.userAddr.Update(ctx, addr.ID.String(), addr, nil); err != nil {
		return nil, err
	}
	if addr.Prefer {
		if err := s.userAddr.SetPrefer(ctx, addr); err != nil {
			return nil, err
		}
	}
	return &pb.UpdateAddressesReply{}, nil
}

func (s *AccountService) DeleteAddresses(ctx context.Context, req *pb.DeleteAddressRequest) (*pb.DeleteAddressesReply, error) {

	ui, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	addr, err := s.userAddr.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if addr == nil || addr.UserId.String() != ui.GetId() {
		return nil, errors.NotFound("", "")
	}
	if err := s.userAddr.Delete(ctx, addr.ID.String()); err != nil {
		return nil, err
	}
	return &pb.DeleteAddressesReply{}, nil
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
		//update avatar field
		u, err := s.um.FindByID(ctx, user.GetId())
		if err != nil {
			return nil, err
		}
		u.Avatar = &normalizedName
		err = s.um.Update(ctx, u, query.NewField(&fieldmaskpb.FieldMask{Paths: []string{"avatar"}}))
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

func mapAvatar(ctx context.Context, b vfs.Blob, user *biz.User) *blob.BlobFile {
	if user.Avatar == nil {
		return nil
	}

	url, _ := b.PublicUrl(ctx, *user.Avatar)
	return &blob.BlobFile{
		Id:   *user.Avatar,
		Name: "",
		Url:  url.URL,
	}
}

func mapBizUserAddr2Pb(a *biz.UserAddress, b *pb.UserAddress) {
	b.Id = a.ID.String()
	b.Phone = a.Phone
	b.Usage = a.Usage
	b.Prefer = a.Prefer

	b.Address, _ = a.Address.ToPb()
	m, _ := structpb.NewStruct(a.Metadata)
	b.Metadata = m
}

func mapCreateAddr2Biz(a *pb.CreateAddressesRequest, b *biz.UserAddress) error {
	b.Phone = a.Phone
	b.Usage = a.Usage
	b.Prefer = a.Prefer

	addr, err := lbs.NewAddressEntityFromPb(a.Address)
	if err != nil {
		return err
	}
	if addr != nil {
		b.Address = *addr
	}
	if a.Metadata != nil {
		b.Metadata = a.Metadata.AsMap()
	}
	return nil
}
func mapUpdateAddr2Biz(a *pb.UpdateAddress, b *biz.UserAddress) error {
	b.Phone = a.Phone
	b.Usage = a.Usage
	b.Prefer = a.Prefer

	addr, err := lbs.NewAddressEntityFromPb(a.Address)
	if err != nil {
		return err
	}
	if addr != nil {
		b.Address = *addr
	}
	if a.Metadata != nil {
		b.Metadata = a.Metadata.AsMap()
	}
	return nil
}
