package wallet

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/currency"
	"xr-game-server/core/event"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/entity/recharge"
	"xr-game-server/gameevent"
)

const defaultTierFirstBonusRatio = 20.0

// onUsdIncomeArrived 订阅美金入账 → 发币 → 发布充值金币到账。
// 普通用户:可走账号首充/档位首充加赠(档位来自 recharge_cfgs)。
// 币商:独立档位表 coin_merchant_recharge_cfgs,无首充/档位首充概念,只按订单金币发放。
func onUsdIncomeArrived(val any) {
	data, ok := val.(*gameevent.UsdIncomeArrivedEventData)
	if !ok || data == nil || data.Order == nil {
		g.Log().Errorf(gctx.New(), "UsdIncomeArrivedEvent payload type error: %T", val)
		return
	}
	order := data.Order
	if order.UserId == 0 || order.Gold <= 0 {
		return
	}
	if order.Status == entity.RechargeOrderStatusCompleted {
		return
	}

	if isCoinMerchantRechargeOrder(order) {
		grantCoinMerchantRechargeGold(data, order)
		return
	}
	grantNormalUserRechargeGold(data, order)
}

func isCoinMerchantRechargeOrder(order *entity.RechargeOrder) bool {
	return order != nil && order.PayChannel == entity.RechargeCfgTypeCoinMerchant
}

// grantCoinMerchantRechargeGold 币商充值:无账号首充/档位首充,不写首充标记,不加赠
func grantCoinMerchantRechargeGold(data *gameevent.UsdIncomeArrivedEventData, order *entity.RechargeOrder) {
	baseGold := order.Gold
	logReason := data.Reason
	if logReason == currency.ReasonRecharge || logReason == 0 {
		logReason = currency.ReasonRechargeCoinMerchant
	}
	after, err := GoldAdd(order.UserId, baseGold, logReason)
	if err != nil {
		g.Log().Errorf(gctx.New(), "coin merchant recharge gold grant failed order=%d user=%d gold=%v err=%v",
			order.ID, order.UserId, baseGold, err)
		return
	}
	markOrderPaid(order, baseGold)
	addUserRechargeGoldStat(order.UserId, baseGold)
	fillGranted(data, baseGold, baseGold, after, false, false)
	publishRechargeGoldEvents(data, order, baseGold, baseGold, after, false, false)
}

// grantNormalUserRechargeGold 普通用户充值:可档位首充加赠,并标记账号/档位首充
func grantNormalUserRechargeGold(data *gameevent.UsdIncomeArrivedEventData, order *entity.RechargeOrder) {
	baseGold := order.Gold
	isTierFirst := order.CfgId > 0 && userinfodao.IsRechargeCfgFirstRecharge(order.UserId, order.CfgId)

	creditedGold := baseGold
	if isTierFirst {
		creditedGold = applyTierFirstBonus(baseGold)
	}

	logReason := data.Reason
	if logReason == currency.ReasonRecharge || logReason == 0 {
		logReason = resolveNormalRechargeGoldReason(order, isTierFirst)
	}

	after, err := GoldAdd(order.UserId, creditedGold, logReason)
	if err != nil {
		g.Log().Errorf(gctx.New(), "recharge gold grant failed order=%d user=%d gold=%v err=%v",
			order.ID, order.UserId, creditedGold, err)
		return
	}

	markOrderPaid(order, creditedGold)

	isTierFirstDone := false
	isAccountFirstDone := false
	// 仅普通充值档位(recharge_cfgs);币商档位 ID 空间独立,禁止写入本表
	if order.CfgId > 0 && userinfodao.MarkRechargeCfgFirstRechargeDone(order.UserId, order.CfgId) {
		isTierFirstDone = true
		if userinfodao.MarkFirstRechargeDone(order.UserId) {
			isAccountFirstDone = true
		}
	}

	fillGranted(data, baseGold, creditedGold, after, isTierFirstDone, isAccountFirstDone)
	addUserRechargeGoldStat(order.UserId, creditedGold)
	publishRechargeGoldEvents(data, order, baseGold, creditedGold, after, isTierFirstDone, isAccountFirstDone)
}

func markOrderPaid(order *entity.RechargeOrder, creditedGold float64) {
	if creditedGold != order.Gold {
		order.SetGold(creditedGold)
	}
	paidAt := time.Now()
	order.SetStatus(entity.RechargeOrderStatusCompleted)
	order.SetPaidAt(paidAt)
	order.SetUpdatedAt(paidAt)
}

func fillGranted(
	data *gameevent.UsdIncomeArrivedEventData,
	baseGold, creditedGold, after float64,
	isTierFirst, isAccountFirst bool,
) {
	data.Granted = true
	data.BaseGold = baseGold
	data.CreditedGold = creditedGold
	data.GoldBalance = after
	data.IsTierFirst = isTierFirst
	data.IsAccountFirst = isAccountFirst
}

// addUserRechargeGoldStat 发币成功后写入用户累计充值金币/次数(须在发布金币到账事件之前)
func addUserRechargeGoldStat(userId uint64, creditedGold float64) {
	if userId == 0 || creditedGold <= 0 {
		return
	}
	stat := userinfodao.GetUserCumulativeStatByUserId(userId)
	stat.AddTotalPayCount(1)
	stat.AddTotalRechargeGold(creditedGold)
	userinfodao.PublishUserCumulativeStat(stat)
}

func publishRechargeGoldEvents(
	data *gameevent.UsdIncomeArrivedEventData,
	order *entity.RechargeOrder,
	baseGold, creditedGold, after float64,
	isTierFirst, isAccountFirst bool,
) {
	goldEvt := gameevent.NewRechargeGoldArrivedEventData(
		order, data.Kind, baseGold, creditedGold, after, isTierFirst, isAccountFirst,
	)
	event.Pub(gameevent.RechargeGoldArrivedEvent, goldEvt)
}

// applyTierFirstBonus 普通用户档位首充加赠(读首充活动配置;币商不走此逻辑)
func applyTierFirstBonus(baseGold float64) float64 {
	if baseGold <= 0 {
		return baseGold
	}
	ratio := defaultTierFirstBonusRatio
	if row := cfgdao.LoadFirstRechargeActivityCfg(); row != nil {
		ratio = row.FirstRechargeRatio
	}
	if ratio < 0 {
		ratio = 0
	}
	if ratio <= 0 {
		return baseGold
	}
	return baseGold * (1 + ratio/100)
}

func resolveNormalRechargeGoldReason(order *entity.RechargeOrder, isTierFirst bool) currency.Reason {
	if order == nil {
		return currency.ReasonRecharge
	}
	if order.Source == entity.RechargeOrderSourceManual {
		return currency.ReasonRechargeManual
	}
	if userinfodao.IsRechargeWhitelist(order.UserId) {
		return currency.ReasonRechargeWhitelist
	}
	if order.CfgId == 0 {
		return currency.ReasonRechargeCustom
	}
	if isTierFirst {
		return currency.ReasonRechargeFirstBonus
	}
	switch order.PayChannel {
	case entity.RechargeCfgTypeGoogle:
		return currency.ReasonRechargeGoogle
	case entity.RechargeCfgTypeIOS:
		return currency.ReasonRechargeIOS
	case entity.RechargeCfgTypeChannel:
		return currency.ReasonRechargeChannel
	default:
		return currency.ReasonRechargeTier
	}
}
