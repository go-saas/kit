package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/goxiaoy/go-saas-kit/pkg/authn"
	"github.com/goxiaoy/go-saas/common"
	shttp "github.com/goxiaoy/go-saas/common/http"
)

const (
	userKey       = "user"
	clientKey     = "client"
	tenantInfoKey = "tenant.info"
)

type SaasPropagator struct {
	hmtOpt *shttp.WebMultiTenancyOption
	l      *log.Helper
}

var _ Propagator = (*SaasPropagator)(nil)

func NewSaasContributor(hmtOpt *shttp.WebMultiTenancyOption, logger log.Logger) *SaasPropagator {
	return &SaasPropagator{
		hmtOpt: hmtOpt,
		l:      log.NewHelper(log.With(logger, "module", "SaasPropagator")),
	}
}

func (s *SaasPropagator) Extract(ctx context.Context, headers Header) (context.Context, error) {
	if !headers.HasKey(s.hmtOpt.TenantKey) {
		return ctx, nil
	}
	if headers.HasKey(s.hmtOpt.TenantKey) {
		tenantId := headers.Get(s.hmtOpt.TenantKey)
		s.l.Infof("recover tenant: %s", tenantId)
		ctx = common.NewCurrentTenant(ctx, tenantId, "")
	}
	if headers.HasKey(tenantInfoKey) {
		if infoJson, err := base64.StdEncoding.DecodeString(headers.Get(tenantInfoKey)); err == nil {
			tenantConfig := &common.TenantConfig{}
			if err := json.Unmarshal(infoJson, tenantConfig); err == nil {
				s.l.Infof("recover tenant config: %s", infoJson)
				ctx = common.NewTenantConfigContext(ctx, tenantConfig.ID, tenantConfig)
			}
		}
	}
	return ctx, nil
}

func (s *SaasPropagator) Inject(ctx context.Context, headers Header) error {
	ti, _ := common.FromCurrentTenant(ctx)
	headers.Set(s.hmtOpt.TenantKey, ti.GetId())
	if tenantConfig, ok := common.FromTenantConfigContext(ctx, ti.GetId()); ok {
		b, _ := json.Marshal(tenantConfig)
		headers.Set(tenantInfoKey, base64.StdEncoding.EncodeToString(b))
	}
	return nil
}

func (s *SaasPropagator) Fields() []string {
	return []string{
		s.hmtOpt.TenantKey,
		tenantInfoKey,
	}
}

type UserPropagator struct {
	l *log.Helper
}

var _ Propagator = (*UserPropagator)(nil)

func NewUserContributor(logger log.Logger) *UserPropagator {
	return &UserPropagator{l: log.NewHelper(log.With(logger, "module", "UserPropagator"))}
}

func (u *UserPropagator) Extract(ctx context.Context, headers Header) (context.Context, error) {
	if !headers.HasKey(userKey) {
		return ctx, nil
	}
	user := headers.Get(userKey)
	u.l.Infof("recover user: %s", user)
	return authn.NewUserContext(ctx, authn.NewUserInfo(user)), nil
}

func (u *UserPropagator) Inject(ctx context.Context, carrier Header) error {
	if userInfo, ok := authn.FromUserContext(ctx); ok {
		carrier.Set(userKey, userInfo.GetId())
	}
	return nil
}

func (u *UserPropagator) Fields() []string {
	return []string{
		userKey,
	}
}

type ClientPropagator struct {
	recoverOnly bool
	l           *log.Helper
}

var _ Propagator = (*ClientPropagator)(nil)

func NewClientContributor(recoverOnly bool, logger log.Logger) *ClientPropagator {
	return &ClientPropagator{
		recoverOnly: recoverOnly,
		l:           log.NewHelper(log.With(logger, "module", "ClientPropagator")),
	}
}

func (u *ClientPropagator) Extract(ctx context.Context, headers Header) (context.Context, error) {
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

func (u *ClientPropagator) Inject(ctx context.Context, carrier Header) error {
	if u.recoverOnly {
		return nil
	}
	if client, ok := authn.FromClientContext(ctx); ok {
		if len(client) == 0 {
			//can not set empty value header in gateway
			client = "-"
		}
		carrier.Set(clientKey, client)
	}
	return nil
}

func (u *ClientPropagator) Fields() []string {
	return []string{
		clientKey,
	}
}
