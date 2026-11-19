package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyEmailAccountStat db.TbName = "daily_email_account_stats"
)

const DailyEmailAccountStatVal db.TbCol = "val"

// DailyEmailAccountStat 每日邮箱绑定或邮箱注册人数(合并计数,每次成功 +1)
type DailyEmailAccountStat struct {
	ID  string `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Val uint64 `gorm:"default:0;comment:人数" json:"val"`
}

func NewDailyEmailAccountStat(date string) *DailyEmailAccountStat {
	return &DailyEmailAccountStat{ID: date, Val: 0}
}

func (r *DailyEmailAccountStat) AddVal(n uint64) {
	r.Val = math.Add(r.Val, n)
	syndb.AddData(TbDailyEmailAccountStat, DailyEmailAccountStatVal, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Val,
	})
}

func initDailyEmailAccountStat() {
	syndb.RegLazy(TbDailyEmailAccountStat, DailyEmailAccountStatVal)
	migrate.AutoMigrate(&DailyEmailAccountStat{})
}