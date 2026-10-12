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
	TbMonthlyShortVideoPlayStartStat db.TbName = "monthly_short_video_play_start_stats"
)

const MonthlyShortVideoPlayStartStatCount db.TbCol = "count"

type MonthlyShortVideoPlayStartStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识" json:"month"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func FormatMonthlyShortVideoPlayStartStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyShortVideoPlayStartStat(month string) *MonthlyShortVideoPlayStartStat {
	return &MonthlyShortVideoPlayStartStat{ID: month, Count: 0}
}

func (r *MonthlyShortVideoPlayStartStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyShortVideoPlayStartStat, MonthlyShortVideoPlayStartStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyShortVideoPlayStartStat() {
	syndb.RegLazy(TbMonthlyShortVideoPlayStartStat, MonthlyShortVideoPlayStartStatCount)
	migrate.AutoMigrate(&MonthlyShortVideoPlayStartStat{})
}
