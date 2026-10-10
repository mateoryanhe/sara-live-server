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
	TbMonthlyHotLiveRoomLeaveStat db.TbName = "monthly_hot_live_room_leave_stats"
)

const MonthlyHotLiveRoomLeaveStatCount db.TbCol = "count"

// MonthlyHotLiveRoomLeaveStat 每月退出 Hot 直播间次数(主键=YYYY-MM)
type MonthlyHotLiveRoomLeaveStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识" json:"month"`
	Count uint64 `gorm:"default:0;comment:退出次数" json:"count"`
}

func FormatMonthlyHotLiveRoomLeaveStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyHotLiveRoomLeaveStat(month string) *MonthlyHotLiveRoomLeaveStat {
	return &MonthlyHotLiveRoomLeaveStat{ID: month, Count: 0}
}

func (r *MonthlyHotLiveRoomLeaveStat) AddLeaveCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyHotLiveRoomLeaveStat, MonthlyHotLiveRoomLeaveStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyHotLiveRoomLeaveStat() {
	syndb.RegLazy(TbMonthlyHotLiveRoomLeaveStat, MonthlyHotLiveRoomLeaveStatCount)
	migrate.AutoMigrate(&MonthlyHotLiveRoomLeaveStat{})
}
