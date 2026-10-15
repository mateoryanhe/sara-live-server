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
	TbWeeklyOfflineGameRoomPlayStat db.TbName = "weekly_offline_game_room_play_stats"
)

const WeeklyOfflineGameRoomPlayStatCount db.TbCol = "count"

type WeeklyOfflineGameRoomPlayStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func FormatWeeklyOfflineGameRoomPlayStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyOfflineGameRoomPlayStat(week string) *WeeklyOfflineGameRoomPlayStat {
	return &WeeklyOfflineGameRoomPlayStat{ID: week, Count: 0}
}

func (r *WeeklyOfflineGameRoomPlayStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyOfflineGameRoomPlayStat, WeeklyOfflineGameRoomPlayStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyOfflineGameRoomPlayStat() {
	syndb.RegLazy(TbWeeklyOfflineGameRoomPlayStat, WeeklyOfflineGameRoomPlayStatCount)
	migrate.AutoMigrate(&WeeklyOfflineGameRoomPlayStat{})
}
