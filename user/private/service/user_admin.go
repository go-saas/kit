package service

import (
	"context"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/user/api"
	pb "github.com/go-saas/kit/user/api/user/v1"
	"github.com/go-saas/kit/user/private/biz"
	"github.com/mennanov/fmutils"
)

func (s *UserService) ListUsersAdmin(ctx context.Context, req *pb.AdminListUsersRequest) (*pb.AdminListUsersResponse, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceAdminUser, "*"), authz.ReadAction); err != nil {
		return nil, err
	}
	ret := &pb.AdminListUsersResponse{}
	totalCount, filterCount, err := s.um.CountAdmin(ctx, req)
	if err != nil {
		return nil, err
	}
	ret.TotalSize = int32(totalCount)
	ret.FilterSize = int32(filterCount)

	items, err := s.um.ListAdmin(ctx, req)
	if err != nil {
		return ret, err
	}
	rItems := make([]*pb.User, len(items))
	for index, u := range items {
		res := MapBizUserToApi(ctx, u, s.blob)
		if req.Fields != nil {
			fmutils.Filter(res, req.Fields.Paths)
		}
		rItems[index] = res
	}
	ret.Items = rItems
	return ret, nil
}

func (s *UserService) GetUserAdmin(ctx context.Context, req *pb.AdminGetUserRequest) (*pb.User, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceAdminUser, req.Id), authz.ReadAction); err != nil {
		return nil, err
	}
	u, err := s.um.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, pb.ErrorUserNotFoundLocalized(ctx, nil, nil)
	}
	res := MapBizUserToApi(ctx, u, s.blob)
	return res, nil
}

func (s *UserService) CreateUserAdmin(ctx context.Context, req *pb.AdminCreateUserRequest) (*pb.User, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceAdminUser, "*"), authz.CreateAction); err != nil {
		return nil, err
	}
	// check confirm password
	if req.Password != "" {
		if req.ConfirmPassword != req.Password {
			return nil, pb.ErrorConfirmPasswordMismatchLocalized(ctx, nil, nil)
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
		u.SetEmail(req.Email.Value, false)
	}
	if req.Phone != nil {
		u.SetPhone(req.Phone.Value, false)
	}
	if req.Birthday != nil {
		b := req.Birthday.AsTime()
		u.Birthday = &b
	}
	gender := req.Gender.String()
	if len(req.Avatar) > 0 {
		u.Avatar = &req.Avatar
	}

	u.Gender = &gender
	var err error
	if req.Password != "" {
		err = s.um.CreateWithPassword(ctx, &u, req.Password, true)
	} else {
		err = s.um.Create(ctx, &u)
	}
	if err != nil {
		return nil, err
	}

	res := MapBizUserToApi(ctx, &u, s.blob)
	return res, nil
}

func (s *UserService) UpdateUserAdmin(ctx context.Context, req *pb.AdminUpdateUserRequest) (*pb.User, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceAdminUser, req.User.Id), authz.UpdateAction); err != nil {
		return nil, err
	}
	// check confirm password
	if req.User.Password != "" {
		if req.User.ConfirmPassword != req.User.Password {
			return nil, pb.ErrorConfirmPasswordMismatchLocalized(ctx, nil, nil)
		}
	}

	if req.UpdateMask != nil {
		fmutils.Filter(req, req.UpdateMask.Paths)
	}
	u, err := s.um.FindByID(ctx, req.User.Id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, pb.ErrorUserNotFoundLocalized(ctx, nil, nil)
	}

	if req.User.Password != "" {
		//reset password
		if err := s.um.UpdatePassword(ctx, u, req.User.Password); err != nil {
			return nil, err
		}
	}
	if req.GetUser().GetUsername() != nil {
		v := req.GetUser().GetUsername().Value
		u.Username = &v
	}

	if req.GetUser().GetName() != nil {
		v := req.GetUser().GetName().Value
		u.Name = &v
	}
	if req.GetUser().GetPhone() != nil {
		v := req.GetUser().GetPhone().Value
		u.SetPhone(v, false)
	}
	if req.GetUser().GetEmail() != nil {
		v := req.GetUser().GetEmail().Value
		u.SetEmail(v, false)
	}
	if req.GetUser().GetBirthday() != nil {
		v := req.GetUser().GetBirthday().AsTime()
		u.Birthday = &v
	}
	if len(req.User.Avatar) > 0 {
		u.Avatar = &req.User.Avatar
	}
	g := req.GetUser().Gender.Enum().String()
	u.Gender = &g
	if err := s.um.Update(ctx, u, nil); err != nil {
		return nil, err
	}
	res := MapBizUserToApi(ctx, u, s.blob)
	return res, nil
}

func (s *UserService) DeleteUserAdmin(ctx context.Context, req *pb.AdminDeleteUserRequest) (*pb.AdminDeleteUserResponse, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceUser, req.Id), authz.DeleteAction); err != nil {
		return nil, err
	}

	u, err := s.um.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, pb.ErrorUserNotFoundLocalized(ctx, nil, nil)
	}
	if err := s.um.Delete(ctx, u); err != nil {
		return nil, err
	}
	return &pb.AdminDeleteUserResponse{}, nil
}
