package current

import "context"

type userKey struct{}

func NewUserContext(ctx context.Context, user UserInfo) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

func FromUserContext(ctx context.Context) (user UserInfo, ok bool) {
	user, ok = ctx.Value(userKey{}).(UserInfo)
	return
}
