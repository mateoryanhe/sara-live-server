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
	TbWeeklyCall1v1InitiateStat db.TbName = "weekly_call_1v1_initiate_stats"
)

const WeeklyCall1v1InitiateStatCount db.TbCol = "count"

type WeeklyCall1v1InitiateStat struct {
	ID    string `gorm:"primaryKey;size:16;comment:周标识" json:"week"`
	Count uint64 `gorm:"default:0;comment:发起次数" json:"count"`
}

func FormatWeeklyCall1v1InitiateStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyCall1v1InitiateStat(week string) *WeeklyCall1v1InitiateStat {
	return &WeeklyCall1v1InitiateStat{ID: week, Count: 0}
}

func (r *WeeklyCall1v1InitiateStat) AddCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyCall1v1InitiateStat, WeeklyCall1v1InitiateStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func initWeeklyCall1v1InitiateStat() {
	syndb.RegLazy(TbWeeklyCall1v1InitiateStat, WeeklyCall1v1InitiateStatCount)
	migrate.AutoMigrate(&WeeklyCall1v1InitiateStat{})
}
