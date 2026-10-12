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
	TbWeeklyShortVideoPaywallClickStat db.TbName = "weekly_short_video_paywall_click_stats"
)

const WeeklyShortVideoPaywallClickStatCount db.TbCol = "count"

type WeeklyShortVideoPaywallClickStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func FormatWeeklyShortVideoPaywallClickStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyShortVideoPaywallClickStat(week string) *WeeklyShortVideoPaywallClickStat {
	return &WeeklyShortVideoPaywallClickStat{ID: week, Count: 0}
}

func (r *WeeklyShortVideoPaywallClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyShortVideoPaywallClickStat, WeeklyShortVideoPaywallClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyShortVideoPaywallClickStat() {
	syndb.RegLazy(TbWeeklyShortVideoPaywallClickStat, WeeklyShortVideoPaywallClickStatCount)
	migrate.AutoMigrate(&WeeklyShortVideoPaywallClickStat{})
}
