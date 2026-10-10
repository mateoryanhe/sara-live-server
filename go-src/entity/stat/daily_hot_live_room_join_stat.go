package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyHotLiveRoomJoinStat db.TbName = "daily_hot_live_room_join_stats"
)

const DailyHotLiveRoomJoinStatCount db.TbCol = "count"

// DailyHotLiveRoomJoinStat 每日进入 Hot 直播间次数(主键=YYYY-MM-DD)
type DailyHotLiveRoomJoinStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:进入次数" json:"count"`
}

func NewDailyHotLiveRoomJoinStat(date string) *DailyHotLiveRoomJoinStat {
	return &DailyHotLiveRoomJoinStat{ID: date, Count: 0}
}

func (r *DailyHotLiveRoomJoinStat) AddJoinCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyHotLiveRoomJoinStat, DailyHotLiveRoomJoinStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyHotLiveRoomJoinStat() {
	syndb.RegLazy(TbDailyHotLiveRoomJoinStat, DailyHotLiveRoomJoinStatCount)
	migrate.AutoMigrate(&DailyHotLiveRoomJoinStat{})
}
