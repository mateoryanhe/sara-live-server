package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbMonthlyGameLiveRoomJoinStat db.TbName = "monthly_game_live_room_join_stats"
)

const MonthlyGameLiveRoomJoinStatCount db.TbCol = "count"

// MonthlyGameLiveRoomJoinStat 每月进入游戏直播类直播间次数(主键=YYYY-MM)
type MonthlyGameLiveRoomJoinStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:月标识(YYYY-MM)" json:"month"`
	Count uint64 `gorm:"default:0;comment:进入次数" json:"count"`
}

func FormatMonthlyGameLiveRoomJoinStatKey(t time.Time) string {
	return t.Format("2006-01")
}

func NewMonthlyGameLiveRoomJoinStat(month string) *MonthlyGameLiveRoomJoinStat {
	return &MonthlyGameLiveRoomJoinStat{ID: month, Count: 0}
}

func (r *MonthlyGameLiveRoomJoinStat) AddJoinCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyGameLiveRoomJoinStat, MonthlyGameLiveRoomJoinStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initMonthlyGameLiveRoomJoinStat() {
	syndb.RegLazy(TbMonthlyGameLiveRoomJoinStat, MonthlyGameLiveRoomJoinStatCount)
	migrate.AutoMigrate(&MonthlyGameLiveRoomJoinStat{})
}
