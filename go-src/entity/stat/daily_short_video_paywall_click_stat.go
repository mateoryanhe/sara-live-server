package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyShortVideoPaywallClickStat db.TbName = "daily_short_video_paywall_click_stats"
)

const DailyShortVideoPaywallClickStatCount db.TbCol = "count"

type DailyShortVideoPaywallClickStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func NewDailyShortVideoPaywallClickStat(date string) *DailyShortVideoPaywallClickStat {
	return &DailyShortVideoPaywallClickStat{ID: date, Count: 0}
}

func (r *DailyShortVideoPaywallClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyShortVideoPaywallClickStat, DailyShortVideoPaywallClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyShortVideoPaywallClickStat() {
	syndb.RegLazy(TbDailyShortVideoPaywallClickStat, DailyShortVideoPaywallClickStatCount)
	migrate.AutoMigrate(&DailyShortVideoPaywallClickStat{})
}
