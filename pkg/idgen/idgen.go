package idgen

import "context"

type Generator interface {
	Gen(ctx context.Context) (string, error)
}

type GeneratorFunc func(ctx context.Context) (string, error)

func (g GeneratorFunc) Gen(ctx context.Context) (string, error) {
	return g(ctx)
}
