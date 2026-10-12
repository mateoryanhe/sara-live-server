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
	TbWeeklyChargePanelEnterStat db.TbName = "weekly_charge_panel_enter_stats"
)

const WeeklyChargePanelEnterStatCount db.TbCol = "count"

type WeeklyChargePanelEnterStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func FormatWeeklyChargePanelEnterStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyChargePanelEnterStat(week string) *WeeklyChargePanelEnterStat {
	return &WeeklyChargePanelEnterStat{ID: week, Count: 0}
}

func (r *WeeklyChargePanelEnterStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyChargePanelEnterStat, WeeklyChargePanelEnterStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyChargePanelEnterStat() {
	syndb.RegLazy(TbWeeklyChargePanelEnterStat, WeeklyChargePanelEnterStatCount)
	migrate.AutoMigrate(&WeeklyChargePanelEnterStat{})
}
