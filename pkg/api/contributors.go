package api

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/goxiaoy/go-saas-kit/pkg/authn"
	"github.com/goxiaoy/go-saas/common"
	shttp "github.com/goxiaoy/go-saas/common/http"
)

type SaasContributor struct {
	hmtOpt *shttp.WebMultiTenancyOption
	l      *log.Helper
}

var _ Contributor = (*SaasContributor)(nil)

func NewSaasContributor(hmtOpt *shttp.WebMultiTenancyOption, logger log.Logger) *SaasContributor {
	return &SaasContributor{
		hmtOpt: hmtOpt,
		l:      log.NewHelper(log.With(logger, "module", "SaasContributor")),
	}
}

func (s *SaasContributor) RecoverContext(ctx context.Context, headers Header) (context.Context, error) {
	if !headers.HasKey(s.hmtOpt.TenantKey) {
		return ctx, nil
	}
	tenantId := headers.Get(s.hmtOpt.TenantKey)
	s.l.Infof("recover tenant: %s", tenantId)
	return common.NewCurrentTenant(ctx, tenantId, ""), nil
}

func (s *SaasContributor) CreateHeader(ctx context.Context) map[string]string {
	ti, _ := common.FromCurrentTenant(ctx)
	return map[string]string{
		s.hmtOpt.TenantKey: ti.GetId(),
	}
}

const (
	userKey   = "user"
	clientKey = "client"
)

type UserContributor struct {
	l *log.Helper
}

var _ Contributor = (*UserContributor)(nil)

func NewUserContributor(logger log.Logger) *UserContributor {
	return &UserContributor{l: log.NewHelper(log.With(logger, "module", "UserContributor"))}
}

func (u *UserContributor) RecoverContext(ctx context.Context, headers Header) (context.Context, error) {
	if !headers.HasKey(userKey) {
		return ctx, nil
	}
	user := headers.Get(userKey)
	u.l.Infof("recover user: %s", user)
	return authn.NewUserContext(ctx, authn.NewUserInfo(user)), nil
}

func (u *UserContributor) CreateHeader(ctx context.Context) map[string]string {
	res := map[string]string{}
	if userInfo, ok := authn.FromUserContext(ctx); ok {
		res[userKey] = userInfo.GetId()
	}
	return res
}

type ClientContributor struct {
	recoverOnly bool
	l           *log.Helper
}

var _ Contributor = (*ClientContributor)(nil)

func NewClientContributor(recoverOnly bool, logger log.Logger) *ClientContributor {
	return &ClientContributor{
		recoverOnly: recoverOnly,
		l:           log.NewHelper(log.With(logger, "module", "ClientContributor")),
	}
}

func (u *ClientContributor) RecoverContext(ctx context.Context, headers Header) (context.Context, error) {
	if !headers.HasKey(clientKey) {
		return ctx, nil
	}
	client := headers.Get(clientKey)
	if client == "-" {
		//can not set empty value header in gateway
		client = ""
	}
	u.l.Infof("recover client: %s", client)
	return authn.NewClientContext(ctx, client), nil
}

func (u *ClientContributor) CreateHeader(ctx context.Context) map[string]string {
	if u.recoverOnly {
		return nil
	}
	res := map[string]string{}
	if client, ok := authn.FromClientContext(ctx); ok {
		if len(client) == 0 {
			//can not set empty value header in gateway
			client = "-"
		}
		res[clientKey] = client
	}
	return res
}
