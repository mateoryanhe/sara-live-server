package auth

import (
	"time"
	"xr-game-server/core/cache"
	"xr-game-server/errercode"

	"github.com/gogf/gf/v2/os/gctx"
)

const (
	phoneLoginFailLimit  = 5
	phoneLoginFailExpire = 2 * time.Hour
)

var phoneLoginCacheMgr *cache.CacheMgr

func initPhoneLoginGuard() {
	phoneLoginCacheMgr = cache.NewCacheMgr()
}

func checkPhoneLoginLimit(phoneKey string) error {
	ctx := gctx.New()
	if blocked, _ := phoneLoginCacheMgr.Cache.Contains(ctx, phoneLoginBlockKey(phoneKey)); blocked {
		return errercode.CreateCode(errercode.RequestTooFrequent)
	}
	return nil
}

func markPhoneLoginFailure(phoneKey string) error {
	ctx := gctx.New()
	key := phoneLoginFailKey(phoneKey)
	count := 1

	val, _ := phoneLoginCacheMgr.Cache.Get(ctx, key)
	if val != nil {
		if current, ok := val.Val().(int); ok && current > 0 {
			count = current + 1
		}
	}

	if count >= phoneLoginFailLimit {
		setPhoneLoginCache(phoneLoginBlockKey(phoneKey), time.Now().Unix(), phoneLoginFailExpire)
		_, _ = phoneLoginCacheMgr.Cache.Remove(ctx, key)
		return errercode.CreateCode(errercode.RequestTooFrequent)
	}

	setPhoneLoginCache(key, count, phoneLoginFailExpire)
	return nil
}

func clearPhoneLoginFailure(phoneKey string) {
	ctx := gctx.New()
	_, _ = phoneLoginCacheMgr.Cache.Remove(ctx, phoneLoginFailKey(phoneKey))
	_, _ = phoneLoginCacheMgr.Cache.Remove(ctx, phoneLoginBlockKey(phoneKey))
}

func phoneLoginFailKey(phoneKey string) string {
	return "phone_login_fail:" + phoneKey
}

func phoneLoginBlockKey(phoneKey string) string {
	return "phone_login_block:" + phoneKey
}

func setPhoneLoginCache(key string, data any, ttl time.Duration) {
	_ = phoneLoginCacheMgr.Cache.Set(gctx.New(), key, data, ttl)
}
