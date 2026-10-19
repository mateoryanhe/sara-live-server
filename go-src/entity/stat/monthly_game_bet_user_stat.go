package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbMonthlyGameBetUserStat db.TbName = "monthly_game_bet_user_stats"
)

const MonthlyGameBetUserStatUserCount db.TbCol = "user_count"

// MonthlyGameBetUserStat 每月游戏下注人数(去重,主键=YYYY-MM)
type MonthlyGameBetUserStat struct {
	ID        string `gorm:"primaryKey;size:10;comment:月标识(YYYY-MM)" json:"month"`
	UserCount uint64 `gorm:"default:0;comment:下注人数(去重)" json:"userCount"`
}

func FormatMonthlyGameBetUserStatKey(t time.Time) string {
	return t.Format("2006-01")
}

func NewMonthlyGameBetUserStat(month string) *MonthlyGameBetUserStat {
	return &MonthlyGameBetUserStat{ID: month, UserCount: 0}
}

func (r *MonthlyGameBetUserStat) AddGameBetUserCount(n uint64) {
	r.UserCount = math.Add(r.UserCount, n)
	syndb.AddData(TbMonthlyGameBetUserStat, MonthlyGameBetUserStatUserCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.UserCount,
	})
}

func initMonthlyGameBetUserStat() {
	syndb.RegLazy(TbMonthlyGameBetUserStat, MonthlyGameBetUserStatUserCount)
	migrate.AutoMigrate(&MonthlyGameBetUserStat{})
}
