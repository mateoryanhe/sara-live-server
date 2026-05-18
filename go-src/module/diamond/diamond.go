package diamond

import (
	"xr-game-server/constants/currency"
	"xr-game-server/core/event"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
)

// Add 给指定用户增加钻石,amount 必须为正数,reason 流水原因枚举
func Add(userId uint64, amount float64, reason currency.Reason) (float64, error) {
	if amount <= 0 {
		return 0, errercode.CreateCode(errercode.DiamondAmountInvalid)
	}
	data := userinfodao.GetUserInfoByUserId(userId)
	before := data.Diamond
	after := before + amount
	data.SetDiamond(after)
	event.Pub(gameevent.CurrencyChangeEvent, gameevent.NewCurrencyChangeEventData(
		userId, gameevent.CurrencyTypeDiamond, gameevent.CurrencyActionAdd,
		amount, before, after, reason,
	))
	return after, nil
}

// Sub 扣减指定用户钻石,amount 必须为正数,余额不足返回错误,reason 流水原因枚举
func Sub(userId uint64, amount float64, reason currency.Reason) (float64, error) {
	if amount <= 0 {
		return 0, errercode.CreateCode(errercode.DiamondAmountInvalid)
	}
	data := userinfodao.GetUserInfoByUserId(userId)
	if data.Diamond < amount {
		return data.Diamond, errercode.CreateCode(errercode.DiamondNotEnough)
	}
	before := data.Diamond
	after := before - amount
	data.SetDiamond(after)
	event.Pub(gameevent.CurrencyChangeEvent, gameevent.NewCurrencyChangeEventData(
		userId, gameevent.CurrencyTypeDiamond, gameevent.CurrencyActionSub,
		amount, before, after, reason,
	))
	return after, nil
}
