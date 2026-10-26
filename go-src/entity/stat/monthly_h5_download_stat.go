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
	TbMonthlyH5DownloadStat db.TbName = "monthly_h5_download_stats"
)

const MonthlyH5DownloadStatCount db.TbCol = "count"

type MonthlyH5DownloadStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识" json:"month"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func FormatMonthlyH5DownloadStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyH5DownloadStat(month string) *MonthlyH5DownloadStat {
	return &MonthlyH5DownloadStat{ID: month, Count: 0}
}

func (r *MonthlyH5DownloadStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyH5DownloadStat, MonthlyH5DownloadStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyH5DownloadStat() {
	syndb.RegLazy(TbMonthlyH5DownloadStat, MonthlyH5DownloadStatCount)
	migrate.AutoMigrate(&MonthlyH5DownloadStat{})
}
