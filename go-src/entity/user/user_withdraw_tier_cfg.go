package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const TbUserWithdrawTierCfg db.TbName = "user_withdraw_tier_cfgs"

const (
	UserWithdrawTierCfgName           db.TbCol = "name"
	UserWithdrawTierCfgWithdrawUsd    db.TbCol = "withdraw_usd"
	UserWithdrawTierCfgMinGameConsumeGold db.TbCol = "min_game_consume_gold"
	UserWithdrawTierCfgStatus         db.TbCol = "status"
)

const (
	UserWithdrawTierCfgStatusOffShelf uint8 = 0
	UserWithdrawTierCfgStatusOnShelf  uint8 = 1
)

// UserWithdrawTierCfg 用户提现档位(CMS 管理,上架项缓存)
type UserWithdrawTierCfg struct {
	migrate.OneModel
	Name           string  `gorm:"size:64;uniqueIndex;comment:档位名称" json:"name"`
	WithdrawUsd    float64 `gorm:"type:decimal(16,4);default:0;comment:提现金额(USD)" json:"withdrawUsd"`
	MinGameConsumeGold float64 `gorm:"type:decimal(20,2);default:0;comment:游戏消费累计金币门槛(配置,逻辑待接)" json:"minGameConsumeGold"`
	Status         uint8   `gorm:"default:0;comment:状态(0下架,1上架)" json:"status"`
}

func (r *UserWithdrawTierCfg) queue(column db.TbCol, value any) {
	syndb.AddData(TbUserWithdrawTierCfg, column, &syndb.ColData{IdVal: r.ID, ColVal: value})
}

func (r *UserWithdrawTierCfg) touchUpdatedAt() {
	r.UpdatedAt = time.Now()
	r.queue(db.UpdatedAtName, r.UpdatedAt)
}

func (r *UserWithdrawTierCfg) SetName(v string) {
	r.Name = v
	r.touchUpdatedAt()
	r.queue(UserWithdrawTierCfgName, v)
}

func (r *UserWithdrawTierCfg) SetWithdrawUsd(v float64) {
	r.WithdrawUsd = v
	r.touchUpdatedAt()
	r.queue(UserWithdrawTierCfgWithdrawUsd, v)
}

func (r *UserWithdrawTierCfg) SetMinGameConsumeGold(v float64) {
	r.MinGameConsumeGold = v
	r.touchUpdatedAt()
	r.queue(UserWithdrawTierCfgMinGameConsumeGold, v)
}

func (r *UserWithdrawTierCfg) SetStatus(v uint8) {
	r.Status = v
	r.touchUpdatedAt()
	r.queue(UserWithdrawTierCfgStatus, v)
}

func initUserWithdrawTierCfg() {
	syndb.RegQuick(TbUserWithdrawTierCfg, db.CreatedAtName)
	syndb.RegQuick(TbUserWithdrawTierCfg, db.UpdatedAtName)
	syndb.RegQuick(TbUserWithdrawTierCfg, UserWithdrawTierCfgName)
	syndb.RegQuick(TbUserWithdrawTierCfg, UserWithdrawTierCfgWithdrawUsd)
	syndb.RegQuick(TbUserWithdrawTierCfg, UserWithdrawTierCfgMinGameConsumeGold)
	syndb.RegQuick(TbUserWithdrawTierCfg, UserWithdrawTierCfgStatus)
	migrate.AutoMigrate(&UserWithdrawTierCfg{})
}
