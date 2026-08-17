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
	SettlementSalary float64 `gorm:"type:decimal(16,4);default:0;comment:结算薪资" json:"settlementSalary"`
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

func initGuildIncomeSettled() {
	regLiveRoomIncomeCols(TbGuildIncomeSettled)
	syndb.RegQuick(TbGuildIncomeSettled, LiveRoomIncomeSettlementSalary)
	migrate.AutoMigrate(&GuildIncomeSettled{})
}
