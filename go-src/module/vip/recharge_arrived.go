package vip

import (
	"math"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/gameevent"
)

const rechargeCurrencyUSD = "USD"

func onRechargeArrived(val any) {
	data, ok := val.(*gameevent.RechargeArrivedEventData)
	if !ok || data == nil {
		g.Log().Errorf(gctx.New(), "RechargeArrivedEvent payload type error: %T", val)
		return
	}
	if data.UserId == 0 || data.Price == 0 {
		return
	}

	stat := userinfodao.GetUserCumulativeStatByUserId(data.UserId)
	totalRechargeCents := toRechargeCents(stat.TotalRecharge)
	targetLevel := calcTargetVipLevel(totalRechargeCents)
	if targetLevel == 0 {
		return
	}

	user := userinfodao.GetUserInfoByUserId(data.UserId)
	if targetLevel <= user.VipLevel {
		return
	}

	maxLevel := getMaxEnabledVipLevel()
	if maxLevel > 0 && targetLevel > maxLevel {
		targetLevel = maxLevel
	}
	if targetLevel <= user.VipLevel {
		return
	}

	user.SetVipLevel(targetLevel)
	pushVipLevelToApp(data.UserId, targetLevel)
	g.Log().Infof(gctx.New(), "vip upgrade userId=%d level=%d totalRechargeCents=%d",
		data.UserId, targetLevel, totalRechargeCents)
}

// calcTargetVipLevel 根据累计充值(美分)计算应达到的VIP等级
func calcTargetVipLevel(totalRechargeCents uint64) uint32 {
	var target uint32
	for _, cfg := range GetAllVipCfgFromMemory() {
		//if cfg.Status != entity.VipCfgStatusEnabled {
		//	continue
		//}
		if totalRechargeCents >= cfg.UpgradeRechargeLimit && cfg.Level > target {
			target = cfg.Level
		}
	}
	return target
}

// getMaxEnabledVipLevel 获取已开启配置中的最高VIP等级
func getMaxEnabledVipLevel() uint32 {
	var maxLevel uint32
	for _, cfg := range GetAllVipCfgFromMemory() {
		//if cfg.Status != entity.VipCfgStatusEnabled {
		//	continue
		//}
		if cfg.Level > maxLevel {
			maxLevel = cfg.Level
		}
	}
	return maxLevel
}

func toRechargeCents(totalRecharge float64) uint64 {
	if totalRecharge <= 0 {
		return 0
	}
	return uint64(math.Floor(totalRecharge))
}
