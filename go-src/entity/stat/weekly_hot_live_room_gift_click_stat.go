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
	TbWeeklyHotLiveRoomGiftClickStat db.TbName = "weekly_hot_live_room_gift_click_stats"
)

const WeeklyHotLiveRoomGiftClickStatCount db.TbCol = "count"

type WeeklyHotLiveRoomGiftClickStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func FormatWeeklyHotLiveRoomGiftClickStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyHotLiveRoomGiftClickStat(week string) *WeeklyHotLiveRoomGiftClickStat {
	return &WeeklyHotLiveRoomGiftClickStat{ID: week, Count: 0}
}

func (r *WeeklyHotLiveRoomGiftClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyHotLiveRoomGiftClickStat, WeeklyHotLiveRoomGiftClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyHotLiveRoomGiftClickStat() {
	syndb.RegLazy(TbWeeklyHotLiveRoomGiftClickStat, WeeklyHotLiveRoomGiftClickStatCount)
	migrate.AutoMigrate(&WeeklyHotLiveRoomGiftClickStat{})
}
