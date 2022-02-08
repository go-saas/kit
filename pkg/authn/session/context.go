package session

import "context"

type clientStateKey struct {
}

type clientStateWriterKey struct {
}

func NewClientStateContext(ctx context.Context, state ClientState) context.Context {
	return context.WithValue(ctx, clientStateKey{}, state)
}

func FromClientStateContext(ctx context.Context) (state ClientState, ok bool) {
	state, ok = ctx.Value(clientStateKey{}).(ClientState)
	return
}

func NewClientStateWriterContext(ctx context.Context, sw ClientStateWriter) context.Context {
	return context.WithValue(ctx, clientStateWriterKey{}, sw)
}

func FromClientStateWriterContext(ctx context.Context) (sw ClientStateWriter, ok bool) {
	sw, ok = ctx.Value(clientStateWriterKey{}).(ClientStateWriter)
	return
}
