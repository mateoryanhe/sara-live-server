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
	TbMonthlyMiniGameRoundStartStat db.TbName = "monthly_mini_game_round_start_stats"
)

const MonthlyMiniGameRoundStartStatCount db.TbCol = "count"

type MonthlyMiniGameRoundStartStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识" json:"month"`
	Count uint64 `gorm:"default:0;comment:开始游戏次数" json:"count"`
}

func FormatMonthlyMiniGameRoundStartStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyMiniGameRoundStartStat(month string) *MonthlyMiniGameRoundStartStat {
	return &MonthlyMiniGameRoundStartStat{ID: month, Count: 0}
}

func (r *MonthlyMiniGameRoundStartStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyMiniGameRoundStartStat, MonthlyMiniGameRoundStartStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyMiniGameRoundStartStat() {
	syndb.RegLazy(TbMonthlyMiniGameRoundStartStat, MonthlyMiniGameRoundStartStatCount)
	migrate.AutoMigrate(&MonthlyMiniGameRoundStartStat{})
}
