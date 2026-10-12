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
	TbWeeklyShortVideoPaySuccessStat db.TbName = "weekly_short_video_pay_success_stats"
)

const WeeklyShortVideoPaySuccessStatCount db.TbCol = "count"

type WeeklyShortVideoPaySuccessStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func FormatWeeklyShortVideoPaySuccessStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyShortVideoPaySuccessStat(week string) *WeeklyShortVideoPaySuccessStat {
	return &WeeklyShortVideoPaySuccessStat{ID: week, Count: 0}
}

func (r *WeeklyShortVideoPaySuccessStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyShortVideoPaySuccessStat, WeeklyShortVideoPaySuccessStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyShortVideoPaySuccessStat() {
	syndb.RegLazy(TbWeeklyShortVideoPaySuccessStat, WeeklyShortVideoPaySuccessStatCount)
	migrate.AutoMigrate(&WeeklyShortVideoPaySuccessStat{})
}
