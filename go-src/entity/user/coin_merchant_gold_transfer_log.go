package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const TbCoinMerchantGoldTransferLog db.TbName = "coin_merchant_gold_transfer_logs"

const (
	CoinMerchantGoldTransferLogMerchantUserId   db.TbCol = "coin_merchant_user_id"
	CoinMerchantGoldTransferLogTargetUserId     db.TbCol = "target_user_id"
	CoinMerchantGoldTransferLogAmount           db.TbCol = "amount"
	CoinMerchantGoldTransferLogSenderGoldBefore db.TbCol = "sender_gold_before"
	CoinMerchantGoldTransferLogSenderGoldAfter  db.TbCol = "sender_gold_after"
	CoinMerchantGoldTransferLogTargetGoldBefore db.TbCol = "target_gold_before"
	CoinMerchantGoldTransferLogTargetGoldAfter  db.TbCol = "target_gold_after"
)

// CoinMerchantGoldTransferLog 币商向普通用户转赠金币的独立业务流水。
// 各字段保留数据库默认值，允许 syndb 按列缓冲写入并最终补齐整行。
type CoinMerchantGoldTransferLog struct {
	migrate.OneModel
	CreatedAt          time.Time `gorm:"index:idx_cmgtl_merchant_created,priority:2" json:"-"`
	CoinMerchantUserId uint64    `gorm:"index:idx_cmgtl_merchant_created,priority:1;default:0;comment:转出币商用户ID" json:"coinMerchantUserId"`
	TargetUserId       uint64    `gorm:"index;default:0;comment:目标用户ID" json:"targetUserId"`
	Amount             float64   `gorm:"type:decimal(20,2);default:0;comment:转赠金币数量" json:"amount"`
	SenderGoldBefore   float64   `gorm:"type:decimal(20,2);default:0;comment:币商转出前金币" json:"senderGoldBefore"`
	SenderGoldAfter    float64   `gorm:"type:decimal(20,2);default:0;comment:币商转出后金币" json:"senderGoldAfter"`
	TargetGoldBefore   float64   `gorm:"type:decimal(20,2);default:0;comment:目标用户转入前金币" json:"targetGoldBefore"`
	TargetGoldAfter    float64   `gorm:"type:decimal(20,2);default:0;comment:目标用户转入后金币" json:"targetGoldAfter"`
}

func NewCoinMerchantGoldTransferLog(
	merchantUserId, targetUserId uint64,
	amount, senderGoldBefore, senderGoldAfter, targetGoldBefore, targetGoldAfter float64,
) *CoinMerchantGoldTransferLog {
	row := &CoinMerchantGoldTransferLog{}
	row.ID = snowflake.GetId()
	now := time.Now()
	row.SetCreatedAt(now)
	row.SetUpdatedAt(now)
	row.SetCoinMerchantUserId(merchantUserId)
	row.SetTargetUserId(targetUserId)
	row.SetAmount(amount)
	row.SetSenderGoldBefore(senderGoldBefore)
	row.SetSenderGoldAfter(senderGoldAfter)
	row.SetTargetGoldBefore(targetGoldBefore)
	row.SetTargetGoldAfter(targetGoldAfter)
	return row
}

func (r *CoinMerchantGoldTransferLog) queue(column db.TbCol, value any) {
	syndb.AddData(TbCoinMerchantGoldTransferLog, column, &syndb.ColData{IdVal: r.ID, ColVal: value})
}

func (r *CoinMerchantGoldTransferLog) SetCreatedAt(value time.Time) {
	r.CreatedAt = value
	r.queue(db.CreatedAtName, value)
}

func (r *CoinMerchantGoldTransferLog) SetUpdatedAt(value time.Time) {
	r.UpdatedAt = value
	r.queue(db.UpdatedAtName, value)
}

func (r *CoinMerchantGoldTransferLog) SetCoinMerchantUserId(value uint64) {
	r.CoinMerchantUserId = value
	r.queue(CoinMerchantGoldTransferLogMerchantUserId, value)
}

func (r *CoinMerchantGoldTransferLog) SetTargetUserId(value uint64) {
	r.TargetUserId = value
	r.queue(CoinMerchantGoldTransferLogTargetUserId, value)
}

func (r *CoinMerchantGoldTransferLog) SetAmount(value float64) {
	r.Amount = value
	r.queue(CoinMerchantGoldTransferLogAmount, value)
}

func (r *CoinMerchantGoldTransferLog) SetSenderGoldBefore(value float64) {
	r.SenderGoldBefore = value
	r.queue(CoinMerchantGoldTransferLogSenderGoldBefore, value)
}

func (r *CoinMerchantGoldTransferLog) SetSenderGoldAfter(value float64) {
	r.SenderGoldAfter = value
	r.queue(CoinMerchantGoldTransferLogSenderGoldAfter, value)
}

func (r *CoinMerchantGoldTransferLog) SetTargetGoldBefore(value float64) {
	r.TargetGoldBefore = value
	r.queue(CoinMerchantGoldTransferLogTargetGoldBefore, value)
}

func (r *CoinMerchantGoldTransferLog) SetTargetGoldAfter(value float64) {
	r.TargetGoldAfter = value
	r.queue(CoinMerchantGoldTransferLogTargetGoldAfter, value)
}

func initCoinMerchantGoldTransferLog() {
	syndb.RegQuick(TbCoinMerchantGoldTransferLog, db.CreatedAtName)
	syndb.RegLazy(TbCoinMerchantGoldTransferLog, db.UpdatedAtName)
	syndb.RegQuick(TbCoinMerchantGoldTransferLog, CoinMerchantGoldTransferLogMerchantUserId)
	syndb.RegQuick(TbCoinMerchantGoldTransferLog, CoinMerchantGoldTransferLogTargetUserId)
	syndb.RegQuick(TbCoinMerchantGoldTransferLog, CoinMerchantGoldTransferLogAmount)
	syndb.RegQuick(TbCoinMerchantGoldTransferLog, CoinMerchantGoldTransferLogSenderGoldBefore)
	syndb.RegQuick(TbCoinMerchantGoldTransferLog, CoinMerchantGoldTransferLogSenderGoldAfter)
	syndb.RegQuick(TbCoinMerchantGoldTransferLog, CoinMerchantGoldTransferLogTargetGoldBefore)
	syndb.RegQuick(TbCoinMerchantGoldTransferLog, CoinMerchantGoldTransferLogTargetGoldAfter)
	migrate.AutoMigrate(&CoinMerchantGoldTransferLog{})
}
