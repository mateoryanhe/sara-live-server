package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyShortVideoPaySuccessStat db.TbName = "daily_short_video_pay_success_stats"
)

const DailyShortVideoPaySuccessStatCount db.TbCol = "count"

type DailyShortVideoPaySuccessStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func NewDailyShortVideoPaySuccessStat(date string) *DailyShortVideoPaySuccessStat {
	return &DailyShortVideoPaySuccessStat{ID: date, Count: 0}
}

func (r *DailyShortVideoPaySuccessStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyShortVideoPaySuccessStat, DailyShortVideoPaySuccessStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyShortVideoPaySuccessStat() {
	syndb.RegLazy(TbDailyShortVideoPaySuccessStat, DailyShortVideoPaySuccessStatCount)
	migrate.AutoMigrate(&DailyShortVideoPaySuccessStat{})
}
