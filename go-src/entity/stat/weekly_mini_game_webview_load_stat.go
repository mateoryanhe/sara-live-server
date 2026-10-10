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
	TbWeeklyMiniGameWebViewLoadStat db.TbName = "weekly_mini_game_webview_load_stats"
)

const WeeklyMiniGameWebViewLoadStatCount db.TbCol = "count"

type WeeklyMiniGameWebViewLoadStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:加载成功次数" json:"count"`
}

func FormatWeeklyMiniGameWebViewLoadStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyMiniGameWebViewLoadStat(week string) *WeeklyMiniGameWebViewLoadStat {
	return &WeeklyMiniGameWebViewLoadStat{ID: week, Count: 0}
}

func (r *WeeklyMiniGameWebViewLoadStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyMiniGameWebViewLoadStat, WeeklyMiniGameWebViewLoadStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyMiniGameWebViewLoadStat() {
	syndb.RegLazy(TbWeeklyMiniGameWebViewLoadStat, WeeklyMiniGameWebViewLoadStatCount)
	migrate.AutoMigrate(&WeeklyMiniGameWebViewLoadStat{})
}
