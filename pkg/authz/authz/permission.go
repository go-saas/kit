package authz

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
	RemoveGrant(ctx context.Context, resource Resource, action Action, subject Subject, tenantID string, effects []Effect) error
}

func EnsureGrant(ctx context.Context, mgr PermissionManagementService, checker PermissionChecker, resource Resource, action Action, subject Subject, tenantID string) error {
	eff, err := checker.IsGrantTenant(ctx, resource, action, tenantID, subject)
	if err != nil {
		return err
	}
	if eff != EffectGrant {
		err = mgr.AddGrant(ctx, resource, action, subject, tenantID, EffectGrant)
		if err != nil {
			return err
		}
	}
	return nil
}

func EnsureForbidden(ctx context.Context, mgr PermissionManagementService, checker PermissionChecker, resource Resource, action Action, subject Subject, tenantID string) error {
	eff, err := checker.IsGrantTenant(ctx, resource, action, tenantID, subject)
	if err != nil {
		return err
	}
	if eff != EffectForbidden {
		err = mgr.AddGrant(ctx, resource, action, subject, tenantID, EffectForbidden)
		if err != nil {
			return err
		}
	}
	return nil
}

type PermissionService struct {
	v   []PermissionBean
	mux sync.Mutex
	log *log.Helper
}

var _ PermissionManagementService = (*PermissionService)(nil)

func NewPermissionService(logger log.Logger) *PermissionService {
	return &PermissionService{log: log.NewHelper(log.With(logger, "module", "authz.permission"))}
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
		p.v = append(p.v, NewPermissionBean(permission.Resource, permission.Action, subject, NormalizeTenantId(ctx, permission.TenantID), permission.Effect))
	}
	return nil
}

func (p *PermissionService) RemoveGrant(ctx context.Context, resource Resource, action Action, subject Subject, tenantID string, effects []Effect) error {
	p.mux.Lock()
	defer p.mux.Unlock()
	tenantID = NormalizeTenantId(ctx, tenantID)
	var v []PermissionBean
	if len(effects) == 0 {
		effects = []Effect{EffectGrant, EffectForbidden, EffectUnknown}
	}
	for _, bean := range p.v {
		preserved := true
		if (bean.Namespace == resource.GetNamespace()) &&
			(bean.Subject == subject.GetIdentity()) &&
			(bean.Resource == resource.GetIdentity()) &&
			(bean.Action == action.GetIdentity()) &&
			(bean.TenantID == tenantID) {
			for _, e := range effects {
				if bean.Effect == e {
					preserved = false
				}
			}
		}
		if preserved {
			v = append(v, bean)
		}
	}
	p.v = v
	return nil
}

func (p *PermissionService) AddGrant(ctx context.Context, resource Resource, action Action, subject Subject, tenantID string, effect Effect) error {
	p.mux.Lock()
	defer p.mux.Unlock()
	tenantID = NormalizeTenantId(ctx, tenantID)
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

func NormalizeTenantId(ctx context.Context, tenantId string) string {
	ti := common.FromCurrentTenant(ctx)
	if ti.GetId() == "" {
		//host side
		return tenantId
	}
	return ti.GetId()
}

var PermissionProviderSet = wire.NewSet(
	NewPermissionService,
	wire.Bind(new(PermissionManagementService), new(*PermissionService)),
	wire.Bind(new(PermissionChecker), new(*PermissionService)),
)
