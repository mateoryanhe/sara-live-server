package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbMonthlyMiniGamePlusClickStat db.TbName = "monthly_mini_game_plus_click_stats"
)

const MonthlyMiniGamePlusClickStatCount db.TbCol = "count"

type MonthlyMiniGamePlusClickStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:月标识(YYYY-MM)" json:"month"`
	Count uint64 `gorm:"default:0;comment:点击游戏手柄次数" json:"count"`
}

func FormatMonthlyMiniGamePlusClickStatKey(t time.Time) string {
	return t.Format("2006-01")
}

func NewMonthlyMiniGamePlusClickStat(month string) *MonthlyMiniGamePlusClickStat {
	return &MonthlyMiniGamePlusClickStat{ID: month, Count: 0}
}

func (r *MonthlyMiniGamePlusClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyMiniGamePlusClickStat, MonthlyMiniGamePlusClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyMiniGamePlusClickStat() {
	syndb.RegLazy(TbMonthlyMiniGamePlusClickStat, MonthlyMiniGamePlusClickStatCount)
	migrate.AutoMigrate(&MonthlyMiniGamePlusClickStat{})
}
