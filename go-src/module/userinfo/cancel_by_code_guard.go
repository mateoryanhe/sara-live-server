package userinfo

import (
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/errercode"
)

const (
	cancelByCodeFailLimit  = 3
	cancelByCodeFailExpire = 2 * time.Hour
)

var cancelByCodeCacheMgr *cache.CacheMgr

func initCancelByCodeGuard() {
	cancelByCodeCacheMgr = cache.NewCacheMgr()
}

func checkCancelByCodeIPLimit(ip string) error {
	if ip == "" {
		return nil
	}
	ctx := gctx.New()
	if blocked, _ := cancelByCodeCacheMgr.Cache.Contains(ctx, cancelByCodeBlockKey(ip)); blocked {
		return errercode.CreateCode(errercode.RequestTooFrequent)
	}
	return nil
}

func markCancelByCodeFailure(ip string, code errercode.XRCode) error {
	if ip == "" {
		return errercode.CreateCode(code)
	}
	ctx := gctx.New()
	key := cancelByCodeFailKey(ip)
	count := 1

	val, _ := cancelByCodeCacheMgr.Cache.Get(ctx, key)
	if val != nil {
		if current, ok := val.Val().(int); ok && current > 0 {
			count = current + 1
		}
	}

	if count >= cancelByCodeFailLimit {
		setCancelByCodeCache(cancelByCodeBlockKey(ip), time.Now().Unix(), cancelByCodeFailExpire)
		_, _ = cancelByCodeCacheMgr.Cache.Remove(ctx, key)
		return errercode.CreateCode(errercode.RequestTooFrequent)
	}

	setCancelByCodeCache(key, count, cancelByCodeFailExpire)
	return errercode.CreateCode(code)
}

func clearCancelByCodeFailure(ip string) {
	if ip == "" {
		return
	}
	ctx := gctx.New()
	_, _ = cancelByCodeCacheMgr.Cache.Remove(ctx, cancelByCodeFailKey(ip))
	_, _ = cancelByCodeCacheMgr.Cache.Remove(ctx, cancelByCodeBlockKey(ip))
}

func cancelByCodeFailKey(ip string) string {
	return "cancel_by_code_fail:" + ip
}

func cancelByCodeBlockKey(ip string) string {
	return "cancel_by_code_block:" + ip
}

func setCancelByCodeCache(key string, data any, ttl time.Duration) {
	_ = cancelByCodeCacheMgr.Cache.Set(gctx.New(), key, data, ttl)
}
