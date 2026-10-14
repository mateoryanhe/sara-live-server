package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyOneToOneRoomListStat db.TbName = "daily_one_to_one_room_list_stats"
)

const DailyOneToOneRoomListStatCount db.TbCol = "count"

// DailyOneToOneRoomListStat 每日 1v1 房间列表查询次数(主键=YYYY-MM-DD)
type DailyOneToOneRoomListStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:查询次数" json:"count"`
}

func NewDailyOneToOneRoomListStat(date string) *DailyOneToOneRoomListStat {
	return &DailyOneToOneRoomListStat{ID: date, Count: 0}
}

func (r *DailyOneToOneRoomListStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyOneToOneRoomListStat, DailyOneToOneRoomListStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyOneToOneRoomListStat() {
	syndb.RegLazy(TbDailyOneToOneRoomListStat, DailyOneToOneRoomListStatCount)
	migrate.AutoMigrate(&DailyOneToOneRoomListStat{})
}
