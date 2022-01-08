package jwt

import "context"

type jwtKey struct{}
type claimKey struct{}

func NewClaimsContext(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimKey{}, claims)
}

func FromClaimsContext(ctx context.Context) (claims *Claims, ok bool) {
	claims, ok = ctx.Value(claimKey{}).(*Claims)
	return
}

func NewJWTContext(ctx context.Context, jwt string) context.Context {
	return context.WithValue(ctx, jwtKey{}, jwt)
}

func FromJWTContext(ctx context.Context) (jwt string, ok bool) {
	jwt, ok = ctx.Value(jwtKey{}).(string)
	return
}
