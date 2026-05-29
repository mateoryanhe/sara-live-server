package vip

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/entity"
)

func onRechargeArrived(val any) {
	data, ok := val.(*entity.RechargeOrder)
	if !ok || data == nil {
		g.Log().Errorf(gctx.New(), "RechargeArrivedEvent payload type error: %T", val)
		return
	}
	order := data

	//充值成功完成
	stat := userinfodao.GetUserCumulativeStatByUserId(order.UserId)
	stat.AddTotalRecharge(order.Price)
	stat.AddTotalPayCount(1)

	totalRecharge := stat.TotalRecharge

	targetLevel := calcTargetVipLevel(totalRecharge)
	if targetLevel == 0 {
		return
	}

	user := userinfodao.GetUserInfoByUserId(order.UserId)
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
	pushVipLevelToApp(order.UserId, targetLevel)
	g.Log().Infof(gctx.New(), "vip upgrade userId=%d level=%d totalRecharge=%.4f",
		order.UserId, targetLevel, totalRecharge)
}

// calcTargetVipLevel 根据累计充值(USD)计算应达到的VIP等级
func calcTargetVipLevel(totalRecharge float64) uint32 {
	var target uint32
	for _, cfg := range GetAllVipCfgFromMemory() {
		if totalRecharge >= cfg.UpgradeRechargeLimit && cfg.Level > target {
			target = cfg.Level
		}
	}
	return target
}

// getMaxEnabledVipLevel 获取已开启配置中的最高VIP等级
func getMaxEnabledVipLevel() uint32 {
	var maxLevel uint32
	for _, cfg := range GetAllVipCfgFromMemory() {
		if cfg.Level > maxLevel {
			maxLevel = cfg.Level
		}
	}
	return maxLevel
}
