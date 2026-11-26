package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyGoogleRechargeResultStat db.TbName = "daily_google_recharge_result_stats"
)

const DailyGoogleRechargeResultStatCount db.TbCol = "count"

type DailyGoogleRechargeResultStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func NewDailyGoogleRechargeResultStat(date string) *DailyGoogleRechargeResultStat {
	return &DailyGoogleRechargeResultStat{ID: date, Count: 0}
}

func (r *DailyGoogleRechargeResultStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyGoogleRechargeResultStat, DailyGoogleRechargeResultStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyGoogleRechargeResultStat() {
	syndb.RegLazy(TbDailyGoogleRechargeResultStat, DailyGoogleRechargeResultStatCount)
	migrate.AutoMigrate(&DailyGoogleRechargeResultStat{})
}
