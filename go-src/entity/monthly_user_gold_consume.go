package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const TbMonthlyUserGoldConsume db.TbName = "monthly_user_gold_consumes"

const (
	MonthlyUserGoldConsumeUserId       db.TbCol = "user_id"
	MonthlyUserGoldConsumeConsumeMonth db.TbCol = "consume_month"
)

type MonthlyUserGoldConsume struct {
	ID           string     `gorm:"primaryKey;size:32" json:"id"`
	UserId       uint64     `gorm:"index;default:0" json:"userId"`
	ConsumeMonth string     `gorm:"size:10;index" json:"consumeMonth"`
	CreatedAt    *time.Time `json:"createdAt"`
}

func BuildMonthlyUserGoldConsumeId(month string, userId uint64) string {
	return fmt.Sprintf("%s_%d", month, userId)
}

func NewMonthlyUserGoldConsume(month string, userId uint64) *MonthlyUserGoldConsume {
	r := &MonthlyUserGoldConsume{}
	r.ID = BuildMonthlyUserGoldConsumeId(month, userId)
	r.SetUserId(userId)
	r.SetConsumeMonth(month)
	return r
}

func (r *MonthlyUserGoldConsume) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddDataToLazyChan(TbMonthlyUserGoldConsume, MonthlyUserGoldConsumeUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *MonthlyUserGoldConsume) SetConsumeMonth(v string) {
	r.ConsumeMonth = v
	syndb.AddDataToLazyChan(TbMonthlyUserGoldConsume, MonthlyUserGoldConsumeConsumeMonth, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *MonthlyUserGoldConsume) SetCreatedAt(v *time.Time) {
	r.CreatedAt = v
	syndb.AddDataToLazyChan(TbMonthlyUserGoldConsume, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initMonthlyUserGoldConsume() {
	syndb.RegLazyWithMiddle(TbMonthlyUserGoldConsume, db.CreatedAtName)
	syndb.RegLazyWithMiddle(TbMonthlyUserGoldConsume, MonthlyUserGoldConsumeUserId)
	syndb.RegLazyWithMiddle(TbMonthlyUserGoldConsume, MonthlyUserGoldConsumeConsumeMonth)
	migrate.AutoMigrate(&MonthlyUserGoldConsume{})
}
