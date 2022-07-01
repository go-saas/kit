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
}

type option struct {
	group   *singleflight.Group
	options []store.Option
}

type Option func(*option)

func WithGroup(g ...*singleflight.Group) Option {
	return func(o *option) {
		if len(g) > 0 {
			o.group = g[0]
		} else {
			o.group = &singleflight.Group{}
		}
	}
}

func WithStoreOption(opt ...store.Option) Option {
	return func(o *option) {
		o.options = opt
	}
}

func (h *Helper[T]) GetOrSet(ctx context.Context, key any, fn func(ctx context.Context) (T, error), opts ...Option) (v T, err error, set bool) {
	opt := &option{}
	for _, o := range opts {
		o(opt)
	}
	run := func() (v T, err error, set bool) {
		v, err = h.Get(ctx, key)
		if err == nil {
			return
		}
		if (store.NotFound{}).Is(err) {
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
		}
		return
	}

	if opt.group == nil {
		return run()
	}
	//run into group
	keyStr := h.getCacheKey(key)

	_, err, _ = opt.group.Do(keyStr, func() (interface{}, error) {
		v, err, set = run()
		return v, err
	})
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
