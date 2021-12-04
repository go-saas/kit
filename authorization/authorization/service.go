package authorization

import (
	"context"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/auth/current"
)

type Service interface {
	CheckForSubjects(ctx context.Context, resource Resource, action Action, subject ...Subject) (Result, error)
	Check(ctx context.Context, resource Resource, action Action) (Result, error)
}

type SubjectContributor interface {
	Process(ctx context.Context, subject Subject) ([]Subject, error)
}

type Option struct {
	SubjectContributorList []SubjectContributor
}

func NewAuthorizationOption(subjectContributorList ...SubjectContributor) *Option {
	return &Option{SubjectContributorList: subjectContributorList}
}

type DefaultAuthorizationService struct {
	opt     *Option
	checker PermissionChecker
}

var _ Service = (*DefaultAuthorizationService)(nil)

func NewDefaultAuthorizationService(opt *Option, checker PermissionChecker) *DefaultAuthorizationService {
	return &DefaultAuthorizationService{opt: opt, checker: checker}
}

func (a *DefaultAuthorizationService) CheckForSubjects(ctx context.Context, resource Resource, action Action, subject ...Subject) (Result, error) {
	if always, ok := FromAlwaysAuthorizationContext(ctx); ok {
		if always {
			return NewAllowAuthorizationResult(), nil
		} else {
			return NewDisallowAuthorizationResult(nil), nil
		}
	}
	var subjectList []Subject

	addIfNotPresent := func(subject Subject) {
		for _, s := range subjectList {
			if s.GetIdentity() == subject.GetIdentity() {
				return
			}
		}
		subjectList = append(subjectList, subject)
	}
	for _, s := range subject {
		addIfNotPresent(s)
		for _, contributor := range a.opt.SubjectContributorList {
			if subjects, err := contributor.Process(ctx, s); err != nil {
				return NewDisallowAuthorizationResult(nil), err
			} else {
				for _, s2 := range subjects {
					addIfNotPresent(s2)
				}
			}
		}
	}
	var anyAllow bool
	for _, s := range subjectList {
		grantType, err := a.checker.IsGrant(ctx, resource, action, s)
		if err != nil {
			return NewDisallowAuthorizationResult(nil), err
		}
		if grantType == EffectForbidden {
			return NewDisallowAuthorizationResult(nil), err
		}
		if grantType == EffectGrant {
			anyAllow = true
		}
	}
	if anyAllow {
		return NewAllowAuthorizationResult(), nil
	}
	return NewDisallowAuthorizationResult(nil), nil
}

func (a *DefaultAuthorizationService) Check(ctx context.Context, resource Resource, action Action) (Result, error) {
	var subjects []Subject
	var userId string
	if userInfo, ok := current.FromUserContext(ctx); ok {
		userId = userInfo.GetId()
		subjects = append(subjects, NewUserSubject(userId))
	}
	if clientId, ok := current.FromClientContext(ctx); ok {
		subjects = append(subjects, NewClientSubject(clientId))
	}
	return a.CheckForSubjects(ctx, resource, action, subjects...)
}

var ProviderSet = wire.NewSet(NewDefaultAuthorizationService, wire.Bind(new(Service), new(*DefaultAuthorizationService)), NewPermissionService,
	wire.Bind(new(PermissionManagementService), new(*PermissionService)), wire.Bind(new(PermissionChecker), new(*PermissionService)))
