package service

import (
	"context"
	"github.com/go-saas/kit/oidc/api"
	"github.com/go-saas/kit/pkg/authz/authz"
	client "github.com/ory/hydra-client-go"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"time"

	pb "github.com/go-saas/kit/oidc/api/client/v1"
	"google.golang.org/protobuf/types/known/emptypb"
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
	resp, raw, err := s.client.AdminApi.ListOAuth2Clients(ctx).ClientName(req.ClientName).Limit(req.Limit).Offset(req.Offset).Owner(req.Owner).Execute()
	if err != nil {
		return nil, err
	}
	total, _ := strconv.Atoi(raw.Header.Get("X-Total-Count"))

	return &pb.OAuth2ClientList{TotalCount: int32(total), Items: lo.Map(resp, func(t client.OAuth2Client, _ int) *pb.OAuth2Client {
		return mapClients(t)
	})}, nil
}
func (s *ClientService) GetOAuth2Client(ctx context.Context, req *pb.GetOAuth2ClientRequest) (*pb.OAuth2Client, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceClient, "*"), authz.ReadAction); err != nil {
		return nil, err
	}
	resp, _, err := s.client.AdminApi.GetOAuth2Client(ctx, req.Id).Execute()
	if err != nil {
		return nil, err
	}
	c := mapClients(*resp)
	return c, nil
}
func (s *ClientService) CreateOAuth2Client(ctx context.Context, req *pb.OAuth2Client) (*pb.OAuth2Client, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceClient, "*"), authz.CreateAction); err != nil {
		return nil, err
	}
	c := mapOAuthClients(req)
	resp, _, err := s.client.AdminApi.CreateOAuth2Client(ctx).OAuth2Client(c).Execute()
	if err != nil {
		return nil, err
	}
	return mapClients(*resp), nil
}
func (s *ClientService) DeleteOAuth2Client(ctx context.Context, req *pb.DeleteOAuth2ClientRequest) (*emptypb.Empty, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceClient, "*"), authz.DeleteAction); err != nil {
		return nil, err
	}
	_, err := s.client.AdminApi.DeleteOAuth2Client(ctx, req.Id).Execute()
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func (s *ClientService) PatchOAuth2Client(ctx context.Context, req *pb.PatchOAuth2ClientRequest) (*pb.OAuth2Client, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceClient, "*"), authz.UpdateAction); err != nil {
		return nil, err
	}
	resp, _, err := s.client.AdminApi.PatchOAuth2Client(ctx, req.Id).PatchDocument(lo.Map(req.Client, func(t *pb.PatchOAuth2Client, _ int) client.PatchDocument {
		return client.PatchDocument{
			From:  t.From,
			Op:    t.Op,
			Path:  t.Path,
			Value: mapStruct(t.Value),
		}
	})).Execute()
	if err != nil {
		return nil, err
	}
	return mapClients(*resp), nil
}

func (s *ClientService) UpdateOAuth2Client(ctx context.Context, req *pb.UpdateOAuth2ClientRequest) (*pb.OAuth2Client, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceClient, "*"), authz.UpdateAction); err != nil {
		return nil, err
	}
	resp, _, err := s.client.AdminApi.UpdateOAuth2Client(ctx, req.Id).OAuth2Client(mapOAuthClients(req.Client)).Execute()
	if err != nil {
		return nil, err
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
		CreatedAt:                         mapTime(c.CreatedAt),
		FrontchannelLogoutSessionRequired: c.FrontchannelLogoutSessionRequired,
		FrontchannelLogoutUri:             c.FrontchannelLogoutUri,
		GrantTypes:                        c.GrantTypes,
		Jwks:                              mapInterface(c.Jwks),
		JwksUri:                           c.JwksUri,
		LogoUri:                           c.LogoUri,
		Metadata:                          mapInterface(c.Metadata),
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
		UpdatedAt:                         mapTime(c.UpdatedAt),
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
		CreatedAt:                         mapPbTime(c.CreatedAt),
		FrontchannelLogoutSessionRequired: c.FrontchannelLogoutSessionRequired,
		FrontchannelLogoutUri:             c.FrontchannelLogoutUri,
		GrantTypes:                        c.GrantTypes,
		Jwks:                              mapStruct(c.Jwks),
		JwksUri:                           c.JwksUri,
		LogoUri:                           c.LogoUri,
		Metadata:                          mapStruct(c.Metadata),
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
		UpdatedAt:                         mapPbTime(c.UpdatedAt),
		UserinfoSignedResponseAlg:         c.UserinfoSignedResponseAlg,
	}
	return ret
}

func mapTime(time *time.Time) *timestamppb.Timestamp {
	if time == nil {
		return nil
	}
	return timestamppb.New(*time)
}
func mapInterface(m map[string]interface{}) *structpb.Struct {
	if m == nil {
		return nil
	}
	r, _ := structpb.NewStruct(m)
	return r
}

func mapStruct(m *structpb.Struct) map[string]interface{} {
	if m == nil {
		return nil
	}
	return m.AsMap()
}
func mapPbTime(t *timestamppb.Timestamp) *time.Time {
	if t == nil {
		return nil
	}
	ret := t.AsTime()
	return &ret
}
