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
	TbWeeklyLiveFirstFrameRenderStat db.TbName = "weekly_live_first_frame_render_stats"
)

const WeeklyLiveFirstFrameRenderStatCount db.TbCol = "count"

// WeeklyLiveFirstFrameRenderStat 每周直播首帧渲染完成次数(主键=YYYY-WW)
type WeeklyLiveFirstFrameRenderStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func FormatWeeklyLiveFirstFrameRenderStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyLiveFirstFrameRenderStat(week string) *WeeklyLiveFirstFrameRenderStat {
	return &WeeklyLiveFirstFrameRenderStat{ID: week, Count: 0}
}

func (r *WeeklyLiveFirstFrameRenderStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyLiveFirstFrameRenderStat, WeeklyLiveFirstFrameRenderStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyLiveFirstFrameRenderStat() {
	syndb.RegLazy(TbWeeklyLiveFirstFrameRenderStat, WeeklyLiveFirstFrameRenderStatCount)
	migrate.AutoMigrate(&WeeklyLiveFirstFrameRenderStat{})
}
