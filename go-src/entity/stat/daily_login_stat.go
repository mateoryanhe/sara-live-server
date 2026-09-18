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
	DailyLoginStatCount                      db.TbCol = "count"
	DailyLoginStatRegisterCount              db.TbCol = "register_count"
	DailyLoginStatRechargeAmount             db.TbCol = "recharge_amount"
	DailyLoginStatNormalUserRechargeAmount   db.TbCol = "normal_user_recharge_amount"
	DailyLoginStatCoinMerchantRechargeAmount db.TbCol = "coin_merchant_recharge_amount"
	DailyLoginStatGoldConsumeAmount          db.TbCol = "gold_consume_amount"
	DailyLoginStatDiamondConsumeAmount       db.TbCol = "diamond_consume_amount"
	DailyLoginStatRechargeUserCount          db.TbCol = "recharge_user_count"
	DailyLoginStatGoldConsumeUserCount       db.TbCol = "gold_consume_user_count"
	DailyLoginStatDiamondConsumeUserCount    db.TbCol = "diamond_consume_user_count"
	DailyLoginStatAudienceUserCount          db.TbCol = "audience_user_count"
)

// DailyLoginStat 每日登录统计(主键ID即日期 YYYY-MM-DD)
type DailyLoginStat struct {
	ID                         string  `gorm:"primaryKey;size:10;comment:日期(YYYY-MM-DD)" json:"date"`
	Count                      uint64  `gorm:"default:0;comment:登录数量" json:"count"`
	RegisterCount              uint64  `gorm:"default:0;comment:注册人数" json:"registerCount"`
	RechargeAmount             float64 `gorm:"type:decimal(10,4);default:0;comment:全部真实美金入账(USD)" json:"rechargeAmount"`
	NormalUserRechargeAmount   float64 `gorm:"type:decimal(10,4);default:0;comment:普通用户真实美金入账(USD)" json:"normalUserRechargeAmount"`
	CoinMerchantRechargeAmount float64 `gorm:"type:decimal(10,4);default:0;comment:币商真实美金入账(USD)" json:"coinMerchantRechargeAmount"`
	GoldConsumeAmount          float64 `gorm:"default:0;comment:金币消费金额" json:"goldConsumeAmount"`
	DiamondConsumeAmount       float64 `gorm:"default:0;comment:钻石消费金额" json:"diamondConsumeAmount"`
	RechargeUserCount          uint64  `gorm:"default:0;comment:充值人数(去重)" json:"rechargeUserCount"`
	GoldConsumeUserCount       uint64  `gorm:"default:0;comment:金币消费人数(去重)" json:"goldConsumeUserCount"`
	DiamondConsumeUserCount    uint64  `gorm:"default:0;comment:钻石消费人数(去重)" json:"diamondConsumeUserCount"`
	AudienceUserCount          uint64  `gorm:"default:0;comment:有效观众人数(去重,跨直播间)" json:"audienceUserCount"`
}

// FormatDailyLoginStatDate 格式化统计日期
func FormatDailyLoginStatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func NewDailyLoginStat(date string) *DailyLoginStat {
	return &DailyLoginStat{
		ID:                         date,
		Count:                      0,
		RegisterCount:              0,
		RechargeAmount:             0,
		NormalUserRechargeAmount:   0,
		CoinMerchantRechargeAmount: 0,
		GoldConsumeAmount:          0,
		DiamondConsumeAmount:       0,
		RechargeUserCount:          0,
		GoldConsumeUserCount:       0,
		DiamondConsumeUserCount:    0,
		AudienceUserCount:          0,
	}
}

func (receiver *DailyLoginStat) AddLoginCount(n uint64) {
	receiver.Count = math.Add(receiver.Count, n)
	syndb.AddData(TbDailyLoginStat, DailyLoginStatCount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.Count,
	})
}

func (receiver *DailyLoginStat) AddRegisterCount(n uint64) {
	receiver.RegisterCount = math.Add(receiver.RegisterCount, n)
	syndb.AddData(TbDailyLoginStat, DailyLoginStatRegisterCount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.RegisterCount,
	})
}

func (receiver *DailyLoginStat) AddRechargeAmount(val float64) {
	receiver.RechargeAmount = math.AddFloat64(receiver.RechargeAmount, val)
	syndb.AddData(TbDailyLoginStat, DailyLoginStatRechargeAmount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.RechargeAmount,
	})
}

func (receiver *DailyLoginStat) AddNormalUserRechargeAmount(val float64) {
	receiver.NormalUserRechargeAmount = math.AddFloat64(receiver.NormalUserRechargeAmount, val)
	syndb.AddData(TbDailyLoginStat, DailyLoginStatNormalUserRechargeAmount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.NormalUserRechargeAmount,
	})
}

func (receiver *DailyLoginStat) AddCoinMerchantRechargeAmount(val float64) {
	receiver.CoinMerchantRechargeAmount = math.AddFloat64(receiver.CoinMerchantRechargeAmount, val)
	syndb.AddData(TbDailyLoginStat, DailyLoginStatCoinMerchantRechargeAmount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.CoinMerchantRechargeAmount,
	})
}

func (receiver *DailyLoginStat) AddGoldConsumeAmount(val float64) {
	receiver.GoldConsumeAmount = math.AddFloat64(receiver.GoldConsumeAmount, val)
	syndb.AddData(TbDailyLoginStat, DailyLoginStatGoldConsumeAmount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.GoldConsumeAmount,
	})
}

func (receiver *DailyLoginStat) AddDiamondConsumeAmount(val float64) {
	receiver.DiamondConsumeAmount = math.AddFloat64(receiver.DiamondConsumeAmount, val)
	syndb.AddData(TbDailyLoginStat, DailyLoginStatDiamondConsumeAmount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.DiamondConsumeAmount,
	})
}

func (receiver *DailyLoginStat) AddRechargeUserCount(n uint64) {
	receiver.RechargeUserCount = math.Add(receiver.RechargeUserCount, n)
	syndb.AddData(TbDailyLoginStat, DailyLoginStatRechargeUserCount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.RechargeUserCount,
	})
}

func (receiver *DailyLoginStat) AddGoldConsumeUserCount(n uint64) {
	receiver.GoldConsumeUserCount = math.Add(receiver.GoldConsumeUserCount, n)
	syndb.AddData(TbDailyLoginStat, DailyLoginStatGoldConsumeUserCount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.GoldConsumeUserCount,
	})
}

func (receiver *DailyLoginStat) AddDiamondConsumeUserCount(n uint64) {
	receiver.DiamondConsumeUserCount = math.Add(receiver.DiamondConsumeUserCount, n)
	syndb.AddData(TbDailyLoginStat, DailyLoginStatDiamondConsumeUserCount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.DiamondConsumeUserCount,
	})
}

func (receiver *DailyLoginStat) AddAudienceUserCount(n uint64) {
	receiver.AudienceUserCount = math.Add(receiver.AudienceUserCount, n)
	syndb.AddData(TbDailyLoginStat, DailyLoginStatAudienceUserCount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.AudienceUserCount,
	})
}

func initDailyLoginStat() {
	syndb.RegLazy(TbDailyLoginStat, DailyLoginStatCount)
	syndb.RegLazy(TbDailyLoginStat, DailyLoginStatRegisterCount)
	syndb.RegLazy(TbDailyLoginStat, DailyLoginStatRechargeAmount)
	syndb.RegLazy(TbDailyLoginStat, DailyLoginStatNormalUserRechargeAmount)
	syndb.RegLazy(TbDailyLoginStat, DailyLoginStatCoinMerchantRechargeAmount)
	syndb.RegLazy(TbDailyLoginStat, DailyLoginStatGoldConsumeAmount)
	syndb.RegLazy(TbDailyLoginStat, DailyLoginStatDiamondConsumeAmount)
	syndb.RegLazy(TbDailyLoginStat, DailyLoginStatRechargeUserCount)
	syndb.RegLazy(TbDailyLoginStat, DailyLoginStatGoldConsumeUserCount)
	syndb.RegLazy(TbDailyLoginStat, DailyLoginStatDiamondConsumeUserCount)
	syndb.RegLazy(TbDailyLoginStat, DailyLoginStatAudienceUserCount)
	migrate.AutoMigrate(&DailyLoginStat{})
}
