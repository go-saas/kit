package service

import (
	"context"

	pb "github.com/go-saas/kit/oidc/api/key/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type KeyService struct {
	pb.UnimplementedKeyServiceServer
}

func NewKeyService() *KeyService {
	return &KeyService{}
}

func (s *KeyService) DeleteJsonWebKeySet(ctx context.Context, req *pb.DeleteJsonWebKeySetRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
func (s *KeyService) GetJsonWebKeySet(ctx context.Context, req *pb.GetJsonWebKeySetRequest) (*pb.JSONWebKeySet, error) {
	return &pb.JSONWebKeySet{}, nil
}
func (s *KeyService) CreateJsonWebKeySet(ctx context.Context, req *pb.CreateJsonWebKeySetRequest) (*pb.JSONWebKeySet, error) {
	return &pb.JSONWebKeySet{}, nil
}
func (s *KeyService) UpdateJsonWebKeySet(ctx context.Context, req *pb.UpdateJsonWebKeySetRequest) (*pb.JSONWebKeySet, error) {
	return &pb.JSONWebKeySet{}, nil
}
func (s *KeyService) DeleteJsonWebKey(ctx context.Context, req *pb.DeleteJsonWebKeyRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
func (s *KeyService) GetJsonWebKey(ctx context.Context, req *pb.GetJsonWebKeyRequest) (*pb.JSONWebKeySet, error) {
	return &pb.JSONWebKeySet{}, nil
}
func (s *KeyService) UpdateJsonWebKey(ctx context.Context, req *pb.UpdateJsonWebKeyRequest) (*pb.JSONWebKey, error) {
	return &pb.JSONWebKey{}, nil
}
