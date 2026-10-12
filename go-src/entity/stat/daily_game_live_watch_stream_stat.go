package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyGameLiveWatchStreamStat db.TbName = "daily_game_live_watch_stream_stats"
)

const DailyGameLiveWatchStreamStatCount db.TbCol = "count"

// DailyGameLiveWatchStreamStat 每日游戏直播流首帧加载次数(主键=YYYY-MM-DD)
type DailyGameLiveWatchStreamStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func NewDailyGameLiveWatchStreamStat(date string) *DailyGameLiveWatchStreamStat {
	return &DailyGameLiveWatchStreamStat{ID: date, Count: 0}
}

func (r *DailyGameLiveWatchStreamStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyGameLiveWatchStreamStat, DailyGameLiveWatchStreamStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyGameLiveWatchStreamStat() {
	syndb.RegLazy(TbDailyGameLiveWatchStreamStat, DailyGameLiveWatchStreamStatCount)
	migrate.AutoMigrate(&DailyGameLiveWatchStreamStat{})
}
