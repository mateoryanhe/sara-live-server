package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyOfflineGameRoomPlayStat db.TbName = "daily_offline_game_room_play_stats"
)

const DailyOfflineGameRoomPlayStatCount db.TbCol = "count"

// DailyOfflineGameRoomPlayStat 每日离线游戏房开始玩游戏次数(主键=YYYY-MM-DD)
type DailyOfflineGameRoomPlayStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func NewDailyOfflineGameRoomPlayStat(date string) *DailyOfflineGameRoomPlayStat {
	return &DailyOfflineGameRoomPlayStat{ID: date, Count: 0}
}

func (r *DailyOfflineGameRoomPlayStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyOfflineGameRoomPlayStat, DailyOfflineGameRoomPlayStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyOfflineGameRoomPlayStat() {
	syndb.RegLazy(TbDailyOfflineGameRoomPlayStat, DailyOfflineGameRoomPlayStatCount)
	migrate.AutoMigrate(&DailyOfflineGameRoomPlayStat{})
}
