package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const TbWeeklyUserGameBet db.TbName = "weekly_user_game_bets"

const (
	WeeklyUserGameBetUserId db.TbCol = "user_id"
	WeeklyUserGameBetWeek   db.TbCol = "bet_week"
)

type WeeklyUserGameBet struct {
	ID       string     `gorm:"primaryKey;size:32" json:"id"`
	UserId   uint64     `gorm:"index;default:0" json:"userId"`
	BetWeek  string     `gorm:"size:16;index" json:"betWeek"`
	CreatedAt *time.Time `json:"createdAt"`
}

func BuildWeeklyUserGameBetId(week string, userId uint64) string {
	return fmt.Sprintf("%s_%d", week, userId)
}

func NewWeeklyUserGameBet(week string, userId uint64) *WeeklyUserGameBet {
	r := &WeeklyUserGameBet{}
	r.ID = BuildWeeklyUserGameBetId(week, userId)
	r.SetUserId(userId)
	r.SetBetWeek(week)
	return r
}

func (r *WeeklyUserGameBet) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddData(TbWeeklyUserGameBet, WeeklyUserGameBetUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *WeeklyUserGameBet) SetBetWeek(v string) {
	r.BetWeek = v
	syndb.AddData(TbWeeklyUserGameBet, WeeklyUserGameBetWeek, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *WeeklyUserGameBet) SetCreatedAt(v *time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbWeeklyUserGameBet, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initWeeklyUserGameBet() {
	syndb.RegLazy(TbWeeklyUserGameBet, db.CreatedAtName)
	syndb.RegLazy(TbWeeklyUserGameBet, WeeklyUserGameBetUserId)
	syndb.RegLazy(TbWeeklyUserGameBet, WeeklyUserGameBetWeek)
	migrate.AutoMigrate(&WeeklyUserGameBet{})
}
