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
	TbWeeklyOneToOneRoomListStat db.TbName = "weekly_one_to_one_room_list_stats"
)

const WeeklyOneToOneRoomListStatCount db.TbCol = "count"

type WeeklyOneToOneRoomListStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:查询次数" json:"count"`
}

func FormatWeeklyOneToOneRoomListStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyOneToOneRoomListStat(week string) *WeeklyOneToOneRoomListStat {
	return &WeeklyOneToOneRoomListStat{ID: week, Count: 0}
}

func (r *WeeklyOneToOneRoomListStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyOneToOneRoomListStat, WeeklyOneToOneRoomListStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyOneToOneRoomListStat() {
	syndb.RegLazy(TbWeeklyOneToOneRoomListStat, WeeklyOneToOneRoomListStatCount)
	migrate.AutoMigrate(&WeeklyOneToOneRoomListStat{})
}
