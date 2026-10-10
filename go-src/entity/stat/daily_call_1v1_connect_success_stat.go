package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyCall1v1ConnectSuccessStat db.TbName = "daily_call_1v1_connect_success_stats"
)

const DailyCall1v1ConnectSuccessStatCount db.TbCol = "count"

// DailyCall1v1ConnectSuccessStat 每日双方成功接通的 1v1 视频通话次数(主键=YYYY-MM-DD)
type DailyCall1v1ConnectSuccessStat struct {
	ID    string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count uint64 `gorm:"default:0;comment:接通次数" json:"count"`
}

func NewDailyCall1v1ConnectSuccessStat(date string) *DailyCall1v1ConnectSuccessStat {
	return &DailyCall1v1ConnectSuccessStat{ID: date, Count: 0}
}

func (r *DailyCall1v1ConnectSuccessStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbDailyCall1v1ConnectSuccessStat, DailyCall1v1ConnectSuccessStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initDailyCall1v1ConnectSuccessStat() {
	syndb.RegLazy(TbDailyCall1v1ConnectSuccessStat, DailyCall1v1ConnectSuccessStatCount)
	migrate.AutoMigrate(&DailyCall1v1ConnectSuccessStat{})
}
