package authz

type Result struct {
	Allowed      bool
	Requirements []*Requirement
}

func NewAllowAuthorizationResult() *Result {
	return &Result{Allowed: true}
}

func NewDisallowAuthorizationResult(requirements ...*Requirement) *Result {
	return &Result{Allowed: false, Requirements: requirements}
}
