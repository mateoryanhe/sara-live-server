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
	TbMonthlyOfficialSiteDownloadStat db.TbName = "monthly_official_site_download_stats"
)

const MonthlyOfficialSiteDownloadStatCount db.TbCol = "count"

type MonthlyOfficialSiteDownloadStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识" json:"month"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func FormatMonthlyOfficialSiteDownloadStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyOfficialSiteDownloadStat(month string) *MonthlyOfficialSiteDownloadStat {
	return &MonthlyOfficialSiteDownloadStat{ID: month, Count: 0}
}

func (r *MonthlyOfficialSiteDownloadStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyOfficialSiteDownloadStat, MonthlyOfficialSiteDownloadStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyOfficialSiteDownloadStat() {
	syndb.RegLazy(TbMonthlyOfficialSiteDownloadStat, MonthlyOfficialSiteDownloadStatCount)
	migrate.AutoMigrate(&MonthlyOfficialSiteDownloadStat{})
}
