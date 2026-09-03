package vip

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/entity/recharge"
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
	if order.Gold > 0 {
		stat.AddTotalRechargeGold(order.Gold)
	}
	userinfodao.PublishUserCumulativeStat(stat)

	totalRechargeGold := stat.TotalRechargeGold

	targetLevel := calcTargetVipLevel(totalRechargeGold)
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
	userinfodao.PublishUserInfo(user)
	pushVipLevelToApp(order.UserId, targetLevel)

}

// calcTargetVipLevel 根据累计充值到账金币计算应达到的VIP等级。
// 配置中每级 UpgradeRechargeLimit 为该等级累计充值金币上限(如 L1=1000,L2=5000):
// 累计金币 < L1上限 → L1; 累计金币 >= 最高级上限 → 最高等级。
func calcTargetVipLevel(totalRechargeGold float64) uint32 {
	if totalRechargeGold <= 0 {
		return 0
	}
	cfgs := GetAllVipCfgFromMemory()
	if len(cfgs) == 0 {
		return 0
	}
	for _, cfg := range cfgs {
		if cfg.UpgradeRechargeLimit > totalRechargeGold {
			return cfg.Level
		}
	}
	return cfgs[len(cfgs)-1].Level
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
