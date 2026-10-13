package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const TbUserWithdrawLog db.TbName = "user_withdraw_logs"

const (
	UserWithdrawLogUserId               db.TbCol = "user_id"
	UserWithdrawLogGoldAmount           db.TbCol = "gold_amount"
	UserWithdrawLogWithdrawUsd          db.TbCol = "withdraw_usd"
	UserWithdrawLogUsdToGoldRate        db.TbCol = "usd_to_gold_rate"
	UserWithdrawLogStatus               db.TbCol = "status"
	UserWithdrawLogTransferOrderId      db.TbCol = "transfer_order_id"
	UserWithdrawLogTransferPlatformNo   db.TbCol = "transfer_platform_no"
	UserWithdrawLogTransferCurrency     db.TbCol = "transfer_currency"
	UserWithdrawLogTransferLocalAmount  db.TbCol = "transfer_local_amount"
	UserWithdrawLogTransferAt           db.TbCol = "transfer_at"
	UserWithdrawLogTransferFailMsg      db.TbCol = "transfer_fail_msg"
	UserWithdrawLogRejectReason         db.TbCol = "reject_reason"
)

const (
	UserWithdrawStatusPending      uint8 = 0 // 待审核
	UserWithdrawStatusApproved     uint8 = 1 // 审核通过,待代付
	UserWithdrawStatusTransferred  uint8 = 2 // 代付成功
	UserWithdrawStatusTransferring uint8 = 3 // 已提交代付
	UserWithdrawStatusRejected     uint8 = 4 // 审核拒绝(已退金币)
)

// UserWithdrawLog 用户提现单
type UserWithdrawLog struct {
	migrate.OneModel
	CreatedAt              time.Time  `gorm:"index:idx_uwl_user_status_created,priority:3" json:"-"`
	UserId                 uint64     `gorm:"index:idx_uwl_user_status_created,priority:1;default:0;comment:用户ID" json:"userId"`
	GoldAmount             float64    `gorm:"type:decimal(20,2);default:0;comment:提现金币数" json:"goldAmount"`
	WithdrawUsd            float64    `gorm:"type:decimal(16,4);default:0;comment:折算USD" json:"withdrawUsd"`
	UsdToGoldRate          int        `gorm:"default:0;comment:申请时1USD兑换金币数快照" json:"usdToGoldRate"`
	Status                 uint8      `gorm:"index:idx_uwl_user_status_created,priority:2;default:0;comment:状态" json:"status"`
	TransferOrderId        string     `gorm:"size:64;default:'';comment:HaiPay orderId" json:"transferOrderId"`
	TransferPlatformNo     string     `gorm:"size:64;default:'';comment:HaiPay orderNo" json:"transferPlatformNo"`
	TransferCurrency       string     `gorm:"size:16;default:'';comment:代付币种" json:"transferCurrency"`
	TransferLocalAmount    float64    `gorm:"type:decimal(20,4);default:0;comment:代付本地金额" json:"transferLocalAmount"`
	TransferAt             *time.Time `gorm:"comment:代付成功时间" json:"transferAt"`
	TransferFailMsg        string     `gorm:"size:255;default:'';comment:代付失败原因" json:"transferFailMsg"`
	RejectReason           string     `gorm:"size:255;default:'';comment:审核拒绝原因" json:"rejectReason"`
}

func NewUserWithdrawLog(userId uint64, goldAmount, withdrawUsd float64, usdToGoldRate int) *UserWithdrawLog {
	now := time.Now()
	return &UserWithdrawLog{
		OneModel: migrate.OneModel{
			ID:        snowflake.GetId(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserId:        userId,
		GoldAmount:    goldAmount,
		WithdrawUsd:   withdrawUsd,
		UsdToGoldRate: usdToGoldRate,
		Status:        UserWithdrawStatusPending,
	}
}

func (r *UserWithdrawLog) queue(column db.TbCol, value any) {
	syndb.AddData(TbUserWithdrawLog, column, &syndb.ColData{IdVal: r.ID, ColVal: value})
}

func (r *UserWithdrawLog) SetCreatedAt(value time.Time) {
	r.CreatedAt = value
	r.queue(db.CreatedAtName, value)
}

func (r *UserWithdrawLog) SetUpdatedAt(value time.Time) {
	r.UpdatedAt = value
	r.queue(db.UpdatedAtName, value)
}

func (r *UserWithdrawLog) SetUserId(value uint64) {
	r.UserId = value
	r.queue(UserWithdrawLogUserId, value)
}

func (r *UserWithdrawLog) SetGoldAmount(value float64) {
	r.GoldAmount = value
	r.queue(UserWithdrawLogGoldAmount, value)
}

func (r *UserWithdrawLog) SetWithdrawUsd(value float64) {
	r.WithdrawUsd = value
	r.queue(UserWithdrawLogWithdrawUsd, value)
}

func (r *UserWithdrawLog) SetUsdToGoldRate(value int) {
	r.UsdToGoldRate = value
	r.queue(UserWithdrawLogUsdToGoldRate, value)
}

func (r *UserWithdrawLog) SetStatus(value uint8) {
	r.Status = value
	r.queue(UserWithdrawLogStatus, value)
}

func (r *UserWithdrawLog) SetTransferPayout(orderId, platformNo, currency string, localAmount float64) {
	r.TransferOrderId = orderId
	r.TransferPlatformNo = platformNo
	r.TransferCurrency = currency
	r.TransferLocalAmount = localAmount
	r.queue(UserWithdrawLogTransferOrderId, orderId)
	r.queue(UserWithdrawLogTransferPlatformNo, platformNo)
	r.queue(UserWithdrawLogTransferCurrency, currency)
	r.queue(UserWithdrawLogTransferLocalAmount, localAmount)
}

func (r *UserWithdrawLog) SetTransferAt(value *time.Time) {
	r.TransferAt = value
	r.queue(UserWithdrawLogTransferAt, value)
}

func (r *UserWithdrawLog) SetTransferFailMsg(value string) {
	r.TransferFailMsg = value
	r.queue(UserWithdrawLogTransferFailMsg, value)
}

func (r *UserWithdrawLog) SetRejectReason(value string) {
	r.RejectReason = value
	r.queue(UserWithdrawLogRejectReason, value)
}

func initUserWithdrawLog() {
	syndb.RegQuick(TbUserWithdrawLog, db.CreatedAtName)
	syndb.RegLazy(TbUserWithdrawLog, db.UpdatedAtName)
	syndb.RegQuick(TbUserWithdrawLog, UserWithdrawLogUserId)
	syndb.RegQuick(TbUserWithdrawLog, UserWithdrawLogGoldAmount)
	syndb.RegQuick(TbUserWithdrawLog, UserWithdrawLogWithdrawUsd)
	syndb.RegQuick(TbUserWithdrawLog, UserWithdrawLogUsdToGoldRate)
	syndb.RegQuick(TbUserWithdrawLog, UserWithdrawLogStatus)
	syndb.RegLazy(TbUserWithdrawLog, UserWithdrawLogTransferOrderId)
	syndb.RegLazy(TbUserWithdrawLog, UserWithdrawLogTransferPlatformNo)
	syndb.RegLazy(TbUserWithdrawLog, UserWithdrawLogTransferCurrency)
	syndb.RegLazy(TbUserWithdrawLog, UserWithdrawLogTransferLocalAmount)
	syndb.RegLazy(TbUserWithdrawLog, UserWithdrawLogTransferAt)
	syndb.RegLazy(TbUserWithdrawLog, UserWithdrawLogTransferFailMsg)
	syndb.RegLazy(TbUserWithdrawLog, UserWithdrawLogRejectReason)
	migrate.AutoMigrate(&UserWithdrawLog{})
}
