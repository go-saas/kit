
package service

import(
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/types/known/timestamppb"
	"github.com/goxiaoy/go-saas-kit/auth/current"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/user/v1"
	"github.com/goxiaoy/go-saas-kit/user/internal_/biz"

	pb "github.com/goxiaoy/go-saas-kit/user/api/account/v1"
)

type AccountService struct {
	pb.UnimplementedAccountServer
	um *biz.UserManager
}

func NewAccountService(um *biz.UserManager) *AccountService {
	return &AccountService{
		um: um,
	}
}

func (s *AccountService) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	userinfo,ok := current.FromUserContext(ctx)
	if !ok {
		return nil,errors.Unauthorized("","")
	}
	u ,err:=s.um.FindByID(ctx,userinfo.Id)
	if err!=nil {
		return nil,errors.Forbidden("","")
	}
	res := &pb.GetProfileResponse{
		Id:u.ID.String(),
		Username: &wrappers.StringValue{Value: *u.Username} ,
	}
	if u.Name!=nil{
		res.Name =  &wrappers.StringValue{Value: *u.Name }
	}
	if u.Phone!=nil{
		res.Phone =  &wrappers.StringValue{Value: *u.Phone }
	}
	if u.Email!=nil{
		res.Email = &wrappers.StringValue{Value: *u.Email }
	}
	if u.Birthday!=nil{
		res.Birthday =  timestamppb.New(*u.Birthday)
	}
	if u.Gender!=nil{
		v,ok:= v1.Gender_value[*u.Gender]
		if ok{
			res.Gender = *v1.Gender(v).Enum()
		}
	}
	for _, role := range u.Roles {
		res.Roles = append(res.Roles, &v1.Role{
			Id:   role.ID.String(),
			Name: role.Name,
		})
	}
	return res, nil
}
func (s *AccountService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	return &pb.UpdateProfileResponse{}, nil
}
func (s *AccountService) GetSettings(ctx context.Context, req *pb.GetSettingsRequest) (*pb.GetSettingsResponse, error) {
	return &pb.GetSettingsResponse{}, nil
}
func (s *AccountService) UpdateSettings(ctx context.Context, req *pb.UpdateSettingsRequest) (*pb.UpdateSettingsResponse, error) {
	return &pb.UpdateSettingsResponse{}, nil
}
func (s *AccountService) GetAddresses(ctx context.Context, req *pb.GetAddressesRequest) (*pb.GetAddressesReply, error) {
	return &pb.GetAddressesReply{}, nil
}
func (s *AccountService) UpdateAddresses(ctx context.Context, req *pb.UpdateAddressesRequest) (*pb.UpdateAddressesReply, error) {
	return &pb.UpdateAddressesReply{}, nil
}
