package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbMonthlyHotLiveRoomJoinStat db.TbName = "monthly_hot_live_room_join_stats"
)

const MonthlyHotLiveRoomJoinStatCount db.TbCol = "count"

// MonthlyHotLiveRoomJoinStat 每月进入 Hot 直播间次数(主键=YYYY-MM)
type MonthlyHotLiveRoomJoinStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:月标识(YYYY-MM)" json:"month"`
	Count uint64 `gorm:"default:0;comment:进入次数" json:"count"`
}

func FormatMonthlyHotLiveRoomJoinStatKey(t time.Time) string {
	return t.Format("2006-01")
}

func NewMonthlyHotLiveRoomJoinStat(month string) *MonthlyHotLiveRoomJoinStat {
	return &MonthlyHotLiveRoomJoinStat{ID: month, Count: 0}
}

func (r *MonthlyHotLiveRoomJoinStat) AddJoinCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyHotLiveRoomJoinStat, MonthlyHotLiveRoomJoinStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyHotLiveRoomJoinStat() {
	syndb.RegLazy(TbMonthlyHotLiveRoomJoinStat, MonthlyHotLiveRoomJoinStatCount)
	migrate.AutoMigrate(&MonthlyHotLiveRoomJoinStat{})
}
