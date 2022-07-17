package di

import (
	"fmt"
	"github.com/goava/di"
	"reflect"
)

func NewSet(constructor ...interface{}) di.Option {
	var opts []di.Option
	for _, i := range constructor {
		if opt, ok := i.(di.Option); ok {
			opts = append(opts, opt)
		} else if p, ok := i.(*Provider); ok {
			opts = append(opts, di.Provide(p.constructor, p.opts...))
		} else {
			ctype := reflect.TypeOf(i)
			if ctype == nil {
				panic("can't provide an untyped nil")
			}
			if ctype.Kind() == reflect.Func {
				opts = append(opts, di.Provide(i))
			} else {
				panic(fmt.Errorf("can not resolve %v (type %v), you probably need wrap with di.Value", i, ctype))
			}
		}
	}
	return di.Options(opts...)

}

func Value[T any](v T, opts ...di.ProvideOption) di.Option {
	return di.Provide(func() T { return v }, opts...)
}

type Provider struct {
	constructor interface{}
	opts        []di.ProvideOption
}

func NewProvider(constructor interface{}, opts ...di.ProvideOption) *Provider {
	return &Provider{
		constructor: constructor,
		opts:        opts,
	}
}
