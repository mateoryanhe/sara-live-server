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
	TbWeeklyGameLiveWatchStreamStat db.TbName = "weekly_game_live_watch_stream_stats"
)

const WeeklyGameLiveWatchStreamStatCount db.TbCol = "count"

// WeeklyGameLiveWatchStreamStat 每周游戏直播流首帧加载次数(主键=YYYY-WW)
type WeeklyGameLiveWatchStreamStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func FormatWeeklyGameLiveWatchStreamStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyGameLiveWatchStreamStat(week string) *WeeklyGameLiveWatchStreamStat {
	return &WeeklyGameLiveWatchStreamStat{ID: week, Count: 0}
}

func (r *WeeklyGameLiveWatchStreamStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyGameLiveWatchStreamStat, WeeklyGameLiveWatchStreamStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyGameLiveWatchStreamStat() {
	syndb.RegLazy(TbWeeklyGameLiveWatchStreamStat, WeeklyGameLiveWatchStreamStatCount)
	migrate.AutoMigrate(&WeeklyGameLiveWatchStreamStat{})
}
