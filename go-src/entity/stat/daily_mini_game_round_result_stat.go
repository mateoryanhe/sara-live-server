package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyMiniGameRoundResultStat db.TbName = "daily_mini_game_round_result_stats"
)

const DailyMiniGameRoundResultStatCount db.TbCol = "count"

// DailyMiniGameRoundResultStat 每日单局游戏结算次数(主键=YYYY-MM-DD)
type DailyMiniGameRoundResultStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:结算次数" json:"count"`
}

func NewDailyMiniGameRoundResultStat(date string) *DailyMiniGameRoundResultStat {
	return &DailyMiniGameRoundResultStat{ID: date, Count: 0}
}

func (r *DailyMiniGameRoundResultStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyMiniGameRoundResultStat, DailyMiniGameRoundResultStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyMiniGameRoundResultStat() {
	syndb.RegLazy(TbDailyMiniGameRoundResultStat, DailyMiniGameRoundResultStatCount)
	migrate.AutoMigrate(&DailyMiniGameRoundResultStat{})
}
