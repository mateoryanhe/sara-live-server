package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbWeeklyUserRecharge db.TbName = "weekly_user_recharges"
)

const (
	WeeklyUserRechargeUserId       db.TbCol = "user_id"
	WeeklyUserRechargeRechargeWeek db.TbCol = "recharge_week"
)

// WeeklyUserRecharge 每周用户充值去重记录(主键 ID = week_userId)
type WeeklyUserRecharge struct {
	ID           string     `gorm:"primaryKey;size:32;comment:复合ID(week_userId)" json:"id"`
	UserId       uint64     `gorm:"index;default:0;comment:用户ID" json:"userId"`
	RechargeWeek string     `gorm:"size:16;index;comment:充值周标识" json:"rechargeWeek"`
	CreatedAt    *time.Time `json:"createdAt"`
}

func BuildWeeklyUserRechargeId(week string, userId uint64) string {
	return fmt.Sprintf("%s_%d", week, userId)
}

func NewWeeklyUserRecharge(week string, userId uint64) *WeeklyUserRecharge {
	r := &WeeklyUserRecharge{}
	r.ID = BuildWeeklyUserRechargeId(week, userId)
	r.SetUserId(userId)
	r.SetRechargeWeek(week)
	return r
}

func (r *WeeklyUserRecharge) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddDataToLazyChan(TbWeeklyUserRecharge, WeeklyUserRechargeUserId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *WeeklyUserRecharge) SetRechargeWeek(v string) {
	r.RechargeWeek = v
	syndb.AddDataToLazyChan(TbWeeklyUserRecharge, WeeklyUserRechargeRechargeWeek, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *WeeklyUserRecharge) SetCreatedAt(v *time.Time) {
	r.CreatedAt = v
	syndb.AddDataToLazyChan(TbWeeklyUserRecharge, db.CreatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func initWeeklyUserRecharge() {
	syndb.RegLazyWithMiddle(TbWeeklyUserRecharge, db.CreatedAtName)
	syndb.RegLazyWithMiddle(TbWeeklyUserRecharge, WeeklyUserRechargeUserId)
	syndb.RegLazyWithMiddle(TbWeeklyUserRecharge, WeeklyUserRechargeRechargeWeek)
	migrate.AutoMigrate(&WeeklyUserRecharge{})
}
