package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyChargeClickStat db.TbName = "daily_charge_click_stats"
)

const DailyChargeClickStatCount db.TbCol = "count"

type DailyChargeClickStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func NewDailyChargeClickStat(date string) *DailyChargeClickStat {
	return &DailyChargeClickStat{ID: date, Count: 0}
}

func (r *DailyChargeClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyChargeClickStat, DailyChargeClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyChargeClickStat() {
	syndb.RegLazy(TbDailyChargeClickStat, DailyChargeClickStatCount)
	migrate.AutoMigrate(&DailyChargeClickStat{})
}
