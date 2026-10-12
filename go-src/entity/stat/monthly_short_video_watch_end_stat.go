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
	TbMonthlyShortVideoWatchEndStat db.TbName = "monthly_short_video_watch_end_stats"
)

const MonthlyShortVideoWatchEndStatCount db.TbCol = "count"

type MonthlyShortVideoWatchEndStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识" json:"month"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func FormatMonthlyShortVideoWatchEndStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyShortVideoWatchEndStat(month string) *MonthlyShortVideoWatchEndStat {
	return &MonthlyShortVideoWatchEndStat{ID: month, Count: 0}
}

func (r *MonthlyShortVideoWatchEndStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyShortVideoWatchEndStat, MonthlyShortVideoWatchEndStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyShortVideoWatchEndStat() {
	syndb.RegLazy(TbMonthlyShortVideoWatchEndStat, MonthlyShortVideoWatchEndStatCount)
	migrate.AutoMigrate(&MonthlyShortVideoWatchEndStat{})
}
