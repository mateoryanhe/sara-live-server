package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbMonthlyMiniGameWebViewLoadStat db.TbName = "monthly_mini_game_webview_load_stats"
)

const MonthlyMiniGameWebViewLoadStatCount db.TbCol = "count"

type MonthlyMiniGameWebViewLoadStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:月标识(YYYY-MM)" json:"month"`
	Count uint64 `gorm:"default:0;comment:加载成功次数" json:"count"`
}

func FormatMonthlyMiniGameWebViewLoadStatKey(t time.Time) string {
	return t.Format("2006-01")
}

func NewMonthlyMiniGameWebViewLoadStat(month string) *MonthlyMiniGameWebViewLoadStat {
	return &MonthlyMiniGameWebViewLoadStat{ID: month, Count: 0}
}

func (r *MonthlyMiniGameWebViewLoadStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyMiniGameWebViewLoadStat, MonthlyMiniGameWebViewLoadStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyMiniGameWebViewLoadStat() {
	syndb.RegLazy(TbMonthlyMiniGameWebViewLoadStat, MonthlyMiniGameWebViewLoadStatCount)
	migrate.AutoMigrate(&MonthlyMiniGameWebViewLoadStat{})
}
