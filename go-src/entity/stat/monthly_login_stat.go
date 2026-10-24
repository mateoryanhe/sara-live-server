package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbMonthlyLoginStat db.TbName = "monthly_login_stats"
)

const (
	MonthlyLoginStatCount                      db.TbCol = "count"
	MonthlyLoginStatRegisterCount              db.TbCol = "register_count"
	MonthlyLoginStatRechargeAmount             db.TbCol = "recharge_amount"
	MonthlyLoginStatNormalUserRechargeAmount   db.TbCol = "normal_user_recharge_amount"
	MonthlyLoginStatCoinMerchantRechargeAmount db.TbCol = "coin_merchant_recharge_amount"
	MonthlyLoginStatAppDeviceRechargeAmount    db.TbCol = "app_device_recharge_amount"
	MonthlyLoginStatH5DeviceRechargeAmount     db.TbCol = "h5_device_recharge_amount"
	MonthlyLoginStatAppDeviceLoginCount        db.TbCol = "app_device_login_count"
	MonthlyLoginStatH5DeviceLoginCount         db.TbCol = "h5_device_login_count"
	MonthlyLoginStatGoldConsumeAmount          db.TbCol = "gold_consume_amount"
	MonthlyLoginStatDiamondConsumeAmount       db.TbCol = "diamond_consume_amount"
	MonthlyLoginStatRechargeUserCount          db.TbCol = "recharge_user_count"
	MonthlyLoginStatAudienceUserCount          db.TbCol = "audience_user_count"
)

// MonthlyLoginStat 每月登录统计(主键ID即月标识,如 2026-05)
type MonthlyLoginStat struct {
	ID                         string  `gorm:"primaryKey;size:10;comment:月标识(YYYY-MM)" json:"month"`
	Count                      uint64  `gorm:"default:0;comment:登录数量" json:"count"`
	RegisterCount              uint64  `gorm:"default:0;comment:注册人数" json:"registerCount"`
	RechargeAmount             float64 `gorm:"type:decimal(10,4);default:0;comment:全部真实美金入账(USD)" json:"rechargeAmount"`
	NormalUserRechargeAmount   float64 `gorm:"type:decimal(10,4);default:0;comment:普通用户真实美金入账(USD)" json:"normalUserRechargeAmount"`
	CoinMerchantRechargeAmount float64 `gorm:"type:decimal(10,4);default:0;comment:币商真实美金入账(USD)" json:"coinMerchantRechargeAmount"`
	AppDeviceRechargeAmount    float64 `gorm:"type:decimal(10,4);default:0;comment:App设备码用户真实美金入账(USD)" json:"appDeviceRechargeAmount"`
	H5DeviceRechargeAmount     float64 `gorm:"type:decimal(10,4);default:0;comment:H5设备码用户真实美金入账(USD)" json:"h5DeviceRechargeAmount"`
	AppDeviceLoginCount        uint64  `gorm:"default:0;comment:App设备码用户登录人数(去重)" json:"appDeviceLoginCount"`
	H5DeviceLoginCount         uint64  `gorm:"default:0;comment:H5设备码用户登录人数(去重)" json:"h5DeviceLoginCount"`
	GoldConsumeAmount          float64 `gorm:"default:0;comment:金币消费金额" json:"goldConsumeAmount"`
	DiamondConsumeAmount       float64 `gorm:"default:0;comment:钻石消费金额" json:"diamondConsumeAmount"`
	RechargeUserCount          uint64  `gorm:"default:0;comment:充值人数(去重)" json:"rechargeUserCount"`
	AudienceUserCount          uint64  `gorm:"default:0;comment:有效观众人数(去重,跨直播间)" json:"audienceUserCount"`
}

// FormatMonthlyLoginStatKey 格式化月统计标识
func FormatMonthlyLoginStatKey(t time.Time) string {
	return t.Format("2006-01")
}

func NewMonthlyLoginStat(month string) *MonthlyLoginStat {
	return &MonthlyLoginStat{
		ID:                         month,
		Count:                      0,
		RegisterCount:              0,
		RechargeAmount:             0,
		NormalUserRechargeAmount:   0,
		CoinMerchantRechargeAmount: 0,
		AppDeviceRechargeAmount:    0,
		H5DeviceRechargeAmount:     0,
		AppDeviceLoginCount:        0,
		H5DeviceLoginCount:         0,
		GoldConsumeAmount:          0,
		DiamondConsumeAmount:       0,
		RechargeUserCount:          0,
		AudienceUserCount:          0,
	}
}

func (r *MonthlyLoginStat) AddLoginCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbMonthlyLoginStat, MonthlyLoginStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func (r *MonthlyLoginStat) AddRegisterCount(n uint64) {
	r.RegisterCount = math.Add(r.RegisterCount, n)
	syndb.AddData(TbMonthlyLoginStat, MonthlyLoginStatRegisterCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.RegisterCount,
	})
}

func (r *MonthlyLoginStat) AddRechargeAmount(val float64) {
	r.RechargeAmount = math.AddFloat64(r.RechargeAmount, val)
	syndb.AddData(TbMonthlyLoginStat, MonthlyLoginStatRechargeAmount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.RechargeAmount,
	})
}

func (r *MonthlyLoginStat) AddNormalUserRechargeAmount(val float64) {
	r.NormalUserRechargeAmount = math.AddFloat64(r.NormalUserRechargeAmount, val)
	syndb.AddData(TbMonthlyLoginStat, MonthlyLoginStatNormalUserRechargeAmount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.NormalUserRechargeAmount,
	})
}

func (r *MonthlyLoginStat) AddCoinMerchantRechargeAmount(val float64) {
	r.CoinMerchantRechargeAmount = math.AddFloat64(r.CoinMerchantRechargeAmount, val)
	syndb.AddData(TbMonthlyLoginStat, MonthlyLoginStatCoinMerchantRechargeAmount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.CoinMerchantRechargeAmount,
	})
}

func (r *MonthlyLoginStat) AddAppDeviceRechargeAmount(val float64) {
	r.AppDeviceRechargeAmount = math.AddFloat64(r.AppDeviceRechargeAmount, val)
	syndb.AddData(TbMonthlyLoginStat, MonthlyLoginStatAppDeviceRechargeAmount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.AppDeviceRechargeAmount,
	})
}

func (r *MonthlyLoginStat) AddH5DeviceRechargeAmount(val float64) {
	r.H5DeviceRechargeAmount = math.AddFloat64(r.H5DeviceRechargeAmount, val)
	syndb.AddData(TbMonthlyLoginStat, MonthlyLoginStatH5DeviceRechargeAmount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.H5DeviceRechargeAmount,
	})
}

func (r *MonthlyLoginStat) AddAppDeviceLoginCount(n uint64) {
	r.AppDeviceLoginCount = math.Add(r.AppDeviceLoginCount, n)
	syndb.AddData(TbMonthlyLoginStat, MonthlyLoginStatAppDeviceLoginCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.AppDeviceLoginCount,
	})
}

func (r *MonthlyLoginStat) AddH5DeviceLoginCount(n uint64) {
	r.H5DeviceLoginCount = math.Add(r.H5DeviceLoginCount, n)
	syndb.AddData(TbMonthlyLoginStat, MonthlyLoginStatH5DeviceLoginCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.H5DeviceLoginCount,
	})
}

func (r *MonthlyLoginStat) AddGoldConsumeAmount(val float64) {
	r.GoldConsumeAmount = math.AddFloat64(r.GoldConsumeAmount, val)
	syndb.AddData(TbMonthlyLoginStat, MonthlyLoginStatGoldConsumeAmount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.GoldConsumeAmount,
	})
}

func (r *MonthlyLoginStat) AddDiamondConsumeAmount(val float64) {
	r.DiamondConsumeAmount = math.AddFloat64(r.DiamondConsumeAmount, val)
	syndb.AddData(TbMonthlyLoginStat, MonthlyLoginStatDiamondConsumeAmount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.DiamondConsumeAmount,
	})
}

func (r *MonthlyLoginStat) AddRechargeUserCount(n uint64) {
	r.RechargeUserCount = math.Add(r.RechargeUserCount, n)
	syndb.AddData(TbMonthlyLoginStat, MonthlyLoginStatRechargeUserCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.RechargeUserCount,
	})
}

func (r *MonthlyLoginStat) AddAudienceUserCount(n uint64) {
	r.AudienceUserCount = math.Add(r.AudienceUserCount, n)
	syndb.AddData(TbMonthlyLoginStat, MonthlyLoginStatAudienceUserCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.AudienceUserCount,
	})
}

func initMonthlyLoginStat() {
	syndb.RegLazy(TbMonthlyLoginStat, MonthlyLoginStatCount)
	syndb.RegLazy(TbMonthlyLoginStat, MonthlyLoginStatRegisterCount)
	syndb.RegLazy(TbMonthlyLoginStat, MonthlyLoginStatRechargeAmount)
	syndb.RegLazy(TbMonthlyLoginStat, MonthlyLoginStatNormalUserRechargeAmount)
	syndb.RegLazy(TbMonthlyLoginStat, MonthlyLoginStatCoinMerchantRechargeAmount)
	syndb.RegLazy(TbMonthlyLoginStat, MonthlyLoginStatAppDeviceRechargeAmount)
	syndb.RegLazy(TbMonthlyLoginStat, MonthlyLoginStatH5DeviceRechargeAmount)
	syndb.RegLazy(TbMonthlyLoginStat, MonthlyLoginStatAppDeviceLoginCount)
	syndb.RegLazy(TbMonthlyLoginStat, MonthlyLoginStatH5DeviceLoginCount)
	syndb.RegLazy(TbMonthlyLoginStat, MonthlyLoginStatGoldConsumeAmount)
	syndb.RegLazy(TbMonthlyLoginStat, MonthlyLoginStatDiamondConsumeAmount)
	syndb.RegLazy(TbMonthlyLoginStat, MonthlyLoginStatRechargeUserCount)
	syndb.RegLazy(TbMonthlyLoginStat, MonthlyLoginStatAudienceUserCount)
	migrate.AutoMigrate(&MonthlyLoginStat{})
}
