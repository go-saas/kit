package service

import (
	"context"
	"fmt"
	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/dtm-labs/dtm/client/workflow"
	dtmsrv "github.com/go-saas/kit/dtm/service"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/kit/pkg/data"
	"github.com/go-saas/kit/pkg/query"
	kithttp "github.com/go-saas/kit/pkg/server/http"
	"github.com/go-saas/kit/pkg/utils"
	v12 "github.com/go-saas/kit/saas/api/plan/v1"
	conf2 "github.com/go-saas/kit/saas/private/conf"
	v1 "github.com/go-saas/kit/user/api/user/v1"
	"github.com/go-saas/saas"
	shttp "github.com/go-saas/saas/http"
	"github.com/go-saas/sessions"
	"github.com/goxiaoy/vfs"
	"github.com/segmentio/ksuid"
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
	"google.golang.org/protobuf/types/known/timestamppb"
)

var wfCreateTenantName = "saas_create_tenant"

type TenantService struct {
	useCase         *biz.TenantUseCase
	auth            authz.Service
	blob            vfs.Blob
	trusted         sapi.TrustedContextValidator
	app             *conf.AppConfig
	saasConf        *conf2.SaasConf
	webConf         *shttp.WebMultiTenancyOption
	txhelper        *dtmsrv.Helper
	userInternalSrv v1.UserInternalServiceServer
}

var _ pb.TenantServiceServer = (*TenantService)(nil)

func NewTenantService(
	useCase *biz.TenantUseCase,
	auth authz.Service,
	trusted sapi.TrustedContextValidator,
	blob vfs.Blob,
	app *conf.AppConfig,
	saasConf *conf2.SaasConf,
	wenConf *shttp.WebMultiTenancyOption,
	txhelper *dtmsrv.Helper,
	userInternalSrv v1.UserInternalServiceServer,
) *TenantService {
	s := &TenantService{
		useCase:         useCase,
		auth:            auth,
		trusted:         trusted,
		blob:            blob,
		app:             app,
		saasConf:        saasConf,
		webConf:         wenConf,
		txhelper:        txhelper,
		userInternalSrv: userInternalSrv,
	}

	err := s.txhelper.WorkflowRegister2(wfCreateTenantName, func(wf *workflow.Workflow, data []byte) ([]byte, error) {
		var req = &pb.CreateTenantRequest{}
		utils.PbMustUnMarshalJson(data, req)

		var tenantId string
		resp, err := wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			//delete tenant
			return s.txhelper.BarrierUow(wf.Context, bb, biz.ConnName, func(ctx context.Context) error {
				if err := s.useCase.Delete(ctx, tenantId); err != nil {
					return err
				}
				return nil
			})
		}).Do(func(bb *dtmcli.BranchBarrier) ([]byte, error) {
			//create tenant
			resp := &pb.Tenant{}
			err := s.txhelper.BarrierUow(wf.Context, bb, biz.ConnName, func(ctx context.Context) error {
				if len(req.DisplayName) == 0 {
					req.DisplayName = req.Name
				}
				t := &biz.Tenant{
					Name:        req.Name,
					DisplayName: req.DisplayName,
					Region:      req.Region,
					Logo:        req.Logo,
					SeparateDb:  req.SeparateDb,
				}

				if err := s.useCase.Create(ctx, t); err != nil {
					return err
				}
				resp = mapBizTenantToApi(ctx, s.app, s.blob, t)
				return nil
			})
			if err != nil {
				return nil, err
			}
			tenantId = resp.Id
			return utils.PbMustMarshalJson(resp), err
		})
		if err != nil {
			return nil, err
		}

		//call new branch for later gRPC call
		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			//do nothing
			return nil
		})
		_, err = s.userInternalSrv.CreateTenant(wf.Context, &v1.UserInternalCreateTenantRequest{
			TenantId:      tenantId,
			AdminEmail:    req.AdminEmail,
			AdminUsername: req.AdminUsername,
			AdminPassword: req.AdminPassword,
			AdminUserId:   req.AdminUserId,
		})

		if err != nil {
			return nil, fmt.Errorf("%s %w", err.Error(), dtmcli.ErrFailure)
		}
		return resp, nil
	})

	if err != nil {
		panic(err)
	}
	return s
}

func (s *TenantService) CreateTenant(ctx context.Context, req *pb.CreateTenantRequest) (*pb.Tenant, error) {

	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceTenant, "*"), authz.CreateAction); err != nil {
		return nil, err
	}

	if len(req.DisplayName) == 0 {
		req.DisplayName = req.Name
	}

	var err error
	var resp = &pb.Tenant{}
	//Workflow Transaction
	data, err := workflow.ExecuteCtx(ctx, wfCreateTenantName, ksuid.New().String(), utils.PbMustMarshalJson(req))
	if err != nil {
		return nil, err
	}
	utils.PbMustUnMarshalJson(data, resp)
	return resp, err
}

func (s *TenantService) UserCreateTenant(ctx context.Context, req *pb.UserCreateTenantRequest) (*pb.UserCreateTenantReply, error) {
	ui, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	if len(req.DisplayName) == 0 {
		req.DisplayName = req.Name
	}
	uid := ui.GetId()
	var rawReq = &pb.CreateTenantRequest{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Region:      req.Region,
		Logo:        req.Logo,
		AdminUserId: &uid,
	}
	var createTenantResp = &pb.Tenant{}

	//Workflow Transaction
	data, err := workflow.ExecuteCtx(ctx, wfCreateTenantName, ksuid.New().String(), utils.PbMustMarshalJson(rawReq))
	if err != nil {
		return nil, err
	}
	utils.PbMustUnMarshalJson(data, createTenantResp)
	return &pb.UserCreateTenantReply{Tenant: createTenantResp}, err

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
			Value: *data.NewFromDynamicValue(t.Value),
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
	ui, _ := authn.FromUserContext(ctx)
	if len(ti.GetId()) == 0 {
		return &pb.GetCurrentTenantReply{IsHost: true, Tenant: mapBizTenantToInfo(ctx, s.blob, nil, s.app)}, nil
	} else {
		tenant, err := s.useCase.FindByIdOrName(ctx, ti.GetId())
		if err != nil {
			return nil, err
		}
		if tenant == nil {
			return nil, pb.ErrorTenantNotFoundLocalized(ctx, nil, nil)
		}
		info := mapBizTenantToInfo(ctx, s.blob, tenant, s.app)
		if len(ui.GetId()) > 0 {
			info.PlanKey = tenant.PlanKey
			if tenant.Plan != nil {
				info.Plan = &v12.Plan{}
				MapBizPlan2Pb(tenant.Plan, info.Plan)
			}
		}

		return &pb.GetCurrentTenantReply{
			IsHost: false,
			Tenant: info,
		}, nil
	}
}

func (s *TenantService) ListTenant(ctx context.Context, req *pb.ListTenantRequest) (*pb.ListTenantReply, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceTenant, "*"), authz.ReadAction); err != nil {
		return nil, err
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
			Value: con.Value.ToDynamicValue(),
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
		PlanKey:     tenant.PlanKey,
	}
	if tenant.Plan != nil {
		res.Plan = &v12.Plan{}
		MapBizPlan2Pb(tenant.Plan, res.Plan)
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
