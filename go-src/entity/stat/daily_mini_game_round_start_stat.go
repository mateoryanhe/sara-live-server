package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyMiniGameRoundStartStat db.TbName = "daily_mini_game_round_start_stats"
)

const DailyMiniGameRoundStartStatCount db.TbCol = "count"

// DailyMiniGameRoundStartStat 每日点击开始游戏次数(主键=YYYY-MM-DD)
type DailyMiniGameRoundStartStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:开始游戏次数" json:"count"`
}

func NewDailyMiniGameRoundStartStat(date string) *DailyMiniGameRoundStartStat {
	return &DailyMiniGameRoundStartStat{ID: date, Count: 0}
}

func (r *DailyMiniGameRoundStartStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyMiniGameRoundStartStat, DailyMiniGameRoundStartStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyMiniGameRoundStartStat() {
	syndb.RegLazy(TbDailyMiniGameRoundStartStat, DailyMiniGameRoundStartStatCount)
	migrate.AutoMigrate(&DailyMiniGameRoundStartStat{})
}
