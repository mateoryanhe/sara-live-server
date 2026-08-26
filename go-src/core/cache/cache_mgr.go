package cache

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	CacheTime = 30 * time.Minute
)

// Loader 缓存未命中时加载数据.
type Loader[T any] func(ctx context.Context) (T, error)

type cacheCore[T any] struct {
	cache *gcache.Cache
	ttl   time.Duration
}

func (c *cacheCore[T]) ctx(ctx context.Context) context.Context {
	if ctx != nil {
		return ctx
	}
	return gctx.New()
}

func (c *cacheCore[T]) load(ctx context.Context, key any, loader Loader[T]) (T, bool) {
	var zero T
	if c == nil || c.cache == nil || loader == nil {
		return zero, false
	}
	ctx = c.ctx(ctx)
	data, err := c.cache.GetOrSetFuncLock(ctx, key, func(ctx context.Context) (any, error) {
		return loader(ctx)
	}, c.ttl)
	if err != nil || data.IsNil() {
		return zero, false
	}
	if c.ttl > 0 {
		_, _ = c.cache.UpdateExpire(ctx, key, c.ttl)
	}
	typed, ok := data.Val().(T)
	return typed, ok
}

func (c *cacheCore[T]) peek(ctx context.Context, key any) (T, bool) {
	var zero T
	if c == nil || c.cache == nil {
		return zero, false
	}
	v, err := c.cache.Get(c.ctx(ctx), key)
	if err != nil || v.IsNil() {
		return zero, false
	}
	typed, ok := v.Val().(T)
	return typed, ok
}

func (c *cacheCore[T]) publish(ctx context.Context, key any, data T) {
	if c == nil || c.cache == nil {
		return
	}
	ctx = c.ctx(ctx)
	if _, exist, err := c.cache.Update(ctx, key, data); err == nil && exist {
		return
	}
	_ = c.cache.Set(ctx, key, data, c.ttl)
}

func (c *cacheCore[T]) set(ctx context.Context, key any, data T, ttl ...time.Duration) error {
	if c == nil || c.cache == nil {
		return nil
	}
	d := c.ttl
	if len(ttl) > 0 {
		d = ttl[0]
	}
	return c.cache.Set(c.ctx(ctx), key, data, d)
}

func (c *cacheCore[T]) remove(ctx context.Context, keys ...any) {
	if c == nil || c.cache == nil || len(keys) == 0 {
		return
	}
	_, _ = c.cache.Remove(c.ctx(ctx), keys...)
}

func (c *cacheCore[T]) contains(ctx context.Context, key any) bool {
	if c == nil || c.cache == nil {
		return false
	}
	ok, _ := c.cache.Contains(c.ctx(ctx), key)
	return ok
}

func (c *cacheCore[T]) keys(ctx context.Context) ([]any, error) {
	if c == nil || c.cache == nil {
		return nil, nil
	}
	return c.cache.Keys(c.ctx(ctx))
}

func (c *cacheCore[T]) values(ctx context.Context) ([]T, error) {
	if c == nil || c.cache == nil {
		return nil, nil
	}
	vals, err := c.cache.Values(c.ctx(ctx))
	if err != nil {
		return nil, err
	}
	out := make([]T, 0, len(vals))
	for _, v := range vals {
		if v == nil {
			continue
		}
		if typed, ok := v.(T); ok {
			out = append(out, typed)
		}
	}
	return out, nil
}

// RowCache 单行实体缓存. GetRow 返回缓存单例, 原地修改后必须 PublishRow.
type RowCache[T any] struct {
	core *cacheCore[T]
}

func NewRowCache[T any]() *RowCache[T] {
	return &RowCache[T]{
		core: &cacheCore[T]{cache: gcache.New(), ttl: CacheTime},
	}
}

func NewPermanentRowCache[T any]() *RowCache[T] {
	return &RowCache[T]{
		core: &cacheCore[T]{cache: gcache.New(), ttl: 0},
	}
}

func (c *RowCache[T]) GetRow(ctx context.Context, key any, loader Loader[T]) (T, bool) {
	if c == nil || c.core == nil {
		var zero T
		return zero, false
	}
	return c.core.load(ctx, key, loader)
}

func (c *RowCache[T]) MustGetRow(ctx context.Context, key any, loader Loader[T]) T {
	v, _ := c.GetRow(ctx, key, loader)
	return v
}

func (c *RowCache[T]) GetRowCached(ctx context.Context, key any) (T, bool) {
	if c == nil || c.core == nil {
		var zero T
		return zero, false
	}
	return c.core.peek(ctx, key)
}

func (c *RowCache[T]) PublishRow(ctx context.Context, key any, data T) {
	if c == nil || c.core == nil {
		return
	}
	c.core.publish(ctx, key, data)
}

func (c *RowCache[T]) SetRow(ctx context.Context, key any, data T, ttl ...time.Duration) error {
	if c == nil || c.core == nil {
		return nil
	}
	return c.core.set(ctx, key, data, ttl...)
}

func (c *RowCache[T]) RemoveRow(ctx context.Context, keys ...any) {
	if c == nil || c.core == nil {
		return
	}
	c.core.remove(ctx, keys...)
}

func (c *RowCache[T]) ContainsRow(ctx context.Context, key any) bool {
	if c == nil || c.core == nil {
		return false
	}
	return c.core.contains(ctx, key)
}

func (c *RowCache[T]) Keys(ctx context.Context) ([]any, error) {
	if c == nil || c.core == nil {
		return nil, nil
	}
	return c.core.keys(ctx)
}

func (c *RowCache[T]) Values(ctx context.Context) ([]T, error) {
	if c == nil || c.core == nil {
		return nil, nil
	}
	return c.core.values(ctx)
}

// ListCache 列表缓存. GetList 返回缓存列表, 变更后必须 PublishList.
type ListCache[T any] struct {
	core *cacheCore[[]T]
}

func NewListCache[T any]() *ListCache[T] {
	return &ListCache[T]{
		core: &cacheCore[[]T]{cache: gcache.New(), ttl: CacheTime},
	}
}

func NewPermanentListCache[T any]() *ListCache[T] {
	return &ListCache[T]{
		core: &cacheCore[[]T]{cache: gcache.New(), ttl: 0},
	}
}

func (c *ListCache[T]) GetList(ctx context.Context, key any, loader Loader[[]T]) ([]T, bool) {
	if c == nil || c.core == nil {
		return nil, false
	}
	return c.core.load(ctx, key, loader)
}

func (c *ListCache[T]) MustGetList(ctx context.Context, key any, loader Loader[[]T]) []T {
	v, _ := c.GetList(ctx, key, loader)
	return v
}

func (c *ListCache[T]) GetListCached(ctx context.Context, key any) ([]T, bool) {
	if c == nil || c.core == nil {
		return nil, false
	}
	return c.core.peek(ctx, key)
}

func (c *ListCache[T]) PublishList(ctx context.Context, key any, data []T) {
	if c == nil || c.core == nil {
		return
	}
	if data == nil {
		data = []T{}
	}
	c.core.publish(ctx, key, data)
}

func (c *ListCache[T]) SetList(ctx context.Context, key any, data []T, ttl ...time.Duration) error {
	if c == nil || c.core == nil {
		return nil
	}
	if data == nil {
		data = []T{}
	}
	return c.core.set(ctx, key, data, ttl...)
}

func (c *ListCache[T]) RemoveList(ctx context.Context, keys ...any) {
	if c == nil || c.core == nil {
		return
	}
	c.core.remove(ctx, keys...)
}

func (c *ListCache[T]) ContainsList(ctx context.Context, key any) bool {
	if c == nil || c.core == nil {
		return false
	}
	return c.core.contains(ctx, key)
}

func (c *ListCache[T]) Keys(ctx context.Context) ([]any, error) {
	if c == nil || c.core == nil {
		return nil, nil
	}
	return c.core.keys(ctx)
}

func (c *ListCache[T]) Values(ctx context.Context) ([][]T, error) {
	if c == nil || c.core == nil {
		return nil, nil
	}
	vals, err := c.core.values(ctx)
	if err != nil {
		return nil, err
	}
	return vals, nil
}
