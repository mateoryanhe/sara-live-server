package vip

import (
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/dao/userinfodao"
	entity "xr-game-server/entity/user"
)

// getVipProgressGold 返回用于 VIP 判级/进度的累计充值金币。
// 若新字段尚未写入,则从已完成订单回填一次(兼容上线前按 USD 判级的历史数据)。
func getVipProgressGold(stat *entity.UserCumulativeStat) float64 {
	if stat == nil {
		return 0
	}
	if stat.TotalRechargeGold > 0 {
		return stat.TotalRechargeGold
	}
	gold := rechargeorderdao.SumCompletedRechargeGoldByUserId(stat.ID)
	if gold <= 0 {
		return 0
	}
	stat.SetTotalRechargeGold(gold)
	userinfodao.PublishUserCumulativeStat(stat)
	return gold
}
