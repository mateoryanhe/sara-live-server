package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyUserRecharge db.TbName = "daily_user_recharges"
)

const (
	DailyUserRechargeUserId       db.TbCol = "user_id"
	DailyUserRechargeRechargeDate db.TbCol = "recharge_date"
)

// DailyUserRecharge 每日用户充值去重记录(主键 ID = date_userId)
type DailyUserRecharge struct {
	ID           string     `gorm:"primaryKey;size:32;comment:复合ID(date_userId)" json:"id"`
	UserId       uint64     `gorm:"index;default:0;comment:用户ID" json:"userId"`
	RechargeDate string     `gorm:"size:10;index;comment:充值日期YYYY-MM-DD" json:"rechargeDate"`
	CreatedAt    *time.Time `json:"createdAt"`
}

func BuildDailyUserRechargeId(date string, userId uint64) string {
	return fmt.Sprintf("%s_%d", date, userId)
}

func NewDailyUserRecharge(date string, userId uint64) *DailyUserRecharge {
	r := &DailyUserRecharge{}
	r.ID = BuildDailyUserRechargeId(date, userId)
	r.SetUserId(userId)
	r.SetRechargeDate(date)
	return r
}

func (r *DailyUserRecharge) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddDataToLazyChan(TbDailyUserRecharge, DailyUserRechargeUserId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *DailyUserRecharge) SetRechargeDate(v string) {
	r.RechargeDate = v
	syndb.AddDataToLazyChan(TbDailyUserRecharge, DailyUserRechargeRechargeDate, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *DailyUserRecharge) SetCreatedAt(v *time.Time) {
	r.CreatedAt = v
	syndb.AddDataToLazyChan(TbDailyUserRecharge, db.CreatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func initDailyUserRecharge() {
	syndb.RegLazyWithMiddle(TbDailyUserRecharge, db.CreatedAtName)
	syndb.RegLazyWithMiddle(TbDailyUserRecharge, DailyUserRechargeUserId)
	syndb.RegLazyWithMiddle(TbDailyUserRecharge, DailyUserRechargeRechargeDate)
	migrate.AutoMigrate(&DailyUserRecharge{})
}
