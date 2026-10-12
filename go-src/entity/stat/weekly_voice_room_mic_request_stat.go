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
	TbWeeklyVoiceRoomMicRequestStat db.TbName = "weekly_voice_room_mic_request_stats"
)

const WeeklyVoiceRoomMicRequestStatCount db.TbCol = "count"

type WeeklyVoiceRoomMicRequestStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func FormatWeeklyVoiceRoomMicRequestStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyVoiceRoomMicRequestStat(week string) *WeeklyVoiceRoomMicRequestStat {
	return &WeeklyVoiceRoomMicRequestStat{ID: week, Count: 0}
}

func (r *WeeklyVoiceRoomMicRequestStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyVoiceRoomMicRequestStat, WeeklyVoiceRoomMicRequestStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyVoiceRoomMicRequestStat() {
	syndb.RegLazy(TbWeeklyVoiceRoomMicRequestStat, WeeklyVoiceRoomMicRequestStatCount)
	migrate.AutoMigrate(&WeeklyVoiceRoomMicRequestStat{})
}
