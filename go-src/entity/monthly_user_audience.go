package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const TbMonthlyUserAudience db.TbName = "monthly_user_audiences"

const (
	MonthlyUserAudienceUserId        db.TbCol = "user_id"
	MonthlyUserAudienceAudienceMonth db.TbCol = "audience_month"
)

// MonthlyUserAudience 每月有效观众去重(跨直播间)
type MonthlyUserAudience struct {
	ID            string     `gorm:"primaryKey;size:32;comment:复合ID(month_userId)" json:"id"`
	UserId        uint64     `gorm:"index;default:0;comment:用户ID" json:"userId"`
	AudienceMonth string     `gorm:"size:10;index;comment:月标识YYYY-MM" json:"audienceMonth"`
	CreatedAt     *time.Time `json:"createdAt"`
}

func BuildMonthlyUserAudienceId(month string, userId uint64) string {
	return fmt.Sprintf("%s_%d", month, userId)
}

func NewMonthlyUserAudience(month string, userId uint64) *MonthlyUserAudience {
	r := &MonthlyUserAudience{}
	r.ID = BuildMonthlyUserAudienceId(month, userId)
	r.SetUserId(userId)
	r.SetAudienceMonth(month)
	return r
}

func (r *MonthlyUserAudience) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddDataToLazyChan(TbMonthlyUserAudience, MonthlyUserAudienceUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *MonthlyUserAudience) SetAudienceMonth(v string) {
	r.AudienceMonth = v
	syndb.AddDataToLazyChan(TbMonthlyUserAudience, MonthlyUserAudienceAudienceMonth, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *MonthlyUserAudience) SetCreatedAt(v *time.Time) {
	r.CreatedAt = v
	syndb.AddDataToLazyChan(TbMonthlyUserAudience, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initMonthlyUserAudience() {
	syndb.RegLazyWithMiddle(TbMonthlyUserAudience, db.CreatedAtName)
	syndb.RegLazyWithMiddle(TbMonthlyUserAudience, MonthlyUserAudienceUserId)
	syndb.RegLazyWithMiddle(TbMonthlyUserAudience, MonthlyUserAudienceAudienceMonth)
	migrate.AutoMigrate(&MonthlyUserAudience{})
}
