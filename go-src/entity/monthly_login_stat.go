package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbMonthlyLoginStat db.TbName = "monthly_login_stats"
)

const (
	MonthlyLoginStatCount          db.TbCol = "count"
	MonthlyLoginStatRegisterCount  db.TbCol = "register_count"
	MonthlyLoginStatRechargeAmount db.TbCol = "recharge_amount"
)

// MonthlyLoginStat 每月登录统计(主键ID即月标识,如 2026-05)
type MonthlyLoginStat struct {
	ID             string  `gorm:"primaryKey;size:10;comment:月标识(YYYY-MM)" json:"month"`
	Count          uint64  `gorm:"default:0;comment:登录数量" json:"count"`
	RegisterCount  uint64  `gorm:"default:0;comment:注册人数" json:"registerCount"`
	RechargeAmount float64 `gorm:"type:decimal(10,4);default:0;comment:充值金额(USD)" json:"rechargeAmount"`
}

// FormatMonthlyLoginStatKey 格式化月统计标识
func FormatMonthlyLoginStatKey(t time.Time) string {
	return t.Format("2006-01")
}

func NewMonthlyLoginStat(month string) *MonthlyLoginStat {
	return &MonthlyLoginStat{
		ID:             month,
		Count:          0,
		RegisterCount:  0,
		RechargeAmount: 0,
	}
}

func (r *MonthlyLoginStat) AddLoginCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddDataToLazyChan(TbMonthlyLoginStat, MonthlyLoginStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func (r *MonthlyLoginStat) AddRegisterCount(n uint64) {
	r.RegisterCount = math.Add(r.RegisterCount, n)
	syndb.AddDataToLazyChan(TbMonthlyLoginStat, MonthlyLoginStatRegisterCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.RegisterCount,
	})
}

func (r *MonthlyLoginStat) AddRechargeAmount(val float64) {
	r.RechargeAmount = math.AddFloat64(r.RechargeAmount, val)
	syndb.AddDataToLazyChan(TbMonthlyLoginStat, MonthlyLoginStatRechargeAmount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.RechargeAmount,
	})
}

func initMonthlyLoginStat() {
	syndb.RegLazyWithMiddle(TbMonthlyLoginStat, MonthlyLoginStatCount)
	syndb.RegLazyWithMiddle(TbMonthlyLoginStat, MonthlyLoginStatRegisterCount)
	syndb.RegLazyWithMiddle(TbMonthlyLoginStat, MonthlyLoginStatRechargeAmount)
	migrate.AutoMigrate(&MonthlyLoginStat{})
}
