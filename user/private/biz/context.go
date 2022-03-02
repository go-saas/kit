package biz

import "context"

type ignoreUserTenantsKey struct {
}

func NewIgnoreUserTenantsContext(ctx context.Context, ignore bool) context.Context {
	return context.WithValue(ctx, ignoreUserTenantsKey{}, ignore)
}

func FromIgnoreUserTenantsContext(ctx context.Context) (ignore bool) {
	v, ok := ctx.Value(ignoreUserTenantsKey{}).(bool)
	if ok {
		return v
	}
	return false
}
