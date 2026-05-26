package userinfo

import (
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"time"
)

var idCache = gcache.New()

func AddIdToCache(userId uint64) {
	//缓存一分钟
	idCache.Set(gctx.New(), userId, userId, 1*time.Minute)
}

func InCache(userId uint64) bool {
	ok, e := idCache.Contains(gctx.New(), userId)
	if e != nil {
		return false
	}
	return ok
}
