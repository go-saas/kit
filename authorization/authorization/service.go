package authorization

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/auth/current"
	"strings"
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
	log     *log.Helper
}

var _ Service = (*DefaultAuthorizationService)(nil)

func NewDefaultAuthorizationService(opt *Option, checker PermissionChecker, logger log.Logger) *DefaultAuthorizationService {
	return &DefaultAuthorizationService{opt: opt, checker: checker, log: log.NewHelper(log.With(logger, "module", "authorization.service"))}
}

func (a *DefaultAuthorizationService) CheckForSubjects(ctx context.Context, resource Resource, action Action, subject ...Subject) (Result, error) {
	if always, ok := FromAlwaysAuthorizationContext(ctx); ok {
		var subjectStr []string
		for _, s := range subject {
			subjectStr = append(subjectStr, s.GetIdentity())
		}
		if always {
			a.log.Debugf("check permission for subject %s action %s to resource %s granted", strings.Join(subjectStr, ","), action.GetIdentity(), fmt.Sprintf("%s/%s", resource.GetNamespace(), resource.GetIdentity()))
			return NewAllowAuthorizationResult(), nil
		} else {
			a.log.Debugf("check permission for subject %s action %s to resource %s forbidden", strings.Join(subjectStr, ","), action.GetIdentity(), fmt.Sprintf("%s/%s", resource.GetNamespace(), resource.GetIdentity()))
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
	var logStr []string
	for _, s := range subjectList {
		logStr = append(logStr, s.GetIdentity())
	}
	a.log.Debugf("check permission for subject %s action %s to resource %s ", strings.Join(logStr, ","), action.GetIdentity(), fmt.Sprintf("%s/%s", resource.GetNamespace(), resource.GetIdentity()))
	var anyAllow bool
	for _, s := range subjectList {
		grantType, err := a.checker.IsGrant(ctx, resource, action, s)
		if err != nil {
			return NewDisallowAuthorizationResult(nil), err
		}
		if grantType == EffectForbidden {
			a.log.Debugf("check permission for subject %s action %s to resource %s forbidden", strings.Join(logStr, ","), action.GetIdentity(), fmt.Sprintf("%s/%s", resource.GetNamespace(), resource.GetIdentity()))
			return NewDisallowAuthorizationResult(nil), err
		}
		if grantType == EffectGrant {
			anyAllow = true
		}
	}
	if anyAllow {
		a.log.Debugf("check permission for subject %s action %s to resource %s granted", strings.Join(logStr, ","), action.GetIdentity(), fmt.Sprintf("%s/%s", resource.GetNamespace(), resource.GetIdentity()))
		return NewAllowAuthorizationResult(), nil
	}
	a.log.Debugf("check permission for subject %s action %s to resource %s forbidden", strings.Join(logStr, ","), action.GetIdentity(), fmt.Sprintf("%s/%s", resource.GetNamespace(), resource.GetIdentity()))
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
