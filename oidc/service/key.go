package service

import (
	"context"
	"github.com/go-saas/kit/oidc/api"
	"github.com/go-saas/kit/pkg/authz/authz"
	client "github.com/ory/hydra-client-go/v2"
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
	raw, err := s.client.JwkApi.DeleteJsonWebKeySet(ctx, req.Set).Execute()
	if err != nil {
		return nil, TransformHydraErr(raw, err)
	}
	return &emptypb.Empty{}, nil
}
func (s *KeyService) GetJsonWebKeySet(ctx context.Context, req *pb.GetJsonWebKeySetRequest) (*pb.JsonWebKeySet, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceKey, "*"), authz.ReadAction); err != nil {
		return nil, err
	}
	resp, raw, err := s.client.JwkApi.GetJsonWebKeySet(ctx, req.Set).Execute()
	if err != nil {
		return nil, TransformHydraErr(raw, err)
	}
	return mapSet(*resp), nil
}
func (s *KeyService) CreateJsonWebKeySet(ctx context.Context, req *pb.CreateJsonWebKeySetRequest) (*pb.JsonWebKeySet, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceKey, "*"), authz.CreateAction); err != nil {
		return nil, err
	}
	resp, raw, err := s.client.JwkApi.CreateJsonWebKeySet(ctx, req.Set).CreateJsonWebKeySet(client.CreateJsonWebKeySet{
		Alg: req.Keys.Alg,
		Kid: req.Keys.Kid,
		Use: req.Keys.Use,
	}).Execute()
	if err != nil {
		return nil, TransformHydraErr(raw, err)
	}
	return mapSet(*resp), nil
}
func (s *KeyService) UpdateJsonWebKeySet(ctx context.Context, req *pb.UpdateJsonWebKeySetRequest) (*pb.JsonWebKeySet, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceKey, "*"), authz.UpdateAction); err != nil {
		return nil, err
	}
	resp, raw, err := s.client.JwkApi.SetJsonWebKeySet(ctx, req.Set).JsonWebKeySet(mapPbSet(req.Keys)).Execute()
	if err != nil {
		return nil, TransformHydraErr(raw, err)
	}
	return mapSet(*resp), nil
}
func (s *KeyService) DeleteJsonWebKey(ctx context.Context, req *pb.DeleteJsonWebKeyRequest) (*emptypb.Empty, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceKey, "*"), authz.DeleteAction); err != nil {
		return nil, err
	}
	raw, err := s.client.JwkApi.DeleteJsonWebKey(ctx, req.Kid, req.Set).Execute()
	if err != nil {
		return nil, TransformHydraErr(raw, err)
	}
	return &emptypb.Empty{}, nil
}
func (s *KeyService) GetJsonWebKey(ctx context.Context, req *pb.GetJsonWebKeyRequest) (*pb.JsonWebKeySet, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceKey, "*"), authz.ReadAction); err != nil {
		return nil, err
	}
	resp, raw, err := s.client.JwkApi.GetJsonWebKey(ctx, req.Kid, req.Set).Execute()
	if err != nil {
		return nil, TransformHydraErr(raw, err)
	}
	return mapSet(*resp), nil
}

func (s *KeyService) UpdateJsonWebKey(ctx context.Context, req *pb.UpdateJsonWebKeyRequest) (*pb.JsonWebKey, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceKey, "*"), authz.UpdateAction); err != nil {
		return nil, err
	}
	resp, raw, err := s.client.JwkApi.SetJsonWebKey(ctx, req.Kid, req.Set).JsonWebKey(mapPbKey(req.Key)).Execute()
	if err != nil {
		return nil, TransformHydraErr(raw, err)
	}
	return mapKey(*resp), nil
}

func mapKey(a client.JsonWebKey) *pb.JsonWebKey {
	return &pb.JsonWebKey{
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

func mapPbKey(a *pb.JsonWebKey) client.JsonWebKey {
	return client.JsonWebKey{
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
func mapSet(a client.JsonWebKeySet) *pb.JsonWebKeySet {
	return &pb.JsonWebKeySet{Keys: lo.Map(a.Keys, func(t client.JsonWebKey, _ int) *pb.JsonWebKey {
		return mapKey(t)
	})}
}
func mapPbSet(a *pb.JsonWebKeySet) client.JsonWebKeySet {
	return client.JsonWebKeySet{Keys: lo.Map(a.Keys, func(t *pb.JsonWebKey, _ int) client.JsonWebKey {
		return mapPbKey(t)
	})}
}
