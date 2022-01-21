package authorization

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas/common"
	"regexp"
	"strings"
	"sync"
)

type PermissionManagementService interface {
	AddGrant(ctx context.Context, resource Resource, action Action, subject Subject, tenantID string, effect Effect) error
	ListAcl(ctx context.Context, subjects ...Subject) ([]PermissionBean, error)
	UpdateGrant(ctx context.Context, subject Subject, acl []UpdateSubjectPermission) error
}

type PermissionChecker interface {
	IsGrant(ctx context.Context, resource Resource, action Action, subjects ...Subject) (Effect, error)
}

type PermissionService struct {
	v   []PermissionBean
	mux sync.Mutex
	log *log.Helper
}

var _ PermissionManagementService = (*PermissionService)(nil)
var _ PermissionChecker = (*PermissionService)(nil)

func NewPermissionService(logger log.Logger) *PermissionService {
	return &PermissionService{log: log.NewHelper(log.With(logger, "module", "authorization.permission"))}
}

func (p *PermissionService) IsGrant(ctx context.Context, resource Resource, action Action, subjects ...Subject) (Effect, error) {
	p.mux.Lock()
	defer p.mux.Unlock()
	var anyAllow bool
	//TODO host side?
	tenantInfo := common.FromCurrentTenant(ctx)
	for _, subject := range subjects {
		for _, bean := range p.v {
			if match(bean.Namespace, resource.GetNamespace()) &&
				match(bean.Subject, subject.GetIdentity()) &&
				match(bean.Resource, resource.GetIdentity()) &&
				match(bean.Action, action.GetIdentity()) &&
				match(bean.TenantID, tenantInfo.GetId()) {
				if bean.Effect == EffectForbidden {
					p.log.Debugf("Subject %s Action %s to Resource %s forbidden", subject.GetIdentity(), action.GetIdentity(), fmt.Sprintf("%s/%s", resource.GetNamespace(), resource.GetIdentity()))
					return EffectForbidden, nil
				}
				if bean.Effect == EffectGrant {
					anyAllow = true
				}
			}
		}
		if anyAllow {
			p.log.Debugf("Subject %s Action %s to Resource %s grant", subject.GetIdentity(), action.GetIdentity(), fmt.Sprintf("%s/%s", resource.GetNamespace(), resource.GetIdentity()))
			return EffectGrant, nil
		}
		p.log.Debugf("Subject %s Action %s to Resource %s unknown", subject.GetIdentity(), action.GetIdentity(), fmt.Sprintf("%s/%s", resource.GetNamespace(), resource.GetIdentity()))
	}
	return EffectUnknown, nil
}

func (p *PermissionService) ListAcl(ctx context.Context, subjects ...Subject) ([]PermissionBean, error) {
	p.mux.Lock()
	defer p.mux.Unlock()
	//TODO host side?
	tenantInfo := common.FromCurrentTenant(ctx)
	var ret []PermissionBean
	for _, bean := range p.v {
		for _, subject := range subjects {
			if (match(bean.Subject, subject.GetIdentity()) || match(bean.Subject, "")) && match(bean.TenantID, tenantInfo.GetId()) {
				ret = append(ret, bean)
			}
		}
	}
	return ret, nil
}

func (p *PermissionService) UpdateGrant(ctx context.Context, subject Subject, acl []UpdateSubjectPermission) error {
	p.mux.Lock()
	defer p.mux.Unlock()
	//remove previous
	for i := len(p.v) - 1; i >= 0; i-- {
		if subject.GetIdentity() == p.v[i].Subject {
			p.v = append(p.v[:i], p.v[i+1:]...)
		}
	}
	for _, permission := range acl {
		p.v = append(p.v, NewPermissionBean(permission.Resource, permission.Action, subject, permission.TenantID, permission.Effect))
	}
	return nil
}

func (p *PermissionService) AddGrant(ctx context.Context, resource Resource, action Action, subject Subject, tenantID string, effect Effect) error {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.v = append(p.v, NewPermissionBean(resource, action, subject, tenantID, effect))
	p.log.Debugf("add Resource %s Action %s grant %v to Subject %s", fmt.Sprintf("%s/%s", resource.GetNamespace(), resource.GetIdentity()), action.GetIdentity(), effect, subject.GetIdentity())
	return nil
}

// wildCardToRegexp converts a wildcard pattern to a regular expression pattern.
func wildCardToRegexp(pattern string) string {
	var result strings.Builder
	for i, literal := range strings.Split(pattern, "*") {

		// Replace * with .*
		if i > 0 {
			result.WriteString(".*")
		}

		// Quote any regular expression meta characters in the
		// literal text.
		result.WriteString(regexp.QuoteMeta(literal))
	}
	return result.String()
}

func match(pattern string, value string) bool {
	result, _ := regexp.MatchString(wildCardToRegexp(pattern), value)
	return result
}

var MemoryManagerProviderSet = wire.NewSet(
	NewPermissionService,
	wire.Bind(new(PermissionManagementService), new(*PermissionService)),
	wire.Bind(new(PermissionChecker), new(*PermissionService)),
)
