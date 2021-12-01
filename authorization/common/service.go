package common

import (
	"context"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/auth/current"
)

type AuthorizationService interface {
	Check(ctx context.Context, resource Resource, action Action, subject ...Subject) (AuthorizationResult, error)
	CheckCurrent(ctx context.Context, resource Resource, action Action) (AuthorizationResult, error)
}

type SubjectContributor interface {
	Process(subject Subject) ([]Subject, error)
}

type AuthorizationOption struct {
	SubjectContributorList []SubjectContributor
}

func NewAuthorizationOption(subjectContributorList []SubjectContributor) *AuthorizationOption {
	return &AuthorizationOption{SubjectContributorList: subjectContributorList}
}

type DefaultAuthorizationService struct {
	opt     *AuthorizationOption
	checker PermissionChecker
}

var _ AuthorizationService = (*DefaultAuthorizationService)(nil)

func NewDefaultAuthorizationService(opt *AuthorizationOption, checker PermissionChecker) *DefaultAuthorizationService {
	return &DefaultAuthorizationService{opt: opt, checker: checker}
}

func (a *DefaultAuthorizationService) Check(ctx context.Context, resource Resource, action Action, subject ...Subject) (AuthorizationResult, error) {
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
			if s.GetName() == subject.GetName() && s.GetIdentity() == subject.GetIdentity() {
				return
			}
		}
		subjectList = append(subjectList, subject)
	}
	for _, s := range subject {
		addIfNotPresent(s)
		for _, contributor := range a.opt.SubjectContributorList {
			if subjects, err := contributor.Process(s); err != nil {
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
		if grantType == GrantTypeDisallow {
			return NewDisallowAuthorizationResult(nil), err
		}
		if grantType == GrantTypeAllow {
			anyAllow = true
		}
	}
	if anyAllow {
		return NewAllowAuthorizationResult(), nil
	}
	return NewDisallowAuthorizationResult(nil), nil
}

func (a *DefaultAuthorizationService) CheckCurrent(ctx context.Context, resource Resource, action Action) (AuthorizationResult, error) {
	var userId string
	if userInfo, ok := current.FromUserContext(ctx); ok {
		userId = userInfo.GetId()
	}
	return a.Check(ctx, resource, action, NewUserSubject(userId))
}

var ProviderSet = wire.NewSet(NewDefaultAuthorizationService, wire.Bind(new(AuthorizationService), new(*DefaultAuthorizationService)), NewPermissionService,
	wire.Bind(new(PermissionManagementService), new(*PermissionService)), wire.Bind(new(PermissionChecker), new(*PermissionService)))
