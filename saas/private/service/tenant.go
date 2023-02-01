package service

import (
	"context"
	"fmt"
	"github.com/dtm-labs/client/dtmgrpc"
	dtmapi "github.com/go-saas/kit/dtm/api"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/kit/pkg/localize"
	"github.com/go-saas/kit/pkg/query"
	kithttp "github.com/go-saas/kit/pkg/server/http"
	conf2 "github.com/go-saas/kit/saas/private/conf"
	uapi "github.com/go-saas/kit/user/api"
	v1 "github.com/go-saas/kit/user/api/user/v1"
	"github.com/go-saas/saas"
	shttp "github.com/go-saas/saas/http"
	"github.com/go-saas/sessions"
	"github.com/goxiaoy/vfs"
	"github.com/lithammer/shortuuid/v3"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-saas/kit/pkg/authn"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/blob"
	"github.com/go-saas/kit/saas/api"
	pb "github.com/go-saas/kit/saas/api/tenant/v1"
	"github.com/go-saas/kit/saas/private/biz"
	"github.com/google/uuid"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TenantService struct {
	pb.UnimplementedTenantServiceServer
	useCase  *biz.TenantUseCase
	auth     authz.Service
	blob     vfs.Blob
	trusted  sapi.TrustedContextValidator
	app      *conf.AppConfig
	saasConf *conf2.SaasConf
	webConf  *shttp.WebMultiTenancyOption
	tokenMgr sapi.TokenManager
}

func NewTenantService(
	useCase *biz.TenantUseCase,
	auth authz.Service,
	trusted sapi.TrustedContextValidator,
	blob vfs.Blob,
	app *conf.AppConfig,
	saasConf *conf2.SaasConf,
	wenConf *shttp.WebMultiTenancyOption,
	tokenMgr sapi.TokenManager,
) *TenantService {
	return &TenantService{
		useCase:  useCase,
		auth:     auth,
		trusted:  trusted,
		blob:     blob,
		app:      app,
		saasConf: saasConf,
		webConf:  wenConf,
		tokenMgr: tokenMgr,
	}
}

func (s *TenantService) CreateTenant(ctx context.Context, req *pb.CreateTenantRequest) (*pb.Tenant, error) {

	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceTenant, "*"), authz.CreateAction); err != nil {
		return nil, err
	}

	if len(req.DisplayName) == 0 {
		req.DisplayName = req.Name
	}
	uid := uuid.New()
	req.Id = uid.String()

	//XA Transaction
	gid := shortuuid.New()
	var err error
	var createTenantResp *pb.Tenant
	err = dtmgrpc.XaGlobalTransaction(sapi.WithDiscovery(dtmapi.ServiceName), gid, func(xa *dtmgrpc.XaGrpc) error {
		t, err := s.tokenMgr.GetOrGenerateToken(ctx, dtmapi.ClientConf)
		if err != nil {
			return err
		}
		xa.BranchHeaders = map[string]string{
			"Authorization": t,
		}
		//create tenant
		err = xa.CallBranch(req, sapi.WithDiscovery(api.ServiceName)+pb.GrpcOperationTenantInternalServiceCreateTenant, createTenantResp)
		if err != nil {
			return err
		}
		//server id
		req.Id = createTenantResp.Id
		//create user ,seed database
		err = xa.CallBranch(req, sapi.WithDiscovery(uapi.ServiceName)+v1.GrpcOperationUserInternalServiceCreateTenant, &emptypb.Empty{})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return createTenantResp, nil
}

func (s *TenantService) UpdateTenant(ctx context.Context, req *pb.UpdateTenantRequest) (*pb.Tenant, error) {

	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceTenant, req.Tenant.Id), authz.UpdateAction); err != nil {
		return nil, err
	}

	t, err := s.useCase.Get(ctx, req.Tenant.Id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, errors.NotFound("", "")
	}
	t.Name = req.Tenant.Name
	t.DisplayName = req.Tenant.DisplayName
	t.Logo = req.Tenant.Logo

	tenantConn := lo.Map(req.Tenant.Conn, func(t *pb.TenantConnectionString, _ int) biz.TenantConn {
		return biz.TenantConn{
			Key:   t.Key,
			Value: t.Value,
		}
	})

	tenantFeature := lo.Map(req.Tenant.Features, func(t *pb.TenantFeature, _ int) biz.TenantFeature {
		return biz.TenantFeature{
			Key:   t.Key,
			Value: t.Value,
		}
	})
	t.Conn = tenantConn
	t.Features = tenantFeature

	if err := s.useCase.Update(ctx, t, query.NewField(req.UpdateMask)); err != nil {
		return nil, err
	}
	return mapBizTenantToApi(ctx, s.app, s.blob, t), nil
}
func (s *TenantService) DeleteTenant(ctx context.Context, req *pb.DeleteTenantRequest) (*pb.DeleteTenantReply, error) {

	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceTenant, req.Id), authz.DeleteAction); err != nil {
		return nil, err
	}

	if err := s.useCase.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &pb.DeleteTenantReply{}, nil
}

func (s *TenantService) GetTenant(ctx context.Context, req *pb.GetTenantRequest) (*pb.Tenant, error) {
	t, err := s.useCase.FindByIdOrName(ctx, req.IdOrName)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, errors.NotFound("", "")
	}

	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceTenant, t.ID.String()), authz.ReadAction); err != nil {
		return nil, err
	}

	return mapBizTenantToApi(ctx, s.app, s.blob, t), nil
}

// GetTenantPublic return public info of tenant
func (s *TenantService) GetTenantPublic(ctx context.Context, req *pb.GetTenantPublicRequest) (*pb.TenantInfo, error) {
	t, err := s.useCase.FindByIdOrName(ctx, req.IdOrName)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, errors.NotFound("", "")
	}
	return mapBizTenantToInfo(ctx, s.blob, t, s.app), nil
}

func (s *TenantService) GetCurrentTenant(ctx context.Context, req *pb.GetCurrentTenantRequest) (*pb.GetCurrentTenantReply, error) {
	ti, _ := saas.FromCurrentTenant(ctx)
	if len(ti.GetId()) == 0 {
		return &pb.GetCurrentTenantReply{IsHost: true, Tenant: mapBizTenantToInfo(ctx, s.blob, nil, s.app)}, nil
	} else {
		t, err := s.useCase.FindByIdOrName(ctx, ti.GetId())
		if err != nil {
			return nil, err
		}
		if t == nil {
			return nil, pb.ErrorTenantNotFoundLocalized(localize.FromContext(ctx), nil, nil)
		}
		info := mapBizTenantToInfo(ctx, s.blob, t, s.app)
		return &pb.GetCurrentTenantReply{
			IsHost: false,
			Tenant: info,
		}, nil
	}
}

func (s *TenantService) ListTenant(ctx context.Context, req *pb.ListTenantRequest) (*pb.ListTenantReply, error) {

	if authResult, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceTenant, "*"), authz.ReadAction); err != nil {
		return nil, err
	} else if !authResult.Allowed {
		return nil, errors.Forbidden("", "")
	}

	ret := &pb.ListTenantReply{}

	totalCount, filterCount, err := s.useCase.Count(ctx, req)
	ret.TotalSize = int32(totalCount)
	ret.FilterSize = int32(filterCount)
	if err != nil {
		return ret, err
	}
	items, err := s.useCase.List(ctx, req)
	if err != nil {
		return ret, err
	}
	rItems := lo.Map(items, func(g *biz.Tenant, _ int) *pb.Tenant { return mapBizTenantToApi(ctx, s.app, s.blob, g) })
	ret.Items = rItems
	return ret, nil
}

func (s *TenantService) ChangeTenant(ctx context.Context, req *pb.ChangeTenantRequest) (*pb.ChangeTenantReply, error) {
	ret := &pb.ChangeTenantReply{}
	domain := ""
	if s.saasConf != nil && s.saasConf.TenantCookie != nil {
		domain = s.saasConf.TenantCookie.Domain.Value
	}
	if len(req.IdOrName) == 0 || req.IdOrName == "-" {
		ret.IsHost = true
		//clear cookie
		kithttp.SetCookie(ctx, sessions.NewCookie(s.webConf.TenantKey, "", &sessions.Options{MaxAge: -1, Domain: domain, Path: "/"}))
		return ret, nil
	}
	t, err := s.useCase.FindByIdOrName(ctx, req.IdOrName)
	if err != nil {
		return nil, err
	}
	ret.Tenant = mapBizTenantToInfo(ctx, s.blob, t, s.app)
	kithttp.SetCookie(ctx, sessions.NewCookie(s.webConf.TenantKey, ret.Tenant.Name, &sessions.Options{MaxAge: 2147483647, Domain: domain, Path: "/"}))
	return ret, nil
}

func (s *TenantService) UpdateLogo(ctx http.Context) error {
	req := ctx.Request()
	//TODO do not know why should read form file first ...
	if _, _, err := req.FormFile("file"); err != nil {
		return err
	}
	tenantID := req.FormValue("id")
	h := ctx.Middleware(func(ctx context.Context, _ interface{}) (interface{}, error) {
		if len(tenantID) > 0 {
			if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceTenant, tenantID), authz.UpdateAction); err != nil {
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
		normalizedName := fmt.Sprintf("%s/%s%s", biz.TenantLogoPath, uuid.New().String(), ext)

		err = s.blob.MkdirAll(biz.TenantLogoPath, 0755)
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
		if len(tenantID) > 0 {
			//update field
			t, err := s.useCase.FindByIdOrName(ctx, tenantID)
			if err != nil {
				return nil, err
			}
			if t == nil {
				return nil, errors.NotFound("", "")
			}

			t.Logo = normalizedName
			if err := s.useCase.Update(ctx, t, query.NewField(&fieldmaskpb.FieldMask{Paths: []string{"logo"}})); err != nil {
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

func mapBizTenantToApi(ctx context.Context, app *conf.AppConfig, blob vfs.Blob, tenant *biz.Tenant) *pb.Tenant {
	conns := lo.Map(tenant.Conn, func(con biz.TenantConn, _ int) *pb.TenantConnectionString {
		return &pb.TenantConnectionString{
			Key:   con.Key,
			Value: con.Value,
		}
	})

	features := lo.Map(tenant.Features, func(con biz.TenantFeature, _ int) *pb.TenantFeature {
		return &pb.TenantFeature{
			Key:   con.Key,
			Value: con.Value,
		}
	})

	res := &pb.Tenant{
		Id:          tenant.ID.String(),
		Name:        tenant.Name,
		DisplayName: tenant.DisplayName,
		Region:      tenant.Region,
		CreatedAt:   timestamppb.New(tenant.CreatedAt),
		UpdatedAt:   timestamppb.New(tenant.UpdatedAt),
		Conn:        conns,
		Features:    features,
		Logo:        mapLogo(ctx, blob, tenant),
		SeparateDb:  tenant.SeparateDb,
	}
	res.NormalizeHost(ctx, app)
	return res
}

func mapBizTenantToInfo(ctx context.Context, b vfs.Blob, tenant *biz.Tenant, app *conf.AppConfig) *pb.TenantInfo {
	if tenant == nil {
		if app == nil {
			return nil
		}
		return &pb.TenantInfo{
			DisplayName: app.HostDisplayName,
			Logo:        &blob.BlobFile{Url: app.HostLogo},
		}
	}
	res := &pb.TenantInfo{
		Id:          tenant.ID.String(),
		Name:        tenant.Name,
		DisplayName: tenant.DisplayName,
		Region:      tenant.Region,
		Logo:        mapLogo(ctx, b, tenant),
	}
	res.NormalizeHost(ctx, app)
	return res
}

func mapLogo(ctx context.Context, b vfs.Blob, entity *biz.Tenant) *blob.BlobFile {
	if entity.Logo == "" {
		return nil
	}

	url, _ := b.PublicUrl(ctx, entity.Logo)
	return &blob.BlobFile{
		Id:   entity.Logo,
		Name: "",
		Url:  url.URL,
	}
}
