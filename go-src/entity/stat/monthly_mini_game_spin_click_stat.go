package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbMonthlyMiniGameSpinClickStat db.TbName = "monthly_mini_game_spin_click_stats"
)

const MonthlyMiniGameSpinClickStatCount db.TbCol = "count"

type MonthlyMiniGameSpinClickStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:月标识(YYYY-MM)" json:"month"`
	Count uint64 `gorm:"default:0;comment:点击游戏手柄次数" json:"count"`
}

func FormatMonthlyMiniGameSpinClickStatKey(t time.Time) string {
	return t.Format("2006-01")
}

func NewMonthlyMiniGameSpinClickStat(month string) *MonthlyMiniGameSpinClickStat {
	return &MonthlyMiniGameSpinClickStat{ID: month, Count: 0}
}

func (r *MonthlyMiniGameSpinClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyMiniGameSpinClickStat, MonthlyMiniGameSpinClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyMiniGameSpinClickStat() {
	syndb.RegLazy(TbMonthlyMiniGameSpinClickStat, MonthlyMiniGameSpinClickStatCount)
	migrate.AutoMigrate(&MonthlyMiniGameSpinClickStat{})
}
