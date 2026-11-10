package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbLuckyGiftPoolPendingBalance db.TbName = "lucky_gift_pool_pending_balances"
)

const (
	LuckyGiftPoolPendingStatusPending    uint8 = 0 // 未处理:待合并入奖池
	LuckyGiftPoolPendingStatusProcessed  uint8 = 1 // 已处理:已加到 pool_balance
	LuckyGiftPoolPendingStatusProcessing uint8 = 2 // 处理中(防并发双加)
)

const (
	LuckyGiftPoolPendingBalanceAmount db.TbCol = "amount"
	LuckyGiftPoolPendingBalanceStatus db.TbCol = "status"
)

// LuckyGiftPoolPendingBalance 奖池临时入账(冷启/排空时不直加 pool_balance)
type LuckyGiftPoolPendingBalance struct {
	migrate.OneModel
	Amount float64 `gorm:"type:decimal(20,4);default:0;comment:待并入奖金池金额" json:"amount"`
	Status uint8   `gorm:"default:0;index;comment:处理状态(0未处理,1已处理,2处理中)" json:"status"`
}

// NewLuckyGiftPoolPendingBalance 新建待入账并通过 syndb 入库
func NewLuckyGiftPoolPendingBalance(amount float64) *LuckyGiftPoolPendingBalance {
	if amount <= 0 {
		return nil
	}
	r := &LuckyGiftPoolPendingBalance{}
	r.ID = snowflake.GetId()
	now := time.Now()
	r.SetCreatedAt(now)
	r.SetUpdatedAt(now)
	r.SetAmount(amount)
	r.SetStatus(LuckyGiftPoolPendingStatusPending)
	return r
}

func (r *LuckyGiftPoolPendingBalance) queue(column db.TbCol, value any) {
	if r == nil || r.ID == 0 {
		return
	}
	syndb.AddData(TbLuckyGiftPoolPendingBalance, column, &syndb.ColData{IdVal: r.ID, ColVal: value})
}

func (r *LuckyGiftPoolPendingBalance) touchUpdatedAt() {
	r.UpdatedAt = time.Now()
	r.queue(db.UpdatedAtName, r.UpdatedAt)
}

func (r *LuckyGiftPoolPendingBalance) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	r.queue(db.CreatedAtName, v)
}

func (r *LuckyGiftPoolPendingBalance) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	r.queue(db.UpdatedAtName, v)
}

func (r *LuckyGiftPoolPendingBalance) SetAmount(v float64) {
	if v < 0 {
		v = 0
	}
	r.Amount = v
	r.touchUpdatedAt()
	r.queue(LuckyGiftPoolPendingBalanceAmount, v)
}

func (r *LuckyGiftPoolPendingBalance) SetStatus(v uint8) {
	r.Status = v
	r.touchUpdatedAt()
	r.queue(LuckyGiftPoolPendingBalanceStatus, v)
}

func initLuckyGiftPoolPendingBalance() {
	syndb.RegQuick(TbLuckyGiftPoolPendingBalance, db.CreatedAtName)
	syndb.RegQuick(TbLuckyGiftPoolPendingBalance, db.UpdatedAtName)
	syndb.RegQuick(TbLuckyGiftPoolPendingBalance, LuckyGiftPoolPendingBalanceAmount)
	syndb.RegQuick(TbLuckyGiftPoolPendingBalance, LuckyGiftPoolPendingBalanceStatus)
	migrate.AutoMigrate(&LuckyGiftPoolPendingBalance{})
}
