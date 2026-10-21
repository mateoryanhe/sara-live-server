package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyMiniGamePlusClickStat db.TbName = "daily_mini_game_plus_click_stats"
)

const DailyMiniGamePlusClickStatCount db.TbCol = "count"

type DailyMiniGamePlusClickStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:点击游戏手柄次数" json:"count"`
}

func NewDailyMiniGamePlusClickStat(date string) *DailyMiniGamePlusClickStat {
	return &DailyMiniGamePlusClickStat{ID: date, Count: 0}
}

func (r *DailyMiniGamePlusClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyMiniGamePlusClickStat, DailyMiniGamePlusClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyMiniGamePlusClickStat() {
	syndb.RegLazy(TbDailyMiniGamePlusClickStat, DailyMiniGamePlusClickStatCount)
	migrate.AutoMigrate(&DailyMiniGamePlusClickStat{})
}
