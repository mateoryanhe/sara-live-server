package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const TbDailyUserDiamondConsume db.TbName = "daily_user_diamond_consumes"

const (
	DailyUserDiamondConsumeUserId      db.TbCol = "user_id"
	DailyUserDiamondConsumeConsumeDate db.TbCol = "consume_date"
)

type DailyUserDiamondConsume struct {
	ID          string     `gorm:"primaryKey;size:32" json:"id"`
	UserId      uint64     `gorm:"index;default:0" json:"userId"`
	ConsumeDate string     `gorm:"size:10;index" json:"consumeDate"`
	CreatedAt   *time.Time `json:"createdAt"`
}

func BuildDailyUserDiamondConsumeId(date string, userId uint64) string {
	return fmt.Sprintf("%s_%d", date, userId)
}

func NewDailyUserDiamondConsume(date string, userId uint64) *DailyUserDiamondConsume {
	r := &DailyUserDiamondConsume{}
	r.ID = BuildDailyUserDiamondConsumeId(date, userId)
	r.SetUserId(userId)
	r.SetConsumeDate(date)
	return r
}

func (r *DailyUserDiamondConsume) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddDataToLazyChan(TbDailyUserDiamondConsume, DailyUserDiamondConsumeUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *DailyUserDiamondConsume) SetConsumeDate(v string) {
	r.ConsumeDate = v
	syndb.AddDataToLazyChan(TbDailyUserDiamondConsume, DailyUserDiamondConsumeConsumeDate, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *DailyUserDiamondConsume) SetCreatedAt(v *time.Time) {
	r.CreatedAt = v
	syndb.AddDataToLazyChan(TbDailyUserDiamondConsume, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initDailyUserDiamondConsume() {
	syndb.RegLazyWithMiddle(TbDailyUserDiamondConsume, db.CreatedAtName)
	syndb.RegLazyWithMiddle(TbDailyUserDiamondConsume, DailyUserDiamondConsumeUserId)
	syndb.RegLazyWithMiddle(TbDailyUserDiamondConsume, DailyUserDiamondConsumeConsumeDate)
	migrate.AutoMigrate(&DailyUserDiamondConsume{})
}
