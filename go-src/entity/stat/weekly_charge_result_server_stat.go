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
	TbWeeklyChargeResultServerStat db.TbName = "weekly_charge_result_server_stats"
)

const WeeklyChargeResultServerStatCount db.TbCol = "count"

type WeeklyChargeResultServerStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func FormatWeeklyChargeResultServerStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyChargeResultServerStat(week string) *WeeklyChargeResultServerStat {
	return &WeeklyChargeResultServerStat{ID: week, Count: 0}
}

func (r *WeeklyChargeResultServerStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyChargeResultServerStat, WeeklyChargeResultServerStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyChargeResultServerStat() {
	syndb.RegLazy(TbWeeklyChargeResultServerStat, WeeklyChargeResultServerStatCount)
	migrate.AutoMigrate(&WeeklyChargeResultServerStat{})
}
