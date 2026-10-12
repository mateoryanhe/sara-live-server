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
	TbMonthlyChargeResultServerStat db.TbName = "monthly_charge_result_server_stats"
)

const MonthlyChargeResultServerStatCount db.TbCol = "count"

type MonthlyChargeResultServerStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识" json:"month"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func FormatMonthlyChargeResultServerStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyChargeResultServerStat(month string) *MonthlyChargeResultServerStat {
	return &MonthlyChargeResultServerStat{ID: month, Count: 0}
}

func (r *MonthlyChargeResultServerStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyChargeResultServerStat, MonthlyChargeResultServerStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyChargeResultServerStat() {
	syndb.RegLazy(TbMonthlyChargeResultServerStat, MonthlyChargeResultServerStatCount)
	migrate.AutoMigrate(&MonthlyChargeResultServerStat{})
}
