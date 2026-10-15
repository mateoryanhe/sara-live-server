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
	TbMonthlyOfflineGameRoomPlayStat db.TbName = "monthly_offline_game_room_play_stats"
)

const MonthlyOfflineGameRoomPlayStatCount db.TbCol = "count"

type MonthlyOfflineGameRoomPlayStat struct {
	ID    string `gorm:"primaryKey;size:7;comment:月标识" json:"month"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func FormatMonthlyOfflineGameRoomPlayStatKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func NewMonthlyOfflineGameRoomPlayStat(month string) *MonthlyOfflineGameRoomPlayStat {
	return &MonthlyOfflineGameRoomPlayStat{ID: month, Count: 0}
}

func (r *MonthlyOfflineGameRoomPlayStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyOfflineGameRoomPlayStat, MonthlyOfflineGameRoomPlayStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyOfflineGameRoomPlayStat() {
	syndb.RegLazy(TbMonthlyOfflineGameRoomPlayStat, MonthlyOfflineGameRoomPlayStatCount)
	migrate.AutoMigrate(&MonthlyOfflineGameRoomPlayStat{})
}
