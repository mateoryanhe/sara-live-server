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
	TbMonthlyCall1v1InitiateStat db.TbName = "monthly_call_1v1_initiate_stats"
)

const MonthlyCall1v1InitiateStatCount db.TbCol = "count"

type MonthlyCall1v1InitiateStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识" json:"month"`
	Count uint64 `gorm:"default:0;comment:发起次数" json:"count"`
}

func FormatMonthlyCall1v1InitiateStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyCall1v1InitiateStat(month string) *MonthlyCall1v1InitiateStat {
	return &MonthlyCall1v1InitiateStat{ID: month, Count: 0}
}

func (r *MonthlyCall1v1InitiateStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyCall1v1InitiateStat, MonthlyCall1v1InitiateStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyCall1v1InitiateStat() {
	syndb.RegLazy(TbMonthlyCall1v1InitiateStat, MonthlyCall1v1InitiateStatCount)
	migrate.AutoMigrate(&MonthlyCall1v1InitiateStat{})
}
