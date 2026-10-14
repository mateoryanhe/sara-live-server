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
	TbWeeklyMiniGameControllerClickStat db.TbName = "weekly_mini_game_controller_click_stats"
)

const WeeklyMiniGameControllerClickStatCount db.TbCol = "count"

type WeeklyMiniGameControllerClickStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:点击游戏手柄次数" json:"count"`
}

func FormatWeeklyMiniGameControllerClickStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyMiniGameControllerClickStat(week string) *WeeklyMiniGameControllerClickStat {
	return &WeeklyMiniGameControllerClickStat{ID: week, Count: 0}
}

func (r *WeeklyMiniGameControllerClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyMiniGameControllerClickStat, WeeklyMiniGameControllerClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyMiniGameControllerClickStat() {
	syndb.RegLazy(TbWeeklyMiniGameControllerClickStat, WeeklyMiniGameControllerClickStatCount)
	migrate.AutoMigrate(&WeeklyMiniGameControllerClickStat{})
}
