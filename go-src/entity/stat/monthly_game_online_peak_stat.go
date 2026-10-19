package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbMonthlyGameOnlinePeakStat db.TbName = "monthly_game_online_peak_stats"
)

const MonthlyGameOnlinePeakStatPeakCount db.TbCol = "peak_count"

// MonthlyGameOnlinePeakStat 每月游戏在线人数峰值
type MonthlyGameOnlinePeakStat struct {
	ID        string `gorm:"primaryKey;size:7;comment:月(YYYY-MM)" json:"month"`
	PeakCount uint64 `gorm:"default:0;comment:当月峰值在线人数" json:"peakCount"`
}

func NewMonthlyGameOnlinePeakStat(month string) *MonthlyGameOnlinePeakStat {
	return &MonthlyGameOnlinePeakStat{ID: month, PeakCount: 0}
}

func (r *MonthlyGameOnlinePeakStat) RaisePeakCount(n uint64) {
	if n <= r.PeakCount {
		return
	}
	r.PeakCount = n
	syndb.AddData(TbMonthlyGameOnlinePeakStat, MonthlyGameOnlinePeakStatPeakCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.PeakCount,
	})
}

func initMonthlyGameOnlinePeakStat() {
	syndb.RegLazy(TbMonthlyGameOnlinePeakStat, MonthlyGameOnlinePeakStatPeakCount)
	migrate.AutoMigrate(&MonthlyGameOnlinePeakStat{})
}
