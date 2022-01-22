package authorization

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

type Result struct {
	Allowed      bool
	Requirements []Requirement
}

func NewAllowAuthorizationResult() Result {
	return Result{Allowed: true}
}

func NewDisallowAuthorizationResult(requirements []Requirement) Result {
	return Result{Allowed: false, Requirements: requirements}
}

func FormatError(ctx context.Context, result Result, subjects ...Subject) error {
	if result.Allowed {
		return nil
	}
	var authed bool
	for _, sub := range subjects {
		if s, ok := sub.(*UserSubject); ok {
			if len(s.GetUserId()) > 0 {
				authed = true
			}
		}
		if s, ok := sub.(*ClientSubject); ok {
			if len(s.GetClientId()) > 0 {
				authed = true
			}
		}
	}
	if authed {
		//TODO format error
		return errors.Forbidden("", "")
	}
	//no claims
	return errors.Unauthorized("", "")
}
