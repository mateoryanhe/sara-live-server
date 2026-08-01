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
	ttl   time.Duration
}

func NewCacheMgr() *CacheMgr {
	return &CacheMgr{
		Cache: gcache.New(),
		ttl:   CacheTime,
	}
}

// NewPermanentCacheMgr 创建永不过期的缓存管理器
func NewPermanentCacheMgr() *CacheMgr {
	return &CacheMgr{
		Cache: gcache.New(),
		ttl:   0,
	}
}

// GetData 取缓存数据,命中不了缓存,从数据库拉取数据 ttl 30分钟
func (mgr *CacheMgr) GetData(key interface{}, f gcache.Func) any {
	ctx := gctx.New()
	data, _ := mgr.Cache.GetOrSetFuncLock(ctx, key, f, mgr.ttl)
	if mgr.ttl > 0 {
		//延迟缓存失效时间
		_, _ = mgr.Cache.UpdateExpire(ctx, key, mgr.ttl)
	}
	return data.Val()
}

// FlushCache 切片数据变动时候,刷新缓存
func (mgr *CacheMgr) FlushCache(key interface{}, data any) {
	err := mgr.Cache.Set(gctx.New(), key, data, mgr.ttl)
	if err != nil {
		return
	}
}

// GetFromCache 仅从内存缓存读取,未命中不访问数据库
func (mgr *CacheMgr) GetFromCache(key interface{}) any {
	if mgr == nil || mgr.Cache == nil {
		return nil
	}
	v, err := mgr.Cache.Get(gctx.New(), key)
	if err != nil || v.IsNil() {
		return nil
	}
	return v.Val()
}
