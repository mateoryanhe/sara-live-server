package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const TbDailyUserGoldConsume db.TbName = "daily_user_gold_consumes"

const (
	DailyUserGoldConsumeUserId      db.TbCol = "user_id"
	DailyUserGoldConsumeConsumeDate db.TbCol = "consume_date"
)

// DailyUserGoldConsume 每日用户金币消费去重(主键 ID = date_userId)
type DailyUserGoldConsume struct {
	ID          string     `gorm:"primaryKey;size:32;comment:复合ID(date_userId)" json:"id"`
	UserId      uint64     `gorm:"index;default:0;comment:用户ID" json:"userId"`
	ConsumeDate string     `gorm:"size:10;index;comment:消费日期YYYY-MM-DD" json:"consumeDate"`
	CreatedAt   *time.Time `json:"createdAt"`
}

func BuildDailyUserGoldConsumeId(date string, userId uint64) string {
	return fmt.Sprintf("%s_%d", date, userId)
}

func NewDailyUserGoldConsume(date string, userId uint64) *DailyUserGoldConsume {
	r := &DailyUserGoldConsume{}
	r.ID = BuildDailyUserGoldConsumeId(date, userId)
	r.SetUserId(userId)
	r.SetConsumeDate(date)
	return r
}

func (r *DailyUserGoldConsume) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddDataToLazyChan(TbDailyUserGoldConsume, DailyUserGoldConsumeUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *DailyUserGoldConsume) SetConsumeDate(v string) {
	r.ConsumeDate = v
	syndb.AddDataToLazyChan(TbDailyUserGoldConsume, DailyUserGoldConsumeConsumeDate, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *DailyUserGoldConsume) SetCreatedAt(v *time.Time) {
	r.CreatedAt = v
	syndb.AddDataToLazyChan(TbDailyUserGoldConsume, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initDailyUserGoldConsume() {
	syndb.RegLazyWithMiddle(TbDailyUserGoldConsume, db.CreatedAtName)
	syndb.RegLazyWithMiddle(TbDailyUserGoldConsume, DailyUserGoldConsumeUserId)
	syndb.RegLazyWithMiddle(TbDailyUserGoldConsume, DailyUserGoldConsumeConsumeDate)
	migrate.AutoMigrate(&DailyUserGoldConsume{})
}
