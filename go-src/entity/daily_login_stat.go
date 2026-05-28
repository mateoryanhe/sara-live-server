package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyLoginStat db.TbName = "daily_login_stats"
)

const (
	DailyLoginStatCount         db.TbCol = "count"
	DailyLoginStatRegisterCount db.TbCol = "register_count"
)

// DailyLoginStat 每日登录统计(主键ID即日期 YYYY-MM-DD)
type DailyLoginStat struct {
	ID            string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count         uint64 `gorm:"default:0;comment:登录数量" json:"count"`
	RegisterCount uint64 `gorm:"default:0;comment:注册人数" json:"registerCount"`
}

// FormatDailyLoginStatDate 格式化统计日期
func FormatDailyLoginStatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func NewDailyLoginStat(date string) *DailyLoginStat {
	return &DailyLoginStat{
		ID:            date,
		Count:         0,
		RegisterCount: 0,
	}
}

func (receiver *DailyLoginStat) AddLoginCount(n uint64) {
	receiver.Count = math.Add(receiver.Count, n)
	syndb.AddDataToLazyChan(TbDailyLoginStat, DailyLoginStatCount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.Count,
	})
}

func (receiver *DailyLoginStat) AddRegisterCount(n uint64) {
	receiver.RegisterCount = math.Add(receiver.RegisterCount, n)
	syndb.AddDataToLazyChan(TbDailyLoginStat, DailyLoginStatRegisterCount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.RegisterCount,
	})
}

func initDailyLoginStat() {
	syndb.RegLazyWithMiddle(TbDailyLoginStat, DailyLoginStatCount)
	syndb.RegLazyWithMiddle(TbDailyLoginStat, DailyLoginStatRegisterCount)
	migrate.AutoMigrate(&DailyLoginStat{})
}
