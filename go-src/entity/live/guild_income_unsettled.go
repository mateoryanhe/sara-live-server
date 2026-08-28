package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"

	"github.com/gogf/gf/v2/os/gmlock"
)

const (
	TbGuildIncomeUnsettled db.TbName = "guild_income_unsettleds"
)

// GuildIncomeUnsettled 工会未结算收益(主键ID=工会ID,结算时清零)
type GuildIncomeUnsettled struct {
	migrate.OneModel
	LiveRoomIncomeAmounts
}

func NewGuildIncomeUnsettled(guildId uint64) *GuildIncomeUnsettled {
	ret := &GuildIncomeUnsettled{}
	ret.ID = guildId
	now := time.Now()
	ret.CreatedAt = now
	ret.UpdatedAt = now
	syndb.AddData(TbGuildIncomeUnsettled, db.CreatedAtName, &syndb.ColData{IdVal: guildId, ColVal: now})
	syndb.AddData(TbGuildIncomeUnsettled, db.UpdatedAtName, &syndb.ColData{IdVal: guildId, ColVal: now})
	return ret
}

func (r *GuildIncomeUnsettled) AddTotalLiveDuration(v float64) {
	addIncomeAmount(TbGuildIncomeUnsettled, LiveRoomIncomeTotalLiveDuration, r.ID, &r.TotalLiveDuration, v, true, &r.UpdatedAt)
}
func (r *GuildIncomeUnsettled) AddGiftEarn(v float64) {
	addIncomeEarn(TbGuildIncomeUnsettled, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalGiftIncome, &r.TotalGiftIncome)
}
func (r *GuildIncomeUnsettled) AddPaidDanmakuEarn(v float64) {
	addIncomeEarn(TbGuildIncomeUnsettled, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalPaidDanmakuIncome, &r.TotalPaidDanmakuIncome)
}
func (r *GuildIncomeUnsettled) AddPrivateRoomTicketEarn(v float64) {
	addIncomeEarn(TbGuildIncomeUnsettled, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalPrivateRoomTicketIncome, &r.TotalPrivateRoomTicketIncome)
}
func (r *GuildIncomeUnsettled) AddPrivateRoomWatchEarn(v float64) {
	addIncomeEarn(TbGuildIncomeUnsettled, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalPrivateRoomWatchIncome, &r.TotalPrivateRoomWatchIncome)
}

func (r *GuildIncomeUnsettled) AddShortVideoEarn(v float64) {
	addIncomeEarn(TbGuildIncomeUnsettled, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalShortVideoIncome, &r.TotalShortVideoIncome)
}

func (r *GuildIncomeUnsettled) AddGameEarn(goldAmount, incomeDelta float64) {
	addGameEarn(TbGuildIncomeUnsettled, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, goldAmount, incomeDelta)
}

func (r *GuildIncomeUnsettled) AddAmounts(a *LiveRoomIncomeAmounts) {
	if a == nil || a.IsZero() {
		return
	}
	key := liveRoomIncomeLockKey(TbGuildIncomeUnsettled, r.ID)
	gmlock.Lock(key)
	defer gmlock.Unlock(key)
	addIncomeAmountsLocked(TbGuildIncomeUnsettled, r.ID, &r.LiveRoomIncomeAmounts, a)
	touchIncomeUpdatedAt(TbGuildIncomeUnsettled, r.ID, &r.UpdatedAt)
}

func (r *GuildIncomeUnsettled) SnapshotAndClear() LiveRoomIncomeAmounts {
	key := liveRoomIncomeLockKey(TbGuildIncomeUnsettled, r.ID)
	gmlock.Lock(key)
	defer gmlock.Unlock(key)
	snap := r.LiveRoomIncomeAmounts
	clearIncomeAmountsLocked(TbGuildIncomeUnsettled, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt)
	return snap
}

func (r *GuildIncomeUnsettled) Clear() {
	key := liveRoomIncomeLockKey(TbGuildIncomeUnsettled, r.ID)
	gmlock.Lock(key)
	defer gmlock.Unlock(key)
	clearIncomeAmountsLocked(TbGuildIncomeUnsettled, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt)
}

func initGuildIncomeUnsettled() {
	regLiveRoomIncomeCols(TbGuildIncomeUnsettled)
	migrate.AutoMigrate(&GuildIncomeUnsettled{})
}
