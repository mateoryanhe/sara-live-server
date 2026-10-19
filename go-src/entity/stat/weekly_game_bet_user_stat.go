package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbWeeklyGameBetUserStat db.TbName = "weekly_game_bet_user_stats"
)

const WeeklyGameBetUserStatUserCount db.TbCol = "user_count"

// WeeklyGameBetUserStat 每周游戏下注人数(去重,主键=YYYY-WW)
type WeeklyGameBetUserStat struct {
	ID        string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	UserCount uint64 `gorm:"default:0;comment:下注人数(去重)" json:"userCount"`
}

func NewWeeklyGameBetUserStat(week string) *WeeklyGameBetUserStat {
	return &WeeklyGameBetUserStat{ID: week, UserCount: 0}
}

func (r *WeeklyGameBetUserStat) AddGameBetUserCount(n uint64) {
	r.UserCount = math.Add(r.UserCount, n)
	syndb.AddData(TbWeeklyGameBetUserStat, WeeklyGameBetUserStatUserCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.UserCount,
	})
}

func initWeeklyGameBetUserStat() {
	syndb.RegLazy(TbWeeklyGameBetUserStat, WeeklyGameBetUserStatUserCount)
	migrate.AutoMigrate(&WeeklyGameBetUserStat{})
}
