package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyMiniGameExposureStat db.TbName = "daily_mini_game_exposure_stats"
)

const DailyMiniGameExposureStatCount db.TbCol = "count"

type DailyMiniGameExposureStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func NewDailyMiniGameExposureStat(date string) *DailyMiniGameExposureStat {
	return &DailyMiniGameExposureStat{ID: date, Count: 0}
}

func (r *DailyMiniGameExposureStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyMiniGameExposureStat, DailyMiniGameExposureStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyMiniGameExposureStat() {
	syndb.RegLazy(TbDailyMiniGameExposureStat, DailyMiniGameExposureStatCount)
	migrate.AutoMigrate(&DailyMiniGameExposureStat{})
}
