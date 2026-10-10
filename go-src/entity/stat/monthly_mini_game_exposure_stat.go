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
	TbMonthlyMiniGameExposureStat db.TbName = "monthly_mini_game_exposure_stats"
)

const MonthlyMiniGameExposureStatCount db.TbCol = "count"

type MonthlyMiniGameExposureStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识" json:"month"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func FormatMonthlyMiniGameExposureStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyMiniGameExposureStat(month string) *MonthlyMiniGameExposureStat {
	return &MonthlyMiniGameExposureStat{ID: month, Count: 0}
}

func (r *MonthlyMiniGameExposureStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyMiniGameExposureStat, MonthlyMiniGameExposureStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyMiniGameExposureStat() {
	syndb.RegLazy(TbMonthlyMiniGameExposureStat, MonthlyMiniGameExposureStatCount)
	migrate.AutoMigrate(&MonthlyMiniGameExposureStat{})
}
