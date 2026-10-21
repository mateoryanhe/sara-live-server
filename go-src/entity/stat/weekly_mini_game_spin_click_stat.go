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
	TbWeeklyMiniGameSpinClickStat db.TbName = "weekly_mini_game_spin_click_stats"
)

const WeeklyMiniGameSpinClickStatCount db.TbCol = "count"

type WeeklyMiniGameSpinClickStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:点击游戏手柄次数" json:"count"`
}

func FormatWeeklyMiniGameSpinClickStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyMiniGameSpinClickStat(week string) *WeeklyMiniGameSpinClickStat {
	return &WeeklyMiniGameSpinClickStat{ID: week, Count: 0}
}

func (r *WeeklyMiniGameSpinClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyMiniGameSpinClickStat, WeeklyMiniGameSpinClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyMiniGameSpinClickStat() {
	syndb.RegLazy(TbWeeklyMiniGameSpinClickStat, WeeklyMiniGameSpinClickStatCount)
	migrate.AutoMigrate(&WeeklyMiniGameSpinClickStat{})
}
