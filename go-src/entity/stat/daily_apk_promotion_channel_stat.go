package entity

import (
	"fmt"

	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyApkPromotionChannelStat db.TbName = "daily_apk_promotion_channel_stats"
)

const (
	DailyApkPromotionChannelStatStatDate              db.TbCol = "stat_date"
	DailyApkPromotionChannelStatChannelName           db.TbCol = "channel_name"
	DailyApkPromotionChannelStatChannelCode           db.TbCol = "channel_code"
	DailyApkPromotionChannelStatRegisterCount         db.TbCol = "register_count"
	DailyApkPromotionChannelStatDownloadClickCount    db.TbCol = "download_click_count"
	DailyApkPromotionChannelStatOrderCount            db.TbCol = "order_count"
	DailyApkPromotionChannelStatPayUserCount          db.TbCol = "pay_user_count"
	DailyApkPromotionChannelStatOrderTotalAmount      db.TbCol = "order_total_amount"
	DailyApkPromotionChannelStatRechargeSuccessAmount db.TbCol = "recharge_success_amount"
)

// DailyApkPromotionChannelStat APK 推广渠道日统计(主键=YYYY-MM-DD_渠道码)
type DailyApkPromotionChannelStat struct {
	ID                    string  `gorm:"primaryKey;size:80;comment:复合ID(date_channelCode)" json:"id"`
	StatDate              string  `gorm:"size:10;index;comment:统计日期YYYY-MM-DD" json:"statDate"`
	ChannelName           string  `gorm:"size:128;default:'';comment:渠道名称" json:"channelName"`
	ChannelCode           string  `gorm:"size:64;index;comment:渠道编码" json:"channelCode"`
	RegisterCount         uint64  `gorm:"default:0;comment:注册人数" json:"registerCount"`
	DownloadClickCount    uint64  `gorm:"default:0;comment:渠道短链下载点击次数" json:"downloadClickCount"`
	OrderCount            uint64  `gorm:"default:0;comment:下单笔数" json:"orderCount"`
	PayUserCount          uint64  `gorm:"default:0;comment:支付人次(去重)" json:"payUserCount"`
	OrderTotalAmount      float64 `gorm:"type:decimal(10,4);default:0;comment:下单总金额(USD)" json:"orderTotalAmount"`
	RechargeSuccessAmount float64 `gorm:"type:decimal(10,4);default:0;comment:成功充值总金额(USD)" json:"rechargeSuccessAmount"`
}

func BuildDailyApkPromotionChannelStatId(date, channelCode string) string {
	return fmt.Sprintf("%s_%s", date, channelCode)
}

func NewDailyApkPromotionChannelStat(date, channelCode, channelName string) *DailyApkPromotionChannelStat {
	r := &DailyApkPromotionChannelStat{
		ID:          BuildDailyApkPromotionChannelStatId(date, channelCode),
		StatDate:    date,
		ChannelCode: channelCode,
		ChannelName: channelName,
	}
	r.syncMetaCols()
	return r
}

func (r *DailyApkPromotionChannelStat) syncMetaCols() {
	syndb.AddData(TbDailyApkPromotionChannelStat, DailyApkPromotionChannelStatStatDate, &syndb.ColData{
		IdVal: r.ID, ColVal: r.StatDate,
	})
	syndb.AddData(TbDailyApkPromotionChannelStat, DailyApkPromotionChannelStatChannelCode, &syndb.ColData{
		IdVal: r.ID, ColVal: r.ChannelCode,
	})
	syndb.AddData(TbDailyApkPromotionChannelStat, DailyApkPromotionChannelStatChannelName, &syndb.ColData{
		IdVal: r.ID, ColVal: r.ChannelName,
	})
}

func (r *DailyApkPromotionChannelStat) SetChannelName(name string) {
	r.ChannelName = name
	syndb.AddData(TbDailyApkPromotionChannelStat, DailyApkPromotionChannelStatChannelName, &syndb.ColData{
		IdVal: r.ID, ColVal: name,
	})
}

func (r *DailyApkPromotionChannelStat) AddRegisterCount(n uint64) {
	r.RegisterCount = math.Add(r.RegisterCount, n)
	syndb.AddData(TbDailyApkPromotionChannelStat, DailyApkPromotionChannelStatRegisterCount, &syndb.ColData{
		IdVal: r.ID, ColVal: r.RegisterCount,
	})
}

func (r *DailyApkPromotionChannelStat) AddDownloadClickCount(n uint64) {
	r.DownloadClickCount = math.Add(r.DownloadClickCount, n)
	syndb.AddData(TbDailyApkPromotionChannelStat, DailyApkPromotionChannelStatDownloadClickCount, &syndb.ColData{
		IdVal: r.ID, ColVal: r.DownloadClickCount,
	})
}

func (r *DailyApkPromotionChannelStat) AddOrderCount(n uint64) {
	r.OrderCount = math.Add(r.OrderCount, n)
	syndb.AddData(TbDailyApkPromotionChannelStat, DailyApkPromotionChannelStatOrderCount, &syndb.ColData{
		IdVal: r.ID, ColVal: r.OrderCount,
	})
}

func (r *DailyApkPromotionChannelStat) AddPayUserCount(n uint64) {
	r.PayUserCount = math.Add(r.PayUserCount, n)
	syndb.AddData(TbDailyApkPromotionChannelStat, DailyApkPromotionChannelStatPayUserCount, &syndb.ColData{
		IdVal: r.ID, ColVal: r.PayUserCount,
	})
}

func (r *DailyApkPromotionChannelStat) AddOrderTotalAmount(amount float64) {
	r.OrderTotalAmount = math.AddFloat64(r.OrderTotalAmount, amount)
	syndb.AddData(TbDailyApkPromotionChannelStat, DailyApkPromotionChannelStatOrderTotalAmount, &syndb.ColData{
		IdVal: r.ID, ColVal: r.OrderTotalAmount,
	})
}

func (r *DailyApkPromotionChannelStat) AddRechargeSuccessAmount(amount float64) {
	r.RechargeSuccessAmount = math.AddFloat64(r.RechargeSuccessAmount, amount)
	syndb.AddData(TbDailyApkPromotionChannelStat, DailyApkPromotionChannelStatRechargeSuccessAmount, &syndb.ColData{
		IdVal: r.ID, ColVal: r.RechargeSuccessAmount,
	})
}

func initDailyApkPromotionChannelStat() {
	syndb.RegLazy(TbDailyApkPromotionChannelStat, DailyApkPromotionChannelStatStatDate)
	syndb.RegLazy(TbDailyApkPromotionChannelStat, DailyApkPromotionChannelStatChannelName)
	syndb.RegLazy(TbDailyApkPromotionChannelStat, DailyApkPromotionChannelStatChannelCode)
	syndb.RegLazy(TbDailyApkPromotionChannelStat, DailyApkPromotionChannelStatRegisterCount)
	syndb.RegLazy(TbDailyApkPromotionChannelStat, DailyApkPromotionChannelStatDownloadClickCount)
	syndb.RegLazy(TbDailyApkPromotionChannelStat, DailyApkPromotionChannelStatOrderCount)
	syndb.RegLazy(TbDailyApkPromotionChannelStat, DailyApkPromotionChannelStatPayUserCount)
	syndb.RegLazy(TbDailyApkPromotionChannelStat, DailyApkPromotionChannelStatOrderTotalAmount)
	syndb.RegLazy(TbDailyApkPromotionChannelStat, DailyApkPromotionChannelStatRechargeSuccessAmount)
	migrate.AutoMigrate(&DailyApkPromotionChannelStat{})
}
