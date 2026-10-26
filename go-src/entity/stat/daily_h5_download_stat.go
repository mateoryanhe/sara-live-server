package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyH5DownloadStat db.TbName = "daily_h5_download_stats"
)

const DailyH5DownloadStatCount db.TbCol = "count"

type DailyH5DownloadStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func NewDailyH5DownloadStat(date string) *DailyH5DownloadStat {
	return &DailyH5DownloadStat{ID: date, Count: 0}
}

func (r *DailyH5DownloadStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyH5DownloadStat, DailyH5DownloadStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyH5DownloadStat() {
	syndb.RegLazy(TbDailyH5DownloadStat, DailyH5DownloadStatCount)
	migrate.AutoMigrate(&DailyH5DownloadStat{})
}
