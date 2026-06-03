package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const TbWeeklyUserAudience db.TbName = "weekly_user_audiences"

const (
	WeeklyUserAudienceUserId       db.TbCol = "user_id"
	WeeklyUserAudienceAudienceWeek db.TbCol = "audience_week"
)

// WeeklyUserAudience 每周有效观众去重(跨直播间)
type WeeklyUserAudience struct {
	ID           string     `gorm:"primaryKey;size:32;comment:复合ID(week_userId)" json:"id"`
	UserId       uint64     `gorm:"index;default:0;comment:用户ID" json:"userId"`
	AudienceWeek string     `gorm:"size:16;index;comment:周标识" json:"audienceWeek"`
	CreatedAt    *time.Time `json:"createdAt"`
}

func BuildWeeklyUserAudienceId(week string, userId uint64) string {
	return fmt.Sprintf("%s_%d", week, userId)
}

func NewWeeklyUserAudience(week string, userId uint64) *WeeklyUserAudience {
	r := &WeeklyUserAudience{}
	r.ID = BuildWeeklyUserAudienceId(week, userId)
	r.SetUserId(userId)
	r.SetAudienceWeek(week)
	return r
}

func (r *WeeklyUserAudience) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddDataToLazyChan(TbWeeklyUserAudience, WeeklyUserAudienceUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *WeeklyUserAudience) SetAudienceWeek(v string) {
	r.AudienceWeek = v
	syndb.AddDataToLazyChan(TbWeeklyUserAudience, WeeklyUserAudienceAudienceWeek, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *WeeklyUserAudience) SetCreatedAt(v *time.Time) {
	r.CreatedAt = v
	syndb.AddDataToLazyChan(TbWeeklyUserAudience, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initWeeklyUserAudience() {
	syndb.RegLazyWithMiddle(TbWeeklyUserAudience, db.CreatedAtName)
	syndb.RegLazyWithMiddle(TbWeeklyUserAudience, WeeklyUserAudienceUserId)
	syndb.RegLazyWithMiddle(TbWeeklyUserAudience, WeeklyUserAudienceAudienceWeek)
	migrate.AutoMigrate(&WeeklyUserAudience{})
}
