package jwt

import "context"

type jwtKey struct{}

func NewJWTContext(ctx context.Context, jwt string) context.Context {
	return context.WithValue(ctx, jwtKey{}, jwt)
}

// FromJWTContext returns the Transport value stored in ctx, if any.
func FromJWTContext(ctx context.Context) (jwt string, ok bool) {
	jwt, ok = ctx.Value(jwtKey{}).(string)
	return
}
