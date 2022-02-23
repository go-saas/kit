package service

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/sys/private/biz"

	pb "github.com/goxiaoy/go-saas-kit/sys/api/menu/v1"
)

type MenuService struct {
	pb.UnimplementedMenuServiceServer
	auth authz.Service
	repo biz.MenuRepo
}

func NewMenuService(auth authz.Service, repo biz.MenuRepo) *MenuService {
	return &MenuService{auth: auth, repo: repo}
}

func (s *MenuService) ListMenu(ctx context.Context, req *pb.ListMenuRequest) (*pb.ListMenuReply, error) {
	return &pb.ListMenuReply{}, nil
}
func (s *MenuService) GetMenu(ctx context.Context, req *pb.GetMenuRequest) (*pb.Menu, error) {
	return &pb.Menu{}, nil
}
func (s *MenuService) CreateMenu(ctx context.Context, req *pb.CreateMenuRequest) (*pb.Menu, error) {
	return &pb.Menu{}, nil
}
func (s *MenuService) UpdateMenu(ctx context.Context, req *pb.UpdateMenuRequest) (*pb.Menu, error) {
	return &pb.Menu{}, nil
}
func (s *MenuService) DeleteMenu(ctx context.Context, req *pb.DeleteMenuRequest) (*pb.DeleteMenuReply, error) {
	return &pb.DeleteMenuReply{}, nil
}
