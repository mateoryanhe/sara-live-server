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
	TbMonthlyCall1v1ConnectSuccessStat db.TbName = "monthly_call_1v1_connect_success_stats"
)

const MonthlyCall1v1ConnectSuccessStatCount db.TbCol = "count"

type MonthlyCall1v1ConnectSuccessStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识" json:"month"`
	Count uint64 `gorm:"default:0;comment:接通次数" json:"count"`
}

func FormatMonthlyCall1v1ConnectSuccessStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyCall1v1ConnectSuccessStat(month string) *MonthlyCall1v1ConnectSuccessStat {
	return &MonthlyCall1v1ConnectSuccessStat{ID: month, Count: 0}
}

func (r *MonthlyCall1v1ConnectSuccessStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyCall1v1ConnectSuccessStat, MonthlyCall1v1ConnectSuccessStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyCall1v1ConnectSuccessStat() {
	syndb.RegLazy(TbMonthlyCall1v1ConnectSuccessStat, MonthlyCall1v1ConnectSuccessStatCount)
	migrate.AutoMigrate(&MonthlyCall1v1ConnectSuccessStat{})
}
