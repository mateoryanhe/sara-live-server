package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbMonthlyUserRecharge db.TbName = "monthly_user_recharges"
)

const (
	MonthlyUserRechargeUserId        db.TbCol = "user_id"
	MonthlyUserRechargeRechargeMonth db.TbCol = "recharge_month"
)

// MonthlyUserRecharge 每月用户充值去重记录(主键 ID = month_userId)
type MonthlyUserRecharge struct {
	ID            string     `gorm:"primaryKey;size:32;comment:复合ID(month_userId)" json:"id"`
	UserId        uint64     `gorm:"index;default:0;comment:用户ID" json:"userId"`
	RechargeMonth string     `gorm:"size:10;index;comment:充值月标识YYYY-MM" json:"rechargeMonth"`
	CreatedAt     *time.Time `json:"createdAt"`
}

func BuildMonthlyUserRechargeId(month string, userId uint64) string {
	return fmt.Sprintf("%s_%d", month, userId)
}

func NewMonthlyUserRecharge(month string, userId uint64) *MonthlyUserRecharge {
	r := &MonthlyUserRecharge{}
	r.ID = BuildMonthlyUserRechargeId(month, userId)
	r.SetUserId(userId)
	r.SetRechargeMonth(month)
	return r
}

func (r *MonthlyUserRecharge) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddDataToLazyChan(TbMonthlyUserRecharge, MonthlyUserRechargeUserId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *MonthlyUserRecharge) SetRechargeMonth(v string) {
	r.RechargeMonth = v
	syndb.AddDataToLazyChan(TbMonthlyUserRecharge, MonthlyUserRechargeRechargeMonth, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *MonthlyUserRecharge) SetCreatedAt(v *time.Time) {
	r.CreatedAt = v
	syndb.AddDataToLazyChan(TbMonthlyUserRecharge, db.CreatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func initMonthlyUserRecharge() {
	syndb.RegLazyWithMiddle(TbMonthlyUserRecharge, db.CreatedAtName)
	syndb.RegLazyWithMiddle(TbMonthlyUserRecharge, MonthlyUserRechargeUserId)
	syndb.RegLazyWithMiddle(TbMonthlyUserRecharge, MonthlyUserRechargeRechargeMonth)
	migrate.AutoMigrate(&MonthlyUserRecharge{})
}
