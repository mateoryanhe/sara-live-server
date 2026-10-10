package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyMiniGameWebViewLoadStat db.TbName = "daily_mini_game_webview_load_stats"
)

const DailyMiniGameWebViewLoadStatCount db.TbCol = "count"

type DailyMiniGameWebViewLoadStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:加载成功次数" json:"count"`
}

func NewDailyMiniGameWebViewLoadStat(date string) *DailyMiniGameWebViewLoadStat {
	return &DailyMiniGameWebViewLoadStat{ID: date, Count: 0}
}

func (r *DailyMiniGameWebViewLoadStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyMiniGameWebViewLoadStat, DailyMiniGameWebViewLoadStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyMiniGameWebViewLoadStat() {
	syndb.RegLazy(TbDailyMiniGameWebViewLoadStat, DailyMiniGameWebViewLoadStatCount)
	migrate.AutoMigrate(&DailyMiniGameWebViewLoadStat{})
}
