package authz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/go-saas/saas"
	"github.com/goava/di"
	"github.com/samber/lo"
	"strings"
)

type Service interface {
	//CheckForSubjects permission of these subjects directly
	CheckForSubjects(ctx context.Context, resource Resource, action Action, subjects ...Subject) (*Result, error)
	//Check resolve subject from ctx, then check permission of these subjects
	Check(ctx context.Context, resource Resource, action Action) (*Result, error)

	BatchCheckForSubjects(ctx context.Context, requirement RequirementList, subjects ...Subject) (ResultList, error)
	BatchCheck(ctx context.Context, requirement RequirementList) (ResultList, error)

	FormatError(ctx context.Context, requirements RequirementList, result *Result, subjects ...Subject) error
}

type Requirement struct {
	Resource Resource
	Action   Action
}

func NewRequirement(resource Resource, action Action) *Requirement {
	return &Requirement{
		Resource: resource,
		Action:   action,
	}
}

type RequirementList []*Requirement

type SubjectList []Subject

type ResultList []*Result

// SubjectContrib receive one Subject and retrieve as list of subjects
type SubjectContrib interface {
	Process(ctx context.Context, subject Subject) ([]Subject, error)
}

type Option struct {
	SubjectContribList []SubjectContrib
}

func NewAuthorizationOption(subjectContribList ...SubjectContrib) *Option {
	return &Option{SubjectContribList: subjectContribList}
}

type DefaultAuthorizationService struct {
	checker PermissionChecker
	sr      SubjectResolver
	log     *log.Helper
}

var _ Service = (*DefaultAuthorizationService)(nil)

func NewDefaultAuthorizationService(checker PermissionChecker, sr SubjectResolver, logger log.Logger) *DefaultAuthorizationService {
	return &DefaultAuthorizationService{checker: checker, sr: sr, log: log.NewHelper(log.With(logger, "module", "authz.service"))}
}

func (a *DefaultAuthorizationService) CheckForSubjects(ctx context.Context, resource Resource, action Action, subjects ...Subject) (*Result, error) {
	requirements := []*Requirement{NewRequirement(resource, action)}
	resList, err := a.BatchCheckForSubjects(ctx, requirements, subjects...)
	res := NewDisallowAuthorizationResult(requirements...)
	if len(resList) > 0 {
		res = resList[0]
	}
	if err != nil {
		return res, err
	}
	return res, a.FormatError(ctx, requirements, res)
}

func (a *DefaultAuthorizationService) Check(ctx context.Context, resource Resource, action Action) (*Result, error) {
	requirements := []*Requirement{NewRequirement(resource, action)}
	resList, err := a.BatchCheck(ctx, requirements)
	res := NewDisallowAuthorizationResult(requirements...)
	if len(resList) > 0 {
		res = resList[0]
	}
	if err != nil {
		return res, err
	}
	subjects, _ := a.sr.ResolveFromContext(ctx)
	return res, a.FormatError(ctx, requirements, res, subjects...)
}

func (a *DefaultAuthorizationService) check(ctx context.Context, requirements RequirementList, tenant string, subject ...Subject) (ResultList, error) {

	if always, ok := FromAlwaysAuthorizationContext(ctx); ok {
		var subjectStr []string
		for _, s := range subject {
			subjectStr = append(subjectStr, s.GetIdentity())
		}
		if always {
			return lo.Map(requirements, func(t *Requirement, _ int) *Result {
				return NewAllowAuthorizationResult()
			}), nil
		} else {
			return lo.Map(requirements, func(t *Requirement, _ int) *Result {
				return NewDisallowAuthorizationResult(requirements...)
			}), nil
		}
	}

	subjectList, err := a.sr.ResolveProcessed(ctx, subject...)
	if err != nil {
		return nil, err
	}

	grantType, err := a.checker.IsGrantTenant(ctx, requirements, tenant, subjectList...)
	if err != nil {
		return nil, err
	}
	res := lo.Map(grantType, func(effect Effect, i int) *Result {
		if effect == EffectForbidden {
			return NewDisallowAuthorizationResult(requirements[i])
		} else {
			return NewAllowAuthorizationResult()
		}
	})
	return res, nil
}

func (a *DefaultAuthorizationService) BatchCheckForSubjects(ctx context.Context, requirement RequirementList, subjects ...Subject) (ResultList, error) {
	ti, _ := saas.FromCurrentTenant(ctx)
	return a.check(ctx, requirement, ti.GetId(), subjects...)
}

func (a *DefaultAuthorizationService) BatchCheck(ctx context.Context, requirement RequirementList) (ResultList, error) {
	subjects, err := a.sr.ResolveFromContext(ctx)
	if err != nil {
		return nil, err
	}
	ti, _ := saas.FromCurrentTenant(ctx)
	return a.check(ctx, requirement, ti.GetId(), subjects...)
}

func (a *DefaultAuthorizationService) FormatError(ctx context.Context, requirements RequirementList, result *Result, subjects ...Subject) (err error) {
	if len(subjects) == 0 {
		subjects, err = a.sr.ResolveFromContext(ctx)
		if err != nil {
			return
		}
	}
	if result.Allowed {
		return nil
	}
	var authed bool
	for _, sub := range subjects {
		if s, ok := ParseUserSubject(sub); ok {
			if len(s.GetUserId()) > 0 {
				authed = true
			}
		}
		if s, ok := ParseClientSubject(sub); ok {
			if len(s.GetClientId()) > 0 {
				authed = true
			}
		}
	}
	if authed {
		a.log.Warnf("subjects: %s are not allowed to %s", strings.Join(lo.Map(subjects, func(t Subject, _ int) string {
			return t.GetIdentity()
		}), ","), strings.Join(lo.Map(requirements, func(t *Requirement, _ int) string {
			return t.Action.GetIdentity() + " " + t.Resource.GetNamespace() + "/" + t.Resource.GetIdentity()
		}), ","))
		return errors.Forbidden("", "")
	}
	//no claims
	return errors.Unauthorized("", "")
}

var ProviderSet = kitdi.NewSet(
	kitdi.NewProvider(NewDefaultAuthorizationService, di.As(new(Service))),
	kitdi.NewProvider(NewSubjectResolver, di.As(new(SubjectResolver))))
