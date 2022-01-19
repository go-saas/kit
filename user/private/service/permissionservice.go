package service

import (
	"context"

	pb "github.com/goxiaoy/go-saas-kit/user/api/permission/v1"
)

type PermissionServiceService struct {
	pb.UnimplementedPermissionServiceServer
}

func NewPermissionServiceService() *PermissionServiceService {
	return &PermissionServiceService{}
}

func (s *PermissionServiceService) GetCurrent(ctx context.Context, req *pb.GetCurrentPermissionRequest) (*pb.GetCurrentPermissionReply, error) {
	return &pb.GetCurrentPermissionReply{}, nil
}
func (s *PermissionServiceService) CheckCurrent(ctx context.Context, req *pb.CheckPermissionRequest) (*pb.CheckPermissionReply, error) {
	return &pb.CheckPermissionReply{}, nil
}
func (s *PermissionServiceService) UpdateSubjectPermission(ctx context.Context, req *pb.UpdateSubjectPermissionRequest) (*pb.UpdateSubjectPermissionResponse, error) {
	return &pb.UpdateSubjectPermissionResponse{}, nil
}
