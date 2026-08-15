package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const TbWeeklyUserGoldConsume db.TbName = "weekly_user_gold_consumes"

const (
	WeeklyUserGoldConsumeUserId      db.TbCol = "user_id"
	WeeklyUserGoldConsumeConsumeWeek db.TbCol = "consume_week"
)

type WeeklyUserGoldConsume struct {
	ID          string     `gorm:"primaryKey;size:32" json:"id"`
	UserId      uint64     `gorm:"index;default:0" json:"userId"`
	ConsumeWeek string     `gorm:"size:16;index" json:"consumeWeek"`
	CreatedAt   *time.Time `json:"createdAt"`
}

func BuildWeeklyUserGoldConsumeId(week string, userId uint64) string {
	return fmt.Sprintf("%s_%d", week, userId)
}

func NewWeeklyUserGoldConsume(week string, userId uint64) *WeeklyUserGoldConsume {
	r := &WeeklyUserGoldConsume{}
	r.ID = BuildWeeklyUserGoldConsumeId(week, userId)
	r.SetUserId(userId)
	r.SetConsumeWeek(week)
	return r
}

func (r *WeeklyUserGoldConsume) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddData(TbWeeklyUserGoldConsume, WeeklyUserGoldConsumeUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *WeeklyUserGoldConsume) SetConsumeWeek(v string) {
	r.ConsumeWeek = v
	syndb.AddData(TbWeeklyUserGoldConsume, WeeklyUserGoldConsumeConsumeWeek, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *WeeklyUserGoldConsume) SetCreatedAt(v *time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbWeeklyUserGoldConsume, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initWeeklyUserGoldConsume() {
	syndb.RegLazy(TbWeeklyUserGoldConsume, db.CreatedAtName)
	syndb.RegLazy(TbWeeklyUserGoldConsume, WeeklyUserGoldConsumeUserId)
	syndb.RegLazy(TbWeeklyUserGoldConsume, WeeklyUserGoldConsumeConsumeWeek)
	migrate.AutoMigrate(&WeeklyUserGoldConsume{})
}
