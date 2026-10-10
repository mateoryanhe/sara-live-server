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
	TbWeeklyMiniGameExposureStat db.TbName = "weekly_mini_game_exposure_stats"
)

const WeeklyMiniGameExposureStatCount db.TbCol = "count"

type WeeklyMiniGameExposureStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func FormatWeeklyMiniGameExposureStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyMiniGameExposureStat(week string) *WeeklyMiniGameExposureStat {
	return &WeeklyMiniGameExposureStat{ID: week, Count: 0}
}

func (r *WeeklyMiniGameExposureStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyMiniGameExposureStat, WeeklyMiniGameExposureStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyMiniGameExposureStat() {
	syndb.RegLazy(TbWeeklyMiniGameExposureStat, WeeklyMiniGameExposureStatCount)
	migrate.AutoMigrate(&WeeklyMiniGameExposureStat{})
}
