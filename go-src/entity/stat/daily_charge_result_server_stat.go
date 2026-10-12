package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyChargeResultServerStat db.TbName = "daily_charge_result_server_stats"
)

const DailyChargeResultServerStatCount db.TbCol = "count"

type DailyChargeResultServerStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func NewDailyChargeResultServerStat(date string) *DailyChargeResultServerStat {
	return &DailyChargeResultServerStat{ID: date, Count: 0}
}

func (r *DailyChargeResultServerStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyChargeResultServerStat, DailyChargeResultServerStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyChargeResultServerStat() {
	syndb.RegLazy(TbDailyChargeResultServerStat, DailyChargeResultServerStatCount)
	migrate.AutoMigrate(&DailyChargeResultServerStat{})
}
