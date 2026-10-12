package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyVoiceRoomEnterStat db.TbName = "daily_voice_room_enter_stats"
)

const DailyVoiceRoomEnterStatCount db.TbCol = "count"

// DailyVoiceRoomEnterStat 每日进入语聊房次数(主键=YYYY-MM-DD)
type DailyVoiceRoomEnterStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:进入次数" json:"count"`
}

func NewDailyVoiceRoomEnterStat(date string) *DailyVoiceRoomEnterStat {
	return &DailyVoiceRoomEnterStat{ID: date, Count: 0}
}

func (r *DailyVoiceRoomEnterStat) AddJoinCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyVoiceRoomEnterStat, DailyVoiceRoomEnterStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyVoiceRoomEnterStat() {
	syndb.RegLazy(TbDailyVoiceRoomEnterStat, DailyVoiceRoomEnterStatCount)
	migrate.AutoMigrate(&DailyVoiceRoomEnterStat{})
}
