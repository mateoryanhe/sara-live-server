package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"

	"github.com/gogf/gf/v2/os/gmlock"
)

const (
	TbLiveRoomIncomeSettled db.TbName = "live_room_income_settleds"
)

// LiveRoomIncomeSettled 直播间已结算收益(每次结算累加)
type LiveRoomIncomeSettled struct {
	migrate.OneModel
	LiveRoomIncomeAmounts
	SettlementSalary float64 `gorm:"type:decimal(16,4);default:0;comment:结算薪资" json:"settlementSalary"`
}

func NewLiveRoomIncomeSettled(roomId uint64) *LiveRoomIncomeSettled {
	ret := &LiveRoomIncomeSettled{}
	ret.ID = roomId
	now := time.Now()
	ret.CreatedAt = now
	ret.UpdatedAt = now
	syndb.AddData(TbLiveRoomIncomeSettled, db.CreatedAtName, &syndb.ColData{IdVal: roomId, ColVal: now})
	syndb.AddData(TbLiveRoomIncomeSettled, db.UpdatedAtName, &syndb.ColData{IdVal: roomId, ColVal: now})
	return ret
}

func (r *LiveRoomIncomeSettled) AddTotalIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeSettled, LiveRoomIncomeTotalIncome, r.ID, &r.TotalIncome, v, true, &r.UpdatedAt)
}
func (r *LiveRoomIncomeSettled) AddTotalGiftIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeSettled, LiveRoomIncomeTotalGiftIncome, r.ID, &r.TotalGiftIncome, v, true, &r.UpdatedAt)
}
func (r *LiveRoomIncomeSettled) AddTotalPaidDanmakuIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeSettled, LiveRoomIncomeTotalPaidDanmakuIncome, r.ID, &r.TotalPaidDanmakuIncome, v, true, &r.UpdatedAt)
}
func (r *LiveRoomIncomeSettled) AddTotalPrivateRoomTicketIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeSettled, LiveRoomIncomeTotalPrivateRoomTicketIncome, r.ID, &r.TotalPrivateRoomTicketIncome, v, true, &r.UpdatedAt)
}
func (r *LiveRoomIncomeSettled) AddTotalPrivateRoomWatchIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeSettled, LiveRoomIncomeTotalPrivateRoomWatchIncome, r.ID, &r.TotalPrivateRoomWatchIncome, v, true, &r.UpdatedAt)
}
func (r *LiveRoomIncomeSettled) AddTotalVideoCallIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeSettled, LiveRoomIncomeTotalVideoCallIncome, r.ID, &r.TotalVideoCallIncome, v, true, &r.UpdatedAt)
}
func (r *LiveRoomIncomeSettled) AddTotalVideoCallTicketIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeSettled, LiveRoomIncomeTotalVideoCallTicketIncome, r.ID, &r.TotalVideoCallTicketIncome, v, true, &r.UpdatedAt)
}
func (r *LiveRoomIncomeSettled) AddTotalVideoCallBillingIncome(v float64) {
	addIncomeAmount(TbLiveRoomIncomeSettled, LiveRoomIncomeTotalVideoCallBillingIncome, r.ID, &r.TotalVideoCallBillingIncome, v, true, &r.UpdatedAt)
}
func (r *LiveRoomIncomeSettled) AddTotalLiveDuration(v float64) {
	addIncomeAmount(TbLiveRoomIncomeSettled, LiveRoomIncomeTotalLiveDuration, r.ID, &r.TotalLiveDuration, v, true, &r.UpdatedAt)
}

// AddAmounts 结算时将一笔金额累加到已结算表(一次加锁)
func (r *LiveRoomIncomeSettled) AddAmounts(a *LiveRoomIncomeAmounts) {
	if a == nil || a.IsZero() {
		return
	}
	key := liveRoomIncomeLockKey(TbLiveRoomIncomeSettled, r.ID)
	gmlock.Lock(key)
	defer gmlock.Unlock(key)
	addIncomeAmountsLocked(TbLiveRoomIncomeSettled, r.ID, &r.LiveRoomIncomeAmounts, a)
	touchIncomeUpdatedAt(TbLiveRoomIncomeSettled, r.ID, &r.UpdatedAt)
}

// AddSettlementSalary 累加结算薪资
func (r *LiveRoomIncomeSettled) AddSettlementSalary(v float64) {
	if r == nil || v == 0 {
		return
	}
	addIncomeAmount(TbLiveRoomIncomeSettled, LiveRoomIncomeSettlementSalary, r.ID, &r.SettlementSalary, v, false, &r.UpdatedAt)
}

func initLiveRoomIncomeSettled() {
	regLiveRoomIncomeCols(TbLiveRoomIncomeSettled)
	syndb.RegQuick(TbLiveRoomIncomeSettled, LiveRoomIncomeSettlementSalary)
	migrate.AutoMigrate(&LiveRoomIncomeSettled{})
}
