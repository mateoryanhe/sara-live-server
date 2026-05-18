package gameevent

import (
	"xr-game-server/constants/currency"
	"xr-game-server/core/event"
)

const (
	// CurrencyChangeEvent 货币变动事件(金币/钻石 加/减)
	CurrencyChangeEvent event.Type = "CurrencyChangeEvent"
)

// 货币类型
const (
	CurrencyTypeGold    uint8 = 1
	CurrencyTypeDiamond uint8 = 2
)

// 变动动作
const (
	CurrencyActionAdd uint8 = 1
	CurrencyActionSub uint8 = 2
)

// CurrencyChangeEventData 货币流水事件数据
type CurrencyChangeEventData struct {
	UserId uint64
	Type   uint8           // 货币类型:1金币 2钻石
	Action uint8           // 动作:1加 2减
	Amount float64         // 变动数量(正数)
	Before float64         // 变动前余额
	After  float64         // 变动后余额
	Reason currency.Reason // 变动原因(枚举)
}

func NewCurrencyChangeEventData(userId uint64, currencyType, action uint8, amount, before, after float64, reason currency.Reason) *CurrencyChangeEventData {
	return &CurrencyChangeEventData{
		UserId: userId,
		Type:   currencyType,
		Action: action,
		Amount: amount,
		Before: before,
		After:  after,
		Reason: reason,
	}
}
