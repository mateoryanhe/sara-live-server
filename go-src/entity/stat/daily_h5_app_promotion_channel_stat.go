package entity

import (
	"fmt"

	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyH5AppPromotionChannelStat db.TbName = "daily_h5_app_promotion_channel_stats"
)

const (
	DailyH5AppPromotionChannelStatStatDate              db.TbCol = "stat_date"
	DailyH5AppPromotionChannelStatChannelName           db.TbCol = "channel_name"
	DailyH5AppPromotionChannelStatChannelCode           db.TbCol = "channel_code"
	DailyH5AppPromotionChannelStatRegisterCount         db.TbCol = "register_count"
	DailyH5AppPromotionChannelStatDownloadClickCount    db.TbCol = "download_click_count"
	DailyH5AppPromotionChannelStatOrderCount            db.TbCol = "order_count"
	DailyH5AppPromotionChannelStatPayUserCount          db.TbCol = "pay_user_count"
	DailyH5AppPromotionChannelStatOrderTotalAmount      db.TbCol = "order_total_amount"
	DailyH5AppPromotionChannelStatRechargeSuccessAmount db.TbCol = "recharge_success_amount"
)

// DailyH5AppPromotionChannelStat H5下载App推广日统计(主键=YYYY-MM-DD_推广URL渠道编码)
type DailyH5AppPromotionChannelStat struct {
	ID                    string  `gorm:"primaryKey;size:32;comment:复合ID(date_channelCode)" json:"id"`
	StatDate              string  `gorm:"size:10;index;comment:统计日期YYYY-MM-DD" json:"statDate"`
	ChannelName           string  `gorm:"size:128;default:'';comment:渠道名称" json:"channelName"`
	ChannelCode           int     `gorm:"index;comment:推广URL渠道编码" json:"channelCode"`
	RegisterCount         uint64  `gorm:"default:0;comment:注册人数(配对成功)" json:"registerCount"`
	DownloadClickCount    uint64  `gorm:"default:0;comment:H5下载App点击次数" json:"downloadClickCount"`
	OrderCount            uint64  `gorm:"default:0;comment:下单笔数" json:"orderCount"`
	PayUserCount          uint64  `gorm:"default:0;comment:支付人次(去重)" json:"payUserCount"`
	OrderTotalAmount      float64 `gorm:"type:decimal(10,4);default:0;comment:下单总金额(USD)" json:"orderTotalAmount"`
	RechargeSuccessAmount float64 `gorm:"type:decimal(10,4);default:0;comment:成功充值总金额(USD)" json:"rechargeSuccessAmount"`
}

func BuildDailyH5AppPromotionChannelStatId(date string, channelCode int) string {
	return fmt.Sprintf("%s_%d", date, channelCode)
}

func NewDailyH5AppPromotionChannelStat(date string, channelCode int, channelName string) *DailyH5AppPromotionChannelStat {
	r := &DailyH5AppPromotionChannelStat{
		ID:          BuildDailyH5AppPromotionChannelStatId(date, channelCode),
		StatDate:    date,
		ChannelCode: channelCode,
		ChannelName: channelName,
	}
	r.syncMetaCols()
	return r
}

func (r *DailyH5AppPromotionChannelStat) syncMetaCols() {
	syndb.AddData(TbDailyH5AppPromotionChannelStat, DailyH5AppPromotionChannelStatStatDate, &syndb.ColData{
		IdVal: r.ID, ColVal: r.StatDate,
	})
	syndb.AddData(TbDailyH5AppPromotionChannelStat, DailyH5AppPromotionChannelStatChannelCode, &syndb.ColData{
		IdVal: r.ID, ColVal: r.ChannelCode,
	})
	syndb.AddData(TbDailyH5AppPromotionChannelStat, DailyH5AppPromotionChannelStatChannelName, &syndb.ColData{
		IdVal: r.ID, ColVal: r.ChannelName,
	})
}

func (r *DailyH5AppPromotionChannelStat) SetChannelName(name string) {
	r.ChannelName = name
	syndb.AddData(TbDailyH5AppPromotionChannelStat, DailyH5AppPromotionChannelStatChannelName, &syndb.ColData{
		IdVal: r.ID, ColVal: name,
	})
}

func (r *DailyH5AppPromotionChannelStat) AddRegisterCount(n uint64) {
	r.RegisterCount = math.Add(r.RegisterCount, n)
	syndb.AddData(TbDailyH5AppPromotionChannelStat, DailyH5AppPromotionChannelStatRegisterCount, &syndb.ColData{
		IdVal: r.ID, ColVal: r.RegisterCount,
	})
}

func (r *DailyH5AppPromotionChannelStat) AddDownloadClickCount(n uint64) {
	r.DownloadClickCount = math.Add(r.DownloadClickCount, n)
	syndb.AddData(TbDailyH5AppPromotionChannelStat, DailyH5AppPromotionChannelStatDownloadClickCount, &syndb.ColData{
		IdVal: r.ID, ColVal: r.DownloadClickCount,
	})
}

func (r *DailyH5AppPromotionChannelStat) AddOrderCount(n uint64) {
	r.OrderCount = math.Add(r.OrderCount, n)
	syndb.AddData(TbDailyH5AppPromotionChannelStat, DailyH5AppPromotionChannelStatOrderCount, &syndb.ColData{
		IdVal: r.ID, ColVal: r.OrderCount,
	})
}

func (r *DailyH5AppPromotionChannelStat) AddPayUserCount(n uint64) {
	r.PayUserCount = math.Add(r.PayUserCount, n)
	syndb.AddData(TbDailyH5AppPromotionChannelStat, DailyH5AppPromotionChannelStatPayUserCount, &syndb.ColData{
		IdVal: r.ID, ColVal: r.PayUserCount,
	})
}

func (r *DailyH5AppPromotionChannelStat) AddOrderTotalAmount(amount float64) {
	r.OrderTotalAmount = math.AddFloat64(r.OrderTotalAmount, amount)
	syndb.AddData(TbDailyH5AppPromotionChannelStat, DailyH5AppPromotionChannelStatOrderTotalAmount, &syndb.ColData{
		IdVal: r.ID, ColVal: r.OrderTotalAmount,
	})
}

func (r *DailyH5AppPromotionChannelStat) AddRechargeSuccessAmount(amount float64) {
	r.RechargeSuccessAmount = math.AddFloat64(r.RechargeSuccessAmount, amount)
	syndb.AddData(TbDailyH5AppPromotionChannelStat, DailyH5AppPromotionChannelStatRechargeSuccessAmount, &syndb.ColData{
		IdVal: r.ID, ColVal: r.RechargeSuccessAmount,
	})
}

func initDailyH5AppPromotionChannelStat() {
	syndb.RegLazy(TbDailyH5AppPromotionChannelStat, DailyH5AppPromotionChannelStatStatDate)
	syndb.RegLazy(TbDailyH5AppPromotionChannelStat, DailyH5AppPromotionChannelStatChannelName)
	syndb.RegLazy(TbDailyH5AppPromotionChannelStat, DailyH5AppPromotionChannelStatChannelCode)
	syndb.RegLazy(TbDailyH5AppPromotionChannelStat, DailyH5AppPromotionChannelStatRegisterCount)
	syndb.RegLazy(TbDailyH5AppPromotionChannelStat, DailyH5AppPromotionChannelStatDownloadClickCount)
	syndb.RegLazy(TbDailyH5AppPromotionChannelStat, DailyH5AppPromotionChannelStatOrderCount)
	syndb.RegLazy(TbDailyH5AppPromotionChannelStat, DailyH5AppPromotionChannelStatPayUserCount)
	syndb.RegLazy(TbDailyH5AppPromotionChannelStat, DailyH5AppPromotionChannelStatOrderTotalAmount)
	syndb.RegLazy(TbDailyH5AppPromotionChannelStat, DailyH5AppPromotionChannelStatRechargeSuccessAmount)
	migrate.AutoMigrate(&DailyH5AppPromotionChannelStat{})
}
