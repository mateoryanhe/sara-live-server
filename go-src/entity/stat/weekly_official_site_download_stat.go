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
	TbWeeklyOfficialSiteDownloadStat db.TbName = "weekly_official_site_download_stats"
)

const WeeklyOfficialSiteDownloadStatCount db.TbCol = "count"

type WeeklyOfficialSiteDownloadStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func FormatWeeklyOfficialSiteDownloadStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyOfficialSiteDownloadStat(week string) *WeeklyOfficialSiteDownloadStat {
	return &WeeklyOfficialSiteDownloadStat{ID: week, Count: 0}
}

func (r *WeeklyOfficialSiteDownloadStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyOfficialSiteDownloadStat, WeeklyOfficialSiteDownloadStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyOfficialSiteDownloadStat() {
	syndb.RegLazy(TbWeeklyOfficialSiteDownloadStat, WeeklyOfficialSiteDownloadStatCount)
	migrate.AutoMigrate(&WeeklyOfficialSiteDownloadStat{})
}
