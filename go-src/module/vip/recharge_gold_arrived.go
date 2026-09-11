package vip

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/entity/recharge"
	userentity "xr-game-server/entity/user"
	"xr-game-server/gameevent"
)

// onRechargeGoldArrived 普通用户充值金币到账后只读累计充值金币并尝试 VIP 升级。
// TotalRechargeGold 由 wallet 发币时写入;币商无 VIP,跳过。
func onRechargeGoldArrived(val any) {
	data, ok := val.(*gameevent.RechargeGoldArrivedEventData)
	if !ok || data == nil || data.Order == nil {
		g.Log().Errorf(gctx.New(), "RechargeGoldArrivedEvent payload type error: %T", val)
		return
	}
	order := data.Order
	if order.UserId == 0 {
		return
	}
	if isCoinMerchantVipSkip(data.Kind, order) {
		return
	}

	stat := userinfodao.GetUserCumulativeStatByUserId(order.UserId)
	targetLevel := calcTargetVipLevel(stat.TotalRechargeGold)
	if targetLevel == 0 {
		return
	}

	user := userinfodao.GetUserInfoByUserId(order.UserId)
	if user == nil || user.UserType == userentity.UserTypeCoinMerchant {
		return
	}
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

func isCoinMerchantVipSkip(kind gameevent.UsdIncomeKind, order *entity.RechargeOrder) bool {
	if kind == gameevent.UsdIncomeKindCoinMerchant {
		return true
	}
	return order != nil && order.PayChannel == entity.RechargeCfgTypeCoinMerchant
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
