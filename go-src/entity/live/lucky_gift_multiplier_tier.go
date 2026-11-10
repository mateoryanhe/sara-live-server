package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbLuckyGiftMultiplierTier db.TbName = "lucky_gift_multiplier_tiers"
)

const (
	LuckyGiftMultiplierTypeLow    = "low"
	LuckyGiftMultiplierTypeMedium = "medium"
	LuckyGiftMultiplierTypeHigh   = "high"
)

const (
	LuckyGiftMultiplierTierGiftID     db.TbCol = "gift_id"
	LuckyGiftMultiplierTierMultiplier db.TbCol = "multiplier"
	LuckyGiftMultiplierTierWeight     db.TbCol = "weight"
	LuckyGiftMultiplierTierType       db.TbCol = "type"
	LuckyGiftMultiplierTierSort       db.TbCol = "sort"
)

func IsValidLuckyGiftMultiplierType(v string) bool {
	return v == LuckyGiftMultiplierTypeLow || v == LuckyGiftMultiplierTypeMedium || v == LuckyGiftMultiplierTypeHigh
}

// LuckyGiftMultiplierTier 幸运礼物返奖倍率挡位(gift_id=0 为全局; 其余按礼物全量覆盖)。
// syndb 异步入库; 开奖侧读内存缓存。
type LuckyGiftMultiplierTier struct {
	migrate.OneModel
	GiftID     uint64 `gorm:"index;default:0;comment:礼物ID(0=全局默认,否则 live_gifts.id)" json:"giftId"`
	Multiplier uint32 `gorm:"default:1;comment:返奖倍率" json:"multiplier"`
	Weight     uint32 `gorm:"default:1;comment:中奖后抽档权重(越大越易中)" json:"weight"`
	Type       string `gorm:"size:16;default:'low';comment:挡位类型(low/medium/high)" json:"type"`
	Sort       int    `gorm:"default:0;comment:排序(越小越靠前)" json:"sort"`
}

// NewLuckyGiftMultiplierTier 新建挡位并通过 syndb 入库
func NewLuckyGiftMultiplierTier(giftID uint64, multiplier, weight uint32, tierType string, sort int) *LuckyGiftMultiplierTier {
	r := &LuckyGiftMultiplierTier{}
	r.ID = snowflake.GetId()
	now := time.Now()
	r.SetCreatedAt(now)
	r.SetUpdatedAt(now)
	r.SetGiftID(giftID)
	r.SetMultiplier(multiplier)
	r.SetWeight(weight)
	r.SetType(tierType)
	r.SetSort(sort)
	return r
}

func (r *LuckyGiftMultiplierTier) queue(column db.TbCol, value any) {
	if r == nil || r.ID == 0 {
		return
	}
	syndb.AddData(TbLuckyGiftMultiplierTier, column, &syndb.ColData{IdVal: r.ID, ColVal: value})
}

func (r *LuckyGiftMultiplierTier) touchUpdatedAt() {
	r.UpdatedAt = time.Now()
	r.queue(db.UpdatedAtName, r.UpdatedAt)
}

func (r *LuckyGiftMultiplierTier) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	r.queue(db.CreatedAtName, v)
}

func (r *LuckyGiftMultiplierTier) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	r.queue(db.UpdatedAtName, v)
}

func (r *LuckyGiftMultiplierTier) SetGiftID(v uint64) {
	r.GiftID = v
	r.touchUpdatedAt()
	r.queue(LuckyGiftMultiplierTierGiftID, v)
}

func (r *LuckyGiftMultiplierTier) SetMultiplier(v uint32) {
	if v == 0 {
		v = 1
	}
	r.Multiplier = v
	r.touchUpdatedAt()
	r.queue(LuckyGiftMultiplierTierMultiplier, v)
}

func (r *LuckyGiftMultiplierTier) SetWeight(v uint32) {
	if v == 0 {
		v = 1
	}
	r.Weight = v
	r.touchUpdatedAt()
	r.queue(LuckyGiftMultiplierTierWeight, v)
}

func (r *LuckyGiftMultiplierTier) SetType(v string) {
	if !IsValidLuckyGiftMultiplierType(v) {
		v = LuckyGiftMultiplierTypeLow
	}
	r.Type = v
	r.touchUpdatedAt()
	r.queue(LuckyGiftMultiplierTierType, v)
}

func (r *LuckyGiftMultiplierTier) SetSort(v int) {
	r.Sort = v
	r.touchUpdatedAt()
	r.queue(LuckyGiftMultiplierTierSort, v)
}

func initLuckyGiftMultiplierTier() {
	syndb.RegQuick(TbLuckyGiftMultiplierTier, db.CreatedAtName)
	syndb.RegQuick(TbLuckyGiftMultiplierTier, db.UpdatedAtName)
	syndb.RegQuick(TbLuckyGiftMultiplierTier, LuckyGiftMultiplierTierGiftID)
	syndb.RegQuick(TbLuckyGiftMultiplierTier, LuckyGiftMultiplierTierMultiplier)
	syndb.RegQuick(TbLuckyGiftMultiplierTier, LuckyGiftMultiplierTierWeight)
	syndb.RegQuick(TbLuckyGiftMultiplierTier, LuckyGiftMultiplierTierType)
	syndb.RegQuick(TbLuckyGiftMultiplierTier, LuckyGiftMultiplierTierSort)
	migrate.AutoMigrate(&LuckyGiftMultiplierTier{})
}
