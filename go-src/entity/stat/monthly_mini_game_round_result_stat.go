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
	TbMonthlyMiniGameRoundResultStat db.TbName = "monthly_mini_game_round_result_stats"
)

const MonthlyMiniGameRoundResultStatCount db.TbCol = "count"

type MonthlyMiniGameRoundResultStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识" json:"month"`
	Count uint64 `gorm:"default:0;comment:结算次数" json:"count"`
}

func FormatMonthlyMiniGameRoundResultStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyMiniGameRoundResultStat(month string) *MonthlyMiniGameRoundResultStat {
	return &MonthlyMiniGameRoundResultStat{ID: month, Count: 0}
}

func (r *MonthlyMiniGameRoundResultStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyMiniGameRoundResultStat, MonthlyMiniGameRoundResultStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyMiniGameRoundResultStat() {
	syndb.RegLazy(TbMonthlyMiniGameRoundResultStat, MonthlyMiniGameRoundResultStatCount)
	migrate.AutoMigrate(&MonthlyMiniGameRoundResultStat{})
}
