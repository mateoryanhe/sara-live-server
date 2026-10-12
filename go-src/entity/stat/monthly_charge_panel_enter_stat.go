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
	TbMonthlyChargePanelEnterStat db.TbName = "monthly_charge_panel_enter_stats"
)

const MonthlyChargePanelEnterStatCount db.TbCol = "count"

type MonthlyChargePanelEnterStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识" json:"month"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func FormatMonthlyChargePanelEnterStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyChargePanelEnterStat(month string) *MonthlyChargePanelEnterStat {
	return &MonthlyChargePanelEnterStat{ID: month, Count: 0}
}

func (r *MonthlyChargePanelEnterStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyChargePanelEnterStat, MonthlyChargePanelEnterStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyChargePanelEnterStat() {
	syndb.RegLazy(TbMonthlyChargePanelEnterStat, MonthlyChargePanelEnterStatCount)
	migrate.AutoMigrate(&MonthlyChargePanelEnterStat{})
}
