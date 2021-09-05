package service

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/saas/internal_/biz"

	pb "github.com/goxiaoy/go-saas-kit/saas/api/tenant/v1"
)

type TenantService struct {
	pb.UnimplementedTenantServer
	repo biz.TenantRepo
}

func NewTenantService(repo biz.TenantRepo) *TenantService {
	return &TenantService{repo: repo}
}

func (s *TenantService) CreateTenant(ctx context.Context, req *pb.CreateTenantRequest) (*pb.CreateTenantReply, error) {
	return &pb.CreateTenantReply{}, nil
}
func (s *TenantService) UpdateTenant(ctx context.Context, req *pb.UpdateTenantRequest) (*pb.UpdateTenantReply, error) {
	return &pb.UpdateTenantReply{}, nil
}
func (s *TenantService) DeleteTenant(ctx context.Context, req *pb.DeleteTenantRequest) (*pb.DeleteTenantReply, error) {
	return &pb.DeleteTenantReply{}, nil
}
func (s *TenantService) GetTenant(ctx context.Context, req *pb.GetTenantRequest) (*pb.GetTenantReply, error) {
	return &pb.GetTenantReply{}, nil
}
func (s *TenantService) ListTenant(ctx context.Context, req *pb.ListTenantRequest) (*pb.ListTenantReply, error) {
	return &pb.ListTenantReply{}, nil
}
