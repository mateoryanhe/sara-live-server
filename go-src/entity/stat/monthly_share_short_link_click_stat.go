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
	TbMonthlyShareShortLinkClickStat db.TbName = "monthly_share_short_link_click_stats"
)

const MonthlyShareShortLinkClickStatCount db.TbCol = "count"

type MonthlyShareShortLinkClickStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识" json:"month"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func FormatMonthlyShareShortLinkClickStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyShareShortLinkClickStat(month string) *MonthlyShareShortLinkClickStat {
	return &MonthlyShareShortLinkClickStat{ID: month, Count: 0}
}

func (r *MonthlyShareShortLinkClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyShareShortLinkClickStat, MonthlyShareShortLinkClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyShareShortLinkClickStat() {
	syndb.RegLazy(TbMonthlyShareShortLinkClickStat, MonthlyShareShortLinkClickStatCount)
	migrate.AutoMigrate(&MonthlyShareShortLinkClickStat{})
}
