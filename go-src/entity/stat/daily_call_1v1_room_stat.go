package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyCall1v1RoomStat db.TbName = "daily_call_1v1_room_stats"
)

const DailyCall1v1RoomStatCount db.TbCol = "count"

// DailyCall1v1RoomStat 每日 1v1 房间 source=3 视频通话发起次数(主键=YYYY-MM-DD)
type DailyCall1v1RoomStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:发起次数" json:"count"`
}

func NewDailyCall1v1RoomStat(date string) *DailyCall1v1RoomStat {
	return &DailyCall1v1RoomStat{ID: date, Count: 0}
}

func (r *DailyCall1v1RoomStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyCall1v1RoomStat, DailyCall1v1RoomStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyCall1v1RoomStat() {
	syndb.RegLazy(TbDailyCall1v1RoomStat, DailyCall1v1RoomStatCount)
	migrate.AutoMigrate(&DailyCall1v1RoomStat{})
}
