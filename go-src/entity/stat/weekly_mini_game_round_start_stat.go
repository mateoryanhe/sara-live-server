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
	TbWeeklyMiniGameRoundStartStat db.TbName = "weekly_mini_game_round_start_stats"
)

const WeeklyMiniGameRoundStartStatCount db.TbCol = "count"

type WeeklyMiniGameRoundStartStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:开始游戏次数" json:"count"`
}

func FormatWeeklyMiniGameRoundStartStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyMiniGameRoundStartStat(week string) *WeeklyMiniGameRoundStartStat {
	return &WeeklyMiniGameRoundStartStat{ID: week, Count: 0}
}

func (r *WeeklyMiniGameRoundStartStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyMiniGameRoundStartStat, WeeklyMiniGameRoundStartStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyMiniGameRoundStartStat() {
	syndb.RegLazy(TbWeeklyMiniGameRoundStartStat, WeeklyMiniGameRoundStartStatCount)
	migrate.AutoMigrate(&WeeklyMiniGameRoundStartStat{})
}
