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
	TbWeeklyShortVideoWatchEndStat db.TbName = "weekly_short_video_watch_end_stats"
)

const WeeklyShortVideoWatchEndStatCount db.TbCol = "count"

type WeeklyShortVideoWatchEndStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func FormatWeeklyShortVideoWatchEndStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyShortVideoWatchEndStat(week string) *WeeklyShortVideoWatchEndStat {
	return &WeeklyShortVideoWatchEndStat{ID: week, Count: 0}
}

func (r *WeeklyShortVideoWatchEndStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyShortVideoWatchEndStat, WeeklyShortVideoWatchEndStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyShortVideoWatchEndStat() {
	syndb.RegLazy(TbWeeklyShortVideoWatchEndStat, WeeklyShortVideoWatchEndStatCount)
	migrate.AutoMigrate(&WeeklyShortVideoWatchEndStat{})
}
