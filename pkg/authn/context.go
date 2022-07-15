package authn

import "context"

type userKey struct{}

func NewUserContext(ctx context.Context, user UserInfo) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

func FromUserContext(ctx context.Context) (user UserInfo, ok bool) {
	user, ok = ctx.Value(userKey{}).(UserInfo)
	if ok {
		return user, ok
	}
	return NewUserInfo(""), ok
}

type clientKey struct {
}

func NewClientContext(ctx context.Context, clientId string) context.Context {
	return context.WithValue(ctx, clientKey{}, clientId)
}

func FromClientContext(ctx context.Context) (clientId string, ok bool) {
	clientId, ok = ctx.Value(clientKey{}).(string)
	return
}
