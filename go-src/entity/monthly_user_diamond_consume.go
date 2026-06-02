package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const TbMonthlyUserDiamondConsume db.TbName = "monthly_user_diamond_consumes"

const (
	MonthlyUserDiamondConsumeUserId       db.TbCol = "user_id"
	MonthlyUserDiamondConsumeConsumeMonth db.TbCol = "consume_month"
)

type MonthlyUserDiamondConsume struct {
	ID           string     `gorm:"primaryKey;size:32" json:"id"`
	UserId       uint64     `gorm:"index;default:0" json:"userId"`
	ConsumeMonth string     `gorm:"size:10;index" json:"consumeMonth"`
	CreatedAt    *time.Time `json:"createdAt"`
}

func BuildMonthlyUserDiamondConsumeId(month string, userId uint64) string {
	return fmt.Sprintf("%s_%d", month, userId)
}

func NewMonthlyUserDiamondConsume(month string, userId uint64) *MonthlyUserDiamondConsume {
	r := &MonthlyUserDiamondConsume{}
	r.ID = BuildMonthlyUserDiamondConsumeId(month, userId)
	r.SetUserId(userId)
	r.SetConsumeMonth(month)
	return r
}

func (r *MonthlyUserDiamondConsume) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddDataToLazyChan(TbMonthlyUserDiamondConsume, MonthlyUserDiamondConsumeUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *MonthlyUserDiamondConsume) SetConsumeMonth(v string) {
	r.ConsumeMonth = v
	syndb.AddDataToLazyChan(TbMonthlyUserDiamondConsume, MonthlyUserDiamondConsumeConsumeMonth, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *MonthlyUserDiamondConsume) SetCreatedAt(v *time.Time) {
	r.CreatedAt = v
	syndb.AddDataToLazyChan(TbMonthlyUserDiamondConsume, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initMonthlyUserDiamondConsume() {
	syndb.RegLazyWithMiddle(TbMonthlyUserDiamondConsume, db.CreatedAtName)
	syndb.RegLazyWithMiddle(TbMonthlyUserDiamondConsume, MonthlyUserDiamondConsumeUserId)
	syndb.RegLazyWithMiddle(TbMonthlyUserDiamondConsume, MonthlyUserDiamondConsumeConsumeMonth)
	migrate.AutoMigrate(&MonthlyUserDiamondConsume{})
}
