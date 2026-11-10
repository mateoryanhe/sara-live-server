package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbLuckyGiftPoolCfg db.TbName = "lucky_gift_pool_cfgs"
)

// 全局倍率挡位挂在 lucky_gift_multiplier_tiers.gift_id=0
const GlobalLuckyGiftMultiplierGiftID uint64 = 0

const (
	DefaultLuckyGiftPoolBalance       float64 = 10_000_000 // 默认奖金池 1000 万
	DefaultLuckyGiftContributePercent float64 = 0          // 礼物金额进入奖池比例(%)
	DefaultLuckyGiftWinProbability    float64 = 0          // 中奖概率(%)
)

const (
	LuckyGiftPoolCfgPoolBalance       db.TbCol = "pool_balance"
	LuckyGiftPoolCfgContributePercent db.TbCol = "contribute_percent"
	LuckyGiftPoolCfgWinProbability    db.TbCol = "win_probability"
)

// LuckyGiftPoolCfg 幸运礼物奖池全局配置(CMS 单行, syndb 异步入库)。
// 倍率挡位存于 lucky_gift_multiplier_tiers(gift_id=0)，不在本表。
type LuckyGiftPoolCfg struct {
	migrate.OneModel
	PoolBalance       float64 `gorm:"type:decimal(20,4);default:10000000;comment:奖金池余额" json:"poolBalance"`
	ContributePercent float64 `gorm:"type:decimal(10,4);default:0;comment:礼物金额进入奖池比例(%)" json:"contributePercent"`
	WinProbability    float64 `gorm:"type:decimal(10,4);default:0;comment:中奖概率(%)" json:"winProbability"`
}

// NewLuckyGiftPoolCfg 新建奖池配置并通过 syndb 入库
func NewLuckyGiftPoolCfg(poolBalance, contributePercent, winProbability float64) *LuckyGiftPoolCfg {
	r := &LuckyGiftPoolCfg{}
	r.ID = snowflake.GetId()
	now := time.Now()
	r.SetCreatedAt(now)
	r.SetUpdatedAt(now)
	r.SetPoolBalance(poolBalance)
	r.SetContributePercent(contributePercent)
	r.SetWinProbability(winProbability)
	return r
}

func (r *LuckyGiftPoolCfg) queue(column db.TbCol, value any) {
	if r == nil || r.ID == 0 {
		return
	}
	syndb.AddData(TbLuckyGiftPoolCfg, column, &syndb.ColData{IdVal: r.ID, ColVal: value})
}

func (r *LuckyGiftPoolCfg) touchUpdatedAt() {
	r.UpdatedAt = time.Now()
	r.queue(db.UpdatedAtName, r.UpdatedAt)
}

func (r *LuckyGiftPoolCfg) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	r.queue(db.CreatedAtName, v)
}

func (r *LuckyGiftPoolCfg) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	r.queue(db.UpdatedAtName, v)
}

func (r *LuckyGiftPoolCfg) SetPoolBalance(v float64) {
	if v < 0 {
		v = 0
	}
	r.PoolBalance = v
	r.touchUpdatedAt()
	r.queue(LuckyGiftPoolCfgPoolBalance, v)
}

func (r *LuckyGiftPoolCfg) SetContributePercent(v float64) {
	if v < 0 {
		v = 0
	}
	if v > 100 {
		v = 100
	}
	r.ContributePercent = v
	r.touchUpdatedAt()
	r.queue(LuckyGiftPoolCfgContributePercent, v)
}

func (r *LuckyGiftPoolCfg) SetWinProbability(v float64) {
	if v < 0 {
		v = 0
	}
	if v > 100 {
		v = 100
	}
	r.WinProbability = v
	r.touchUpdatedAt()
	r.queue(LuckyGiftPoolCfgWinProbability, v)
}

func initLuckyGiftPoolCfg() {
	syndb.RegQuick(TbLuckyGiftPoolCfg, db.CreatedAtName)
	syndb.RegQuick(TbLuckyGiftPoolCfg, db.UpdatedAtName)
	syndb.RegQuick(TbLuckyGiftPoolCfg, LuckyGiftPoolCfgPoolBalance)
	syndb.RegQuick(TbLuckyGiftPoolCfg, LuckyGiftPoolCfgContributePercent)
	syndb.RegQuick(TbLuckyGiftPoolCfg, LuckyGiftPoolCfgWinProbability)
	migrate.AutoMigrate(&LuckyGiftPoolCfg{})
}
