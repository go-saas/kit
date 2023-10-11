package service

import (
	"context"
	"github.com/eko/gocache/v3/cache"
	"github.com/go-saas/kit/pkg/idp"
	"github.com/go-saas/kit/user/private/biz"
	"github.com/nyaruka/phonenumbers"
	"github.com/samber/lo"
	"strconv"
	"strings"
	"time"

	pb "github.com/go-saas/kit/user/api/auth/v1"
)

const (
	miniprogramLogin biz.TokenPurpose = "wechat_miniprogram_login"
)

type WeChatAuthService struct {
	wechat            *idp.WechatFactory
	um                *biz.UserManager
	authSrv           *AuthService
	loginTwoStepToken *biz.TwoStepTokenProvider[*biz.WeChatMiniProgramLoginTwoStepTokenPayload]
}

var _ pb.WeChatAuthServiceServer = (*WeChatAuthService)(nil)

func NewWeChatAuthService(wechat *idp.WechatFactory, strCache cache.CacheInterface[string], um *biz.UserManager, authSrv *AuthService) *WeChatAuthService {
	return &WeChatAuthService{
		wechat:  wechat,
		um:      um,
		authSrv: authSrv,
		loginTwoStepToken: biz.NewTwoStepTokenProvider(func() *biz.WeChatMiniProgramLoginTwoStepTokenPayload {
			return &biz.WeChatMiniProgramLoginTwoStepTokenPayload{}
		}, strCache),
	}
}

func (s *WeChatAuthService) MiniProgramCode(ctx context.Context, req *pb.WechatMiniProgramCodeReq) (*pb.WeChatLoginReply, error) {
	mini, err := s.wechat.GetMiniProgramByAppID(ctx, req.AppId)
	if err != nil {
		return nil, err
	}
	sess, err := mini.GetAuth().Code2SessionContext(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	//find user by openId and
	user, err := s.um.FindByLogin(ctx, idp.GetWeChatMiniGameOpenIDProvider(req.AppId), sess.OpenID)
	if err != nil {
		return nil, err
	}
	if user == nil && sess.UnionID != "" {
		user, err = s.um.FindByLogin(ctx, idp.GetWeChatUnionIDProvider(), sess.UnionID)
		if err != nil {
			return nil, err
		}
	}
	if user != nil {
		if err != nil {
			return nil, err
		}
		t, err := s.authSrv.generateToken(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		return &pb.WeChatLoginReply{
			Token: &pb.TokenReply{
				AccessToken:  t.accessToken,
				ExpiresIn:    t.expiresIn,
				TokenType:    "Bearer",
				RefreshToken: t.refreshToken,
			}}, nil
	} else {
		nextStep := "phone"
		//user not registered
		token, err := s.loginTwoStepToken.Generate(ctx, miniprogramLogin, &biz.WeChatMiniProgramLoginTwoStepTokenPayload{
			AppId:      req.AppId,
			SessionKey: sess.SessionKey,
			OpenId:     sess.OpenID,
			UnionId:    sess.UnionID,
			Step:       nextStep,
		}, time.Hour*24)
		if err != nil {
			return nil, err
		}
		return &pb.WeChatLoginReply{
			NextToken: token,
			Next:      nextStep,
		}, nil
	}

}
func (s *WeChatAuthService) MiniProgramPhoneCode(ctx context.Context, req *pb.WechatMiniProgramPhoneCodeReq) (*pb.WeChatLoginReply, error) {
	mini, err := s.wechat.GetMiniProgramByAppID(ctx, req.AppId)
	if err != nil {
		return nil, err
	}
	var twoStepPayload *biz.WeChatMiniProgramLoginTwoStepTokenPayload
	if req.NextToken != "" {
		twoStepPayload, err = s.loginTwoStepToken.Retrieve(ctx, miniprogramLogin, req.NextToken)
		if err != nil {
			return nil, err
		}
		if twoStepPayload == nil {
			return nil, pb.ErrorTwoStepFailedLocalized(ctx, nil, nil)
		}
	}

	phoneResp, err := mini.GetAuth().GetPhoneNumberContext(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	countryCode, err := strconv.ParseInt(strings.TrimPrefix(phoneResp.PhoneInfo.CountryCode, "+"), 10, 0)
	if err != nil {
		return nil, err
	}
	pn, err := phonenumbers.Parse(phoneResp.PhoneInfo.PurePhoneNumber, phonenumbers.GetRegionCodeForCountryCode(int(countryCode)))
	if err != nil {
		return nil, err
	}
	formattedNum := phonenumbers.Format(pn, phonenumbers.E164)
	user, err := s.um.FindByPhone(ctx, formattedNum)
	if err != nil {
		return nil, err
	}
	if user == nil {
		//register
		user = &biz.User{}
		user.SetPhone(formattedNum, true)
		if err := s.um.Create(ctx, user); err != nil {
			return nil, err
		}
	}
	//add to login
	if twoStepPayload != nil {
		logins := lo.Filter([]biz.UserLogin{
			{LoginProvider: idp.GetWeChatMiniGameOpenIDProvider(req.AppId), ProviderKey: twoStepPayload.OpenId},
			{LoginProvider: idp.GetWeChatUnionIDProvider(), ProviderKey: twoStepPayload.UnionId},
		}, func(login biz.UserLogin, _ int) bool {
			return len(login.ProviderKey) > 0
		})
		if len(logins) > 0 {
			if err := s.um.AddLogin(ctx, user, logins); err != nil {
				return nil, err
			}
		}
	}

	t, err := s.authSrv.generateToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return &pb.WeChatLoginReply{
		Token: &pb.TokenReply{
			AccessToken:  t.accessToken,
			ExpiresIn:    t.expiresIn,
			TokenType:    "Bearer",
			RefreshToken: t.refreshToken,
		}}, nil
}
