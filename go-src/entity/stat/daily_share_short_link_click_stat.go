package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyShareShortLinkClickStat db.TbName = "daily_share_short_link_click_stats"
)

const DailyShareShortLinkClickStatCount db.TbCol = "count"

type DailyShareShortLinkClickStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func NewDailyShareShortLinkClickStat(date string) *DailyShareShortLinkClickStat {
	return &DailyShareShortLinkClickStat{ID: date, Count: 0}
}

func (r *DailyShareShortLinkClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyShareShortLinkClickStat, DailyShareShortLinkClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyShareShortLinkClickStat() {
	syndb.RegLazy(TbDailyShareShortLinkClickStat, DailyShareShortLinkClickStatCount)
	migrate.AutoMigrate(&DailyShareShortLinkClickStat{})
}
