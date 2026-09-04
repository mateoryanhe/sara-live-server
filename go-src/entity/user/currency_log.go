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
	CurrencyLogUserId       db.TbCol = "user_id"
	CurrencyLogType         db.TbCol = "type"
	CurrencyLogAction       db.TbCol = "action"
	CurrencyLogAmount       db.TbCol = "amount"
	CurrencyLogBefore       db.TbCol = "before"
	CurrencyLogAfter        db.TbCol = "after"
	CurrencyLogReason       db.TbCol = "reason"
	CurrencyLogGameId         db.TbCol = "game_id"
	CurrencyLogGameName       db.TbCol = "game_name"
	CurrencyLogGameCategory   db.TbCol = "game_category"
	CurrencyLogBusinessType   db.TbCol = "business_type"
	CurrencyLogTransactionId  db.TbCol = "transaction_id"
)

// CurrencyLog 金币/钻石流水记录
type CurrencyLog struct {
	migrate.OneModel
	// CreatedAt 覆盖嵌入字段以挂复合索引(富豪榜 type+action+时间; 用户流水 user+type+时间)
	CreatedAt      time.Time `gorm:"index:idx_cl_type_action_created,priority:3;index:idx_cl_user_type_created,priority:3"`
	UserId         uint64    `gorm:"index:idx_cl_user_type_created,priority:1;default:0;comment:用户ID"`
	Type           uint8     `gorm:"index:idx_cl_type_action_created,priority:1;index:idx_cl_user_type_created,priority:2;default:0;comment:货币类型 1金币 2钻石"`
	Action         uint8     `gorm:"index:idx_cl_type_action_created,priority:2;default:0;comment:动作 1加 2减"`
	Amount         float64   `gorm:"default:0;comment:变动数量"`
	Before         float64   `gorm:"default:0;comment:变动前余额"`
	After          float64   `gorm:"default:0;comment:变动后余额"`
	Reason         uint8     `gorm:"default:0;comment:变动原因(枚举,参见 constants/currency.Reason)"`
	GameId         string    `gorm:"size:64;default:'';comment:游戏ID(第三方gameId/游戏编码)"`
	GameName       string    `gorm:"size:128;default:'';comment:游戏名称"`
	GameCategory   string    `gorm:"size:64;default:'';comment:游戏分类"`
	BusinessType   uint8     `gorm:"default:1;comment:商业类型 1社交 2游戏"`
	TransactionId  string    `gorm:"size:64;default:'';index;comment:第三方交易ID"`
}

func NewCurrencyLog(userId uint64, currencyType, action uint8, amount, before, after float64, reason currency.Reason, gameId, gameName, gameCategory string, businessType uint8, transactionId string) *CurrencyLog {
	if businessType == 0 {
		businessType = currency.BusinessTypeSocial
	}
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
	ret.SetGameId(gameId)
	ret.SetGameName(gameName)
	ret.SetGameCategory(gameCategory)
	ret.SetBusinessType(businessType)
	ret.SetTransactionId(transactionId)
	return ret
}

func (receiver *CurrencyLog) SetUserId(userId uint64) {
	receiver.UserId = userId
	syndb.AddData(TbCurrencyLog, CurrencyLogUserId, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: userId,
	})
}

func (receiver *CurrencyLog) SetType(currencyType uint8) {
	receiver.Type = currencyType
	syndb.AddData(TbCurrencyLog, CurrencyLogType, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: currencyType,
	})
}

func (receiver *CurrencyLog) SetAction(action uint8) {
	receiver.Action = action
	syndb.AddData(TbCurrencyLog, CurrencyLogAction, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: action,
	})
}

func (receiver *CurrencyLog) SetAmount(amount float64) {
	receiver.Amount = amount
	syndb.AddData(TbCurrencyLog, CurrencyLogAmount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: amount,
	})
}

func (receiver *CurrencyLog) SetBefore(before float64) {
	receiver.Before = before
	syndb.AddData(TbCurrencyLog, CurrencyLogBefore, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: before,
	})
}

func (receiver *CurrencyLog) SetAfter(after float64) {
	receiver.After = after
	syndb.AddData(TbCurrencyLog, CurrencyLogAfter, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: after,
	})
}

func (receiver *CurrencyLog) SetReason(reason currency.Reason) {
	receiver.Reason = uint8(reason)
	syndb.AddData(TbCurrencyLog, CurrencyLogReason, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: uint8(reason),
	})
}

func (receiver *CurrencyLog) SetGameId(gameId string) {
	receiver.GameId = gameId
	syndb.AddData(TbCurrencyLog, CurrencyLogGameId, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: gameId,
	})
}

func (receiver *CurrencyLog) SetGameName(gameName string) {
	receiver.GameName = gameName
	syndb.AddData(TbCurrencyLog, CurrencyLogGameName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: gameName,
	})
}

func (receiver *CurrencyLog) SetGameCategory(gameCategory string) {
	receiver.GameCategory = gameCategory
	syndb.AddData(TbCurrencyLog, CurrencyLogGameCategory, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: gameCategory,
	})
}

func (receiver *CurrencyLog) SetBusinessType(businessType uint8) {
	if businessType == 0 {
		businessType = currency.BusinessTypeSocial
	}
	receiver.BusinessType = businessType
	syndb.AddData(TbCurrencyLog, CurrencyLogBusinessType, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: businessType,
	})
}

func (receiver *CurrencyLog) SetTransactionId(transactionId string) {
	receiver.TransactionId = transactionId
	syndb.AddData(TbCurrencyLog, CurrencyLogTransactionId, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: transactionId,
	})
}

func (receiver *CurrencyLog) SetCreatedAt(val time.Time) {
	receiver.CreatedAt = val
	syndb.AddData(TbCurrencyLog, db.CreatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *CurrencyLog) SetUpdatedAt(val time.Time) {
	receiver.UpdatedAt = val
	syndb.AddData(TbCurrencyLog, db.UpdatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func initCurrencyLog() {
	syndb.RegQuick(TbCurrencyLog, db.CreatedAtName)
	syndb.RegQuick(TbCurrencyLog, db.UpdatedAtName)

	syndb.RegQuick(TbCurrencyLog, CurrencyLogUserId)
	syndb.RegQuick(TbCurrencyLog, CurrencyLogType)
	syndb.RegQuick(TbCurrencyLog, CurrencyLogAction)
	syndb.RegQuick(TbCurrencyLog, CurrencyLogAmount)
	syndb.RegQuick(TbCurrencyLog, CurrencyLogBefore)
	syndb.RegQuick(TbCurrencyLog, CurrencyLogAfter)
	syndb.RegQuick(TbCurrencyLog, CurrencyLogReason)
	syndb.RegQuick(TbCurrencyLog, CurrencyLogGameId)
	syndb.RegQuick(TbCurrencyLog, CurrencyLogGameName)
	syndb.RegQuick(TbCurrencyLog, CurrencyLogGameCategory)
	syndb.RegQuick(TbCurrencyLog, CurrencyLogBusinessType)
	syndb.RegQuick(TbCurrencyLog, CurrencyLogTransactionId)

	migrate.AutoMigrate(&CurrencyLog{})
}
