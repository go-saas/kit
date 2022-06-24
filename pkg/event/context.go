package event

import "context"

type (
	receiverKey struct{}
	senderKey   struct{}
)

func NewReceiverContext(ctx context.Context, r Receiver) context.Context {
	return context.WithValue(ctx, receiverKey{}, r)
}

func FromReceiverContext(ctx context.Context) (Receiver, bool) {
	v, ok := ctx.Value(receiverKey{}).(Receiver)
	return v, ok
}
func NewSenderContext(ctx context.Context, r Sender) context.Context {
	return context.WithValue(ctx, senderKey{}, r)
}

func FromSenderContext(ctx context.Context) (Sender, bool) {
	v, ok := ctx.Value(senderKey{}).(Sender)
	return v, ok
}
