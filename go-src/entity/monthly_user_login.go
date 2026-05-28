package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbMonthlyUserLogin db.TbName = "monthly_user_logins"
)

const (
	MonthlyUserLoginUserId     db.TbCol = "user_id"
	MonthlyUserLoginLoginMonth db.TbCol = "login_month"
)

// MonthlyUserLogin 每月用户登录去重记录(主键 ID = month_userId)
type MonthlyUserLogin struct {
	ID         string     `gorm:"primaryKey;size:32;comment:复合ID(month_userId)" json:"id"`
	UserId     uint64     `gorm:"index;default:0;comment:用户ID" json:"userId"`
	LoginMonth string     `gorm:"size:10;index;comment:登录月标识YYYY-MM" json:"loginMonth"`
	CreatedAt  *time.Time `json:"createdAt"`
}

func BuildMonthlyUserLoginId(month string, userId uint64) string {
	return fmt.Sprintf("%s_%d", month, userId)
}

func NewMonthlyUserLogin(month string, userId uint64) *MonthlyUserLogin {
	r := &MonthlyUserLogin{}
	r.ID = BuildMonthlyUserLoginId(month, userId)
	r.SetUserId(userId)
	r.SetLoginMonth(month)
	return r
}

func (r *MonthlyUserLogin) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddDataToLazyChan(TbMonthlyUserLogin, MonthlyUserLoginUserId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *MonthlyUserLogin) SetLoginMonth(v string) {
	r.LoginMonth = v
	syndb.AddDataToLazyChan(TbMonthlyUserLogin, MonthlyUserLoginLoginMonth, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *MonthlyUserLogin) SetCreatedAt(v *time.Time) {
	r.CreatedAt = v
	syndb.AddDataToLazyChan(TbMonthlyUserLogin, db.CreatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func initMonthlyUserLogin() {
	syndb.RegLazyWithMiddle(TbMonthlyUserLogin, db.CreatedAtName)
	syndb.RegLazyWithMiddle(TbMonthlyUserLogin, MonthlyUserLoginUserId)
	syndb.RegLazyWithMiddle(TbMonthlyUserLogin, MonthlyUserLoginLoginMonth)
	migrate.AutoMigrate(&MonthlyUserLogin{})
}
