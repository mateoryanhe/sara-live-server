package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyVoiceRoomMicRequestStat db.TbName = "daily_voice_room_mic_request_stats"
)

const DailyVoiceRoomMicRequestStatCount db.TbCol = "count"

type DailyVoiceRoomMicRequestStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func NewDailyVoiceRoomMicRequestStat(date string) *DailyVoiceRoomMicRequestStat {
	return &DailyVoiceRoomMicRequestStat{ID: date, Count: 0}
}

func (r *DailyVoiceRoomMicRequestStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyVoiceRoomMicRequestStat, DailyVoiceRoomMicRequestStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyVoiceRoomMicRequestStat() {
	syndb.RegLazy(TbDailyVoiceRoomMicRequestStat, DailyVoiceRoomMicRequestStatCount)
	migrate.AutoMigrate(&DailyVoiceRoomMicRequestStat{})
}
