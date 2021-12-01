package common

type AuthorizationResult struct {
	Allowed      bool
	Requirements []Requirement
}

func NewAllowAuthorizationResult() AuthorizationResult {
	return AuthorizationResult{Allowed: true}
}

func NewDisallowAuthorizationResult(requirements []Requirement) AuthorizationResult {
	return AuthorizationResult{Allowed: false, Requirements: requirements}
}
