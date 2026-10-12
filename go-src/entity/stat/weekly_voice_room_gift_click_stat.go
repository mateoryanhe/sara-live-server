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
	TbWeeklyVoiceRoomGiftClickStat db.TbName = "weekly_voice_room_gift_click_stats"
)

const WeeklyVoiceRoomGiftClickStatCount db.TbCol = "count"

type WeeklyVoiceRoomGiftClickStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func FormatWeeklyVoiceRoomGiftClickStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyVoiceRoomGiftClickStat(week string) *WeeklyVoiceRoomGiftClickStat {
	return &WeeklyVoiceRoomGiftClickStat{ID: week, Count: 0}
}

func (r *WeeklyVoiceRoomGiftClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyVoiceRoomGiftClickStat, WeeklyVoiceRoomGiftClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyVoiceRoomGiftClickStat() {
	syndb.RegLazy(TbWeeklyVoiceRoomGiftClickStat, WeeklyVoiceRoomGiftClickStatCount)
	migrate.AutoMigrate(&WeeklyVoiceRoomGiftClickStat{})
}
