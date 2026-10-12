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
	TbMonthlyVoiceRoomMicRequestStat db.TbName = "monthly_voice_room_mic_request_stats"
)

const MonthlyVoiceRoomMicRequestStatCount db.TbCol = "count"

type MonthlyVoiceRoomMicRequestStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识" json:"month"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func FormatMonthlyVoiceRoomMicRequestStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyVoiceRoomMicRequestStat(month string) *MonthlyVoiceRoomMicRequestStat {
	return &MonthlyVoiceRoomMicRequestStat{ID: month, Count: 0}
}

func (r *MonthlyVoiceRoomMicRequestStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyVoiceRoomMicRequestStat, MonthlyVoiceRoomMicRequestStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyVoiceRoomMicRequestStat() {
	syndb.RegLazy(TbMonthlyVoiceRoomMicRequestStat, MonthlyVoiceRoomMicRequestStatCount)
	migrate.AutoMigrate(&MonthlyVoiceRoomMicRequestStat{})
}
