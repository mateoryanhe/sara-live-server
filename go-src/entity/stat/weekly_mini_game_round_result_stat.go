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
	TbWeeklyMiniGameRoundResultStat db.TbName = "weekly_mini_game_round_result_stats"
)

const WeeklyMiniGameRoundResultStatCount db.TbCol = "count"

type WeeklyMiniGameRoundResultStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:结算次数" json:"count"`
}

func FormatWeeklyMiniGameRoundResultStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyMiniGameRoundResultStat(week string) *WeeklyMiniGameRoundResultStat {
	return &WeeklyMiniGameRoundResultStat{ID: week, Count: 0}
}

func (r *WeeklyMiniGameRoundResultStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyMiniGameRoundResultStat, WeeklyMiniGameRoundResultStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyMiniGameRoundResultStat() {
	syndb.RegLazy(TbWeeklyMiniGameRoundResultStat, WeeklyMiniGameRoundResultStatCount)
	migrate.AutoMigrate(&WeeklyMiniGameRoundResultStat{})
}
