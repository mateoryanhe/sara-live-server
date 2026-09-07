package cache

import (
	"context"
	"testing"
)

func TestCacheCoreCachesNilPointer(t *testing.T) {
	c := NewRowCache[*int]()
	loadCount := 0
	loader := func(ctx context.Context) (*int, error) {
		loadCount++
		return nil, nil
	}
	v1, ok1 := c.GetRow(context.Background(), "k", loader)
	if !ok1 || v1 != nil || loadCount != 1 {
		t.Fatalf("first get: ok=%v v=%v loadCount=%d", ok1, v1, loadCount)
	}
	v2, ok2 := c.GetRow(context.Background(), "k", loader)
	if !ok2 || v2 != nil || loadCount != 1 {
		t.Fatalf("second get should hit cache: ok=%v v=%v loadCount=%d", ok2, v2, loadCount)
	}
	if !c.ContainsRow(context.Background(), "k") {
		t.Fatal("nil pointer miss should still occupy cache key")
	}
}

func TestCacheCoreCachesEmptySlice(t *testing.T) {
	c := NewListCache[int]()
	loadCount := 0
	loader := func(ctx context.Context) ([]int, error) {
		loadCount++
		return nil, nil // loader 返回 nil slice
	}
	v1, ok1 := c.GetList(context.Background(), "k", loader)
	if !ok1 || v1 != nil || loadCount != 1 {
		t.Fatalf("first get: ok=%v v=%v loadCount=%d", ok1, v1, loadCount)
	}
	v2, ok2 := c.GetList(context.Background(), "k", loader)
	if !ok2 || loadCount != 1 {
		t.Fatalf("second get should hit cache: ok=%v loadCount=%d v=%v", ok2, loadCount, v2)
	}
}

func TestCacheCoreCachesZeroUint64(t *testing.T) {
	c := NewRowCache[uint64]()
	loadCount := 0
	loader := func(ctx context.Context) (uint64, error) {
		loadCount++
		return 0, nil
	}
	v1 := c.MustGetRow(context.Background(), "email", loader)
	v2 := c.MustGetRow(context.Background(), "email", loader)
	if v1 != 0 || v2 != 0 || loadCount != 1 {
		t.Fatalf("zero should be cached: v1=%d v2=%d loadCount=%d", v1, v2, loadCount)
	}
}

func TestCacheCorePublishNilPointer(t *testing.T) {
	c := NewRowCache[*string]()
	c.PublishRow(context.Background(), "k", nil)
	if !c.ContainsRow(context.Background(), "k") {
		t.Fatal("PublishRow(nil) must not delete key via gcache Set(nil)")
	}
	v, ok := c.GetRowCached(context.Background(), "k")
	if !ok || v != nil {
		t.Fatalf("peek nil: ok=%v v=%v", ok, v)
	}
}
