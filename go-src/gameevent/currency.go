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

// CurrencyChangeMeta 货币流水附加信息(如游戏消费).
type CurrencyChangeMeta struct {
	GameName     string
	GameCategory string
	BusinessType uint8 // 商业类型:1社交 2游戏,默认社交
}

// CurrencyChangeEventData 货币流水事件数据
type CurrencyChangeEventData struct {
	UserId       uint64
	Type         uint8           // 货币类型:1金币 2钻石
	Action       uint8           // 动作:1加 2减
	Amount       float64         // 变动数量(正数)
	Before       float64         // 变动前余额
	After        float64         // 变动后余额
	Reason       currency.Reason // 变动原因(枚举)
	GameName     string          // 游戏名称(可选)
	GameCategory string          // 游戏分类(可选)
	BusinessType uint8           // 商业类型:1社交 2游戏
}

func NewCurrencyChangeEventData(userId uint64, currencyType, action uint8, amount, before, after float64, reason currency.Reason, meta ...*CurrencyChangeMeta) *CurrencyChangeEventData {
	data := &CurrencyChangeEventData{
		UserId:       userId,
		Type:         currencyType,
		Action:       action,
		Amount:       amount,
		Before:       before,
		After:        after,
		Reason:       reason,
		BusinessType: currency.BusinessTypeSocial,
	}
	if len(meta) > 0 && meta[0] != nil {
		data.GameName = meta[0].GameName
		data.GameCategory = meta[0].GameCategory
		if meta[0].BusinessType != 0 {
			data.BusinessType = meta[0].BusinessType
		}
	}
	return data
}
