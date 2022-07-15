package authz

import "context"

type alwaysKey struct{}

//NewAlwaysAuthorizationContext create a context for always pass or forbidden authorization check. useful for testing
func NewAlwaysAuthorizationContext(ctx context.Context, allow ...bool) context.Context {
	v := true
	if len(allow) > 0 {
		v = allow[0]
	}
	return context.WithValue(ctx, alwaysKey{}, v)
}

func FromAlwaysAuthorizationContext(ctx context.Context) (allow bool, ok bool) {
	allow, ok = ctx.Value(alwaysKey{}).(bool)
	return
}
