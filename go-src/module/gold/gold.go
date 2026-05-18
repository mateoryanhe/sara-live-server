package gold

import (
	"xr-game-server/core/event"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
)

// Add 给指定用户增加金币,amount 必须为正数,reason 流水原因
func Add(userId uint64, amount float64, reason string) (float64, error) {
	if amount <= 0 {
		return 0, errercode.CreateCode(errercode.GoldAmountInvalid)
	}
	data := userinfodao.GetUserInfoByUserId(userId)
	before := data.Gold
	after := before + amount
	data.SetGold(after)
	event.Pub(gameevent.CurrencyChangeEvent, gameevent.NewCurrencyChangeEventData(
		userId, gameevent.CurrencyTypeGold, gameevent.CurrencyActionAdd,
		amount, before, after, reason,
	))
	return after, nil
}

// Sub 扣减指定用户金币,amount 必须为正数,余额不足返回错误,reason 流水原因
func Sub(userId uint64, amount float64, reason string) (float64, error) {
	if amount <= 0 {
		return 0, errercode.CreateCode(errercode.GoldAmountInvalid)
	}
	data := userinfodao.GetUserInfoByUserId(userId)
	if data.Gold < amount {
		return data.Gold, errercode.CreateCode(errercode.GoldNotEnough)
	}
	before := data.Gold
	after := before - amount
	data.SetGold(after)
	event.Pub(gameevent.CurrencyChangeEvent, gameevent.NewCurrencyChangeEventData(
		userId, gameevent.CurrencyTypeGold, gameevent.CurrencyActionSub,
		amount, before, after, reason,
	))
	return after, nil
}
