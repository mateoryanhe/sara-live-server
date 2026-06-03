package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const TbDailyUserAudience db.TbName = "daily_user_audiences"

const (
	DailyUserAudienceUserId       db.TbCol = "user_id"
	DailyUserAudienceAudienceDate db.TbCol = "audience_date"
)

// DailyUserAudience 每日有效观众去重(跨直播间,主键 ID = date_userId)
type DailyUserAudience struct {
	ID           string     `gorm:"primaryKey;size:32;comment:复合ID(date_userId)" json:"id"`
	UserId       uint64     `gorm:"index;default:0;comment:用户ID" json:"userId"`
	AudienceDate string     `gorm:"size:10;index;comment:统计日期YYYY-MM-DD" json:"audienceDate"`
	CreatedAt    *time.Time `json:"createdAt"`
}

func BuildDailyUserAudienceId(date string, userId uint64) string {
	return fmt.Sprintf("%s_%d", date, userId)
}

func NewDailyUserAudience(date string, userId uint64) *DailyUserAudience {
	r := &DailyUserAudience{}
	r.ID = BuildDailyUserAudienceId(date, userId)
	r.SetUserId(userId)
	r.SetAudienceDate(date)
	return r
}

func (r *DailyUserAudience) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddDataToLazyChan(TbDailyUserAudience, DailyUserAudienceUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *DailyUserAudience) SetAudienceDate(v string) {
	r.AudienceDate = v
	syndb.AddDataToLazyChan(TbDailyUserAudience, DailyUserAudienceAudienceDate, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *DailyUserAudience) SetCreatedAt(v *time.Time) {
	r.CreatedAt = v
	syndb.AddDataToLazyChan(TbDailyUserAudience, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initDailyUserAudience() {
	syndb.RegLazyWithMiddle(TbDailyUserAudience, db.CreatedAtName)
	syndb.RegLazyWithMiddle(TbDailyUserAudience, DailyUserAudienceUserId)
	syndb.RegLazyWithMiddle(TbDailyUserAudience, DailyUserAudienceAudienceDate)
	migrate.AutoMigrate(&DailyUserAudience{})
}
