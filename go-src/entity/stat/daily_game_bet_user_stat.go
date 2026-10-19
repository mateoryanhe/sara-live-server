package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyGameBetUserStat db.TbName = "daily_game_bet_user_stats"
)

const DailyGameBetUserStatUserCount db.TbCol = "user_count"

// DailyGameBetUserStat 每日游戏下注人数(去重,主键=YYYY-MM-DD)
type DailyGameBetUserStat struct {
	ID        string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	UserCount uint64 `gorm:"default:0;comment:下注人数(去重)" json:"userCount"`
}

func NewDailyGameBetUserStat(date string) *DailyGameBetUserStat {
	return &DailyGameBetUserStat{ID: date, UserCount: 0}
}

func (r *DailyGameBetUserStat) AddGameBetUserCount(n uint64) {
	r.UserCount = math.Add(r.UserCount, n)
	syndb.AddData(TbDailyGameBetUserStat, DailyGameBetUserStatUserCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.UserCount,
	})
}

func initDailyGameBetUserStat() {
	syndb.RegLazy(TbDailyGameBetUserStat, DailyGameBetUserStatUserCount)
	migrate.AutoMigrate(&DailyGameBetUserStat{})
}
