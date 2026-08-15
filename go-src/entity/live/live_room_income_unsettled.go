package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"

	"github.com/gogf/gf/v2/os/gmlock"
)

const (
	TbLiveRoomIncomeUnsettled db.TbName = "live_room_income_unsettleds"
)

// LiveRoomIncomeUnsettled 直播间未结算收益(结算时清零)
type LiveRoomIncomeUnsettled struct {
	migrate.OneModel
	LiveRoomIncomeAmounts
}

func NewLiveRoomIncomeUnsettled(roomId uint64) *LiveRoomIncomeUnsettled {
	ret := &LiveRoomIncomeUnsettled{}
	ret.ID = roomId
	now := time.Now()
	ret.CreatedAt = now
	ret.UpdatedAt = now
	syndb.AddData(TbLiveRoomIncomeUnsettled, db.CreatedAtName, &syndb.ColData{IdVal: roomId, ColVal: now})
	syndb.AddData(TbLiveRoomIncomeUnsettled, db.UpdatedAtName, &syndb.ColData{IdVal: roomId, ColVal: now})
	return ret
}

func (r *LiveRoomIncomeUnsettled) AddTotalIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeUnsettled, LiveRoomIncomeTotalIncome, r.ID, &r.TotalIncome, v, false, &r.UpdatedAt)
}
func (r *LiveRoomIncomeUnsettled) AddTotalGiftIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeUnsettled, LiveRoomIncomeTotalGiftIncome, r.ID, &r.TotalGiftIncome, v, true, &r.UpdatedAt)
}
func (r *LiveRoomIncomeUnsettled) AddTotalPaidDanmakuIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeUnsettled, LiveRoomIncomeTotalPaidDanmakuIncome, r.ID, &r.TotalPaidDanmakuIncome, v, true, &r.UpdatedAt)
}
func (r *LiveRoomIncomeUnsettled) AddTotalPrivateRoomTicketIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeUnsettled, LiveRoomIncomeTotalPrivateRoomTicketIncome, r.ID, &r.TotalPrivateRoomTicketIncome, v, true, &r.UpdatedAt)
}
func (r *LiveRoomIncomeUnsettled) AddTotalPrivateRoomWatchIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeUnsettled, LiveRoomIncomeTotalPrivateRoomWatchIncome, r.ID, &r.TotalPrivateRoomWatchIncome, v, true, &r.UpdatedAt)
}
func (r *LiveRoomIncomeUnsettled) AddTotalVideoCallIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeUnsettled, LiveRoomIncomeTotalVideoCallIncome, r.ID, &r.TotalVideoCallIncome, v, true, &r.UpdatedAt)
}
func (r *LiveRoomIncomeUnsettled) AddTotalVideoCallTicketIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeUnsettled, LiveRoomIncomeTotalVideoCallTicketIncome, r.ID, &r.TotalVideoCallTicketIncome, v, true, &r.UpdatedAt)
}
func (r *LiveRoomIncomeUnsettled) AddTotalVideoCallBillingIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeUnsettled, LiveRoomIncomeTotalVideoCallBillingIncome, r.ID, &r.TotalVideoCallBillingIncome, v, true, &r.UpdatedAt)
}

// AddGiftEarn 礼物收益(总收益+礼物细分,内部加锁)
func (r *LiveRoomIncomeUnsettled) AddGiftEarn(v float64) {
	addIncomeEarn(TbLiveRoomIncomeUnsettled, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalGiftIncome, &r.TotalGiftIncome)
}

// AddPaidDanmakuEarn 付费弹幕收益(总收益+弹幕细分,内部加锁)
func (r *LiveRoomIncomeUnsettled) AddPaidDanmakuEarn(v float64) {
	addIncomeEarn(TbLiveRoomIncomeUnsettled, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalPaidDanmakuIncome, &r.TotalPaidDanmakuIncome)
}

// AddPrivateRoomTicketEarn 私密房门票收益(总收益+门票细分,内部加锁)
func (r *LiveRoomIncomeUnsettled) AddPrivateRoomTicketEarn(v float64) {
	addIncomeEarn(TbLiveRoomIncomeUnsettled, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalPrivateRoomTicketIncome, &r.TotalPrivateRoomTicketIncome)
}

// AddPrivateRoomWatchEarn 私密房观看收益(总收益+观看细分,内部加锁)
func (r *LiveRoomIncomeUnsettled) AddPrivateRoomWatchEarn(v float64) {
	addIncomeEarn(TbLiveRoomIncomeUnsettled, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalPrivateRoomWatchIncome, &r.TotalPrivateRoomWatchIncome)
}

// AddAmounts 累加一笔金额到未结算表(一次加锁)
func (r *LiveRoomIncomeUnsettled) AddAmounts(a *LiveRoomIncomeAmounts) {
	if a == nil || a.IsZero() {
		return
	}
	key := liveRoomIncomeLockKey(TbLiveRoomIncomeUnsettled, r.ID)
	gmlock.Lock(key)
	defer gmlock.Unlock(key)
	addIncomeAmountsLocked(TbLiveRoomIncomeUnsettled, r.ID, &r.LiveRoomIncomeAmounts, a)
	touchIncomeUpdatedAt(TbLiveRoomIncomeUnsettled, r.ID, &r.UpdatedAt)
}

// SnapshotAndClear 取出当前未结算金额并清零(一次加锁,用于下架归档)
func (r *LiveRoomIncomeUnsettled) SnapshotAndClear() LiveRoomIncomeAmounts {
	key := liveRoomIncomeLockKey(TbLiveRoomIncomeUnsettled, r.ID)
	gmlock.Lock(key)
	defer gmlock.Unlock(key)
	snap := r.LiveRoomIncomeAmounts
	clearIncomeAmountsLocked(TbLiveRoomIncomeUnsettled, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt)
	return snap
}

// Clear 结算后清零未结算收益
func (r *LiveRoomIncomeUnsettled) Clear() {
	key := liveRoomIncomeLockKey(TbLiveRoomIncomeUnsettled, r.ID)
	gmlock.Lock(key)
	defer gmlock.Unlock(key)
	clearIncomeAmountsLocked(TbLiveRoomIncomeUnsettled, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt)
}

func initLiveRoomIncomeUnsettled() {
	regLiveRoomIncomeCols(TbLiveRoomIncomeUnsettled)
	migrate.AutoMigrate(&LiveRoomIncomeUnsettled{})
}
