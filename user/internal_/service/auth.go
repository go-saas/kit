package service

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/auth/jwt"
	pb "github.com/goxiaoy/go-saas-kit/user/api/auth/v1"
	"github.com/goxiaoy/go-saas-kit/user/internal_/biz"
)

type AuthService struct {
	pb.UnimplementedAuthServer

	um     *biz.UserManager
	rm     *biz.RoleManager
	token  jwt.Tokenizer
	config *jwt.TokenizerConfig
	pwdValidator biz.PasswordValidator
}

func NewAuthService(um *biz.UserManager, rm *biz.RoleManager, token jwt.Tokenizer, config *jwt.TokenizerConfig,pwdValidator biz.PasswordValidator) *AuthService {
	return &AuthService{um: um, rm: rm, token: token, config: config,pwdValidator: pwdValidator}
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterAuthRequest) (*pb.RegisterAuthReply, error) {
	return &pb.RegisterAuthReply{}, nil
}
func (s *AuthService) Login(ctx context.Context, req *pb.LoginAuthRequest) (*pb.LoginAuthReply, error) {
	user, err := s.um.FindByName(ctx, req.GetUsername())
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, pb.ErrorInvalidCredentials("")
	}
	// check password
	ok, err := s.um.CheckPassword(ctx, user, req.Password)
	if !ok {
		return nil, pb.ErrorInvalidCredentials("")
	}
	if err != nil {
		return nil, err
	}
	//login success
	t, err := s.token.Issue(user.ID.String())
	if err != nil {
		return nil, err
	}
	return &pb.LoginAuthReply{AccessToken: t, Expires: int32(s.config.ExpireDuration.Seconds()), TokenType: "Bearer"}, nil
}
func (s *AuthService) Token(ctx context.Context, req *pb.LoginAuthRequest) (*pb.LoginAuthReply, error) {
	return &pb.LoginAuthReply{}, nil
}
func (s *AuthService) Refresh(ctx context.Context, req *pb.RefreshTokenAuthRequest) (*pb.RefreshTokenAuthReply, error) {
	return &pb.RefreshTokenAuthReply{}, nil
}
func (s *AuthService) SendPasswordlessToken(ctx context.Context, req *pb.PasswordlessTokenAuthRequest) (*pb.PasswordlessTokenAuthReply, error) {
	return &pb.PasswordlessTokenAuthReply{}, nil
}
func (s *AuthService) LoginPasswordless(ctx context.Context, req *pb.LoginPasswordlessRequest) (*pb.LoginPasswordlessReply, error) {
	return &pb.LoginPasswordlessReply{}, nil
}
func (s *AuthService) SendForgetPasswordToken(ctx context.Context, req *pb.ForgetPasswordTokenRequest) (*pb.ForgetPasswordTokenReply, error) {
	return &pb.ForgetPasswordTokenReply{}, nil
}

func (s *AuthService) ForgetPassword(ctx context.Context, req *pb.ForgetPasswordRequest) (*pb.ForgetPasswordReply, error) {
	return &pb.ForgetPasswordReply{}, nil
}

func (s *AuthService) ValidatePassword(ctx context.Context, req *pb.ValidatePasswordRequest) (*pb.ValidatePasswordReply, error) {
	err:=s.pwdValidator.Validate(ctx,req.Password)
	if err!=nil{
		return nil,err
	}
	return &pb.ValidatePasswordReply{Ok: true},nil
}
