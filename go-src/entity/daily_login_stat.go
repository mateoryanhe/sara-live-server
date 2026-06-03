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
	DailyLoginStatCount                   db.TbCol = "count"
	DailyLoginStatRegisterCount           db.TbCol = "register_count"
	DailyLoginStatRechargeAmount          db.TbCol = "recharge_amount"
	DailyLoginStatGoldConsumeAmount       db.TbCol = "gold_consume_amount"
	DailyLoginStatDiamondConsumeAmount    db.TbCol = "diamond_consume_amount"
	DailyLoginStatRechargeUserCount       db.TbCol = "recharge_user_count"
	DailyLoginStatGoldConsumeUserCount    db.TbCol = "gold_consume_user_count"
	DailyLoginStatDiamondConsumeUserCount db.TbCol = "diamond_consume_user_count"
	DailyLoginStatAudienceUserCount       db.TbCol = "audience_user_count"
)

// DailyLoginStat 每日登录统计(主键ID即日期 YYYY-MM-DD)
type DailyLoginStat struct {
	ID                      string  `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
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

// FormatDailyLoginStatDate 格式化统计日期
func FormatDailyLoginStatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func NewDailyLoginStat(date string) *DailyLoginStat {
	return &DailyLoginStat{
		ID:                      date,
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

func (receiver *DailyLoginStat) AddRechargeAmount(val float64) {
	receiver.RechargeAmount = math.AddFloat64(receiver.RechargeAmount, val)
	syndb.AddDataToLazyChan(TbDailyLoginStat, DailyLoginStatRechargeAmount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.RechargeAmount,
	})
}

func (receiver *DailyLoginStat) AddGoldConsumeAmount(val float64) {
	receiver.GoldConsumeAmount = math.AddFloat64(receiver.GoldConsumeAmount, val)
	syndb.AddDataToLazyChan(TbDailyLoginStat, DailyLoginStatGoldConsumeAmount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.GoldConsumeAmount,
	})
}

func (receiver *DailyLoginStat) AddDiamondConsumeAmount(val float64) {
	receiver.DiamondConsumeAmount = math.AddFloat64(receiver.DiamondConsumeAmount, val)
	syndb.AddDataToLazyChan(TbDailyLoginStat, DailyLoginStatDiamondConsumeAmount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.DiamondConsumeAmount,
	})
}

func (receiver *DailyLoginStat) AddRechargeUserCount(n uint64) {
	receiver.RechargeUserCount = math.Add(receiver.RechargeUserCount, n)
	syndb.AddDataToLazyChan(TbDailyLoginStat, DailyLoginStatRechargeUserCount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.RechargeUserCount,
	})
}

func (receiver *DailyLoginStat) AddGoldConsumeUserCount(n uint64) {
	receiver.GoldConsumeUserCount = math.Add(receiver.GoldConsumeUserCount, n)
	syndb.AddDataToLazyChan(TbDailyLoginStat, DailyLoginStatGoldConsumeUserCount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.GoldConsumeUserCount,
	})
}

func (receiver *DailyLoginStat) AddDiamondConsumeUserCount(n uint64) {
	receiver.DiamondConsumeUserCount = math.Add(receiver.DiamondConsumeUserCount, n)
	syndb.AddDataToLazyChan(TbDailyLoginStat, DailyLoginStatDiamondConsumeUserCount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.DiamondConsumeUserCount,
	})
}

func (receiver *DailyLoginStat) AddAudienceUserCount(n uint64) {
	receiver.AudienceUserCount = math.Add(receiver.AudienceUserCount, n)
	syndb.AddDataToLazyChan(TbDailyLoginStat, DailyLoginStatAudienceUserCount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.AudienceUserCount,
	})
}

func initDailyLoginStat() {
	syndb.RegLazyWithMiddle(TbDailyLoginStat, DailyLoginStatCount)
	syndb.RegLazyWithMiddle(TbDailyLoginStat, DailyLoginStatRegisterCount)
	syndb.RegLazyWithMiddle(TbDailyLoginStat, DailyLoginStatRechargeAmount)
	syndb.RegLazyWithMiddle(TbDailyLoginStat, DailyLoginStatGoldConsumeAmount)
	syndb.RegLazyWithMiddle(TbDailyLoginStat, DailyLoginStatDiamondConsumeAmount)
	syndb.RegLazyWithMiddle(TbDailyLoginStat, DailyLoginStatRechargeUserCount)
	syndb.RegLazyWithMiddle(TbDailyLoginStat, DailyLoginStatGoldConsumeUserCount)
	syndb.RegLazyWithMiddle(TbDailyLoginStat, DailyLoginStatDiamondConsumeUserCount)
	syndb.RegLazyWithMiddle(TbDailyLoginStat, DailyLoginStatAudienceUserCount)
	migrate.AutoMigrate(&DailyLoginStat{})
}
