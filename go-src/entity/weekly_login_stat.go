package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbWeeklyLoginStat db.TbName = "weekly_login_stats"
)

const (
	WeeklyLoginStatCount          db.TbCol = "count"
	WeeklyLoginStatRegisterCount  db.TbCol = "register_count"
	WeeklyLoginStatRechargeAmount db.TbCol = "recharge_amount"
)

// WeeklyLoginStat 每周登录统计(主键ID即周标识,如 2026-W21)
type WeeklyLoginStat struct {
	ID             string  `gorm:"primaryKey;size:16;comment:周标识(YYYY-Www)" json:"week"`
	Count          uint64  `gorm:"default:0;comment:登录数量" json:"count"`
	RegisterCount  uint64  `gorm:"default:0;comment:注册人数" json:"registerCount"`
	RechargeAmount float64 `gorm:"type:decimal(10,4);default:0;comment:充值金额(USD)" json:"rechargeAmount"`
}

// FormatWeeklyLoginStatKey 格式化周统计标识(ISO周)
func FormatWeeklyLoginStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-W%02d", year, week)
}

func NewWeeklyLoginStat(week string) *WeeklyLoginStat {
	return &WeeklyLoginStat{
		ID:             week,
		Count:          0,
		RegisterCount:  0,
		RechargeAmount: 0,
	}
}

func (r *WeeklyLoginStat) AddLoginCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddDataToLazyChan(TbWeeklyLoginStat, WeeklyLoginStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func (r *WeeklyLoginStat) AddRegisterCount(n uint64) {
	r.RegisterCount = math.Add(r.RegisterCount, n)
	syndb.AddDataToLazyChan(TbWeeklyLoginStat, WeeklyLoginStatRegisterCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.RegisterCount,
	})
}

func (r *WeeklyLoginStat) AddRechargeAmount(val float64) {
	r.RechargeAmount = math.AddFloat64(r.RechargeAmount, val)
	syndb.AddDataToLazyChan(TbWeeklyLoginStat, WeeklyLoginStatRechargeAmount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.RechargeAmount,
	})
}

func initWeeklyLoginStat() {
	syndb.RegLazyWithMiddle(TbWeeklyLoginStat, WeeklyLoginStatCount)
	syndb.RegLazyWithMiddle(TbWeeklyLoginStat, WeeklyLoginStatRegisterCount)
	syndb.RegLazyWithMiddle(TbWeeklyLoginStat, WeeklyLoginStatRechargeAmount)
	migrate.AutoMigrate(&WeeklyLoginStat{})
}
