package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const TbDailyUserGameBet db.TbName = "daily_user_game_bets"

const (
	DailyUserGameBetUserId db.TbCol = "user_id"
	DailyUserGameBetDate   db.TbCol = "bet_date"
)

// DailyUserGameBet 每日用户游戏下注去重(主键 ID = date_userId)
type DailyUserGameBet struct {
	ID      string     `gorm:"primaryKey;size:32;comment:复合ID(date_userId)" json:"id"`
	UserId  uint64     `gorm:"index;default:0;comment:用户ID" json:"userId"`
	BetDate string     `gorm:"size:10;index;comment:下注日期YYYY-MM-DD" json:"betDate"`
	CreatedAt *time.Time `json:"createdAt"`
}

func BuildDailyUserGameBetId(date string, userId uint64) string {
	return fmt.Sprintf("%s_%d", date, userId)
}

func NewDailyUserGameBet(date string, userId uint64) *DailyUserGameBet {
	r := &DailyUserGameBet{}
	r.ID = BuildDailyUserGameBetId(date, userId)
	r.SetUserId(userId)
	r.SetBetDate(date)
	return r
}

func (r *DailyUserGameBet) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddData(TbDailyUserGameBet, DailyUserGameBetUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *DailyUserGameBet) SetBetDate(v string) {
	r.BetDate = v
	syndb.AddData(TbDailyUserGameBet, DailyUserGameBetDate, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *DailyUserGameBet) SetCreatedAt(v *time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbDailyUserGameBet, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initDailyUserGameBet() {
	syndb.RegLazy(TbDailyUserGameBet, db.CreatedAtName)
	syndb.RegLazy(TbDailyUserGameBet, DailyUserGameBetUserId)
	syndb.RegLazy(TbDailyUserGameBet, DailyUserGameBetDate)
	migrate.AutoMigrate(&DailyUserGameBet{})
}
