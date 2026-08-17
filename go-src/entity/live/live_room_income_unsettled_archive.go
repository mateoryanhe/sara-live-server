package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbLiveRoomIncomeUnsettledArchive db.TbName = "live_room_income_unsettled_archives"
)

const (
	LiveRoomIncomeUnsettledArchiveRoomId  db.TbCol = "room_id"
	LiveRoomIncomeUnsettledArchiveGuildId db.TbCol = "guild_id"
)

// LiveRoomIncomeUnsettledArchive 下架时未结算收益快照(每次下架一条新记录,仅历史留存不还原)
type LiveRoomIncomeUnsettledArchive struct {
	migrate.OneModel
	RoomId  uint64 `gorm:"index;default:0;comment:直播间ID" json:"roomId"`
	GuildId uint64 `gorm:"index;default:0;comment:工会ID" json:"guildId"`
	LiveRoomIncomeAmounts
}

// NewLiveRoomIncomeUnsettledArchive 新建一条下架未结算归档并入库
func NewLiveRoomIncomeUnsettledArchive(roomId, guildId uint64, a *LiveRoomIncomeAmounts) *LiveRoomIncomeUnsettledArchive {
	if a == nil {
		a = &LiveRoomIncomeAmounts{}
	}
	ret := &LiveRoomIncomeUnsettledArchive{}
	ret.ID = snowflake.GetId()
	now := time.Now()
	ret.CreatedAt = now
	ret.UpdatedAt = now
	ret.RoomId = roomId
	ret.GuildId = guildId
	ret.LiveRoomIncomeAmounts = *a

	syndb.AddData(TbLiveRoomIncomeUnsettledArchive, db.CreatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(TbLiveRoomIncomeUnsettledArchive, db.UpdatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(TbLiveRoomIncomeUnsettledArchive, LiveRoomIncomeUnsettledArchiveRoomId, &syndb.ColData{IdVal: ret.ID, ColVal: roomId})
	syndb.AddData(TbLiveRoomIncomeUnsettledArchive, LiveRoomIncomeUnsettledArchiveGuildId, &syndb.ColData{IdVal: ret.ID, ColVal: guildId})
	writeIncomeAmountLocked(TbLiveRoomIncomeUnsettledArchive, LiveRoomIncomeTotalIncome, ret.ID, a.TotalIncome)
	writeIncomeAmountLocked(TbLiveRoomIncomeUnsettledArchive, LiveRoomIncomeTotalGiftIncome, ret.ID, a.TotalGiftIncome)
	writeIncomeAmountLocked(TbLiveRoomIncomeUnsettledArchive, LiveRoomIncomeTotalPaidDanmakuIncome, ret.ID, a.TotalPaidDanmakuIncome)
	writeIncomeAmountLocked(TbLiveRoomIncomeUnsettledArchive, LiveRoomIncomeTotalPrivateRoomTicketIncome, ret.ID, a.TotalPrivateRoomTicketIncome)
	writeIncomeAmountLocked(TbLiveRoomIncomeUnsettledArchive, LiveRoomIncomeTotalPrivateRoomWatchIncome, ret.ID, a.TotalPrivateRoomWatchIncome)
	writeIncomeAmountLocked(TbLiveRoomIncomeUnsettledArchive, LiveRoomIncomeTotalVideoCallIncome, ret.ID, a.TotalVideoCallIncome)
	writeIncomeAmountLocked(TbLiveRoomIncomeUnsettledArchive, LiveRoomIncomeTotalVideoCallTicketIncome, ret.ID, a.TotalVideoCallTicketIncome)
	writeIncomeAmountLocked(TbLiveRoomIncomeUnsettledArchive, LiveRoomIncomeTotalVideoCallBillingIncome, ret.ID, a.TotalVideoCallBillingIncome)
	writeIncomeAmountLocked(TbLiveRoomIncomeUnsettledArchive, LiveRoomIncomeTotalLiveDuration, ret.ID, a.TotalLiveDuration)
	writeIncomeCountLocked(TbLiveRoomIncomeUnsettledArchive, LiveRoomIncomeEffectiveLiveCount, ret.ID, a.EffectiveLiveCount)
	return ret
}

func initLiveRoomIncomeUnsettledArchive() {
	syndb.RegQuick(TbLiveRoomIncomeUnsettledArchive, db.CreatedAtName)
	syndb.RegQuick(TbLiveRoomIncomeUnsettledArchive, db.UpdatedAtName)
	syndb.RegQuick(TbLiveRoomIncomeUnsettledArchive, LiveRoomIncomeUnsettledArchiveRoomId)
	syndb.RegQuick(TbLiveRoomIncomeUnsettledArchive, LiveRoomIncomeUnsettledArchiveGuildId)
	regLiveRoomIncomeCols(TbLiveRoomIncomeUnsettledArchive)
	migrate.AutoMigrate(&LiveRoomIncomeUnsettledArchive{})
}
