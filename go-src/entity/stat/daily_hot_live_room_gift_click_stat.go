package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyHotLiveRoomGiftClickStat db.TbName = "daily_hot_live_room_gift_click_stats"
)

const DailyHotLiveRoomGiftClickStatCount db.TbCol = "count"

type DailyHotLiveRoomGiftClickStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func NewDailyHotLiveRoomGiftClickStat(date string) *DailyHotLiveRoomGiftClickStat {
	return &DailyHotLiveRoomGiftClickStat{ID: date, Count: 0}
}

func (r *DailyHotLiveRoomGiftClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyHotLiveRoomGiftClickStat, DailyHotLiveRoomGiftClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyHotLiveRoomGiftClickStat() {
	syndb.RegLazy(TbDailyHotLiveRoomGiftClickStat, DailyHotLiveRoomGiftClickStatCount)
	migrate.AutoMigrate(&DailyHotLiveRoomGiftClickStat{})
}
