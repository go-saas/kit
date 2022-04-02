package biz

import "context"

type enableUserTenantsKey struct{}

func NewEnableUserTenantsContext(ctx context.Context, enable ...bool) context.Context {
	if len(enable) > 0 {
		return context.WithValue(ctx, enableUserTenantsKey{}, enable[0])
	} else {
		return context.WithValue(ctx, enableUserTenantsKey{}, true)
	}
}

func FromEnableUserTenantContext(ctx context.Context) bool {
	v, ok := ctx.Value(enableUserTenantsKey{}).(bool)
	if ok {
		return v
	}
	return false
}
