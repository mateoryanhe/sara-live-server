package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyOfficialSiteDownloadStat db.TbName = "daily_official_site_download_stats"
)

const DailyOfficialSiteDownloadStatCount db.TbCol = "count"

type DailyOfficialSiteDownloadStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func NewDailyOfficialSiteDownloadStat(date string) *DailyOfficialSiteDownloadStat {
	return &DailyOfficialSiteDownloadStat{ID: date, Count: 0}
}

func (r *DailyOfficialSiteDownloadStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyOfficialSiteDownloadStat, DailyOfficialSiteDownloadStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyOfficialSiteDownloadStat() {
	syndb.RegLazy(TbDailyOfficialSiteDownloadStat, DailyOfficialSiteDownloadStatCount)
	migrate.AutoMigrate(&DailyOfficialSiteDownloadStat{})
}
