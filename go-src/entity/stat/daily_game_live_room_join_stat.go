package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyGameLiveRoomJoinStat db.TbName = "daily_game_live_room_join_stats"
)

const DailyGameLiveRoomJoinStatCount db.TbCol = "count"

// DailyGameLiveRoomJoinStat 每日进入游戏直播类直播间次数(主键=YYYY-MM-DD)
type DailyGameLiveRoomJoinStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:进入次数" json:"count"`
}

func NewDailyGameLiveRoomJoinStat(date string) *DailyGameLiveRoomJoinStat {
	return &DailyGameLiveRoomJoinStat{ID: date, Count: 0}
}

func (r *DailyGameLiveRoomJoinStat) AddJoinCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyGameLiveRoomJoinStat, DailyGameLiveRoomJoinStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyGameLiveRoomJoinStat() {
	syndb.RegLazy(TbDailyGameLiveRoomJoinStat, DailyGameLiveRoomJoinStatCount)
	migrate.AutoMigrate(&DailyGameLiveRoomJoinStat{})
}
