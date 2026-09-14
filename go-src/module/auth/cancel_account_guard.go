package auth

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/dao/accountdao"
	"xr-game-server/errercode"
	"xr-game-server/module/accountcfg"
)

const appCancelDailyKeyLayout = "2006-01-02"

var appCancelGuardCacheMgr *cache.RowCache[int]

func initAppCancelGuard() {
	appCancelGuardCacheMgr = cache.NewRowCache[int]()
}

// checkAppCancelAccountGuard 注销侧:风控开启时校验每日注销次数
func checkAppCancelAccountGuard(openId string, channel uint) error {
	openId = accountdao.LogicalOpenId(openId)
	if openId == "" || channel == 0 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	if !accountcfg.IsDeviceRegisterRiskEnabled() {
		return nil
	}
	limit := accountcfg.GetDeviceCancelDailyLimit()
	if limit <= 0 {
		return nil
	}
	if getAppCancelDailyCount(openId, channel) >= limit {
		return errercode.CreateCode(errercode.RequestTooFrequent)
	}
	return nil
}

// checkDeviceRegisterAccountLimit 设备码新注册前:风控开启时校验同设备账号数上限(含已注销)
func checkDeviceRegisterAccountLimit(openId string, channel uint) error {
	openId = accountdao.LogicalOpenId(openId)
	if openId == "" || channel == 0 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	if !accountcfg.IsDeviceRegisterRiskEnabled() {
		return nil
	}
	maxCount := accountcfg.GetDeviceAccountMaxCount()
	if maxCount <= 0 {
		return nil
	}
	list := accountdao.GetAccountList(openId, channel)
	if len(list) >= maxCount {
		return errercode.CreateCode(errercode.DeviceAccountRegisterLimit)
	}
	return nil
}

func recordAppCancelAccountSuccess(openId string, channel uint) {
	openId = accountdao.LogicalOpenId(openId)
	if openId == "" || channel == 0 {
		return
	}
	if !accountcfg.IsDeviceRegisterRiskEnabled() {
		return
	}
	ctx := gctx.New()
	key := appCancelDailyKey(openId, channel)
	count := 1
	if current, ok := appCancelGuardCacheMgr.GetRowCached(ctx, key); ok && current > 0 {
		count = current + 1
	}
	setAppCancelGuardCache(key, count, ttlUntilNextDay())
}

func getAppCancelDailyCount(openId string, channel uint) int {
	ctx := gctx.New()
	count, ok := appCancelGuardCacheMgr.GetRowCached(ctx, appCancelDailyKey(openId, channel))
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

func setAppCancelGuardCache(key string, data int, ttl time.Duration) {
	_ = appCancelGuardCacheMgr.SetRow(gctx.New(), key, data, ttl)
}
