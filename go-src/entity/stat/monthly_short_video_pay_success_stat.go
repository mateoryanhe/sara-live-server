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
	TbMonthlyShortVideoPaySuccessStat db.TbName = "monthly_short_video_pay_success_stats"
)

const MonthlyShortVideoPaySuccessStatCount db.TbCol = "count"

type MonthlyShortVideoPaySuccessStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识" json:"month"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func FormatMonthlyShortVideoPaySuccessStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyShortVideoPaySuccessStat(month string) *MonthlyShortVideoPaySuccessStat {
	return &MonthlyShortVideoPaySuccessStat{ID: month, Count: 0}
}

func (r *MonthlyShortVideoPaySuccessStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyShortVideoPaySuccessStat, MonthlyShortVideoPaySuccessStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyShortVideoPaySuccessStat() {
	syndb.RegLazy(TbMonthlyShortVideoPaySuccessStat, MonthlyShortVideoPaySuccessStatCount)
	migrate.AutoMigrate(&MonthlyShortVideoPaySuccessStat{})
}
