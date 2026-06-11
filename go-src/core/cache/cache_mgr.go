package cache

import (
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"time"
)

const (
	CacheTime = 30 * time.Minute
)

type CacheMgr struct {
	Cache *gcache.Cache
}

func NewCacheMgr() *CacheMgr {
	return &CacheMgr{
		Cache: gcache.New(),
	}
}

// GetData 取缓存数据,命中不了缓存,从数据库拉取数据 ttl 30分钟
func (mgr *CacheMgr) GetData(key interface{}, f gcache.Func) any {
	ctx := gctx.New()
	//缓存30分钟，让变更数据,有时间同步到数据库
	data, _ := mgr.Cache.GetOrSetFuncLock(ctx, key, f, CacheTime)
	//延迟缓存失效时间
	_, _ = mgr.Cache.UpdateExpire(ctx, key, CacheTime)
	return data.Val()
}

// FlushCache 切片数据变动时候,刷新缓存
func (mgr *CacheMgr) FlushCache(key interface{}, data any) {
	err := mgr.Cache.Set(gctx.New(), key, data, CacheTime)
	if err != nil {
		return
	}
}
