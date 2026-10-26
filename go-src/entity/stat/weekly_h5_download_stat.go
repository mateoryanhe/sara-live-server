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
	TbWeeklyH5DownloadStat db.TbName = "weekly_h5_download_stats"
)

const WeeklyH5DownloadStatCount db.TbCol = "count"

type WeeklyH5DownloadStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func FormatWeeklyH5DownloadStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyH5DownloadStat(week string) *WeeklyH5DownloadStat {
	return &WeeklyH5DownloadStat{ID: week, Count: 0}
}

func (r *WeeklyH5DownloadStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyH5DownloadStat, WeeklyH5DownloadStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyH5DownloadStat() {
	syndb.RegLazy(TbWeeklyH5DownloadStat, WeeklyH5DownloadStatCount)
	migrate.AutoMigrate(&WeeklyH5DownloadStat{})
}
