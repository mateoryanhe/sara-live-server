package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbGuildIncomeUnsettledArchive db.TbName = "guild_income_unsettled_archives"
)

const (
	GuildIncomeUnsettledArchiveGuildId db.TbCol = "guild_id"
)

// GuildIncomeUnsettledArchive 工会未结算收益归档(每次结算/归档一条新记录)
type GuildIncomeUnsettledArchive struct {
	migrate.OneModel
	GuildId uint64 `gorm:"index;default:0;comment:工会ID" json:"guildId"`
	LiveRoomIncomeAmounts
}

func NewGuildIncomeUnsettledArchive(guildId uint64, a *LiveRoomIncomeAmounts) *GuildIncomeUnsettledArchive {
	if a == nil {
		a = &LiveRoomIncomeAmounts{}
	}
	ret := &GuildIncomeUnsettledArchive{}
	ret.ID = snowflake.GetId()
	now := time.Now()
	ret.CreatedAt = now
	ret.UpdatedAt = now
	ret.GuildId = guildId
	ret.LiveRoomIncomeAmounts = *a

	syndb.AddData(TbGuildIncomeUnsettledArchive, db.CreatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(TbGuildIncomeUnsettledArchive, db.UpdatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(TbGuildIncomeUnsettledArchive, GuildIncomeUnsettledArchiveGuildId, &syndb.ColData{IdVal: ret.ID, ColVal: guildId})
	writeIncomeAmountLocked(TbGuildIncomeUnsettledArchive, LiveRoomIncomeTotalIncome, ret.ID, a.TotalIncome)
	writeIncomeAmountLocked(TbGuildIncomeUnsettledArchive, LiveRoomIncomeTotalGiftIncome, ret.ID, a.TotalGiftIncome)
	writeIncomeAmountLocked(TbGuildIncomeUnsettledArchive, LiveRoomIncomeTotalPaidDanmakuIncome, ret.ID, a.TotalPaidDanmakuIncome)
	writeIncomeAmountLocked(TbGuildIncomeUnsettledArchive, LiveRoomIncomeTotalPrivateRoomTicketIncome, ret.ID, a.TotalPrivateRoomTicketIncome)
	writeIncomeAmountLocked(TbGuildIncomeUnsettledArchive, LiveRoomIncomeTotalPrivateRoomWatchIncome, ret.ID, a.TotalPrivateRoomWatchIncome)
	writeIncomeAmountLocked(TbGuildIncomeUnsettledArchive, LiveRoomIncomeTotalVideoCallIncome, ret.ID, a.TotalVideoCallIncome)
	writeIncomeAmountLocked(TbGuildIncomeUnsettledArchive, LiveRoomIncomeTotalVideoCallTicketIncome, ret.ID, a.TotalVideoCallTicketIncome)
	writeIncomeAmountLocked(TbGuildIncomeUnsettledArchive, LiveRoomIncomeTotalVideoCallBillingIncome, ret.ID, a.TotalVideoCallBillingIncome)
	writeIncomeAmountLocked(TbGuildIncomeUnsettledArchive, LiveRoomIncomeTotalLiveDuration, ret.ID, a.TotalLiveDuration)
	return ret
}

func initGuildIncomeUnsettledArchive() {
	syndb.RegQuick(TbGuildIncomeUnsettledArchive, db.CreatedAtName)
	syndb.RegQuick(TbGuildIncomeUnsettledArchive, db.UpdatedAtName)
	syndb.RegQuick(TbGuildIncomeUnsettledArchive, GuildIncomeUnsettledArchiveGuildId)
	regLiveRoomIncomeCols(TbGuildIncomeUnsettledArchive)
	migrate.AutoMigrate(&GuildIncomeUnsettledArchive{})
}
