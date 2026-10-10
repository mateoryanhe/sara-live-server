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
	TbWeeklyCall1v1ConnectSuccessStat db.TbName = "weekly_call_1v1_connect_success_stats"
)

const WeeklyCall1v1ConnectSuccessStatCount db.TbCol = "count"

type WeeklyCall1v1ConnectSuccessStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:接通次数" json:"count"`
}

func FormatWeeklyCall1v1ConnectSuccessStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyCall1v1ConnectSuccessStat(week string) *WeeklyCall1v1ConnectSuccessStat {
	return &WeeklyCall1v1ConnectSuccessStat{ID: week, Count: 0}
}

func (r *WeeklyCall1v1ConnectSuccessStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyCall1v1ConnectSuccessStat, WeeklyCall1v1ConnectSuccessStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyCall1v1ConnectSuccessStat() {
	syndb.RegLazy(TbWeeklyCall1v1ConnectSuccessStat, WeeklyCall1v1ConnectSuccessStatCount)
	migrate.AutoMigrate(&WeeklyCall1v1ConnectSuccessStat{})
}
