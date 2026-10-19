package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyGameOnlinePeakStat db.TbName = "daily_game_online_peak_stats"
)

const DailyGameOnlinePeakStatPeakCount db.TbCol = "peak_count"

// DailyGameOnlinePeakStat 每日游戏在线人数峰值(连续下注视为在线, 1 分钟无下注离线)
type DailyGameOnlinePeakStat struct {
	ID        string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	PeakCount uint64 `gorm:"default:0;comment:当日峰值在线人数" json:"peakCount"`
}

func NewDailyGameOnlinePeakStat(date string) *DailyGameOnlinePeakStat {
	return &DailyGameOnlinePeakStat{ID: date, PeakCount: 0}
}

func (r *DailyGameOnlinePeakStat) RaisePeakCount(n uint64) {
	if n <= r.PeakCount {
		return
	}
	r.PeakCount = n
	syndb.AddData(TbDailyGameOnlinePeakStat, DailyGameOnlinePeakStatPeakCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.PeakCount,
	})
}

func initDailyGameOnlinePeakStat() {
	syndb.RegLazy(TbDailyGameOnlinePeakStat, DailyGameOnlinePeakStatPeakCount)
	migrate.AutoMigrate(&DailyGameOnlinePeakStat{})
}
