package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbLiveRoomIncomeTotal db.TbName = "live_room_income_totals"
)

// LiveRoomIncomeTotal 直播间生涯累计收益(只增不结算,用于展示主播获得了多少)
type LiveRoomIncomeTotal struct {
	migrate.OneModel
	LiveRoomIncomeAmounts
	SettlementSalary      float64 `gorm:"type:decimal(16,4);default:0;comment:结算薪资" json:"settlementSalary"`
	SettlementShareAmount float64 `gorm:"type:decimal(16,4);default:0;comment:结算分佣金额" json:"settlementShareAmount"`
}

func NewLiveRoomIncomeTotal(roomId uint64) *LiveRoomIncomeTotal {
	ret := &LiveRoomIncomeTotal{}
	ret.ID = roomId
	now := time.Now()
	ret.CreatedAt = now
	ret.UpdatedAt = now
	syndb.AddData(TbLiveRoomIncomeTotal, db.CreatedAtName, &syndb.ColData{IdVal: roomId, ColVal: now})
	syndb.AddData(TbLiveRoomIncomeTotal, db.UpdatedAtName, &syndb.ColData{IdVal: roomId, ColVal: now})
	return ret
}

func (r *LiveRoomIncomeTotal) AddTotalIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeTotal, LiveRoomIncomeTotalIncome, r.ID, &r.TotalIncome, v, false, &r.UpdatedAt)
}
func (r *LiveRoomIncomeTotal) AddTotalGiftIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeTotal, LiveRoomIncomeTotalGiftIncome, r.ID, &r.TotalGiftIncome, v, true, &r.UpdatedAt)
}
func (r *LiveRoomIncomeTotal) AddTotalPaidDanmakuIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeTotal, LiveRoomIncomeTotalPaidDanmakuIncome, r.ID, &r.TotalPaidDanmakuIncome, v, true, &r.UpdatedAt)
}
func (r *LiveRoomIncomeTotal) AddTotalPrivateRoomTicketIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeTotal, LiveRoomIncomeTotalPrivateRoomTicketIncome, r.ID, &r.TotalPrivateRoomTicketIncome, v, true, &r.UpdatedAt)
}
func (r *LiveRoomIncomeTotal) AddTotalPrivateRoomWatchIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeTotal, LiveRoomIncomeTotalPrivateRoomWatchIncome, r.ID, &r.TotalPrivateRoomWatchIncome, v, true, &r.UpdatedAt)
}
func (r *LiveRoomIncomeTotal) AddTotalVideoCallIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeTotal, LiveRoomIncomeTotalVideoCallIncome, r.ID, &r.TotalVideoCallIncome, v, true, &r.UpdatedAt)
}
func (r *LiveRoomIncomeTotal) AddTotalVideoCallTicketIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeTotal, LiveRoomIncomeTotalVideoCallTicketIncome, r.ID, &r.TotalVideoCallTicketIncome, v, true, &r.UpdatedAt)
}
func (r *LiveRoomIncomeTotal) AddTotalVideoCallBillingIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeTotal, LiveRoomIncomeTotalVideoCallBillingIncome, r.ID, &r.TotalVideoCallBillingIncome, v, true, &r.UpdatedAt)
}
func (r *LiveRoomIncomeTotal) AddTotalLiveDuration(v float64) {
	addIncomeAmount(TbLiveRoomIncomeTotal, LiveRoomIncomeTotalLiveDuration, r.ID, &r.TotalLiveDuration, v, true, &r.UpdatedAt)
}

// AddGiftEarn 礼物收益(总收益+礼物细分,内部加锁)
func (r *LiveRoomIncomeTotal) AddGiftEarn(v float64) {
	addIncomeEarn(TbLiveRoomIncomeTotal, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalGiftIncome, &r.TotalGiftIncome)
}

// AddPaidDanmakuEarn 付费弹幕收益(总收益+弹幕细分,内部加锁)
func (r *LiveRoomIncomeTotal) AddPaidDanmakuEarn(v float64) {
	addIncomeEarn(TbLiveRoomIncomeTotal, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalPaidDanmakuIncome, &r.TotalPaidDanmakuIncome)
}

// AddPrivateRoomTicketEarn 私密房门票收益(总收益+门票细分,内部加锁)
func (r *LiveRoomIncomeTotal) AddPrivateRoomTicketEarn(v float64) {
	addIncomeEarn(TbLiveRoomIncomeTotal, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalPrivateRoomTicketIncome, &r.TotalPrivateRoomTicketIncome)
}

// AddPrivateRoomWatchEarn 私密房观看收益(总收益+观看细分,内部加锁)
func (r *LiveRoomIncomeTotal) AddPrivateRoomWatchEarn(v float64) {
	addIncomeEarn(TbLiveRoomIncomeTotal, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalPrivateRoomWatchIncome, &r.TotalPrivateRoomWatchIncome)
}

// AddShortVideoEarn 短视频付费观看收益(总收益+短视频细分,内部加锁)
func (r *LiveRoomIncomeTotal) AddShortVideoEarn(v float64) {
	addIncomeEarn(TbLiveRoomIncomeTotal, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalShortVideoIncome, &r.TotalShortVideoIncome)
}

// AddSettlementSalary 累加结算薪资
func (r *LiveRoomIncomeTotal) AddSettlementSalary(v float64) {
	if r == nil || v == 0 {
		return
	}
	addIncomeAmount(TbLiveRoomIncomeTotal, LiveRoomIncomeSettlementSalary, r.ID, &r.SettlementSalary, v, false, &r.UpdatedAt)
}

// AddSettlementShareAmount 累加结算分佣金额
func (r *LiveRoomIncomeTotal) AddSettlementShareAmount(v float64) {
	if r == nil || v == 0 {
		return
	}
	addIncomeAmount(TbLiveRoomIncomeTotal, LiveRoomIncomeSettlementShareAmount, r.ID, &r.SettlementShareAmount, v, false, &r.UpdatedAt)
}

func initLiveRoomIncomeTotal() {
	regLiveRoomIncomeCols(TbLiveRoomIncomeTotal)
	syndb.RegQuick(TbLiveRoomIncomeTotal, LiveRoomIncomeSettlementSalary)
	syndb.RegQuick(TbLiveRoomIncomeTotal, LiveRoomIncomeSettlementShareAmount)
	migrate.AutoMigrate(&LiveRoomIncomeTotal{})
}
