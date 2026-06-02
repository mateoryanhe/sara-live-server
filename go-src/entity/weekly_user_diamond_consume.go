package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const TbWeeklyUserDiamondConsume db.TbName = "weekly_user_diamond_consumes"

const (
	WeeklyUserDiamondConsumeUserId      db.TbCol = "user_id"
	WeeklyUserDiamondConsumeConsumeWeek db.TbCol = "consume_week"
)

type WeeklyUserDiamondConsume struct {
	ID          string     `gorm:"primaryKey;size:32" json:"id"`
	UserId      uint64     `gorm:"index;default:0" json:"userId"`
	ConsumeWeek string     `gorm:"size:16;index" json:"consumeWeek"`
	CreatedAt   *time.Time `json:"createdAt"`
}

func BuildWeeklyUserDiamondConsumeId(week string, userId uint64) string {
	return fmt.Sprintf("%s_%d", week, userId)
}

func NewWeeklyUserDiamondConsume(week string, userId uint64) *WeeklyUserDiamondConsume {
	r := &WeeklyUserDiamondConsume{}
	r.ID = BuildWeeklyUserDiamondConsumeId(week, userId)
	r.SetUserId(userId)
	r.SetConsumeWeek(week)
	return r
}

func (r *WeeklyUserDiamondConsume) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddDataToLazyChan(TbWeeklyUserDiamondConsume, WeeklyUserDiamondConsumeUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *WeeklyUserDiamondConsume) SetConsumeWeek(v string) {
	r.ConsumeWeek = v
	syndb.AddDataToLazyChan(TbWeeklyUserDiamondConsume, WeeklyUserDiamondConsumeConsumeWeek, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *WeeklyUserDiamondConsume) SetCreatedAt(v *time.Time) {
	r.CreatedAt = v
	syndb.AddDataToLazyChan(TbWeeklyUserDiamondConsume, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initWeeklyUserDiamondConsume() {
	syndb.RegLazyWithMiddle(TbWeeklyUserDiamondConsume, db.CreatedAtName)
	syndb.RegLazyWithMiddle(TbWeeklyUserDiamondConsume, WeeklyUserDiamondConsumeUserId)
	syndb.RegLazyWithMiddle(TbWeeklyUserDiamondConsume, WeeklyUserDiamondConsumeConsumeWeek)
	migrate.AutoMigrate(&WeeklyUserDiamondConsume{})
}
