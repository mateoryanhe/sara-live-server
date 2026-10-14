package entity

import (
	"fmt"

	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbDailyPromotionChannelStat db.TbName = "daily_promotion_channel_stats"
)

const (
	DailyPromotionChannelStatStatDate              db.TbCol = "stat_date"
	DailyPromotionChannelStatChannelName           db.TbCol = "channel_name"
	DailyPromotionChannelStatChannelCode           db.TbCol = "channel_code"
	DailyPromotionChannelStatRegisterCount         db.TbCol = "register_count"
	DailyPromotionChannelStatOrderCount            db.TbCol = "order_count"
	DailyPromotionChannelStatPayUserCount          db.TbCol = "pay_user_count"
	DailyPromotionChannelStatOrderTotalAmount      db.TbCol = "order_total_amount"
	DailyPromotionChannelStatRechargeSuccessAmount db.TbCol = "recharge_success_amount"
)

// DailyPromotionChannelStat 推广渠道日统计(主键=YYYY-MM-DD_渠道编码)
type DailyPromotionChannelStat struct {
	ID                    string  `gorm:"primaryKey;size:32;comment:复合ID(date_channelCode)" json:"id"`
	StatDate              string  `gorm:"size:10;index;comment:统计日期YYYY-MM-DD" json:"statDate"`
	ChannelName           string  `gorm:"size:128;default:'';comment:渠道名称" json:"channelName"`
	ChannelCode           int     `gorm:"index;comment:渠道编码" json:"channelCode"`
	RegisterCount         uint64  `gorm:"default:0;comment:注册人数" json:"registerCount"`
	OrderCount            uint64  `gorm:"default:0;comment:下单笔数" json:"orderCount"`
	PayUserCount          uint64  `gorm:"default:0;comment:支付人次(去重)" json:"payUserCount"`
	OrderTotalAmount      float64 `gorm:"type:decimal(10,4);default:0;comment:下单总金额(USD)" json:"orderTotalAmount"`
	RechargeSuccessAmount float64 `gorm:"type:decimal(10,4);default:0;comment:成功充值总金额(USD)" json:"rechargeSuccessAmount"`
}

func BuildDailyPromotionChannelStatId(date string, channelCode int) string {
	return fmt.Sprintf("%s_%d", date, channelCode)
}

func NewDailyPromotionChannelStat(date string, channelCode int, channelName string) *DailyPromotionChannelStat {
	r := &DailyPromotionChannelStat{
		ID:          BuildDailyPromotionChannelStatId(date, channelCode),
		StatDate:    date,
		ChannelCode: channelCode,
		ChannelName: channelName,
	}
	r.syncMetaCols()
	return r
}

func (r *DailyPromotionChannelStat) syncMetaCols() {
	syndb.AddData(TbDailyPromotionChannelStat, DailyPromotionChannelStatStatDate, &syndb.ColData{
		IdVal: r.ID, ColVal: r.StatDate,
	})
	syndb.AddData(TbDailyPromotionChannelStat, DailyPromotionChannelStatChannelCode, &syndb.ColData{
		IdVal: r.ID, ColVal: r.ChannelCode,
	})
	syndb.AddData(TbDailyPromotionChannelStat, DailyPromotionChannelStatChannelName, &syndb.ColData{
		IdVal: r.ID, ColVal: r.ChannelName,
	})
}

func (r *DailyPromotionChannelStat) SetChannelName(name string) {
	r.ChannelName = name
	syndb.AddData(TbDailyPromotionChannelStat, DailyPromotionChannelStatChannelName, &syndb.ColData{
		IdVal: r.ID, ColVal: name,
	})
}

func (r *DailyPromotionChannelStat) AddRegisterCount(n uint64) {
	r.RegisterCount = math.Add(r.RegisterCount, n)
	syndb.AddData(TbDailyPromotionChannelStat, DailyPromotionChannelStatRegisterCount, &syndb.ColData{
		IdVal: r.ID, ColVal: r.RegisterCount,
	})
}

func (r *DailyPromotionChannelStat) AddOrderCount(n uint64) {
	r.OrderCount = math.Add(r.OrderCount, n)
	syndb.AddData(TbDailyPromotionChannelStat, DailyPromotionChannelStatOrderCount, &syndb.ColData{
		IdVal: r.ID, ColVal: r.OrderCount,
	})
}

func (r *DailyPromotionChannelStat) AddPayUserCount(n uint64) {
	r.PayUserCount = math.Add(r.PayUserCount, n)
	syndb.AddData(TbDailyPromotionChannelStat, DailyPromotionChannelStatPayUserCount, &syndb.ColData{
		IdVal: r.ID, ColVal: r.PayUserCount,
	})
}

func (r *DailyPromotionChannelStat) AddOrderTotalAmount(amount float64) {
	r.OrderTotalAmount = math.AddFloat64(r.OrderTotalAmount, amount)
	syndb.AddData(TbDailyPromotionChannelStat, DailyPromotionChannelStatOrderTotalAmount, &syndb.ColData{
		IdVal: r.ID, ColVal: r.OrderTotalAmount,
	})
}

func (r *DailyPromotionChannelStat) AddRechargeSuccessAmount(amount float64) {
	r.RechargeSuccessAmount = math.AddFloat64(r.RechargeSuccessAmount, amount)
	syndb.AddData(TbDailyPromotionChannelStat, DailyPromotionChannelStatRechargeSuccessAmount, &syndb.ColData{
		IdVal: r.ID, ColVal: r.RechargeSuccessAmount,
	})
}

func initDailyPromotionChannelStat() {
	syndb.RegLazy(TbDailyPromotionChannelStat, DailyPromotionChannelStatStatDate)
	syndb.RegLazy(TbDailyPromotionChannelStat, DailyPromotionChannelStatChannelName)
	syndb.RegLazy(TbDailyPromotionChannelStat, DailyPromotionChannelStatChannelCode)
	syndb.RegLazy(TbDailyPromotionChannelStat, DailyPromotionChannelStatRegisterCount)
	syndb.RegLazy(TbDailyPromotionChannelStat, DailyPromotionChannelStatOrderCount)
	syndb.RegLazy(TbDailyPromotionChannelStat, DailyPromotionChannelStatPayUserCount)
	syndb.RegLazy(TbDailyPromotionChannelStat, DailyPromotionChannelStatOrderTotalAmount)
	syndb.RegLazy(TbDailyPromotionChannelStat, DailyPromotionChannelStatRechargeSuccessAmount)
	migrate.AutoMigrate(&DailyPromotionChannelStat{})
}
