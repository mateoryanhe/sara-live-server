package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbWeeklyUserLogin db.TbName = "weekly_user_logins"
)

const (
	WeeklyUserLoginUserId    db.TbCol = "user_id"
	WeeklyUserLoginLoginWeek db.TbCol = "login_week"
)

// WeeklyUserLogin 每周用户登录去重记录(主键 ID = week_userId)
type WeeklyUserLogin struct {
	ID        string     `gorm:"primaryKey;size:32;comment:复合ID(week_userId)" json:"id"`
	UserId    uint64     `gorm:"index;default:0;comment:用户ID" json:"userId"`
	LoginWeek string     `gorm:"size:16;index;comment:登录周标识YYYY-Www" json:"loginWeek"`
	CreatedAt *time.Time `json:"createdAt"`
}

func BuildWeeklyUserLoginId(week string, userId uint64) string {
	return fmt.Sprintf("%s_%d", week, userId)
}

func NewWeeklyUserLogin(week string, userId uint64) *WeeklyUserLogin {
	r := &WeeklyUserLogin{}
	r.ID = BuildWeeklyUserLoginId(week, userId)
	r.SetUserId(userId)
	r.SetLoginWeek(week)
	return r
}

func (r *WeeklyUserLogin) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddDataToLazyChan(TbWeeklyUserLogin, WeeklyUserLoginUserId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *WeeklyUserLogin) SetLoginWeek(v string) {
	r.LoginWeek = v
	syndb.AddDataToLazyChan(TbWeeklyUserLogin, WeeklyUserLoginLoginWeek, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *WeeklyUserLogin) SetCreatedAt(v *time.Time) {
	r.CreatedAt = v
	syndb.AddDataToLazyChan(TbWeeklyUserLogin, db.CreatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func initWeeklyUserLogin() {
	syndb.RegLazyWithMiddle(TbWeeklyUserLogin, db.CreatedAtName)
	syndb.RegLazyWithMiddle(TbWeeklyUserLogin, WeeklyUserLoginUserId)
	syndb.RegLazyWithMiddle(TbWeeklyUserLogin, WeeklyUserLoginLoginWeek)
	migrate.AutoMigrate(&WeeklyUserLogin{})
}
