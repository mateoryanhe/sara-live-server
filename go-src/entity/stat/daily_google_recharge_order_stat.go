package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyGoogleRechargeOrderStat db.TbName = "daily_google_recharge_order_stats"
)

const DailyGoogleRechargeOrderStatCount db.TbCol = "count"

type DailyGoogleRechargeOrderStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func NewDailyGoogleRechargeOrderStat(date string) *DailyGoogleRechargeOrderStat {
	return &DailyGoogleRechargeOrderStat{ID: date, Count: 0}
}

func (r *DailyGoogleRechargeOrderStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyGoogleRechargeOrderStat, DailyGoogleRechargeOrderStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyGoogleRechargeOrderStat() {
	syndb.RegLazy(TbDailyGoogleRechargeOrderStat, DailyGoogleRechargeOrderStatCount)
	migrate.AutoMigrate(&DailyGoogleRechargeOrderStat{})
}
