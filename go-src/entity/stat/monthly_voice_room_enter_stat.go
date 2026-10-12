package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbMonthlyVoiceRoomEnterStat db.TbName = "monthly_voice_room_enter_stats"
)

const MonthlyVoiceRoomEnterStatCount db.TbCol = "count"

// MonthlyVoiceRoomEnterStat 每月进入语聊房次数(主键=YYYY-MM)
type MonthlyVoiceRoomEnterStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:月标识(YYYY-MM)" json:"month"`
	Count uint64 `gorm:"default:0;comment:进入次数" json:"count"`
}

func FormatMonthlyVoiceRoomEnterStatKey(t time.Time) string {
	return t.Format("2006-01")
}

func NewMonthlyVoiceRoomEnterStat(month string) *MonthlyVoiceRoomEnterStat {
	return &MonthlyVoiceRoomEnterStat{ID: month, Count: 0}
}

func (r *MonthlyVoiceRoomEnterStat) AddJoinCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyVoiceRoomEnterStat, MonthlyVoiceRoomEnterStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyVoiceRoomEnterStat() {
	syndb.RegLazy(TbMonthlyVoiceRoomEnterStat, MonthlyVoiceRoomEnterStatCount)
	migrate.AutoMigrate(&MonthlyVoiceRoomEnterStat{})
}
