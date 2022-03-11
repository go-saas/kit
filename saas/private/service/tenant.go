package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/uuid"
	"github.com/goxiaoy/go-saas-kit/pkg/authn"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	pb "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
	"github.com/goxiaoy/go-saas-kit/saas/private/biz"
	"github.com/goxiaoy/go-saas/common"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"os"
	"path/filepath"
)

type TenantService struct {
	pb.UnimplementedTenantServiceServer
	useCase *biz.TenantUseCase
	auth    authz.Service
	blob    blob.Factory
}

func NewTenantService(useCase *biz.TenantUseCase, auth authz.Service, blob blob.Factory) *TenantService {
	return &TenantService{useCase: useCase, auth: auth, blob: blob}
}

func (s *TenantService) CreateTenant(ctx context.Context, req *pb.CreateTenantRequest) (*pb.Tenant, error) {

	if _, err := authz.CheckForHostOnly(ctx, s.auth, authz.NewEntityResource("saas.tenant", "*"), authz.CreateAction); err != nil {
		return nil, err
	}

	disPlayName := req.Name
	if req.DisplayName != "" {
		disPlayName = req.DisplayName
	}
	t := &biz.Tenant{
		Name:        req.Name,
		DisplayName: disPlayName,
		Region:      req.Region,
		Logo:        req.Logo,
	}
	if err := s.useCase.Create(ctx, t); err != nil {
		return nil, err
	}

	return mapBizTenantToApi(ctx, s.blob, t), nil
}
func (s *TenantService) UpdateTenant(ctx context.Context, req *pb.UpdateTenantRequest) (*pb.Tenant, error) {

	if _, err := authz.CheckForHostOnly(ctx, s.auth, authz.NewEntityResource("saas.tenant", req.Tenant.Id), authz.UpdateAction); err != nil {
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

	if err := s.useCase.Update(ctx, t, req.UpdateMask); err != nil {
		return nil, err
	}
	return mapBizTenantToApi(ctx, s.blob, t), nil
}
func (s *TenantService) DeleteTenant(ctx context.Context, req *pb.DeleteTenantRequest) (*pb.DeleteTenantReply, error) {

	if _, err := authz.CheckForHostOnly(ctx, s.auth, authz.NewEntityResource("saas.tenant", req.Id), authz.DeleteAction); err != nil {
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
	//TODO separate this check???
	if claim, ok := jwt.FromClaimsContext(ctx); ok && len(claim.ClientId) > 0 {
		//internal api call
		return mapBizTenantToApi(ctx, s.blob, t), nil
	}
	if _, err := s.auth.Check(ctx, authz.NewEntityResource("saas.tenant", t.ID), authz.GetAction); err != nil {
		return nil, err
	}

	return mapBizTenantToApi(ctx, s.blob, t), nil
}

func (s *TenantService) GetCurrentTenant(ctx context.Context, req *pb.GetCurrentTenantRequest) (*pb.GetCurrentTenantReply, error) {
	ti, _ := common.FromCurrentTenant(ctx)
	if len(ti.GetId()) == 0 {
		return &pb.GetCurrentTenantReply{IsHost: true}, nil
	} else {
		t, err := s.useCase.FindByIdOrName(ctx, ti.GetId())
		if err != nil {
			return nil, err
		}
		if t == nil {
			return nil, pb.ErrorTenantNotFound("")
		}
		return &pb.GetCurrentTenantReply{
			Id:          t.ID,
			Name:        t.Name,
			DisplayName: t.DisplayName,
			Region:      t.Region,
			Logo:        mapLogo(ctx, s.blob, t),
		}, nil
	}
}

func (s *TenantService) ListTenant(ctx context.Context, req *pb.ListTenantRequest) (*pb.ListTenantReply, error) {

	if authResult, err := authz.CheckForHostOnly(ctx, s.auth, authz.NewEntityResource("saas.tenant", "*"), authz.ListAction); err != nil {
		return nil, err
	} else if !authResult.Allowed {
		return nil, errors.Forbidden("", "")
	}

	ret := &pb.ListTenantReply{}

	totalCount, filterCount, err := s.useCase.Count(ctx, req.Search, req.Filter)
	ret.TotalSize = int32(totalCount)
	ret.FilterSize = int32(filterCount)
	if err != nil {
		return ret, err
	}
	items, err := s.useCase.List(ctx, req)
	if err != nil {
		return ret, err
	}
	rItems := lo.Map(items, func(g *biz.Tenant, _ int) *pb.Tenant { return mapBizTenantToApi(ctx, s.blob, g) })
	ret.Items = rItems
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
			if _, err := s.auth.Check(ctx, authz.NewEntityResource("saas.tenant", tenantID), authz.UpdateAction); err != nil {
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
		normalizedName := fmt.Sprintf("tenant/logo/%s%s", uuid.New().String(), ext)
		logoBlob := biz.LogoBlob(ctx, s.blob)
		a := logoBlob.GetAfero()
		err = a.MkdirAll("tenant/logo", 0755)
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
			if err := s.useCase.Update(ctx, t, &fieldmaskpb.FieldMask{Paths: []string{"logo"}}); err != nil {
				return nil, err
			}
		}
		profile := biz.LogoBlob(ctx, s.blob)

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

func mapBizTenantToApi(ctx context.Context, blob blob.Factory, tenant *biz.Tenant) *pb.Tenant {
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
		Id:          tenant.ID,
		Name:        tenant.Name,
		DisplayName: tenant.DisplayName,
		Region:      tenant.Region,
		CreatedAt:   timestamppb.New(tenant.CreatedAt),
		UpdatedAt:   timestamppb.New(tenant.UpdatedAt),
		Conn:        conns,
		Features:    features,
		Logo:        mapLogo(ctx, blob, tenant),
	}

	return res
}

func mapLogo(ctx context.Context, factory blob.Factory, entity *biz.Tenant) *blob.BlobFile {
	if entity.Logo == "" {
		return nil
	}
	profile := biz.LogoBlob(ctx, factory)

	url, _ := profile.GeneratePublicUrl(entity.Logo)
	return &blob.BlobFile{
		Id:   entity.Logo,
		Name: "",
		Url:  url,
	}
}
