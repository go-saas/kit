package authz

import (
	"context"
	"github.com/go-saas/kit/pkg/authn"
)

type SubjectResolver interface {
	//ResolveFromContext extract subjects from current ctx
	ResolveFromContext(ctx context.Context) ([]Subject, error)
	//ResolveProcessed recursively find related subjects. (RBAC)
	ResolveProcessed(ctx context.Context, subjects ...Subject) ([]Subject, error)
}

type SubjectResolverImpl struct {
	opt *Option
}

var _ SubjectResolver = (*SubjectResolverImpl)(nil)

func NewSubjectResolver(opt *Option) *SubjectResolverImpl {
	return &SubjectResolverImpl{opt: opt}
}

func (s *SubjectResolverImpl) ResolveFromContext(ctx context.Context) ([]Subject, error) {
	var subjects []Subject
	var userId string
	userInfo, _ := authn.FromUserContext(ctx)
	userId = userInfo.GetId()
	//append empty user
	subjects = append(subjects, NewUserSubject(userId))
	if clientId, ok := authn.FromClientContext(ctx); ok && len(clientId) > 0 {
		//do not append empty client
		subjects = append(subjects, NewClientSubject(clientId))
	}
	return subjects, nil
}

func (s *SubjectResolverImpl) ResolveProcessed(ctx context.Context, subjects ...Subject) ([]Subject, error) {
	//use contributor
	var subjectList []Subject

	addIfNotPresent := func(subject Subject) bool {
		for _, s := range subjectList {
			if s.GetIdentity() == subject.GetIdentity() {
				return false
			}
		}
		subjectList = append(subjectList, subject)
		return true
	}
	for _, s := range subjects {
		addIfNotPresent(s)
	}
	i := 0
	for {
		if i == len(subjectList) {
			break
		}
		for _, contributor := range s.opt.SubjectContribList {
			if subjects, err := contributor.Process(ctx, subjectList[i]); err != nil {
				return nil, err
			} else {
				for _, s2 := range subjects {
					addIfNotPresent(s2)
				}
			}
		}
		i++
	}
	return subjectList, nil
}
