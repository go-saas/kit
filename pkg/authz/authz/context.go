package authz

import "context"

type alwaysKey struct{}

//NewAlwaysAuthorizationContext create a context for always pass or forbidden authorization check. useful for testing
func NewAlwaysAuthorizationContext(ctx context.Context, allow bool) context.Context {
	return context.WithValue(ctx, alwaysKey{}, allow)
}

func FromAlwaysAuthorizationContext(ctx context.Context) (allow bool, ok bool) {
	allow, ok = ctx.Value(alwaysKey{}).(bool)
	return
}
