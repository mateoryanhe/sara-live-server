package xrtoken

import (
	"sort"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
)

// QueryAppTokens 从内存缓存分页查询App Token,可按用户ID过滤
func QueryAppTokens(userId uint64, pageIndex, pageSize int) (int, []AppTokenCacheItem) {
	ctx := gctx.New()
	keys, err := appCache.Keys(ctx)
	if err != nil || len(keys) == 0 {
		return 0, nil
	}

	items := make([]AppTokenCacheItem, 0, len(keys))
	for _, key := range keys {
		uid := gconv.Uint64(key)
		if userId > 0 && uid != userId {
			continue
		}
		tokenVar, err := appCache.Get(ctx, key)
		if err != nil || tokenVar.IsNil() {
			continue
		}
		ttl, err := appCache.GetExpire(ctx, key)
		if err != nil || ttl <= 0 {
			continue
		}
		items = append(items, AppTokenCacheItem{
			UserId:   uid,
			Token:    tokenVar.String(),
			ExpireAt: time.Now().Add(ttl),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].UserId > items[j].UserId
	})

	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	total := len(items)
	start := (pageIndex - 1) * pageSize
	if start >= total {
		return total, []AppTokenCacheItem{}
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return total, items[start:end]
}
