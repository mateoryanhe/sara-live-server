package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbMonthlyMiniGameControllerClickStat db.TbName = "monthly_mini_game_controller_click_stats"
)

const MonthlyMiniGameControllerClickStatCount db.TbCol = "count"

type MonthlyMiniGameControllerClickStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:月标识(YYYY-MM)" json:"month"`
	Count uint64 `gorm:"default:0;comment:点击游戏手柄次数" json:"count"`
}

func FormatMonthlyMiniGameControllerClickStatKey(t time.Time) string {
	return t.Format("2006-01")
}

func NewMonthlyMiniGameControllerClickStat(month string) *MonthlyMiniGameControllerClickStat {
	return &MonthlyMiniGameControllerClickStat{ID: month, Count: 0}
}

func (r *MonthlyMiniGameControllerClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyMiniGameControllerClickStat, MonthlyMiniGameControllerClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyMiniGameControllerClickStat() {
	syndb.RegLazy(TbMonthlyMiniGameControllerClickStat, MonthlyMiniGameControllerClickStatCount)
	migrate.AutoMigrate(&MonthlyMiniGameControllerClickStat{})
}
