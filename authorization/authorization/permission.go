package authorization

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"regexp"
	"strings"
	"sync"
)

type PermissionManagementService interface {
	AddGrant(ctx context.Context, resource Resource, action Action, subject Subject, effect Effect) error
}

type PermissionChecker interface {
	IsGrant(ctx context.Context, resource Resource, action Action, subject Subject) (Effect, error)
}

type permissionBean struct {
	namespace string
	resource  string
	action    string
	subject   string
	effect    Effect
}

func newPermissionBean(resource Resource, action Action, subject Subject, effect Effect) permissionBean {
	return permissionBean{
		namespace: resource.GetNamespace(),
		resource:  resource.GetIdentity(),
		action:    action.GetIdentity(),
		subject:   subject.GetIdentity(),
		effect:    effect,
	}
}

type PermissionService struct {
	v   []permissionBean
	mux sync.Mutex
	log *log.Helper
}

var _ PermissionManagementService = (*PermissionService)(nil)
var _ PermissionChecker = (*PermissionService)(nil)

func NewPermissionService(logger log.Logger) *PermissionService {
	return &PermissionService{log: log.NewHelper(logger)}
}

func (p *PermissionService) IsGrant(ctx context.Context, resource Resource, action Action, subject Subject) (Effect, error) {
	p.mux.Lock()
	defer p.mux.Unlock()
	var anyAllow bool
	for _, bean := range p.v {
		//TODO regex match
		if match(bean.namespace, resource.GetNamespace()) && match(bean.subject, subject.GetIdentity()) && match(bean.resource, resource.GetIdentity()) && match(bean.action, action.GetIdentity()) {
			if bean.effect == EffectForbidden {
				p.log.Debugf("subject %s action %s to resource %s forbidden", subject.GetIdentity(), action.GetIdentity(), resource.GetIdentity())
				return EffectForbidden, nil
			}
			if bean.effect == EffectGrant {
				anyAllow = true
			}
		}
	}
	if anyAllow {
		p.log.Debugf("subject %s action %s to resource %s grant", subject.GetIdentity(), action.GetIdentity(), resource.GetIdentity())
		return EffectGrant, nil
	}
	p.log.Debugf("subject %s action %s to resource %s unknown", subject.GetIdentity(), action.GetIdentity(), resource.GetIdentity())
	return EffectUnknown, nil
}

func (p *PermissionService) AddGrant(ctx context.Context, resource Resource, action Action, subject Subject, effect Effect) error {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.v = append(p.v, newPermissionBean(resource, action, subject, effect))
	p.log.Debugf("add resource %s action %s grant %v to subject %s", resource.GetIdentity(), action.GetIdentity(), effect, subject.GetIdentity())
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
