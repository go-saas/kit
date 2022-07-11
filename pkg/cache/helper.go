package cache

import (
	"context"
	"crypto"
	"fmt"
	"github.com/eko/gocache/v3/cache"
	"github.com/eko/gocache/v3/store"
	"golang.org/x/sync/singleflight"
	"reflect"
)

type Helper[T any] struct {
	cache.CacheInterface[T]
	group singleflight.Group
}

func NewHelper[T any](c cache.CacheInterface[T]) *Helper[T] {
	return &Helper[T]{CacheInterface: c}
}

type option struct {
	group   *singleflight.Group
	options []store.Option
}

type Option func(*option)

// WithGroup pass nil to disable singleflight.Group
func WithGroup(g ...*singleflight.Group) Option {
	return func(o *option) {
		if len(g) > 0 {
			o.group = g[0]
		}
	}
}

func WithStoreOption(opt ...store.Option) Option {
	return func(o *option) {
		o.options = opt
	}
}

func (h *Helper[T]) GetOrSet(ctx context.Context, key any, fn func(ctx context.Context) (T, error), opts ...Option) (v T, err error, set bool) {
	opt := &option{group: &h.group}
	for _, o := range opts {
		o(opt)
	}
	v, err = h.Get(ctx, key)
	if err == nil {
		//resolve from cache
		return
	}
	if !(store.NotFound{}).Is(err) {
		//cache error
		return
	}
	run := func() (v T, err error, set bool) {
		//use factory
		v, err = fn(ctx)
		if err != nil {
			return
		}
		//push back
		err = h.Set(ctx, key, v, opt.options...)
		if err != nil {
			return
		}
		set = true
		return
	}
	if opt.group == nil {
		return run()
	}
	//run into group
	keyStr := h.getCacheKey(key)
	var value interface{}
	value, err, _ = opt.group.Do(keyStr, func() (interface{}, error) {
		v, err, set = run()
		return v, err
	})
	v = value.(T)
	return
}

func (h *Helper[T]) getCacheKey(key any) string {
	switch key.(type) {
	case string:
		return key.(string)
	case cache.CacheKeyGenerator:
		return key.(cache.CacheKeyGenerator).GetCacheKey()
	default:
		return checksum(key)
	}
}

// checksum hashes a given object into a string
func checksum(object any) string {
	digester := crypto.MD5.New()
	fmt.Fprint(digester, reflect.TypeOf(object))
	fmt.Fprint(digester, object)
	hash := digester.Sum(nil)

	return fmt.Sprintf("%x", hash)
}
