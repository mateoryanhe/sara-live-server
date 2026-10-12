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
	TbMonthlyChargeClickStat db.TbName = "monthly_charge_click_stats"
)

const MonthlyChargeClickStatCount db.TbCol = "count"

type MonthlyChargeClickStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识" json:"month"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func FormatMonthlyChargeClickStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyChargeClickStat(month string) *MonthlyChargeClickStat {
	return &MonthlyChargeClickStat{ID: month, Count: 0}
}

func (r *MonthlyChargeClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyChargeClickStat, MonthlyChargeClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyChargeClickStat() {
	syndb.RegLazy(TbMonthlyChargeClickStat, MonthlyChargeClickStatCount)
	migrate.AutoMigrate(&MonthlyChargeClickStat{})
}
