package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const TbMonthlyUserGameBet db.TbName = "monthly_user_game_bets"

const (
	MonthlyUserGameBetUserId db.TbCol = "user_id"
	MonthlyUserGameBetMonth  db.TbCol = "bet_month"
)

type MonthlyUserGameBet struct {
	ID        string     `gorm:"primaryKey;size:32" json:"id"`
	UserId    uint64     `gorm:"index;default:0" json:"userId"`
	BetMonth  string     `gorm:"size:10;index" json:"betMonth"`
	CreatedAt *time.Time `json:"createdAt"`
}

func BuildMonthlyUserGameBetId(month string, userId uint64) string {
	return fmt.Sprintf("%s_%d", month, userId)
}

func NewMonthlyUserGameBet(month string, userId uint64) *MonthlyUserGameBet {
	r := &MonthlyUserGameBet{}
	r.ID = BuildMonthlyUserGameBetId(month, userId)
	r.SetUserId(userId)
	r.SetBetMonth(month)
	return r
}

func (r *MonthlyUserGameBet) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddData(TbMonthlyUserGameBet, MonthlyUserGameBetUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *MonthlyUserGameBet) SetBetMonth(v string) {
	r.BetMonth = v
	syndb.AddData(TbMonthlyUserGameBet, MonthlyUserGameBetMonth, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *MonthlyUserGameBet) SetCreatedAt(v *time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbMonthlyUserGameBet, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initMonthlyUserGameBet() {
	syndb.RegLazy(TbMonthlyUserGameBet, db.CreatedAtName)
	syndb.RegLazy(TbMonthlyUserGameBet, MonthlyUserGameBetUserId)
	syndb.RegLazy(TbMonthlyUserGameBet, MonthlyUserGameBetMonth)
	migrate.AutoMigrate(&MonthlyUserGameBet{})
}
