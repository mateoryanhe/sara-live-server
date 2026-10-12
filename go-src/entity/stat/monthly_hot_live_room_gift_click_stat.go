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
	TbMonthlyHotLiveRoomGiftClickStat db.TbName = "monthly_hot_live_room_gift_click_stats"
)

const MonthlyHotLiveRoomGiftClickStatCount db.TbCol = "count"

type MonthlyHotLiveRoomGiftClickStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识" json:"month"`
	Count uint64 `gorm:"default:0;comment:曝光次数" json:"count"`
}

func FormatMonthlyHotLiveRoomGiftClickStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyHotLiveRoomGiftClickStat(month string) *MonthlyHotLiveRoomGiftClickStat {
	return &MonthlyHotLiveRoomGiftClickStat{ID: month, Count: 0}
}

func (r *MonthlyHotLiveRoomGiftClickStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyHotLiveRoomGiftClickStat, MonthlyHotLiveRoomGiftClickStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyHotLiveRoomGiftClickStat() {
	syndb.RegLazy(TbMonthlyHotLiveRoomGiftClickStat, MonthlyHotLiveRoomGiftClickStatCount)
	migrate.AutoMigrate(&MonthlyHotLiveRoomGiftClickStat{})
}
