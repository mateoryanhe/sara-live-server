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
	TbWeeklyChargeClickStat db.TbName = "weekly_charge_click_stats"
)

const WeeklyChargeClickStatCount db.TbCol = "count"

type WeeklyChargeClickStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func FormatWeeklyChargeClickStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyChargeClickStat(week string) *WeeklyChargeClickStat {
	return &WeeklyChargeClickStat{ID: week, Count: 0}
}

func (r *WeeklyChargeClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyChargeClickStat, WeeklyChargeClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyChargeClickStat() {
	syndb.RegLazy(TbWeeklyChargeClickStat, WeeklyChargeClickStatCount)
	migrate.AutoMigrate(&WeeklyChargeClickStat{})
}
