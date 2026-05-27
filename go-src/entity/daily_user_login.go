package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyUserLogin db.TbName = "daily_user_logins"
)

const (
	DailyUserLoginUserId    db.TbCol = "user_id"
	DailyUserLoginLoginDate db.TbCol = "login_date"
)

// DailyUserLogin 每日用户登录去重记录(主键 ID = date_userId)
type DailyUserLogin struct {
	ID        string     `gorm:"primaryKey;size:32;comment:复合ID(date_userId)" json:"id"`
	UserId    uint64     `gorm:"index;default:0;comment:用户ID" json:"userId"`
	LoginDate string     `gorm:"size:10;index;comment:登录日期YYYY-MM-DD" json:"loginDate"`
	CreatedAt *time.Time `json:"createdAt"`
}

func BuildDailyUserLoginId(date string, userId uint64) string {
	return fmt.Sprintf("%s_%d", date, userId)
}

func NewDailyUserLogin(date string, userId uint64) *DailyUserLogin {
	r := &DailyUserLogin{}
	r.ID = BuildDailyUserLoginId(date, userId)

	r.SetUserId(userId)
	r.SetLoginDate(date)
	return r
}

func (r *DailyUserLogin) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddDataToLazyChan(TbDailyUserLogin, DailyUserLoginUserId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *DailyUserLogin) SetLoginDate(v string) {
	r.LoginDate = v
	syndb.AddDataToLazyChan(TbDailyUserLogin, DailyUserLoginLoginDate, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *DailyUserLogin) SetCreatedAt(v *time.Time) {
	r.CreatedAt = v
	syndb.AddDataToLazyChan(TbDailyUserLogin, db.CreatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func initDailyUserLogin() {
	syndb.RegLazyWithMiddle(TbDailyUserLogin, db.CreatedAtName)
	syndb.RegLazyWithMiddle(TbDailyUserLogin, DailyUserLoginUserId)
	syndb.RegLazyWithMiddle(TbDailyUserLogin, DailyUserLoginLoginDate)
	migrate.AutoMigrate(&DailyUserLogin{})
}
