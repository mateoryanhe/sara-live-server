package recharge

import (
	"xr-game-server/dao/userinfodao"
	"xr-game-server/entity/recharge"
	"xr-game-server/gameevent"
)

// resolveUsdIncomeKind 解析美金入账类型。白名单优先于渠道/人工。
func resolveUsdIncomeKind(order *entity.RechargeOrder) gameevent.UsdIncomeKind {
	if order == nil || order.UserId == 0 {
		return gameevent.UsdIncomeKindUnknown
	}
	if userinfodao.IsRechargeWhitelist(order.UserId) {
		return gameevent.UsdIncomeKindWhitelist
	}
	if order.Source == entity.RechargeOrderSourceManual {
		return gameevent.UsdIncomeKindManual
	}
	switch order.PayChannel {
	case entity.RechargeCfgTypeGoogle:
		return gameevent.UsdIncomeKindGoogle
	case entity.RechargeCfgTypeIOS:
		return gameevent.UsdIncomeKindIOS
	case entity.RechargeCfgTypeChannel:
		return gameevent.UsdIncomeKindChannel
	case entity.RechargeCfgTypeCoinMerchant:
		return gameevent.UsdIncomeKindCoinMerchant
	default:
		return gameevent.UsdIncomeKindUnknown
	}
}
