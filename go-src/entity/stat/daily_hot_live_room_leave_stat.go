package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyHotLiveRoomLeaveStat db.TbName = "daily_hot_live_room_leave_stats"
)

const DailyHotLiveRoomLeaveStatCount db.TbCol = "count"

// DailyHotLiveRoomLeaveStat 每日退出 Hot 直播间次数(主键=YYYY-MM-DD)
type DailyHotLiveRoomLeaveStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:退出次数" json:"count"`
}

func NewDailyHotLiveRoomLeaveStat(date string) *DailyHotLiveRoomLeaveStat {
	return &DailyHotLiveRoomLeaveStat{ID: date, Count: 0}
}

func (r *DailyHotLiveRoomLeaveStat) AddLeaveCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyHotLiveRoomLeaveStat, DailyHotLiveRoomLeaveStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyHotLiveRoomLeaveStat() {
	syndb.RegLazy(TbDailyHotLiveRoomLeaveStat, DailyHotLiveRoomLeaveStatCount)
	migrate.AutoMigrate(&DailyHotLiveRoomLeaveStat{})
}
