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
	TbMonthlyLiveFirstFrameRenderStat db.TbName = "monthly_live_first_frame_render_stats"
)

const MonthlyLiveFirstFrameRenderStatCount db.TbCol = "count"

// MonthlyLiveFirstFrameRenderStat 每月直播首帧渲染完成次数(主键=YYYY-MM)
type MonthlyLiveFirstFrameRenderStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识" json:"month"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func FormatMonthlyLiveFirstFrameRenderStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyLiveFirstFrameRenderStat(month string) *MonthlyLiveFirstFrameRenderStat {
	return &MonthlyLiveFirstFrameRenderStat{ID: month, Count: 0}
}

func (r *MonthlyLiveFirstFrameRenderStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyLiveFirstFrameRenderStat, MonthlyLiveFirstFrameRenderStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyLiveFirstFrameRenderStat() {
	syndb.RegLazy(TbMonthlyLiveFirstFrameRenderStat, MonthlyLiveFirstFrameRenderStatCount)
	migrate.AutoMigrate(&MonthlyLiveFirstFrameRenderStat{})
}
