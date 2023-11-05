package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-saas/kit/pkg/authn"
	"github.com/go-saas/saas"
	"github.com/go-saas/saas/data"
	"strconv"
)

const (
	InternalKeyPrefix = "internal-"
	TenantKey         = InternalKeyPrefix + "tenant"
	TenantInfoKey     = InternalKeyPrefix + "tenant-info"
	TenantFilterKey   = InternalKeyPrefix + "tenant-data-filter"
	TenantAutoSetKey  = InternalKeyPrefix + "tenant-data-autoset"
	UserKey           = InternalKeyPrefix + "user"
	ClientKey         = InternalKeyPrefix + "client"
)

type SaasPropagator struct {
	l *log.Helper
}

var _ Propagator = (*SaasPropagator)(nil)

func NewSaasPropagator(logger log.Logger) *SaasPropagator {
	return &SaasPropagator{
		l: log.NewHelper(log.With(logger, "module", "SaasPropagator")),
	}
}

func (s *SaasPropagator) Extract(ctx context.Context, headers Header) (context.Context, error) {
	if headers.HasKey(TenantKey) {
		tenantId := headers.Get(TenantKey)
		s.l.Debugf("recover tenant: %s", tenantId)
		ctx = saas.NewCurrentTenant(ctx, tenantId, "")
	}
	if headers.HasKey(TenantInfoKey) {
		if infoJson, err := base64.StdEncoding.DecodeString(headers.Get(TenantInfoKey)); err == nil {
			tenantConfig := &saas.TenantConfig{}
			if err := json.Unmarshal(infoJson, tenantConfig); err == nil {
				s.l.Debugf("recover tenant config: %s", infoJson)
				ctx = saas.NewTenantConfigContext(ctx, tenantConfig.ID, tenantConfig)
			}
		}
	}
	if headers.HasKey(TenantFilterKey) {
		if v, err := strconv.ParseBool(headers.Get(TenantFilterKey)); err == nil {
			s.l.Debugf("recover tenant data filter: %s", v)
			ctx = data.NewMultiTenancyDataFilter(ctx, v)
		}
	}
	if headers.HasKey(TenantAutoSetKey) {
		if v, err := strconv.ParseBool(headers.Get(TenantAutoSetKey)); err == nil {
			s.l.Debugf("recover tenant auto set: %s", v)
			ctx = data.NewAutoSetTenantId(ctx, v)
		}
	}
	return ctx, nil
}

func (s *SaasPropagator) Inject(ctx context.Context, headers Header) error {
	ti, _ := saas.FromCurrentTenant(ctx)
	headers.Set(TenantKey, ti.GetId())
	if tenantConfig, ok := saas.FromTenantConfigContext(ctx, ti.GetId()); ok {
		b, _ := json.Marshal(tenantConfig)
		headers.Set(TenantInfoKey, base64.StdEncoding.EncodeToString(b))
	}
	filter := data.FromMultiTenancyDataFilter(ctx)
	headers.Set(TenantFilterKey, strconv.FormatBool(filter))
	autoset := data.FromAutoSetTenantId(ctx)
	headers.Set(TenantAutoSetKey, strconv.FormatBool(autoset))
	for _, k := range headers.Keys() {
		s.l.Debugf("inject header key:%s value: %s", k, headers.Get(k))
	}
	return nil
}

func (s *SaasPropagator) Fields() []string {
	return []string{
		TenantKey,
		TenantInfoKey,
		TenantFilterKey,
		TenantAutoSetKey,
	}
}

type UserPropagator struct {
	l *log.Helper
}

var _ Propagator = (*UserPropagator)(nil)

func NewUserPropagator(logger log.Logger) *UserPropagator {
	return &UserPropagator{l: log.NewHelper(log.With(logger, "module", "UserPropagator"))}
}

func (u *UserPropagator) Extract(ctx context.Context, headers Header) (context.Context, error) {
	if !headers.HasKey(UserKey) {
		return ctx, nil
	}
	user := headers.Get(UserKey)
	u.l.Debugf("recover user: %s", user)
	return authn.NewUserContext(ctx, authn.NewUserInfo(user)), nil
}

func (u *UserPropagator) Inject(ctx context.Context, carrier Header) error {
	if userInfo, ok := authn.FromUserContext(ctx); ok {
		carrier.Set(UserKey, userInfo.GetId())
	}
	return nil
}

func (u *UserPropagator) Fields() []string {
	return []string{
		UserKey,
	}
}

type ClientPropagator struct {
	extractOnly bool
	l           *log.Helper
}

var _ Propagator = (*ClientPropagator)(nil)

func NewClientPropagator(extractOnly bool, logger log.Logger) *ClientPropagator {
	return &ClientPropagator{
		extractOnly: extractOnly,
		l:           log.NewHelper(log.With(logger, "module", "ClientPropagator")),
	}
}

func (u *ClientPropagator) Extract(ctx context.Context, headers Header) (context.Context, error) {
	if !headers.HasKey(ClientKey) {
		return ctx, nil
	}
	client := headers.Get(ClientKey)
	if client == "-" {
		//can not set empty value header in gateway
		client = ""
	}
	u.l.Debugf("recover client: %s", client)
	return authn.NewClientContext(ctx, client), nil
}

func (u *ClientPropagator) Inject(ctx context.Context, carrier Header) error {
	if u.extractOnly {
		return nil
	}
	if client, ok := authn.FromClientContext(ctx); ok {
		if len(client) == 0 {
			//can not set empty value header in gateway
			client = "-"
		}
		carrier.Set(ClientKey, client)
	}
	return nil
}

func (u *ClientPropagator) Fields() []string {
	return []string{
		ClientKey,
	}
}
