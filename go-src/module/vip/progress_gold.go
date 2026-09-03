package vip

import (
	entity "xr-game-server/entity/user"
)

// getVipProgressGold 返回用于 VIP 判级/进度的累计充值金币(仅读 total_recharge_gold,不做订单汇总)。
func getVipProgressGold(stat *entity.UserCumulativeStat) float64 {
	if stat == nil {
		return 0
	}
	return stat.TotalRechargeGold
}
