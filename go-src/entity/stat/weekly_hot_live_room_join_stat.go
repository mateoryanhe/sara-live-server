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
	TbWeeklyHotLiveRoomJoinStat db.TbName = "weekly_hot_live_room_join_stats"
)

const WeeklyHotLiveRoomJoinStatCount db.TbCol = "count"

// WeeklyHotLiveRoomJoinStat 每周进入 Hot 直播间次数(主键=YYYY-WW,与 weekly_login_stats 一致)
type WeeklyHotLiveRoomJoinStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:进入次数" json:"count"`
}

func FormatWeeklyHotLiveRoomJoinStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyHotLiveRoomJoinStat(week string) *WeeklyHotLiveRoomJoinStat {
	return &WeeklyHotLiveRoomJoinStat{ID: week, Count: 0}
}

func (r *WeeklyHotLiveRoomJoinStat) AddJoinCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyHotLiveRoomJoinStat, WeeklyHotLiveRoomJoinStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyHotLiveRoomJoinStat() {
	syndb.RegLazy(TbWeeklyHotLiveRoomJoinStat, WeeklyHotLiveRoomJoinStatCount)
	migrate.AutoMigrate(&WeeklyHotLiveRoomJoinStat{})
}
