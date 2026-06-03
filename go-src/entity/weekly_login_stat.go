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
	WeeklyLoginStatCount                   db.TbCol = "count"
	WeeklyLoginStatRegisterCount           db.TbCol = "register_count"
	WeeklyLoginStatRechargeAmount          db.TbCol = "recharge_amount"
	WeeklyLoginStatGoldConsumeAmount       db.TbCol = "gold_consume_amount"
	WeeklyLoginStatDiamondConsumeAmount    db.TbCol = "diamond_consume_amount"
	WeeklyLoginStatRechargeUserCount       db.TbCol = "recharge_user_count"
	WeeklyLoginStatGoldConsumeUserCount    db.TbCol = "gold_consume_user_count"
	WeeklyLoginStatDiamondConsumeUserCount db.TbCol = "diamond_consume_user_count"
	WeeklyLoginStatAudienceUserCount       db.TbCol = "audience_user_count"
)

// WeeklyLoginStat 每周登录统计(主键ID即周标识,如 2026-W21)
type WeeklyLoginStat struct {
	ID                      string  `gorm:"primaryKey;size:16;comment:周标识(YYYY-Www)" json:"week"`
	Count                   uint64  `gorm:"default:0;comment:登录数量" json:"count"`
	RegisterCount           uint64  `gorm:"default:0;comment:注册人数" json:"registerCount"`
	RechargeAmount          float64 `gorm:"type:decimal(10,4);default:0;comment:充值金额(USD)" json:"rechargeAmount"`
	GoldConsumeAmount       float64 `gorm:"default:0;comment:金币消费金额" json:"goldConsumeAmount"`
	DiamondConsumeAmount    float64 `gorm:"default:0;comment:钻石消费金额" json:"diamondConsumeAmount"`
	RechargeUserCount       uint64  `gorm:"default:0;comment:充值人数(去重)" json:"rechargeUserCount"`
	GoldConsumeUserCount    uint64  `gorm:"default:0;comment:金币消费人数(去重)" json:"goldConsumeUserCount"`
	DiamondConsumeUserCount uint64  `gorm:"default:0;comment:钻石消费人数(去重)" json:"diamondConsumeUserCount"`
	AudienceUserCount       uint64  `gorm:"default:0;comment:有效观众人数(去重,跨直播间)" json:"audienceUserCount"`
}

// FormatWeeklyLoginStatKey 格式化周统计标识(ISO周)
func FormatWeeklyLoginStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyLoginStat(week string) *WeeklyLoginStat {
	return &WeeklyLoginStat{
		ID:                      week,
		Count:                   0,
		RegisterCount:           0,
		RechargeAmount:          0,
		GoldConsumeAmount:       0,
		DiamondConsumeAmount:    0,
		RechargeUserCount:       0,
		GoldConsumeUserCount:    0,
		DiamondConsumeUserCount: 0,
		AudienceUserCount:       0,
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

func (r *WeeklyLoginStat) AddGoldConsumeAmount(val float64) {
	r.GoldConsumeAmount = math.AddFloat64(r.GoldConsumeAmount, val)
	syndb.AddDataToLazyChan(TbWeeklyLoginStat, WeeklyLoginStatGoldConsumeAmount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.GoldConsumeAmount,
	})
}

func (r *WeeklyLoginStat) AddDiamondConsumeAmount(val float64) {
	r.DiamondConsumeAmount = math.AddFloat64(r.DiamondConsumeAmount, val)
	syndb.AddDataToLazyChan(TbWeeklyLoginStat, WeeklyLoginStatDiamondConsumeAmount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.DiamondConsumeAmount,
	})
}

func (r *WeeklyLoginStat) AddRechargeUserCount(n uint64) {
	r.RechargeUserCount = math.Add(r.RechargeUserCount, n)
	syndb.AddDataToLazyChan(TbWeeklyLoginStat, WeeklyLoginStatRechargeUserCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.RechargeUserCount,
	})
}

func (r *WeeklyLoginStat) AddGoldConsumeUserCount(n uint64) {
	r.GoldConsumeUserCount = math.Add(r.GoldConsumeUserCount, n)
	syndb.AddDataToLazyChan(TbWeeklyLoginStat, WeeklyLoginStatGoldConsumeUserCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.GoldConsumeUserCount,
	})
}

func (r *WeeklyLoginStat) AddDiamondConsumeUserCount(n uint64) {
	r.DiamondConsumeUserCount = math.Add(r.DiamondConsumeUserCount, n)
	syndb.AddDataToLazyChan(TbWeeklyLoginStat, WeeklyLoginStatDiamondConsumeUserCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.DiamondConsumeUserCount,
	})
}

func (r *WeeklyLoginStat) AddAudienceUserCount(n uint64) {
	r.AudienceUserCount = math.Add(r.AudienceUserCount, n)
	syndb.AddDataToLazyChan(TbWeeklyLoginStat, WeeklyLoginStatAudienceUserCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.AudienceUserCount,
	})
}

func initWeeklyLoginStat() {
	syndb.RegLazyWithMiddle(TbWeeklyLoginStat, WeeklyLoginStatCount)
	syndb.RegLazyWithMiddle(TbWeeklyLoginStat, WeeklyLoginStatRegisterCount)
	syndb.RegLazyWithMiddle(TbWeeklyLoginStat, WeeklyLoginStatRechargeAmount)
	syndb.RegLazyWithMiddle(TbWeeklyLoginStat, WeeklyLoginStatGoldConsumeAmount)
	syndb.RegLazyWithMiddle(TbWeeklyLoginStat, WeeklyLoginStatDiamondConsumeAmount)
	syndb.RegLazyWithMiddle(TbWeeklyLoginStat, WeeklyLoginStatRechargeUserCount)
	syndb.RegLazyWithMiddle(TbWeeklyLoginStat, WeeklyLoginStatGoldConsumeUserCount)
	syndb.RegLazyWithMiddle(TbWeeklyLoginStat, WeeklyLoginStatDiamondConsumeUserCount)
	syndb.RegLazyWithMiddle(TbWeeklyLoginStat, WeeklyLoginStatAudienceUserCount)
	migrate.AutoMigrate(&WeeklyLoginStat{})
}
