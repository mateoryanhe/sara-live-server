package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyShortVideoPlayStartStat db.TbName = "daily_short_video_play_start_stats"
)

const DailyShortVideoPlayStartStatCount db.TbCol = "count"

type DailyShortVideoPlayStartStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func NewDailyShortVideoPlayStartStat(date string) *DailyShortVideoPlayStartStat {
	return &DailyShortVideoPlayStartStat{ID: date, Count: 0}
}

func (r *DailyShortVideoPlayStartStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyShortVideoPlayStartStat, DailyShortVideoPlayStartStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyShortVideoPlayStartStat() {
	syndb.RegLazy(TbDailyShortVideoPlayStartStat, DailyShortVideoPlayStartStatCount)
	migrate.AutoMigrate(&DailyShortVideoPlayStartStat{})
}
