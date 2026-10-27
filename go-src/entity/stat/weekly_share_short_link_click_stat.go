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
	TbWeeklyShareShortLinkClickStat db.TbName = "weekly_share_short_link_click_stats"
)

const WeeklyShareShortLinkClickStatCount db.TbCol = "count"

type WeeklyShareShortLinkClickStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func FormatWeeklyShareShortLinkClickStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyShareShortLinkClickStat(week string) *WeeklyShareShortLinkClickStat {
	return &WeeklyShareShortLinkClickStat{ID: week, Count: 0}
}

func (r *WeeklyShareShortLinkClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyShareShortLinkClickStat, WeeklyShareShortLinkClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyShareShortLinkClickStat() {
	syndb.RegLazy(TbWeeklyShareShortLinkClickStat, WeeklyShareShortLinkClickStatCount)
	migrate.AutoMigrate(&WeeklyShareShortLinkClickStat{})
}
