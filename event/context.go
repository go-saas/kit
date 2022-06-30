package event

import "context"

type (
	consumerKey struct{}
	producerKey struct{}
)

func NewConsumerContext(ctx context.Context, r Consumer) context.Context {
	return context.WithValue(ctx, consumerKey{}, r)
}

func FromConsumerContext(ctx context.Context) (Consumer, bool) {
	v, ok := ctx.Value(consumerKey{}).(Consumer)
	return v, ok
}
func NewProducerContext(ctx context.Context, r Producer) context.Context {
	return context.WithValue(ctx, producerKey{}, r)
}

func FromProducerContext(ctx context.Context) (Producer, bool) {
	v, ok := ctx.Value(producerKey{}).(Producer)
	return v, ok
}
