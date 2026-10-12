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
	TbWeeklyVoiceRoomEnterStat db.TbName = "weekly_voice_room_enter_stats"
)

const WeeklyVoiceRoomEnterStatCount db.TbCol = "count"

// WeeklyVoiceRoomEnterStat 每周进入语聊房次数(主键=YYYY-WW)
type WeeklyVoiceRoomEnterStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:进入次数" json:"count"`
}

func FormatWeeklyVoiceRoomEnterStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyVoiceRoomEnterStat(week string) *WeeklyVoiceRoomEnterStat {
	return &WeeklyVoiceRoomEnterStat{ID: week, Count: 0}
}

func (r *WeeklyVoiceRoomEnterStat) AddJoinCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyVoiceRoomEnterStat, WeeklyVoiceRoomEnterStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyVoiceRoomEnterStat() {
	syndb.RegLazy(TbWeeklyVoiceRoomEnterStat, WeeklyVoiceRoomEnterStatCount)
	migrate.AutoMigrate(&WeeklyVoiceRoomEnterStat{})
}
