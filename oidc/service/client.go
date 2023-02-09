package service

import (
	"context"
	"github.com/go-saas/kit/oidc/api"
	pb "github.com/go-saas/kit/oidc/api/client/v1"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/utils"
	client "github.com/ory/hydra-client-go/v2"
	"github.com/peterhellberg/link"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"net/url"
	"strconv"
)

type ClientService struct {
	pb.UnimplementedClientServiceServer
	client *client.APIClient
	auth   authz.Service
}

func NewClientService(client *client.APIClient, auth authz.Service) *ClientService {
	return &ClientService{client: client, auth: auth}
}

func (s *ClientService) ListOAuth2Clients(ctx context.Context, req *pb.ListClientRequest) (*pb.OAuth2ClientList, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceClient, "*"), authz.ReadAction); err != nil {
		return nil, err
	}
	rreq := s.client.OAuth2Api.ListOAuth2Clients(ctx).ClientName(req.ClientName).PageSize(req.Limit).Owner(req.Owner)
	if len(req.AfterPageToken) > 0 {
		rreq = rreq.PageToken(req.AfterPageToken)
	}
	if len(req.BeforePageToken) > 0 {
		rreq = rreq.PageToken(req.BeforePageToken)
	}
	if len(req.AfterPageToken) > 0 {
		rreq = rreq.PageToken(req.AfterPageToken)
	}
	resp, raw, err := rreq.Execute()
	if err != nil {
		return nil, TransformHydraErr(raw, err)
	}
	total, _ := strconv.Atoi(raw.Header.Get("X-Total-Count"))

	ret := &pb.OAuth2ClientList{TotalSize: int32(total), Items: lo.Map(resp, func(t client.OAuth2Client, _ int) *pb.OAuth2Client {
		return mapClients(t)
	})}
	respLink := raw.Header.Get("Link")
	parsePageToken := func(us string) string {
		u, err := url.Parse(us)
		if err != nil {
			return ""
		}
		return u.Query().Get("page_token")
	}
	for _, l := range link.Parse(respLink) {
		if l.Rel == "next" {
			t := parsePageToken(l.URI)
			ret.NextAfterPageToken = &t
		}
		if l.Rel == "prev" {
			t := parsePageToken(l.URI)
			ret.NextBeforePageToken = &t
		}
	}
	return ret, nil
}
func (s *ClientService) GetOAuth2Client(ctx context.Context, req *pb.GetOAuth2ClientRequest) (*pb.OAuth2Client, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceClient, "*"), authz.ReadAction); err != nil {
		return nil, err
	}
	resp, raw, err := s.client.OAuth2Api.GetOAuth2Client(ctx, req.Id).Execute()
	if err != nil {
		return nil, TransformHydraErr(raw, err)
	}
	c := mapClients(*resp)
	return c, nil
}
func (s *ClientService) CreateOAuth2Client(ctx context.Context, req *pb.OAuth2Client) (*pb.OAuth2Client, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceClient, "*"), authz.CreateAction); err != nil {
		return nil, err
	}
	c := mapOAuthClients(req)
	resp, raw, err := s.client.OAuth2Api.CreateOAuth2Client(ctx).OAuth2Client(c).Execute()
	if err != nil {
		return nil, TransformHydraErr(raw, err)
	}
	return mapClients(*resp), nil
}
func (s *ClientService) DeleteOAuth2Client(ctx context.Context, req *pb.DeleteOAuth2ClientRequest) (*emptypb.Empty, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceClient, "*"), authz.DeleteAction); err != nil {
		return nil, err
	}
	raw, err := s.client.OAuth2Api.DeleteOAuth2Client(ctx, req.Id).Execute()
	if err != nil {
		return nil, TransformHydraErr(raw, err)
	}
	return &emptypb.Empty{}, nil
}
func (s *ClientService) PatchOAuth2Client(ctx context.Context, req *pb.PatchOAuth2ClientRequest) (*pb.OAuth2Client, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceClient, "*"), authz.UpdateAction); err != nil {
		return nil, err
	}
	resp, raw, err := s.client.OAuth2Api.PatchOAuth2Client(ctx, req.Id).JsonPatch(lo.Map(req.Client, func(t *pb.PatchOAuth2Client, _ int) client.JsonPatch {
		return client.JsonPatch{
			From:  t.From,
			Op:    t.Op,
			Path:  t.Path,
			Value: utils.Structpb2Map(t.Value),
		}
	})).Execute()
	if err != nil {
		return nil, TransformHydraErr(raw, err)
	}
	return mapClients(*resp), nil
}

func (s *ClientService) UpdateOAuth2Client(ctx context.Context, req *pb.UpdateOAuth2ClientRequest) (*pb.OAuth2Client, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceClient, "*"), authz.UpdateAction); err != nil {
		return nil, err
	}
	resp, raw, err := s.client.OAuth2Api.SetOAuth2Client(ctx, req.Id).OAuth2Client(mapOAuthClients(req.Client)).Execute()
	if err != nil {
		return nil, TransformHydraErr(raw, err)
	}
	return mapClients(*resp), nil
}

func mapClients(c client.OAuth2Client) *pb.OAuth2Client {
	ret := &pb.OAuth2Client{
		AllowedCorsOrigins:                c.AllowedCorsOrigins,
		Audience:                          c.Audience,
		BackchannelLogoutSessionRequired:  c.BackchannelLogoutSessionRequired,
		BackchannelLogoutUri:              c.BackchannelLogoutUri,
		ClientId:                          c.ClientId,
		ClientName:                        c.ClientName,
		ClientSecret:                      c.ClientSecret,
		ClientSecretExpiresAt:             c.ClientSecretExpiresAt,
		ClientUri:                         c.ClientUri,
		Contacts:                          c.Contacts,
		CreatedAt:                         utils.Time2Timepb(c.CreatedAt),
		FrontchannelLogoutSessionRequired: c.FrontchannelLogoutSessionRequired,
		FrontchannelLogoutUri:             c.FrontchannelLogoutUri,
		GrantTypes:                        c.GrantTypes,
		Jwks:                              safeConvert(c.Jwks),
		JwksUri:                           c.JwksUri,
		LogoUri:                           c.LogoUri,
		Metadata:                          safeConvert(c.Metadata),
		Owner:                             c.Owner,
		PolicyUri:                         c.PolicyUri,
		PostLogoutRedirectUris:            c.PostLogoutRedirectUris,
		RedirectUris:                      c.RedirectUris,
		RegistrationAccessToken:           c.RegistrationAccessToken,
		RegistrationClientUri:             c.RegistrationClientUri,
		RequestObjectSigningAlg:           c.RequestObjectSigningAlg,
		RequestUris:                       c.RequestUris,
		ResponseTypes:                     c.ResponseTypes,
		Scope:                             c.Scope,
		SectorIdentifierUri:               c.SectorIdentifierUri,
		SubjectType:                       c.SubjectType,
		TokenEndpointAuthMethod:           c.TokenEndpointAuthMethod,
		TokenEndpointAuthSigningAlg:       c.TokenEndpointAuthSigningAlg,
		TosUri:                            c.TosUri,
		UpdatedAt:                         utils.Time2Timepb(c.UpdatedAt),
		UserinfoSignedResponseAlg:         c.UserinfoSignedResponseAlg,
	}
	return ret
}

func mapOAuthClients(c *pb.OAuth2Client) client.OAuth2Client {
	ret := client.OAuth2Client{
		AllowedCorsOrigins:                c.AllowedCorsOrigins,
		Audience:                          c.Audience,
		BackchannelLogoutSessionRequired:  c.BackchannelLogoutSessionRequired,
		BackchannelLogoutUri:              c.BackchannelLogoutUri,
		ClientId:                          c.ClientId,
		ClientName:                        c.ClientName,
		ClientSecret:                      c.ClientSecret,
		ClientSecretExpiresAt:             c.ClientSecretExpiresAt,
		ClientUri:                         c.ClientUri,
		Contacts:                          c.Contacts,
		CreatedAt:                         utils.Timepb2Time(c.CreatedAt),
		FrontchannelLogoutSessionRequired: c.FrontchannelLogoutSessionRequired,
		FrontchannelLogoutUri:             c.FrontchannelLogoutUri,
		GrantTypes:                        c.GrantTypes,
		Jwks:                              utils.Structpb2Map(c.Jwks),
		JwksUri:                           c.JwksUri,
		LogoUri:                           c.LogoUri,
		Metadata:                          utils.Structpb2Map(c.Metadata),
		Owner:                             c.Owner,
		PolicyUri:                         c.PolicyUri,
		PostLogoutRedirectUris:            c.PostLogoutRedirectUris,
		RedirectUris:                      c.RedirectUris,
		RegistrationAccessToken:           c.RegistrationAccessToken,
		RegistrationClientUri:             c.RegistrationClientUri,
		RequestObjectSigningAlg:           c.RequestObjectSigningAlg,
		RequestUris:                       c.RequestUris,
		ResponseTypes:                     c.ResponseTypes,
		Scope:                             c.Scope,
		SectorIdentifierUri:               c.SectorIdentifierUri,
		SubjectType:                       c.SubjectType,
		TokenEndpointAuthMethod:           c.TokenEndpointAuthMethod,
		TokenEndpointAuthSigningAlg:       c.TokenEndpointAuthSigningAlg,
		TosUri:                            c.TosUri,
		UpdatedAt:                         utils.Timepb2Time(c.UpdatedAt),
		UserinfoSignedResponseAlg:         c.UserinfoSignedResponseAlg,
	}
	return ret
}

func safeConvert(a interface{}) *structpb.Struct {
	if a == nil {
		return nil
	}
	return utils.Map2Structpb(a.(map[string]interface{}))
}
