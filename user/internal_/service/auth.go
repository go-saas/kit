package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/uuid"
	"github.com/goxiaoy/go-saas-kit/auth/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	pb "github.com/goxiaoy/go-saas-kit/user/api/auth/v1"
	"github.com/goxiaoy/go-saas-kit/user/internal_/biz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"time"
)

type AuthService struct {
	pb.UnimplementedAuthServer

	um               *biz.UserManager
	rm               *biz.RoleManager
	token            jwt.Tokenizer
	config           *jwt.TokenizerConfig
	pwdValidator     biz.PasswordValidator
	refreshTokenRepo biz.RefreshTokenRepo
	security         *conf.Security
}

func NewAuthService(um *biz.UserManager,
	rm *biz.RoleManager,
	token jwt.Tokenizer,
	config *jwt.TokenizerConfig,
	pwdValidator biz.PasswordValidator,
	refreshTokenRepo biz.RefreshTokenRepo,
	security *conf.Security) *AuthService {
	return &AuthService{um: um, rm: rm, token: token, config: config, pwdValidator: pwdValidator, refreshTokenRepo: refreshTokenRepo, security: security}
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
	t, err := s.generateToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return &pb.LoginAuthReply{AccessToken: t.accessToken, Expires: t.expiresIn, ExpiresIn: t.expiresIn, TokenType: "Bearer", RefreshToken: t.refreshToken}, nil
}

func (s *AuthService) Token(ctx context.Context, req *pb.TokenRequest) (*pb.TokenReply, error) {
	if req.GrantType == "" || req.GrantType == "password" {
		res, err := s.Login(ctx, &pb.LoginAuthRequest{Username: req.Username, Password: req.Password})
		if err != nil {
			return nil, err
		}
		return &pb.TokenReply{
			AccessToken:  res.AccessToken,
			TokenType:    res.TokenType,
			RefreshToken: res.RefreshToken,
			ExpiresIn:    res.ExpiresIn,
		}, nil
	}
	if req.GrantType == "refresh_token" {
		//refresh
		//find token
		if req.RefreshToken == "" {
			return nil, errors.BadRequest("", "refreshToken can not be empty")
		}
		t, err := s.refresh(ctx, req.RefreshToken)
		if err != nil {
			return nil, err
		}
		return &pb.TokenReply{AccessToken: t.accessToken, ExpiresIn: t.expiresIn, TokenType: "Bearer", RefreshToken: t.refreshToken}, nil
	}
	return nil, status.Errorf(codes.Unimplemented, "not implemented")

}

func (s *AuthService) Refresh(ctx context.Context, req *pb.RefreshTokenAuthRequest) (*pb.RefreshTokenAuthReply, error) {
	t, err := s.refresh(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &pb.RefreshTokenAuthReply{AccessToken: t.accessToken, ExpiresIn: t.expiresIn, TokenType: "Bearer", RefreshToken: t.refreshToken}, nil
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
	err := s.pwdValidator.Validate(ctx, req.Password)
	if err != nil {
		return nil, err
	}
	return &pb.ValidatePasswordReply{Ok: true}, nil
}

type tokenModel struct {
	accessToken  string
	refreshToken string
	expiresIn    int32
}

func (s *AuthService) generateToken(ctx context.Context, userId uuid.UUID) (*tokenModel, error) {
	var duration time.Duration = 0
	if s.security.Jwt.RefreshTokenExpireIn != nil {
		duration = s.security.Jwt.RefreshTokenExpireIn.AsDuration()
	}

	refreshToken := biz.NewRefreshToken(userId, duration, server.ClientUserAgent(ctx), server.ClientIP(ctx))
	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, err
	}
	//access token
	t, err := s.token.Issue(jwt.NewUserClaim(userId.String()), 0)
	if err != nil {
		return nil, err
	}
	exp := int32(s.config.ExpireDuration.Seconds())
	return &tokenModel{
		accessToken:  t,
		refreshToken: refreshToken.Token,
		expiresIn:    exp,
	}, nil
}

func (s *AuthService) refresh(ctx context.Context, refreshToken string) (*tokenModel, error) {
	//find
	token, err := s.refreshTokenRepo.Find(ctx, refreshToken, true)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, errors.BadRequest("", "refreshToken invalid")
	}
	if token.Valid() {
		//token valid, regenerate

		//find user again
		user, err := s.um.FindByID(ctx, token.UserId.String())
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, errors.BadRequest("", "refreshToken invalid")
		}

		t, err := s.generateToken(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		//delete pre
		if err := s.refreshTokenRepo.Revoke(ctx, refreshToken); err != nil {
			return nil, err
		}
		return &tokenModel{accessToken: t.accessToken, expiresIn: t.expiresIn, refreshToken: t.refreshToken}, nil
	}
	return nil, errors.BadRequest("", "refreshToken invalid")
}
