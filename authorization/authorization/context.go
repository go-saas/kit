package authorization

import "context"

type alwaysKey struct{}

func NewAlwaysAuthorizationContext(ctx context.Context, allow bool) context.Context {
	return context.WithValue(ctx, alwaysKey{}, allow)
}

func FromAlwaysAuthorizationContext(ctx context.Context) (allow bool, ok bool) {
	allow, ok = ctx.Value(alwaysKey{}).(bool)
	return
}
