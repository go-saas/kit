package service

import (
	"context"

	pb "github.com/go-saas/kit/oidc/api/client/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClientService struct {
	pb.UnimplementedClientServiceServer
}

func NewClientService() *ClientService {
	return &ClientService{}
}

func (s *ClientService) ListOAuth2Clients(ctx context.Context, req *pb.ListClientRequest) (*pb.OAuth2ClientList, error) {
	return &pb.OAuth2ClientList{}, nil
}
func (s *ClientService) GetOAuth2Client(ctx context.Context, req *pb.GetOAuth2ClientRequest) (*pb.OAuth2Client, error) {
	return &pb.OAuth2Client{}, nil
}
func (s *ClientService) CreateOAuth2Client(ctx context.Context, req *pb.OAuth2Client) (*pb.OAuth2Client, error) {
	return &pb.OAuth2Client{}, nil
}
func (s *ClientService) DeleteOAuth2Client(ctx context.Context, req *pb.DeleteOAuth2ClientRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
func (s *ClientService) PatchOAuth2Client(ctx context.Context, req *pb.PatchOAuth2ClientRequest) (*pb.OAuth2Client, error) {
	return &pb.OAuth2Client{}, nil
}
func (s *ClientService) UpdateOAuth2Client(ctx context.Context, req *pb.PatchOAuth2ClientRequest) (*pb.OAuth2Client, error) {
	return &pb.OAuth2Client{}, nil
}
