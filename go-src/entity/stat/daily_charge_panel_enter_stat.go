package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyChargePanelEnterStat db.TbName = "daily_charge_panel_enter_stats"
)

const DailyChargePanelEnterStatCount db.TbCol = "count"

type DailyChargePanelEnterStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func NewDailyChargePanelEnterStat(date string) *DailyChargePanelEnterStat {
	return &DailyChargePanelEnterStat{ID: date, Count: 0}
}

func (r *DailyChargePanelEnterStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyChargePanelEnterStat, DailyChargePanelEnterStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyChargePanelEnterStat() {
	syndb.RegLazy(TbDailyChargePanelEnterStat, DailyChargePanelEnterStatCount)
	migrate.AutoMigrate(&DailyChargePanelEnterStat{})
}
