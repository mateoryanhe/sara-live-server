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
	TbWeeklyGameLiveRoomJoinStat db.TbName = "weekly_game_live_room_join_stats"
)

const WeeklyGameLiveRoomJoinStatCount db.TbCol = "count"

// WeeklyGameLiveRoomJoinStat 每周进入游戏直播类直播间次数(主键=YYYY-WW)
type WeeklyGameLiveRoomJoinStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:进入次数" json:"count"`
}

func FormatWeeklyGameLiveRoomJoinStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyGameLiveRoomJoinStat(week string) *WeeklyGameLiveRoomJoinStat {
	return &WeeklyGameLiveRoomJoinStat{ID: week, Count: 0}
}

func (r *WeeklyGameLiveRoomJoinStat) AddJoinCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyGameLiveRoomJoinStat, WeeklyGameLiveRoomJoinStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyGameLiveRoomJoinStat() {
	syndb.RegLazy(TbWeeklyGameLiveRoomJoinStat, WeeklyGameLiveRoomJoinStatCount)
	migrate.AutoMigrate(&WeeklyGameLiveRoomJoinStat{})
}
