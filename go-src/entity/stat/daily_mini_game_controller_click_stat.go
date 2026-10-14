package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyMiniGameControllerClickStat db.TbName = "daily_mini_game_controller_click_stats"
)

const DailyMiniGameControllerClickStatCount db.TbCol = "count"

type DailyMiniGameControllerClickStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:点击游戏手柄次数" json:"count"`
}

func NewDailyMiniGameControllerClickStat(date string) *DailyMiniGameControllerClickStat {
	return &DailyMiniGameControllerClickStat{ID: date, Count: 0}
}

func (r *DailyMiniGameControllerClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyMiniGameControllerClickStat, DailyMiniGameControllerClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyMiniGameControllerClickStat() {
	syndb.RegLazy(TbDailyMiniGameControllerClickStat, DailyMiniGameControllerClickStatCount)
	migrate.AutoMigrate(&DailyMiniGameControllerClickStat{})
}
