package service

import (
	"context"
	"github.com/go-saas/kit/oidc/api"
	"github.com/go-saas/kit/pkg/authz/authz"
	client "github.com/ory/hydra-client-go"
	"github.com/samber/lo"

	pb "github.com/go-saas/kit/oidc/api/key/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type KeyService struct {
	pb.UnimplementedKeyServiceServer
	client *client.APIClient
	auth   authz.Service
}

func NewKeyService(client *client.APIClient, auth authz.Service) *KeyService {
	return &KeyService{client: client, auth: auth}
}

func (s *KeyService) DeleteJsonWebKeySet(ctx context.Context, req *pb.DeleteJsonWebKeySetRequest) (*emptypb.Empty, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceKey, "*"), authz.DeleteAction); err != nil {
		return nil, err
	}
	raw, err := s.client.AdminApi.DeleteJsonWebKeySet(ctx, req.Set).Execute()
	if err != nil {
		return nil, TransformHydraErr(raw, err)
	}
	return &emptypb.Empty{}, nil
}
func (s *KeyService) GetJsonWebKeySet(ctx context.Context, req *pb.GetJsonWebKeySetRequest) (*pb.JSONWebKeySet, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceKey, "*"), authz.ReadAction); err != nil {
		return nil, err
	}
	resp, raw, err := s.client.AdminApi.GetJsonWebKeySet(ctx, req.Set).Execute()
	if err != nil {
		return nil, TransformHydraErr(raw, err)
	}
	return mapSet(*resp), nil
}
func (s *KeyService) CreateJsonWebKeySet(ctx context.Context, req *pb.CreateJsonWebKeySetRequest) (*pb.JSONWebKeySet, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceKey, "*"), authz.CreateAction); err != nil {
		return nil, err
	}
	resp, raw, err := s.client.AdminApi.CreateJsonWebKeySet(ctx, req.Set).JsonWebKeySetGeneratorRequest(client.JsonWebKeySetGeneratorRequest{
		Alg: req.Keys.Alg,
		Kid: req.Keys.Kid,
		Use: req.Keys.Use,
	}).Execute()
	if err != nil {
		return nil, TransformHydraErr(raw, err)
	}
	return mapSet(*resp), nil
}
func (s *KeyService) UpdateJsonWebKeySet(ctx context.Context, req *pb.UpdateJsonWebKeySetRequest) (*pb.JSONWebKeySet, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceKey, "*"), authz.UpdateAction); err != nil {
		return nil, err
	}
	resp, raw, err := s.client.AdminApi.UpdateJsonWebKeySet(ctx, req.Set).JSONWebKeySet(mapPbSet(req.Keys)).Execute()
	if err != nil {
		return nil, TransformHydraErr(raw, err)
	}
	return mapSet(*resp), nil
}
func (s *KeyService) DeleteJsonWebKey(ctx context.Context, req *pb.DeleteJsonWebKeyRequest) (*emptypb.Empty, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceKey, "*"), authz.DeleteAction); err != nil {
		return nil, err
	}
	raw, err := s.client.AdminApi.DeleteJsonWebKey(ctx, req.Kid, req.Set).Execute()
	if err != nil {
		return nil, TransformHydraErr(raw, err)
	}
	return &emptypb.Empty{}, nil
}
func (s *KeyService) GetJsonWebKey(ctx context.Context, req *pb.GetJsonWebKeyRequest) (*pb.JSONWebKeySet, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceKey, "*"), authz.ReadAction); err != nil {
		return nil, err
	}
	resp, raw, err := s.client.AdminApi.GetJsonWebKey(ctx, req.Kid, req.Set).Execute()
	if err != nil {
		return nil, TransformHydraErr(raw, err)
	}
	return mapSet(*resp), nil
}

func (s *KeyService) UpdateJsonWebKey(ctx context.Context, req *pb.UpdateJsonWebKeyRequest) (*pb.JSONWebKey, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceKey, "*"), authz.UpdateAction); err != nil {
		return nil, err
	}
	resp, raw, err := s.client.AdminApi.UpdateJsonWebKey(ctx, req.Kid, req.Set).JSONWebKey(mapPbKey(req.Key)).Execute()
	if err != nil {
		return nil, TransformHydraErr(raw, err)
	}
	return mapKey(*resp), nil
}

func mapKey(a client.JSONWebKey) *pb.JSONWebKey {
	return &pb.JSONWebKey{
		Alg: a.Alg,
		Crv: a.Crv,
		D:   a.D,
		Dp:  a.Dp,
		Dq:  a.Dq,
		E:   a.E,
		K:   a.K,
		Kid: a.Kid,
		Kty: a.Kty,
		N:   a.N,
		P:   a.P,
		Q:   a.Q,
		Qi:  a.Qi,
		Use: a.Use,
		X:   a.X,
		X5C: a.X5c,
		Y:   a.Y,
	}
}

func mapPbKey(a *pb.JSONWebKey) client.JSONWebKey {
	return client.JSONWebKey{
		Alg: a.Alg,
		Crv: a.Crv,
		D:   a.D,
		Dp:  a.Dp,
		Dq:  a.Dq,
		E:   a.E,
		K:   a.K,
		Kid: a.Kid,
		Kty: a.Kty,
		N:   a.N,
		P:   a.P,
		Q:   a.Q,
		Qi:  a.Qi,
		Use: a.Use,
		X:   a.X,
		X5c: a.X5C,
		Y:   a.Y,
	}
}
func mapSet(a client.JSONWebKeySet) *pb.JSONWebKeySet {
	return &pb.JSONWebKeySet{Keys: lo.Map(a.Keys, func(t client.JSONWebKey, _ int) *pb.JSONWebKey {
		return mapKey(t)
	})}
}
func mapPbSet(a *pb.JSONWebKeySet) client.JSONWebKeySet {
	return client.JSONWebKeySet{Keys: lo.Map(a.Keys, func(t *pb.JSONWebKey, _ int) client.JSONWebKey {
		return mapPbKey(t)
	})}
}
