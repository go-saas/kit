package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	api2 "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authn"
	"github.com/go-saas/kit/pkg/authn/jwt"
	"github.com/go-saas/kit/pkg/authn/session"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/conf"
	kithttp "github.com/go-saas/kit/pkg/server/http"
	pb "github.com/go-saas/kit/user/api/auth/v1"
	v1 "github.com/go-saas/kit/user/api/user/v1"
	"github.com/go-saas/kit/user/private/biz"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
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
	logger           *klog.Helper
	emailer          biz.EmailSender
	auth             authz.Service
	trust            api2.TrustedContextValidator
	signIn           *biz.SignInManager
}

func NewAuthService(um *biz.UserManager,
	rm *biz.RoleManager,
	token jwt.Tokenizer,
	config *jwt.TokenizerConfig,
	pwdValidator biz.PasswordValidator,
	refreshTokenRepo biz.RefreshTokenRepo,
	emailer biz.EmailSender,
	security *conf.Security,
	signIn *biz.SignInManager,
	authz authz.Service,
	trust api2.TrustedContextValidator,
	logger klog.Logger) *AuthService {
	return &AuthService{
		um:               um,
		rm:               rm,
		token:            token,
		config:           config,
		pwdValidator:     pwdValidator,
		refreshTokenRepo: refreshTokenRepo,
		emailer:          emailer,
		security:         security,
		signIn:           signIn,
		auth:             authz,
		trust:            trust,
		logger:           klog.NewHelper(klog.With(logger, "module", "AuthService")),
	}
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterAuthRequest) (*pb.RegisterAuthReply, error) {
	// check confirm password
	if req.Password != "" {
		if req.ConfirmPassword != req.Password {
			return nil, pb.ErrorConfirmPasswordMismatchLocalized(ctx, nil, nil)
		}
	}
	user := &biz.User{}
	user.Username = &req.Username
	if err := s.um.CreateWithPassword(ctx, user, req.Password, true); err != nil {
		return nil, err
	}
	//login success
	if req.Web {
		if err := s.signIn.SignIn(ctx, user, true); err != nil {
			return nil, err
		}
		return &pb.RegisterAuthReply{}, nil
	}
	t, err := s.generateToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterAuthReply{
		AccessToken:  t.accessToken,
		ExpiresIn:    t.expiresIn,
		TokenType:    "Bearer",
		RefreshToken: t.refreshToken,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginAuthRequest) (*pb.LoginAuthReply, error) {
	user, err := FindUserByUsernameAndValidatePwd(ctx, s.um, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	//login success
	t, err := s.generateToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return &pb.LoginAuthReply{
		AccessToken:  t.accessToken,
		Expires:      t.expiresIn,
		ExpiresIn:    t.expiresIn,
		TokenType:    "Bearer",
		RefreshToken: t.refreshToken,
	}, nil
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
	if req.Phone == nil && req.Email == nil {
		return nil, errors.BadRequest("", "")
	}
	if req.Email != nil {
		token, err := s.um.GenerateEmailLoginPasswordlessToken(ctx, req.Email.Value)
		if err != nil {
			return nil, err
		}
		if err := s.emailer.SendPasswordlessLogin(ctx, req.Email.Value, token); err != nil {
			return nil, err
		}
	}
	if req.Phone != nil {
		token, err := s.um.GeneratePhoneLoginPasswordlessToken(ctx, req.Phone.Value)
		if err != nil {
			return nil, err
		}
		//TODO send token
		s.logger.Infof("send passwordless login token %s to  phone %s", token, req.Phone.Value)
	}
	return &pb.PasswordlessTokenAuthReply{}, nil
}

func (s *AuthService) LoginPasswordless(ctx context.Context, req *pb.LoginPasswordlessRequest) (*pb.LoginPasswordlessReply, error) {
	if req.Phone == nil && req.Email == nil {
		return nil, errors.BadRequest("", "")
	}
	var user *biz.User
	var err error
	if req.Email != nil {
		user, err = s.um.VerifyEmailLoginPasswordlessToken(ctx, req.Email.Value, req.Token)
	}
	if req.Phone != nil {
		user, err = s.um.VerifyPhoneLoginPasswordlessToken(ctx, req.Phone.Value, req.Token)
	}
	if err != nil {
		return nil, err
	}
	if req.Web {
		if err := s.signIn.SignIn(ctx, user, true); err != nil {
			return nil, err
		}
		return &pb.LoginPasswordlessReply{}, nil
	}
	//login success
	t, err := s.generateToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return &pb.LoginPasswordlessReply{
		AccessToken:  t.accessToken,
		ExpiresIn:    t.expiresIn,
		TokenType:    "Bearer",
		RefreshToken: t.refreshToken,
	}, nil

}

func (s *AuthService) SendForgetPasswordToken(ctx context.Context, req *pb.ForgetPasswordTokenRequest) (*pb.ForgetPasswordTokenReply, error) {

	if req.Phone == nil && req.Email == nil {
		return nil, errors.BadRequest("", "")
	}
	if req.Phone != nil {
		user, err := s.um.FindByPhone(ctx, req.Phone.Value)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, v1.ErrorUserNotFoundLocalized(ctx, nil, nil)
		}
		//generate token
		token, err := s.um.GeneratePhoneForgetPasswordToken(ctx, user)
		if err != nil {
			return nil, err
		}
		//TODO send token
		s.logger.Infof("send forget password token %s to %s for user %s", token, *user.Phone, user.ID.String())

	} else if req.Email != nil {
		user, err := s.um.FindByEmail(ctx, req.Email.Value)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, v1.ErrorUserNotFoundLocalized(ctx, nil, nil)
		}
		//generate token
		token, err := s.um.GenerateEmailForgetPasswordToken(ctx, user)
		if err != nil {
			return nil, err
		}
		err = s.emailer.SendForgetPassword(ctx, *user.Email, token)
		if err != nil {
			return nil, err
		}
	}
	return &pb.ForgetPasswordTokenReply{}, nil
}

func (s *AuthService) ForgetPassword(ctx context.Context, req *pb.ForgetPasswordRequest) (*pb.ForgetPasswordReply, error) {

	if req.Phone == nil && req.Email == nil {
		return nil, errors.BadRequest("", "")
	}
	var user *biz.User
	var err error
	if req.Phone != nil {
		user, err = s.um.FindByPhone(ctx, req.Phone.Value)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, v1.ErrorUserNotFoundLocalized(ctx, nil, nil)
		}
		if err := s.um.VerifyPhoneForgetPasswordToken(ctx, req.Phone.Value, req.Token); err != nil {
			return nil, err
		}

	} else if req.Email != nil {
		user, err = s.um.FindByPhone(ctx, req.Phone.Value)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, v1.ErrorUserNotFoundLocalized(ctx, nil, nil)
		}
		if err := s.um.VerifyEmailForgetPasswordToken(ctx, req.Email.Value, req.Token); err != nil {
			return nil, err
		}
	}
	token, err := s.um.GenerateForgetPasswordTwoStepToken(ctx, user)
	if err != nil {
		return nil, err
	}
	return &pb.ForgetPasswordReply{ChangePasswordToken: token}, nil
}

func (s *AuthService) ChangePasswordByForget(ctx context.Context, req *pb.ChangePasswordByForgetRequest) (*pb.ChangePasswordByForgetReply, error) {

	//validate password
	if req.NewPassword != req.ConfirmNewPassword {
		return nil, v1.ErrorConfirmPasswordMismatchLocalized(ctx, nil, nil)
	}
	err := s.um.ChangePasswordByToken(ctx, req.ChangePasswordToken, req.NewPassword)
	if err != nil {
		return nil, err
	}
	return &pb.ChangePasswordByForgetReply{}, nil
}

func (s *AuthService) ValidatePassword(ctx context.Context, req *pb.ValidatePasswordRequest) (*pb.ValidatePasswordReply, error) {

	err := s.pwdValidator.Validate(ctx, req.Password)
	if err != nil {
		return nil, err
	}
	return &pb.ValidatePasswordReply{Ok: true}, nil
}

func (s *AuthService) ChangePasswordByPre(ctx context.Context, req *pb.ChangePasswordByPreRequest) (*pb.ChangePasswordByPreReply, error) {

	ui, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	//validate password
	if req.NewPassword != req.ConfirmNewPassword {
		return nil, v1.ErrorConfirmPasswordMismatchLocalized(ctx, nil, nil)
	}
	user, err := s.um.FindByID(ctx, ui.GetId())
	if err != nil {
		return nil, err
	}
	err = s.um.ChangePassword(ctx, user, req.PrePassword, req.NewPassword)
	if err != nil {
		return nil, err
	}
	return &pb.ChangePasswordByPreReply{}, nil
}

func (s *AuthService) GetCsrfToken(ctx context.Context, req *pb.GetCsrfTokenRequest) (*pb.GetCsrfTokenResponse, error) {
	if t, ok := transport.FromServerContext(ctx); ok {
		if ht, ok := t.(*http.Transport); ok {
			token := csrf.Token(ht.Request())
			t.ReplyHeader().Set("X-CSRF-Token", token)
			return &pb.GetCsrfTokenResponse{CsrfToken: token}, nil
		}
	}
	return nil, pb.ErrorInvalidOperationLocalized(ctx, nil, nil)
}

func (s *AuthService) RefreshRememberToken(ctx context.Context, req *pb.RefreshRememberTokenRequest) (*pb.RefreshRememberTokenReply, error) {
	if ok, err := s.trust.Trusted(ctx); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.Forbidden("", "")
	}
	duration := 0
	if s.security != nil && s.security.RememberCookie != nil && s.security.RememberCookie.MaxAge != nil {
		duration = int(s.security.RememberCookie.MaxAge.Value)
	}
	duration = session.RememberMeExpireSecondsOrDefault(duration)
	user, newToken, err := s.um.RefreshRememberToken(ctx, req.RmToken, time.Duration(duration)*time.Second)
	if err != nil {
		return nil, err
	}
	return &pb.RefreshRememberTokenReply{UserId: user.ID.String(), NewRmToken: newToken}, nil
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

	refreshToken := biz.NewRefreshToken(userId, duration, kithttp.ClientUserAgent(ctx), kithttp.ClientIP(ctx))
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
		return nil, pb.ErrorRefreshTokenInvalidLocalized(ctx, nil, nil)
	}
	if token.Valid() {
		//token valid, regenerate

		//find user again
		user, err := s.um.FindByID(ctx, token.UserId.String())
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, pb.ErrorRefreshTokenInvalidLocalized(ctx, nil, nil)
		}

		t, err := s.generateToken(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		//delete pre
		if err := s.refreshTokenRepo.Revoke(ctx, refreshToken, true); err != nil {
			return nil, err
		}
		return &tokenModel{accessToken: t.accessToken, expiresIn: t.expiresIn, refreshToken: t.refreshToken}, nil
	}
	return nil, pb.ErrorRefreshTokenInvalidLocalized(ctx, nil, nil)
}

func FindUserByUsernameAndValidatePwd(ctx context.Context, um *biz.UserManager, username, password string) (*biz.User, error) {
	user, err := um.FindByIdentity(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, pb.ErrorInvalidCredentialsLocalized(ctx, nil, nil)
	}
	// check password
	err = um.CheckPassword(ctx, user, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
