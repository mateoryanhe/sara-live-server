package entity

import (
	"fmt"
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbWeeklyHotLiveRoomLeaveStat db.TbName = "weekly_hot_live_room_leave_stats"
)

const WeeklyHotLiveRoomLeaveStatCount db.TbCol = "count"

// WeeklyHotLiveRoomLeaveStat 每周退出 Hot 直播间次数(主键=YYYY-WW)
type WeeklyHotLiveRoomLeaveStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:退出次数" json:"count"`
}

func FormatWeeklyHotLiveRoomLeaveStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyHotLiveRoomLeaveStat(week string) *WeeklyHotLiveRoomLeaveStat {
	return &WeeklyHotLiveRoomLeaveStat{ID: week, Count: 0}
}

func (r *WeeklyHotLiveRoomLeaveStat) AddLeaveCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyHotLiveRoomLeaveStat, WeeklyHotLiveRoomLeaveStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyHotLiveRoomLeaveStat() {
	syndb.RegLazy(TbWeeklyHotLiveRoomLeaveStat, WeeklyHotLiveRoomLeaveStatCount)
	migrate.AutoMigrate(&WeeklyHotLiveRoomLeaveStat{})
}
