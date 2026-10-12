package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyShortVideoWatchEndStat db.TbName = "daily_short_video_watch_end_stats"
)

const DailyShortVideoWatchEndStatCount db.TbCol = "count"

type DailyShortVideoWatchEndStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func NewDailyShortVideoWatchEndStat(date string) *DailyShortVideoWatchEndStat {
	return &DailyShortVideoWatchEndStat{ID: date, Count: 0}
}

func (r *DailyShortVideoWatchEndStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyShortVideoWatchEndStat, DailyShortVideoWatchEndStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyShortVideoWatchEndStat() {
	syndb.RegLazy(TbDailyShortVideoWatchEndStat, DailyShortVideoWatchEndStatCount)
	migrate.AutoMigrate(&DailyShortVideoWatchEndStat{})
}
