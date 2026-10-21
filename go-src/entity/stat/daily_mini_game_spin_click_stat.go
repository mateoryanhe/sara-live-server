package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyMiniGameSpinClickStat db.TbName = "daily_mini_game_spin_click_stats"
)

const DailyMiniGameSpinClickStatCount db.TbCol = "count"

type DailyMiniGameSpinClickStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:点击游戏手柄次数" json:"count"`
}

func NewDailyMiniGameSpinClickStat(date string) *DailyMiniGameSpinClickStat {
	return &DailyMiniGameSpinClickStat{ID: date, Count: 0}
}

func (r *DailyMiniGameSpinClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyMiniGameSpinClickStat, DailyMiniGameSpinClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyMiniGameSpinClickStat() {
	syndb.RegLazy(TbDailyMiniGameSpinClickStat, DailyMiniGameSpinClickStatCount)
	migrate.AutoMigrate(&DailyMiniGameSpinClickStat{})
}
