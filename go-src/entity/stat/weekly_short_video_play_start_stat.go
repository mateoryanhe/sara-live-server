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
	TbWeeklyShortVideoPlayStartStat db.TbName = "weekly_short_video_play_start_stats"
)

const WeeklyShortVideoPlayStartStatCount db.TbCol = "count"

type WeeklyShortVideoPlayStartStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func FormatWeeklyShortVideoPlayStartStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyShortVideoPlayStartStat(week string) *WeeklyShortVideoPlayStartStat {
	return &WeeklyShortVideoPlayStartStat{ID: week, Count: 0}
}

func (r *WeeklyShortVideoPlayStartStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyShortVideoPlayStartStat, WeeklyShortVideoPlayStartStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyShortVideoPlayStartStat() {
	syndb.RegLazy(TbWeeklyShortVideoPlayStartStat, WeeklyShortVideoPlayStartStatCount)
	migrate.AutoMigrate(&WeeklyShortVideoPlayStartStat{})
}
