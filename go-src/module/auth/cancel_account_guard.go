package auth

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/dao/accountdao"
	"xr-game-server/errercode"
)

const (
	appCancelListMaxSize    = 100
	appCancelDailyLimit     = 3
	appCancelDailyKeyLayout = "2006-01-02"
)

var appCancelGuardCacheMgr *cache.CacheMgr

func initAppCancelGuard() {
	appCancelGuardCacheMgr = cache.NewCacheMgr()
}

func checkAppCancelAccountGuard(openId string, channel uint) error {
	openId = accountdao.LogicalOpenId(openId)
	if openId == "" || channel == 0 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	list := accountdao.GetAccountList(openId, channel)
	if len(list) >= appCancelListMaxSize {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	if getAppCancelDailyCount(openId, channel) >= appCancelDailyLimit {
		return errercode.CreateCode(errercode.RequestTooFrequent)
	}
	return nil
}

func recordAppCancelAccountSuccess(openId string, channel uint) {
	openId = accountdao.LogicalOpenId(openId)
	if openId == "" || channel == 0 {
		return
	}
	ctx := gctx.New()
	key := appCancelDailyKey(openId, channel)
	count := 1
	if val, _ := appCancelGuardCacheMgr.Cache.Get(ctx, key); val != nil {
		if current, ok := val.Val().(int); ok && current > 0 {
			count = current + 1
		}
	}
	setAppCancelGuardCache(key, count, ttlUntilNextDay())
}

func getAppCancelDailyCount(openId string, channel uint) int {
	ctx := gctx.New()
	val, _ := appCancelGuardCacheMgr.Cache.Get(ctx, appCancelDailyKey(openId, channel))
	if val == nil {
		return 0
	}
	count, ok := val.Val().(int)
	if !ok || count <= 0 {
		return 0
	}
	return count
}

func appCancelDailyKey(openId string, channel uint) string {
	return fmt.Sprintf("app_cancel_daily:%s:%d:%s", time.Now().Format(appCancelDailyKeyLayout), channel, openId)
}

func ttlUntilNextDay() time.Duration {
	now := time.Now()
	nextDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	ttl := nextDay.Sub(now)
	if ttl <= 0 {
		return 24 * time.Hour
	}
	return ttl
}

func setAppCancelGuardCache(key string, data any, ttl time.Duration) {
	_ = appCancelGuardCacheMgr.Cache.Set(gctx.New(), key, data, ttl)
}
