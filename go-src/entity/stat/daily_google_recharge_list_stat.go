package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyGoogleRechargeListStat db.TbName = "daily_google_recharge_list_stats"
)

const DailyGoogleRechargeListStatCount db.TbCol = "count"

type DailyGoogleRechargeListStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:次数" json:"count"`
}

func NewDailyGoogleRechargeListStat(date string) *DailyGoogleRechargeListStat {
	return &DailyGoogleRechargeListStat{ID: date, Count: 0}
}

func (r *DailyGoogleRechargeListStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyGoogleRechargeListStat, DailyGoogleRechargeListStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyGoogleRechargeListStat() {
	syndb.RegLazy(TbDailyGoogleRechargeListStat, DailyGoogleRechargeListStatCount)
	migrate.AutoMigrate(&DailyGoogleRechargeListStat{})
}
