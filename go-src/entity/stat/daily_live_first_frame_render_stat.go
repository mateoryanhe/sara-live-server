package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyLiveFirstFrameRenderStat db.TbName = "daily_live_first_frame_render_stats"
)

const DailyLiveFirstFrameRenderStatCount db.TbCol = "count"

// DailyLiveFirstFrameRenderStat 每日直播首帧渲染完成次数(主键=YYYY-MM-DD)
type DailyLiveFirstFrameRenderStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func NewDailyLiveFirstFrameRenderStat(date string) *DailyLiveFirstFrameRenderStat {
	return &DailyLiveFirstFrameRenderStat{ID: date, Count: 0}
}

func (r *DailyLiveFirstFrameRenderStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyLiveFirstFrameRenderStat, DailyLiveFirstFrameRenderStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyLiveFirstFrameRenderStat() {
	syndb.RegLazy(TbDailyLiveFirstFrameRenderStat, DailyLiveFirstFrameRenderStatCount)
	migrate.AutoMigrate(&DailyLiveFirstFrameRenderStat{})
}
