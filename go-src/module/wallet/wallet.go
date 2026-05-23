package wallet

import (
	"xr-game-server/constants/currency"
	"xr-game-server/core/event"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
)

// DiamondAdd 给指定用户增加钻石,amount 必须为正数,reason 流水原因枚举
func DiamondAdd(userId uint64, amount float64, reason currency.Reason) (float64, error) {
	if amount <= 0 {
		return 0, errercode.CreateCode(errercode.DiamondAmountInvalid)
	}

	data := userinfodao.GetUserInfoByUserId(userId)

	before := data.Diamond
	data.AddDiamond(amount)
	after := data.Diamond
	event.Pub(gameevent.CurrencyChangeEvent, gameevent.NewCurrencyChangeEventData(
		userId, gameevent.CurrencyTypeDiamond, gameevent.CurrencyActionAdd,
		amount, before, after, reason,
	))
	pushDiamondToApp(userId, after)
	return after, nil
}

// DiamondSub 扣减指定用户钻石,amount 必须为正数,余额不足返回错误,reason 流水原因枚举
func DiamondSub(userId uint64, amount float64, reason currency.Reason) (float64, error) {
	if amount <= 0 {
		return 0, errercode.CreateCode(errercode.DiamondAmountInvalid)
	}

	data := userinfodao.GetUserInfoByUserId(userId)
	if amount > data.Diamond {
		return 0, errercode.CreateCode(errercode.DiamondNotEnough)
	}
	before := data.Diamond
	data.SubDiamond(amount)
	after := data.Diamond
	event.Pub(gameevent.CurrencyChangeEvent, gameevent.NewCurrencyChangeEventData(
		userId, gameevent.CurrencyTypeDiamond, gameevent.CurrencyActionSub,
		amount, before, after, reason,
	))
	pushDiamondToApp(userId, after)
	return after, nil
}

// GoldAdd 给指定用户增加金币,amount 必须为正数,reason 流水原因枚举
func GoldAdd(userId uint64, amount float64, reason currency.Reason) (float64, error) {
	if amount <= 0 {
		return 0, errercode.CreateCode(errercode.GoldAmountInvalid)
	}
	data := userinfodao.GetUserInfoByUserId(userId)

	before := data.Gold
	data.AddGold(amount)
	after := data.Gold
	event.Pub(gameevent.CurrencyChangeEvent, gameevent.NewCurrencyChangeEventData(
		userId, gameevent.CurrencyTypeGold, gameevent.CurrencyActionAdd,
		amount, before, after, reason,
	))
	pushGoldToApp(userId, after)
	return after, nil
}

// GoldSub 扣减指定用户金币,amount 必须为正数,余额不足返回错误,reason 流水原因枚举
func GoldSub(userId uint64, amount float64, reason currency.Reason) (float64, error) {
	if amount <= 0 {
		return 0, errercode.CreateCode(errercode.GoldAmountInvalid)
	}
	data := userinfodao.GetUserInfoByUserId(userId)
	if amount > data.Gold {
		return 0, errercode.CreateCode(errercode.GoldNotEnough)
	}
	before := data.Gold
	data.SubGold(amount)
	after := data.Gold
	event.Pub(gameevent.CurrencyChangeEvent, gameevent.NewCurrencyChangeEventData(
		userId, gameevent.CurrencyTypeGold, gameevent.CurrencyActionSub,
		amount, before, after, reason,
	))
	pushGoldToApp(userId, after)
	return after, nil
}
