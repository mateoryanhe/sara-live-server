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
	TbMonthlyOneToOneRoomListStat db.TbName = "monthly_one_to_one_room_list_stats"
)

const MonthlyOneToOneRoomListStatCount db.TbCol = "count"

type MonthlyOneToOneRoomListStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识(YYYY-MM)" json:"month"`
	Count uint64 `gorm:"default:0;comment:查询次数" json:"count"`
}

func FormatMonthlyOneToOneRoomListStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyOneToOneRoomListStat(month string) *MonthlyOneToOneRoomListStat {
	return &MonthlyOneToOneRoomListStat{ID: month, Count: 0}
}

func (r *MonthlyOneToOneRoomListStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyOneToOneRoomListStat, MonthlyOneToOneRoomListStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyOneToOneRoomListStat() {
	syndb.RegLazy(TbMonthlyOneToOneRoomListStat, MonthlyOneToOneRoomListStatCount)
	migrate.AutoMigrate(&MonthlyOneToOneRoomListStat{})
}
