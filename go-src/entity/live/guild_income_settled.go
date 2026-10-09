package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"

	"github.com/gogf/gf/v2/os/gmlock"
)

const (
	TbGuildIncomeSettled db.TbName = "guild_income_settleds"
)

// GuildIncomeSettled 工会已结算收益(主键ID=工会ID,每次结算累加)
type GuildIncomeSettled struct {
	migrate.OneModel
	LiveRoomIncomeAmounts
	SettlementSalary         float64 `gorm:"type:decimal(16,4);default:0;comment:结算薪资" json:"settlementSalary"`
	SettlementShareAmount    float64 `gorm:"type:decimal(16,4);default:0;comment:结算分佣金额" json:"settlementShareAmount"`
	SettlementShareAmountUsd      float64 `gorm:"type:decimal(16,4);default:0;comment:结算分佣金额(USD)" json:"settlementShareAmountUsd"`
	SettlementReceivableUsd       float64 `gorm:"type:decimal(16,4);default:0;comment:累计可收USD合计(工会+主播，代收口径)" json:"settlementReceivableUsd"`
	SettlementGuildReceivableUsd  float64 `gorm:"type:decimal(16,4);default:0;comment:累计工会应得USD" json:"settlementGuildReceivableUsd"`
	SettlementAnchorReceivableUsd float64 `gorm:"type:decimal(16,4);default:0;comment:累计主播应得USD(当前多由工会代收)" json:"settlementAnchorReceivableUsd"`
}

func NewGuildIncomeSettled(guildId uint64) *GuildIncomeSettled {
	ret := &GuildIncomeSettled{}
	ret.ID = guildId
	now := time.Now()
	ret.CreatedAt = now
	ret.UpdatedAt = now
	syndb.AddData(TbGuildIncomeSettled, db.CreatedAtName, &syndb.ColData{IdVal: guildId, ColVal: now})
	syndb.AddData(TbGuildIncomeSettled, db.UpdatedAtName, &syndb.ColData{IdVal: guildId, ColVal: now})
	return ret
}

func (r *GuildIncomeSettled) AddTotalLiveDuration(v float64) {
	addIncomeAmount(TbGuildIncomeSettled, LiveRoomIncomeTotalLiveDuration, r.ID, &r.TotalLiveDuration, v, true, &r.UpdatedAt)
}

func (r *GuildIncomeSettled) AddAmounts(a *LiveRoomIncomeAmounts) {
	if a == nil || a.IsZero() {
		return
	}
	key := liveRoomIncomeLockKey(TbGuildIncomeSettled, r.ID)
	gmlock.Lock(key)
	defer gmlock.Unlock(key)
	addIncomeAmountsLocked(TbGuildIncomeSettled, r.ID, &r.LiveRoomIncomeAmounts, a)
	touchIncomeUpdatedAt(TbGuildIncomeSettled, r.ID, &r.UpdatedAt)
}

// AddSettlementSalary 累加结算薪资
func (r *GuildIncomeSettled) AddSettlementSalary(v float64) {
	if r == nil || v == 0 {
		return
	}
	addIncomeAmount(TbGuildIncomeSettled, LiveRoomIncomeSettlementSalary, r.ID, &r.SettlementSalary, v, false, &r.UpdatedAt)
}

// AddSettlementShareAmount 累加结算分佣金额
func (r *GuildIncomeSettled) AddSettlementShareAmount(v float64) {
	if r == nil || v == 0 {
		return
	}
	addIncomeAmount(TbGuildIncomeSettled, LiveRoomIncomeSettlementShareAmount, r.ID, &r.SettlementShareAmount, v, false, &r.UpdatedAt)
}

// AddSettlementShareAmountUsd 累加结算分佣金额(USD)
func (r *GuildIncomeSettled) AddSettlementShareAmountUsd(v float64) {
	if r == nil || v == 0 {
		return
	}
	addIncomeAmount(TbGuildIncomeSettled, LiveRoomIncomeSettlementShareAmountUsd, r.ID, &r.SettlementShareAmountUsd, v, false, &r.UpdatedAt)
}

// AddSettlementReceivableUsd 累加累计可收 USD 合计(内部同步工会+主播分项)
func (r *GuildIncomeSettled) AddSettlementReceivableUsd(v float64) {
	if r == nil || v == 0 {
		return
	}
	addIncomeAmount(TbGuildIncomeSettled, LiveRoomIncomeSettlementReceivableUsd, r.ID, &r.SettlementReceivableUsd, v, false, &r.UpdatedAt)
}

// AddSettlementGuildReceivableUsd 累加工会自身应得 USD，并同步合计
func (r *GuildIncomeSettled) AddSettlementGuildReceivableUsd(v float64) {
	if r == nil || v == 0 {
		return
	}
	addIncomeAmount(TbGuildIncomeSettled, LiveRoomIncomeSettlementGuildReceivableUsd, r.ID, &r.SettlementGuildReceivableUsd, v, false, &r.UpdatedAt)
	r.AddSettlementReceivableUsd(v)
}

// AddSettlementAnchorReceivableUsd 累加名下主播应得 USD(代收预留)，并同步合计
func (r *GuildIncomeSettled) AddSettlementAnchorReceivableUsd(v float64) {
	if r == nil || v == 0 {
		return
	}
	addIncomeAmount(TbGuildIncomeSettled, LiveRoomIncomeSettlementAnchorReceivableUsd, r.ID, &r.SettlementAnchorReceivableUsd, v, false, &r.UpdatedAt)
	r.AddSettlementReceivableUsd(v)
}

func initGuildIncomeSettled() {
	regLiveRoomIncomeCols(TbGuildIncomeSettled)
	syndb.RegQuick(TbGuildIncomeSettled, LiveRoomIncomeSettlementSalary)
	syndb.RegQuick(TbGuildIncomeSettled, LiveRoomIncomeSettlementShareAmount)
	syndb.RegQuick(TbGuildIncomeSettled, LiveRoomIncomeSettlementShareAmountUsd)
	syndb.RegQuick(TbGuildIncomeSettled, LiveRoomIncomeSettlementReceivableUsd)
	syndb.RegQuick(TbGuildIncomeSettled, LiveRoomIncomeSettlementGuildReceivableUsd)
	syndb.RegQuick(TbGuildIncomeSettled, LiveRoomIncomeSettlementAnchorReceivableUsd)
	migrate.AutoMigrate(&GuildIncomeSettled{})
}
