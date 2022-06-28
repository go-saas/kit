package cache

import (
	"context"
	"github.com/eko/gocache/v3/cache"
	"github.com/eko/gocache/v3/store"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type ProtoCache[T proto.Message] struct {
	creator func() T
	proxy   cache.CacheInterface[string]
}

var _ cache.CacheInterface[*anypb.Any] = (*ProtoCache[*anypb.Any])(nil)

func NewProtoCache[T proto.Message](creator func() T, proxy cache.CacheInterface[string]) *ProtoCache[T] {
	return &ProtoCache[T]{creator: creator, proxy: proxy}
}

func (c *ProtoCache[T]) Get(ctx context.Context, key any) (T, error) {

	t, err := c.proxy.Get(ctx, key)
	if err != nil {
		var n T
		return n, err
	}
	f := c.creator()
	err = protojson.Unmarshal([]byte(t), f)
	if err != nil {
		var n T
		return n, err
	}
	return f, nil
}

func (c *ProtoCache[T]) Set(ctx context.Context, key any, object T, options ...store.Option) error {
	b, err := protojson.Marshal(object)
	if err != nil {
		return err
	}
	return c.proxy.Set(ctx, key, string(b), options...)
}

func (c *ProtoCache[T]) Delete(ctx context.Context, key any) error {
	return c.proxy.Delete(ctx, key)
}

func (c *ProtoCache[T]) Invalidate(ctx context.Context, options ...store.InvalidateOption) error {
	return c.proxy.Invalidate(ctx, options...)
}

func (c *ProtoCache[T]) Clear(ctx context.Context) error {
	return c.proxy.Clear(ctx)
}

func (c *ProtoCache[T]) GetType() string {
	return c.proxy.GetType()
}
