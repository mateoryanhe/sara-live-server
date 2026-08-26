package wallet

import (
	"errors"
	"math"

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

	var before, after float64
	data := userinfodao.GetUserInfoByUserId(userId)
	before = data.Diamond
	data.AddDiamond(amount)
	after = data.Diamond
	userinfodao.PublishUserInfo(data)
	if reason == currency.ReasonRefund {
		stat := userinfodao.GetUserCumulativeStatByUserId(userId)
		stat.AddTotalDiamondConsume(-amount)
		userinfodao.PublishUserCumulativeStat(stat)
	}
	event.Pub(gameevent.CurrencyChangeEvent, gameevent.NewCurrencyChangeEventData(
		userId, gameevent.CurrencyTypeDiamond, gameevent.CurrencyActionAdd,
		amount, before, after, reason,
	))
	pushDiamondToApp(userId, after)
	return after, nil
}

func DiamondNotEnough(userId uint64, amount float64) error {
	data := userinfodao.GetUserInfoByUserId(userId)
	if data.Diamond >= amount {
		return nil
	}
	return errercode.CreateCode(errercode.DiamondNotEnough)
}

// CanPayWithGoldExchange 校验应付金额是否可由钻石+按需兑换金币覆盖
func CanPayWithGoldExchange(userId uint64, amount float64) error {
	if amount <= 0 {
		return nil
	}
	data := userinfodao.GetUserInfoByUserId(userId)
	if data.Diamond >= amount {
		return nil
	}
	shortfall := amount - data.Diamond
	goldNeeded := calcGoldForDiamondShortfall(shortfall)
	if data.Gold >= goldNeeded {
		return nil
	}
	return errercode.CreateCode(errercode.DiamondNotEnough)
}

// DiamondSubWithGoldExchange 扣减钻石;余额不足时按缺口自动用金币兑换钻石(仅兑换所需数量)
func DiamondSubWithGoldExchange(userId uint64, amount float64, reason currency.Reason) (float64, error) {
	after, err := DiamondSub(userId, amount, reason)
	if err == nil {
		return after, nil
	}
	if !isDiamondNotEnough(err) {
		return 0, err
	}

	data := userinfodao.GetUserInfoByUserId(userId)
	shortfall := amount - data.Diamond
	if shortfall <= 0 {
		return DiamondSub(userId, amount, reason)
	}

	goldNeeded := calcGoldForDiamondShortfall(shortfall)
	if data.Gold < goldNeeded {
		return 0, errercode.CreateCode(errercode.DiamondNotEnough)
	}
	if _, _, err = exchangeGoldToDiamondForAutoPay(userId, goldNeeded); err != nil {
		return 0, err
	}
	return DiamondSub(userId, amount, reason)
}

func calcGoldForDiamondShortfall(shortfall float64) float64 {
	if shortfall <= 0 {
		return 0
	}
	return math.Ceil(shortfall / float64(GetGoldToDiamondRate()))
}

func isDiamondNotEnough(err error) bool {
	var bizErr *errercode.XError
	return errors.As(err, &bizErr) && bizErr.Code() == errercode.DiamondNotEnough
}

// DiamondSub 扣减指定用户钻石,amount 必须为正数,余额不足返回错误,reason 流水原因枚举
func DiamondSub(userId uint64, amount float64, reason currency.Reason) (float64, error) {
	if amount <= 0 {
		return 0, errercode.CreateCode(errercode.DiamondAmountInvalid)
	}

	data := userinfodao.GetUserInfoByUserId(userId)
	if amount > data.Diamond {
		pushDiamondToApp(userId, data.Diamond)
		return 0, errercode.CreateCode(errercode.DiamondNotEnough)
	}

	before := data.Diamond
	data.SubDiamond(amount)
	after := data.Diamond
	userinfodao.PublishUserInfo(data)
	stat := userinfodao.GetUserCumulativeStatByUserId(userId)
	stat.AddTotalDiamondConsume(amount)
	userinfodao.PublishUserCumulativeStat(stat)

	event.Pub(gameevent.CurrencyChangeEvent, gameevent.NewCurrencyChangeEventData(
		userId, gameevent.CurrencyTypeDiamond, gameevent.CurrencyActionSub,
		amount, before, after, reason,
	))
	pushDiamondToApp(userId, after)
	return after, nil
}

// GoldAdd 给指定用户增加金币,amount 必须为正数,reason 流水原因枚举
func GoldAdd(userId uint64, amount float64, reason currency.Reason, meta ...*gameevent.CurrencyChangeMeta) (float64, error) {
	if amount <= 0 {
		return 0, errercode.CreateCode(errercode.GoldAmountInvalid)
	}

	var before, after float64
	data := userinfodao.GetUserInfoByUserId(userId)
	before = data.Gold
	data.AddGold(amount)
	after = data.Gold
	userinfodao.PublishUserInfo(data)
	event.Pub(gameevent.CurrencyChangeEvent, gameevent.NewCurrencyChangeEventData(
		userId, gameevent.CurrencyTypeGold, gameevent.CurrencyActionAdd,
		amount, before, after, reason, meta...,
	))
	pushGoldToApp(userId, after)
	return after, nil
}

// GoldSub 扣减指定用户金币,amount 必须为正数,余额不足返回错误,reason 流水原因枚举
func GoldSub(userId uint64, amount float64, reason currency.Reason, meta ...*gameevent.CurrencyChangeMeta) (float64, error) {
	if amount <= 0 {
		return 0, errercode.CreateCode(errercode.GoldAmountInvalid)
	}
	data := userinfodao.GetUserInfoByUserId(userId)
	if amount > data.Gold {
		pushGoldToApp(userId, data.Gold)
		return 0, errercode.CreateCode(errercode.GoldNotEnough)
	}

	before := data.Gold
	data.SubGold(amount)
	after := data.Gold
	userinfodao.PublishUserInfo(data)
	stat := userinfodao.GetUserCumulativeStatByUserId(userId)
	stat.AddTotalGoldConsume(amount)
	userinfodao.PublishUserCumulativeStat(stat)

	event.Pub(gameevent.CurrencyChangeEvent, gameevent.NewCurrencyChangeEventData(
		userId, gameevent.CurrencyTypeGold, gameevent.CurrencyActionSub,
		amount, before, after, reason, meta...,
	))
	pushGoldToApp(userId, after)
	return after, nil
}
