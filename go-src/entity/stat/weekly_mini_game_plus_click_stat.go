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
	TbWeeklyMiniGamePlusClickStat db.TbName = "weekly_mini_game_plus_click_stats"
)

const WeeklyMiniGamePlusClickStatCount db.TbCol = "count"

type WeeklyMiniGamePlusClickStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:点击游戏手柄次数" json:"count"`
}

func FormatWeeklyMiniGamePlusClickStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyMiniGamePlusClickStat(week string) *WeeklyMiniGamePlusClickStat {
	return &WeeklyMiniGamePlusClickStat{ID: week, Count: 0}
}

func (r *WeeklyMiniGamePlusClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyMiniGamePlusClickStat, WeeklyMiniGamePlusClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyMiniGamePlusClickStat() {
	syndb.RegLazy(TbWeeklyMiniGamePlusClickStat, WeeklyMiniGamePlusClickStatCount)
	migrate.AutoMigrate(&WeeklyMiniGamePlusClickStat{})
}
