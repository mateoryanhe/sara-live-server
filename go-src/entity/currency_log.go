package entity

import (
	"time"
	"xr-game-server/constants/currency"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbCurrencyLog db.TbName = "currency_logs"
)

const (
	CurrencyLogUserId db.TbCol = "user_id"
	CurrencyLogType   db.TbCol = "type"
	CurrencyLogAction db.TbCol = "action"
	CurrencyLogAmount db.TbCol = "amount"
	CurrencyLogBefore db.TbCol = "before"
	CurrencyLogAfter  db.TbCol = "after"
	CurrencyLogReason db.TbCol = "reason"
)

// CurrencyLog 金币/钻石流水记录
type CurrencyLog struct {
	migrate.OneModel
	UserId uint64  `gorm:"index;default:0;comment:用户ID"`
	Type   uint8   `gorm:"default:0;comment:货币类型 1金币 2钻石"`
	Action uint8   `gorm:"default:0;comment:动作 1加 2减"`
	Amount float64 `gorm:"default:0;comment:变动数量"`
	Before float64 `gorm:"default:0;comment:变动前余额"`
	After  float64 `gorm:"default:0;comment:变动后余额"`
	Reason uint8   `gorm:"default:0;comment:变动原因(枚举,参见 constants/currency.Reason)"`
}

func NewCurrencyLog(userId uint64, currencyType, action uint8, amount, before, after float64, reason currency.Reason) *CurrencyLog {
	ret := &CurrencyLog{}
	ret.ID = snowflake.GetId()
	ret.SetCreatedAt(time.Now())
	ret.SetUpdatedAt(time.Now())
	ret.SetUserId(userId)
	ret.SetType(currencyType)
	ret.SetAction(action)
	ret.SetAmount(amount)
	ret.SetBefore(before)
	ret.SetAfter(after)
	ret.SetReason(reason)
	return ret
}

func (receiver *CurrencyLog) SetUserId(userId uint64) {
	receiver.UserId = userId
	syndb.AddDataToQuickChan(TbCurrencyLog, CurrencyLogUserId, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: userId,
	})
}

func (receiver *CurrencyLog) SetType(currencyType uint8) {
	receiver.Type = currencyType
	syndb.AddDataToQuickChan(TbCurrencyLog, CurrencyLogType, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: currencyType,
	})
}

func (receiver *CurrencyLog) SetAction(action uint8) {
	receiver.Action = action
	syndb.AddDataToQuickChan(TbCurrencyLog, CurrencyLogAction, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: action,
	})
}

func (receiver *CurrencyLog) SetAmount(amount float64) {
	receiver.Amount = amount
	syndb.AddDataToQuickChan(TbCurrencyLog, CurrencyLogAmount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: amount,
	})
}

func (receiver *CurrencyLog) SetBefore(before float64) {
	receiver.Before = before
	syndb.AddDataToQuickChan(TbCurrencyLog, CurrencyLogBefore, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: before,
	})
}

func (receiver *CurrencyLog) SetAfter(after float64) {
	receiver.After = after
	syndb.AddDataToQuickChan(TbCurrencyLog, CurrencyLogAfter, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: after,
	})
}

func (receiver *CurrencyLog) SetReason(reason currency.Reason) {
	receiver.Reason = uint8(reason)
	syndb.AddDataToQuickChan(TbCurrencyLog, CurrencyLogReason, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: uint8(reason),
	})
}

func (receiver *CurrencyLog) SetCreatedAt(val time.Time) {
	receiver.CreatedAt = val
	syndb.AddDataToQuickChan(TbCurrencyLog, db.CreatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *CurrencyLog) SetUpdatedAt(val time.Time) {
	receiver.UpdatedAt = val
	syndb.AddDataToQuickChan(TbCurrencyLog, db.UpdatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func initCurrencyLog() {
	syndb.RegQuickWithLarge(TbCurrencyLog, db.CreatedAtName)
	syndb.RegQuickWithLarge(TbCurrencyLog, db.UpdatedAtName)

	syndb.RegQuickWithLarge(TbCurrencyLog, CurrencyLogUserId)
	syndb.RegQuickWithLarge(TbCurrencyLog, CurrencyLogType)
	syndb.RegQuickWithLarge(TbCurrencyLog, CurrencyLogAction)
	syndb.RegQuickWithLarge(TbCurrencyLog, CurrencyLogAmount)
	syndb.RegQuickWithLarge(TbCurrencyLog, CurrencyLogBefore)
	syndb.RegQuickWithLarge(TbCurrencyLog, CurrencyLogAfter)
	syndb.RegQuickWithLarge(TbCurrencyLog, CurrencyLogReason)

	migrate.AutoMigrate(&CurrencyLog{})
}
