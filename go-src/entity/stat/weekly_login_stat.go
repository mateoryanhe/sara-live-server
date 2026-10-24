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
	WeeklyLoginStatCount                      db.TbCol = "count"
	WeeklyLoginStatRegisterCount              db.TbCol = "register_count"
	WeeklyLoginStatRechargeAmount             db.TbCol = "recharge_amount"
	WeeklyLoginStatNormalUserRechargeAmount   db.TbCol = "normal_user_recharge_amount"
	WeeklyLoginStatCoinMerchantRechargeAmount db.TbCol = "coin_merchant_recharge_amount"
	WeeklyLoginStatAppDeviceRechargeAmount    db.TbCol = "app_device_recharge_amount"
	WeeklyLoginStatH5DeviceRechargeAmount     db.TbCol = "h5_device_recharge_amount"
	WeeklyLoginStatAppDeviceLoginCount        db.TbCol = "app_device_login_count"
	WeeklyLoginStatH5DeviceLoginCount         db.TbCol = "h5_device_login_count"
	WeeklyLoginStatAnchorPayoutAmount         db.TbCol = "anchor_payout_amount"
	WeeklyLoginStatGoldConsumeAmount          db.TbCol = "gold_consume_amount"
	WeeklyLoginStatDiamondConsumeAmount       db.TbCol = "diamond_consume_amount"
	WeeklyLoginStatRechargeUserCount          db.TbCol = "recharge_user_count"
	WeeklyLoginStatAudienceUserCount          db.TbCol = "audience_user_count"
)

// WeeklyLoginStat 每周登录统计(主键ID即周标识,如 2026-W21)
type WeeklyLoginStat struct {
	ID                         string  `gorm:"primaryKey;size:16;comment:周标识(YYYY-Www)" json:"week"`
	Count                      uint64  `gorm:"default:0;comment:登录数量" json:"count"`
	RegisterCount              uint64  `gorm:"default:0;comment:注册人数" json:"registerCount"`
	RechargeAmount             float64 `gorm:"type:decimal(10,4);default:0;comment:全部真实美金入账(USD)" json:"rechargeAmount"`
	NormalUserRechargeAmount   float64 `gorm:"type:decimal(10,4);default:0;comment:普通用户真实美金入账(USD)" json:"normalUserRechargeAmount"`
	CoinMerchantRechargeAmount float64 `gorm:"type:decimal(10,4);default:0;comment:币商真实美金入账(USD)" json:"coinMerchantRechargeAmount"`
	AppDeviceRechargeAmount    float64 `gorm:"type:decimal(10,4);default:0;comment:App设备码用户真实美金入账(USD)" json:"appDeviceRechargeAmount"`
	H5DeviceRechargeAmount     float64 `gorm:"type:decimal(10,4);default:0;comment:H5设备码用户真实美金入账(USD)" json:"h5DeviceRechargeAmount"`
	AppDeviceLoginCount        uint64  `gorm:"default:0;comment:App设备码用户登录人数(去重)" json:"appDeviceLoginCount"`
	H5DeviceLoginCount         uint64  `gorm:"default:0;comment:H5设备码用户登录人数(去重)" json:"h5DeviceLoginCount"`
	AnchorPayoutAmount         float64 `gorm:"type:decimal(16,4);default:0;comment:主播代付金额(USD)" json:"anchorPayoutAmount"`
	GoldConsumeAmount          float64 `gorm:"default:0;comment:金币消费金额" json:"goldConsumeAmount"`
	DiamondConsumeAmount       float64 `gorm:"default:0;comment:钻石消费金额" json:"diamondConsumeAmount"`
	RechargeUserCount          uint64  `gorm:"default:0;comment:充值人数(去重)" json:"rechargeUserCount"`
	AudienceUserCount          uint64  `gorm:"default:0;comment:有效观众人数(去重,跨直播间)" json:"audienceUserCount"`
}

// FormatWeeklyLoginStatKey 格式化周统计标识(ISO周)
func FormatWeeklyLoginStatKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", year, week)
}

func NewWeeklyLoginStat(week string) *WeeklyLoginStat {
	return &WeeklyLoginStat{
		ID:                         week,
		Count:                      0,
		RegisterCount:              0,
		RechargeAmount:             0,
		NormalUserRechargeAmount:   0,
		CoinMerchantRechargeAmount: 0,
		AppDeviceRechargeAmount:    0,
		H5DeviceRechargeAmount:     0,
		AppDeviceLoginCount:        0,
		H5DeviceLoginCount:         0,
		AnchorPayoutAmount:         0,
		GoldConsumeAmount:          0,
		DiamondConsumeAmount:       0,
		RechargeUserCount:          0,
		AudienceUserCount:          0,
	}
}

func (r *WeeklyLoginStat) AddLoginCount(n uint64) {
	r.Count = math.Add(r.Count, n)
	syndb.AddData(TbWeeklyLoginStat, WeeklyLoginStatCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.Count,
	})
}

func (r *WeeklyLoginStat) AddRegisterCount(n uint64) {
	r.RegisterCount = math.Add(r.RegisterCount, n)
	syndb.AddData(TbWeeklyLoginStat, WeeklyLoginStatRegisterCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.RegisterCount,
	})
}

func (r *WeeklyLoginStat) AddRechargeAmount(val float64) {
	r.RechargeAmount = math.AddFloat64(r.RechargeAmount, val)
	syndb.AddData(TbWeeklyLoginStat, WeeklyLoginStatRechargeAmount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.RechargeAmount,
	})
}

func (r *WeeklyLoginStat) AddNormalUserRechargeAmount(val float64) {
	r.NormalUserRechargeAmount = math.AddFloat64(r.NormalUserRechargeAmount, val)
	syndb.AddData(TbWeeklyLoginStat, WeeklyLoginStatNormalUserRechargeAmount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.NormalUserRechargeAmount,
	})
}

func (r *WeeklyLoginStat) AddCoinMerchantRechargeAmount(val float64) {
	r.CoinMerchantRechargeAmount = math.AddFloat64(r.CoinMerchantRechargeAmount, val)
	syndb.AddData(TbWeeklyLoginStat, WeeklyLoginStatCoinMerchantRechargeAmount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.CoinMerchantRechargeAmount,
	})
}

func (r *WeeklyLoginStat) AddAppDeviceRechargeAmount(val float64) {
	r.AppDeviceRechargeAmount = math.AddFloat64(r.AppDeviceRechargeAmount, val)
	syndb.AddData(TbWeeklyLoginStat, WeeklyLoginStatAppDeviceRechargeAmount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.AppDeviceRechargeAmount,
	})
}

func (r *WeeklyLoginStat) AddH5DeviceRechargeAmount(val float64) {
	r.H5DeviceRechargeAmount = math.AddFloat64(r.H5DeviceRechargeAmount, val)
	syndb.AddData(TbWeeklyLoginStat, WeeklyLoginStatH5DeviceRechargeAmount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.H5DeviceRechargeAmount,
	})
}

func (r *WeeklyLoginStat) AddAppDeviceLoginCount(n uint64) {
	r.AppDeviceLoginCount = math.Add(r.AppDeviceLoginCount, n)
	syndb.AddData(TbWeeklyLoginStat, WeeklyLoginStatAppDeviceLoginCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.AppDeviceLoginCount,
	})
}

func (r *WeeklyLoginStat) AddH5DeviceLoginCount(n uint64) {
	r.H5DeviceLoginCount = math.Add(r.H5DeviceLoginCount, n)
	syndb.AddData(TbWeeklyLoginStat, WeeklyLoginStatH5DeviceLoginCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.H5DeviceLoginCount,
	})
}

func (r *WeeklyLoginStat) AddAnchorPayoutAmount(val float64) {
	r.AnchorPayoutAmount = math.AddFloat64(r.AnchorPayoutAmount, val)
	syndb.AddData(TbWeeklyLoginStat, WeeklyLoginStatAnchorPayoutAmount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.AnchorPayoutAmount,
	})
}

func (r *WeeklyLoginStat) AddGoldConsumeAmount(val float64) {
	r.GoldConsumeAmount = math.AddFloat64(r.GoldConsumeAmount, val)
	syndb.AddData(TbWeeklyLoginStat, WeeklyLoginStatGoldConsumeAmount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.GoldConsumeAmount,
	})
}

func (r *WeeklyLoginStat) AddDiamondConsumeAmount(val float64) {
	r.DiamondConsumeAmount = math.AddFloat64(r.DiamondConsumeAmount, val)
	syndb.AddData(TbWeeklyLoginStat, WeeklyLoginStatDiamondConsumeAmount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.DiamondConsumeAmount,
	})
}

func (r *WeeklyLoginStat) AddRechargeUserCount(n uint64) {
	r.RechargeUserCount = math.Add(r.RechargeUserCount, n)
	syndb.AddData(TbWeeklyLoginStat, WeeklyLoginStatRechargeUserCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.RechargeUserCount,
	})
}

func (r *WeeklyLoginStat) AddAudienceUserCount(n uint64) {
	r.AudienceUserCount = math.Add(r.AudienceUserCount, n)
	syndb.AddData(TbWeeklyLoginStat, WeeklyLoginStatAudienceUserCount, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.AudienceUserCount,
	})
}

func initWeeklyLoginStat() {
	syndb.RegLazy(TbWeeklyLoginStat, WeeklyLoginStatCount)
	syndb.RegLazy(TbWeeklyLoginStat, WeeklyLoginStatRegisterCount)
	syndb.RegLazy(TbWeeklyLoginStat, WeeklyLoginStatRechargeAmount)
	syndb.RegLazy(TbWeeklyLoginStat, WeeklyLoginStatNormalUserRechargeAmount)
	syndb.RegLazy(TbWeeklyLoginStat, WeeklyLoginStatCoinMerchantRechargeAmount)
	syndb.RegLazy(TbWeeklyLoginStat, WeeklyLoginStatAppDeviceRechargeAmount)
	syndb.RegLazy(TbWeeklyLoginStat, WeeklyLoginStatH5DeviceRechargeAmount)
	syndb.RegLazy(TbWeeklyLoginStat, WeeklyLoginStatAppDeviceLoginCount)
	syndb.RegLazy(TbWeeklyLoginStat, WeeklyLoginStatH5DeviceLoginCount)
	syndb.RegLazy(TbWeeklyLoginStat, WeeklyLoginStatAnchorPayoutAmount)
	syndb.RegLazy(TbWeeklyLoginStat, WeeklyLoginStatGoldConsumeAmount)
	syndb.RegLazy(TbWeeklyLoginStat, WeeklyLoginStatDiamondConsumeAmount)
	syndb.RegLazy(TbWeeklyLoginStat, WeeklyLoginStatRechargeUserCount)
	syndb.RegLazy(TbWeeklyLoginStat, WeeklyLoginStatAudienceUserCount)
	migrate.AutoMigrate(&WeeklyLoginStat{})
}
