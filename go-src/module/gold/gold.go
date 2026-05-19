package gold

import (
	"math"
	"xr-game-server/constants/currency"
	"xr-game-server/core/event"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
)

// 货币精度:统一以 1000 为缩放系数,float 仅用于存储,运算使用 int64 避免精度丢失
const currencyScale = 1000

// scale 将 float 金额按 1000 缩放并四舍五入为 int64
func scale(v float64) int64 {
	return int64(math.Round(v * currencyScale))
}

// unscale 将缩放后的 int64 还原为 float(/1000),并四舍五入到 3 位小数
func unscale(v int64) float64 {
	return math.Round(float64(v)) / currencyScale
}

// Add 给指定用户增加金币,amount 必须为正数,reason 流水原因枚举
func Add(userId uint64, amount float64, reason currency.Reason) (float64, error) {
	if amount <= 0 {
		return 0, errercode.CreateCode(errercode.GoldAmountInvalid)
	}
	amountScaled := scale(amount)
	if amountScaled <= 0 {
		return 0, errercode.CreateCode(errercode.GoldAmountInvalid)
	}
	data := userinfodao.GetUserInfoByUserId(userId)
	beforeScaled := scale(data.Gold)
	// 溢出保护:封顶到 math.MaxInt64,超出部分直接丢弃
	var afterScaled int64
	if amountScaled > math.MaxInt64-beforeScaled {
		afterScaled = math.MaxInt64
		amountScaled = afterScaled - beforeScaled // 实际生效的增量
	} else {
		afterScaled = beforeScaled + amountScaled
	}
	before := unscale(beforeScaled)
	after := unscale(afterScaled)
	data.SetGold(after)
	event.Pub(gameevent.CurrencyChangeEvent, gameevent.NewCurrencyChangeEventData(
		userId, gameevent.CurrencyTypeGold, gameevent.CurrencyActionAdd,
		unscale(amountScaled), before, after, reason,
	))
	return after, nil
}

// Sub 扣减指定用户金币,amount 必须为正数,余额不足返回错误,reason 流水原因枚举
func Sub(userId uint64, amount float64, reason currency.Reason) (float64, error) {
	if amount <= 0 {
		return 0, errercode.CreateCode(errercode.GoldAmountInvalid)
	}
	amountScaled := scale(amount)
	if amountScaled <= 0 {
		return 0, errercode.CreateCode(errercode.GoldAmountInvalid)
	}
	data := userinfodao.GetUserInfoByUserId(userId)
	beforeScaled := scale(data.Gold)
	if beforeScaled < amountScaled {
		return unscale(beforeScaled), errercode.CreateCode(errercode.GoldNotEnough)
	}
	afterScaled := beforeScaled - amountScaled
	before := unscale(beforeScaled)
	after := unscale(afterScaled)
	data.SetGold(after)
	event.Pub(gameevent.CurrencyChangeEvent, gameevent.NewCurrencyChangeEventData(
		userId, gameevent.CurrencyTypeGold, gameevent.CurrencyActionSub,
		unscale(amountScaled), before, after, reason,
	))
	return after, nil
}
