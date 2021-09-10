package service

import (
	"context"
	"errors"
	"github.com/ahmetb/go-linq/v3"
	errors2 "github.com/go-kratos/kratos/v2/errors"
	"github.com/golang/protobuf/ptypes/wrappers"
	pb "github.com/goxiaoy/go-saas-kit/user/api/user/v1"
	"github.com/goxiaoy/go-saas-kit/user/internal_/biz"
	"github.com/mennanov/fmutils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	um *biz.UserManager
}

func NewUserService(um *biz.UserManager) *UserService {
	return &UserService{
		um: um,
	}
}

func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	ret := &pb.ListUsersResponse{}

	totalCount, filterCount, err := s.um.Count(ctx, req.Filter)
	if err != nil {
		return nil, err
	}
	ret.TotalSize = totalCount
	ret.FilterSize = filterCount

	items, err := s.um.List(ctx, req)
	if err != nil {
		return ret, err
	}
	rItems := make([]*pb.User, len(items))
	for index, u := range items {
		res := MapBizUserToApi(u)
		if req.Fields!=nil{
			fmutils.Filter(res, req.Fields.Paths)
		}
		rItems[index] = res
	}
	ret.Items = rItems
	return ret, nil

}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	u, err := s.um.FindByID(ctx, req.Id)
	if err!=nil{
		return nil,err
	}
	if u==nil{
		return nil,errors2.NotFound("","")
	}
	res := MapBizUserToApi(u)
	return res, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	// check confirm password
	if req.Password != "" {
		if req.ConfirmPassword != req.Password {
			return nil, pb.ErrorConfirmPasswordMismatch("", "")
		}
	}
	u := biz.User{}
	if req.Name != nil {
		u.Name = &req.Name.Value
	}
	if req.Username != nil {
		u.Username = &req.Username.Value
	}
	if req.Email != nil {
		u.Email = &req.Email.Value
	}
	if req.Phone != nil {
		u.Phone = &req.Phone.Value
	}

	if req.Birthday != nil {
		b := req.Birthday.AsTime()
		u.Birthday = &b
	}
	gender := req.Gender.String()
	u.Gender = &gender
	var err error
	if req.Password != "" {
		err = s.um.CreateWithPassword(ctx, &u, req.Password)
	} else {
		err = s.um.Create(ctx, &u)
	}
	if err != nil {
		if errors.Is(err, biz.ErrInsufficientStrength) {
			return nil, pb.ErrorPasswordInsufficientStrength("", "")
		}
		if errors.Is(err, biz.ErrDuplicateEmail) {
			return nil, pb.ErrorDuplicateEmail("", "")
		}
		if errors.Is(err, biz.ErrDuplicateUsername) {
			return nil, pb.ErrorDuplicateUsername("", "")
		}
		if errors.Is(err, biz.ErrDuplicatePhone) {
			return nil, pb.ErrorDuplicatePhone("", "")
		}
		return nil, err
	}
	res := MapBizUserToApi(&u)
	return res, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	// check confirm password
	if req.User.Password != "" {
		if req.User.ConfirmPassword != req.User.Password {
			return nil, pb.ErrorConfirmPasswordMismatch("", "")
		}
	}

	if req.UpdateMask!=nil{
		fmutils.Filter(req, req.UpdateMask.Paths)
	}
	u,err:=s.um.FindByID(ctx,req.User.Id)
	if err!=nil{
		return nil,err
	}
	if u==nil{
		return nil,errors2.NotFound("","")
	}

	if req.User.Password!=""{
		//reset password
		if err:=s.um.UpdatePassword(ctx,u,req.User.Password);err!=nil{
			return nil,err
		}
	}
	if req.GetUser().GetUsername()!=nil{
		v := req.GetUser().GetUsername().Value
		u.Username = &v
	}

	if req.GetUser().GetName()!=nil{
		v := req.GetUser().GetName().Value
		u.Name = &v
	}
	if req.GetUser().GetPhone()!=nil{
		v := req.GetUser().GetPhone().Value
		u.Phone = &v
	}
	if req.GetUser().GetEmail()!=nil{
		v := req.GetUser().GetEmail().Value
		u.Email = &v
	}
	if req.GetUser().GetBirthday()!=nil{
		v := req.GetUser().GetBirthday().AsTime()
		u.Birthday = &v
	}

	g:=req.GetUser().Gender.Enum().String()
	u.Gender =  &g
	if err := s.um.Update(ctx,u);err!=nil{
		return nil,err
	}
	res :=MapBizUserToApi(u)
	return res, nil
}
func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	return &pb.DeleteUserResponse{}, nil
}

func MapBizUserToApi(u *biz.User ) *pb.User  {
	res := &pb.User{
		Id:       u.ID.String(),
		Roles:    nil,
	}
	if u.Username != nil {
		res.Username = &wrappers.StringValue{Value: *u.Username}
	}
	if u.Name != nil {
		res.Name = &wrappers.StringValue{Value: *u.Name}
	}
	if u.Phone != nil {
		res.Phone = &wrappers.StringValue{Value: *u.Phone}
	}
	if u.Email != nil {
		res.Email = &wrappers.StringValue{Value: *u.Email}
	}
	if u.Birthday != nil {
		res.Birthday = timestamppb.New(*u.Birthday)
	}
	if u.Gender != nil {
		if v, ok := pb.Gender_value[*u.Gender]; ok {
			res.Gender = pb.Gender(v)
		}
	}
	if u.Roles!=nil{
		var returnRoles []*pb.Role
		linq.From(u.Roles).SelectT(func(i biz.Role) *pb.Role{
			return &pb.Role{
				Id:   i.ID.String(),
				Name: i.Name,
			}
		}).ToSlice(&returnRoles)
		res.Roles = returnRoles
	}
	return res
}