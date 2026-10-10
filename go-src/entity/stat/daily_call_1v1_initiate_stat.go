package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyCall1v1InitiateStat db.TbName = "daily_call_1v1_initiate_stats"
)

const DailyCall1v1InitiateStatCount db.TbCol = "count"

// DailyCall1v1InitiateStat 每日 source=1 的 1v1 视频通话发起次数(主键=YYYY-MM-DD)
type DailyCall1v1InitiateStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:发起次数" json:"count"`
}

func NewDailyCall1v1InitiateStat(date string) *DailyCall1v1InitiateStat {
	return &DailyCall1v1InitiateStat{ID: date, Count: 0}
}

func (r *DailyCall1v1InitiateStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyCall1v1InitiateStat, DailyCall1v1InitiateStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyCall1v1InitiateStat() {
	syndb.RegLazy(TbDailyCall1v1InitiateStat, DailyCall1v1InitiateStatCount)
	migrate.AutoMigrate(&DailyCall1v1InitiateStat{})
}
