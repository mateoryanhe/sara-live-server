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
	TbMonthlyGameLiveWatchStreamStat db.TbName = "monthly_game_live_watch_stream_stats"
)

const MonthlyGameLiveWatchStreamStatCount db.TbCol = "count"

// MonthlyGameLiveWatchStreamStat 每月游戏直播流首帧加载次数(主键=YYYY-MM)
type MonthlyGameLiveWatchStreamStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识" json:"month"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func FormatMonthlyGameLiveWatchStreamStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyGameLiveWatchStreamStat(month string) *MonthlyGameLiveWatchStreamStat {
	return &MonthlyGameLiveWatchStreamStat{ID: month, Count: 0}
}

func (r *MonthlyGameLiveWatchStreamStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyGameLiveWatchStreamStat, MonthlyGameLiveWatchStreamStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyGameLiveWatchStreamStat() {
	syndb.RegLazy(TbMonthlyGameLiveWatchStreamStat, MonthlyGameLiveWatchStreamStatCount)
	migrate.AutoMigrate(&MonthlyGameLiveWatchStreamStat{})
}
