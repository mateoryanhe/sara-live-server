package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyVoiceRoomGiftClickStat db.TbName = "daily_voice_room_gift_click_stats"
)

const DailyVoiceRoomGiftClickStatCount db.TbCol = "count"

type DailyVoiceRoomGiftClickStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func NewDailyVoiceRoomGiftClickStat(date string) *DailyVoiceRoomGiftClickStat {
	return &DailyVoiceRoomGiftClickStat{ID: date, Count: 0}
}

func (r *DailyVoiceRoomGiftClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyVoiceRoomGiftClickStat, DailyVoiceRoomGiftClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyVoiceRoomGiftClickStat() {
	syndb.RegLazy(TbDailyVoiceRoomGiftClickStat, DailyVoiceRoomGiftClickStatCount)
	migrate.AutoMigrate(&DailyVoiceRoomGiftClickStat{})
}
