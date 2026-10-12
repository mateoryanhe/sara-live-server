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
	TbMonthlyShortVideoPaywallClickStat db.TbName = "monthly_short_video_paywall_click_stats"
)

const MonthlyShortVideoPaywallClickStatCount db.TbCol = "count"

type MonthlyShortVideoPaywallClickStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识" json:"month"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func FormatMonthlyShortVideoPaywallClickStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyShortVideoPaywallClickStat(month string) *MonthlyShortVideoPaywallClickStat {
	return &MonthlyShortVideoPaywallClickStat{ID: month, Count: 0}
}

func (r *MonthlyShortVideoPaywallClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyShortVideoPaywallClickStat, MonthlyShortVideoPaywallClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyShortVideoPaywallClickStat() {
	syndb.RegLazy(TbMonthlyShortVideoPaywallClickStat, MonthlyShortVideoPaywallClickStatCount)
	migrate.AutoMigrate(&MonthlyShortVideoPaywallClickStat{})
}
