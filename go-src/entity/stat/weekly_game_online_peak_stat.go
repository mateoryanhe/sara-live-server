package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbWeeklyGameOnlinePeakStat db.TbName = "weekly_game_online_peak_stats"
)

const WeeklyGameOnlinePeakStatPeakCount db.TbCol = "peak_count"

// WeeklyGameOnlinePeakStat 每周游戏在线人数峰值
type WeeklyGameOnlinePeakStat struct {
	ID        string `gorm:"primaryKey;size:8;comment:周(YYYY-WW)" json:"week"`
	PeakCount uint64 `gorm:"default:0;comment:当周峰值在线人数" json:"peakCount"`
}

func NewWeeklyGameOnlinePeakStat(week string) *WeeklyGameOnlinePeakStat {
	return &WeeklyGameOnlinePeakStat{ID: week, PeakCount: 0}
}

func (r *WeeklyGameOnlinePeakStat) RaisePeakCount(n uint64) {
	if n <= r.PeakCount {
		return
	}
	r.PeakCount = n
	syndb.AddData(TbWeeklyGameOnlinePeakStat, WeeklyGameOnlinePeakStatPeakCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.PeakCount,
	})
}

func initWeeklyGameOnlinePeakStat() {
	syndb.RegLazy(TbWeeklyGameOnlinePeakStat, WeeklyGameOnlinePeakStatPeakCount)
	migrate.AutoMigrate(&WeeklyGameOnlinePeakStat{})
}
