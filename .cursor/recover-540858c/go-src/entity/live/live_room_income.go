package entity

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/syndb"
)

const (
	LiveRoomIncomeTotalIncome                  db.TbCol = "total_income"
	LiveRoomIncomeTotalGiftIncome              db.TbCol = "total_gift_income"
	LiveRoomIncomeTotalPaidDanmakuIncome       db.TbCol = "total_paid_danmaku_income"
	LiveRoomIncomeTotalPrivateRoomTicketIncome db.TbCol = "total_private_room_ticket_income"
	LiveRoomIncomeTotalPrivateRoomWatchIncome  db.TbCol = "total_private_room_watch_income"
	LiveRoomIncomeTotalVideoCallIncome         db.TbCol = "total_video_call_income"
	LiveRoomIncomeTotalVideoCallTicketIncome   db.TbCol = "total_video_call_ticket_income"
	LiveRoomIncomeTotalVideoCallBillingIncome  db.TbCol = "total_video_call_billing_income"
	LiveRoomIncomeTotalShortVideoIncome        db.TbCol = "total_short_video_income"
	LiveRoomIncomeTotalLiveDuration            db.TbCol = "total_live_duration"
	LiveRoomIncomeSettlementSalary             db.TbCol = "settlement_salary"
	LiveRoomIncomeSettlementShareAmount        db.TbCol = "settlement_share_amount"
)

// LiveRoomIncomeAmounts 直播间/工会收益字段(房间与工会收益表共用结构)
type LiveRoomIncomeAmounts struct {
	TotalIncome                  float64 `gorm:"default:0;comment:直播收益" json:"totalIncome"`
	TotalGiftIncome              float64 `gorm:"default:0;comment:累计礼物收益" json:"totalGiftIncome"`
	TotalPaidDanmakuIncome       float64 `gorm:"default:0;comment:累计付费弹幕收益" json:"totalPaidDanmakuIncome"`
	TotalPrivateRoomTicketIncome float64 `gorm:"default:0;comment:累计私密直播间门票收益" json:"totalPrivateRoomTicketIncome"`
	TotalPrivateRoomWatchIncome  float64 `gorm:"default:0;comment:累计私密房观看收益" json:"totalPrivateRoomWatchIncome"`
	TotalVideoCallIncome         float64 `gorm:"type:decimal(10,4);default:0;comment:累计直播间视频通话收益" json:"totalVideoCallIncome"`
	TotalVideoCallTicketIncome   float64 `gorm:"type:decimal(10,4);default:0;comment:累计直播间视频通话门票收益" json:"totalVideoCallTicketIncome"`
	TotalVideoCallBillingIncome  float64 `gorm:"type:decimal(10,4);default:0;comment:累计直播间视频通话计费收益" json:"totalVideoCallBillingIncome"`
	TotalShortVideoIncome        float64 `gorm:"type:decimal(10,4);default:0;comment:累计短视频付费观看收益" json:"totalShortVideoIncome"`
	TotalLiveDuration            float64 `gorm:"default:0;comment:累计直播时长(秒)" json:"totalLiveDuration"`
}

func liveRoomIncomeLockKey(tb db.TbName, id uint64) string {
	return fmt.Sprintf("live_room_income:%s:%d", tb, id)
}

func (a *LiveRoomIncomeAmounts) clearAmounts() {
	a.TotalIncome = 0
	a.TotalGiftIncome = 0
	a.TotalPaidDanmakuIncome = 0
	a.TotalPrivateRoomTicketIncome = 0
	a.TotalPrivateRoomWatchIncome = 0
	a.TotalVideoCallIncome = 0
	a.TotalVideoCallTicketIncome = 0
	a.TotalVideoCallBillingIncome = 0
	a.TotalShortVideoIncome = 0
	a.TotalLiveDuration = 0
}

// IsZero 是否全部为0
func (a *LiveRoomIncomeAmounts) IsZero() bool {
	if a == nil {
		return true
	}
	return a.TotalIncome == 0 &&
		a.TotalGiftIncome == 0 &&
		a.TotalPaidDanmakuIncome == 0 &&
		a.TotalPrivateRoomTicketIncome == 0 &&
		a.TotalPrivateRoomWatchIncome == 0 &&
		a.TotalVideoCallIncome == 0 &&
		a.TotalVideoCallTicketIncome == 0 &&
		a.TotalVideoCallBillingIncome == 0 &&
		a.TotalShortVideoIncome == 0 &&
		a.TotalLiveDuration == 0
}

// addIncomeAmountsLocked 在已持锁前提下累加各收益字段
func addIncomeAmountsLocked(tb db.TbName, id any, dst *LiveRoomIncomeAmounts, src *LiveRoomIncomeAmounts) {
	if dst == nil || src == nil || src.IsZero() {
		return
	}
	if src.TotalIncome != 0 {
		addIncomeAmountLocked(tb, LiveRoomIncomeTotalIncome, id, &dst.TotalIncome, src.TotalIncome)
	}
	if src.TotalGiftIncome > 0 {
		addIncomeAmountLocked(tb, LiveRoomIncomeTotalGiftIncome, id, &dst.TotalGiftIncome, src.TotalGiftIncome)
	}
	if src.TotalPaidDanmakuIncome > 0 {
		addIncomeAmountLocked(tb, LiveRoomIncomeTotalPaidDanmakuIncome, id, &dst.TotalPaidDanmakuIncome, src.TotalPaidDanmakuIncome)
	}
	if src.TotalPrivateRoomTicketIncome > 0 {
		addIncomeAmountLocked(tb, LiveRoomIncomeTotalPrivateRoomTicketIncome, id, &dst.TotalPrivateRoomTicketIncome, src.TotalPrivateRoomTicketIncome)
	}
	if src.TotalPrivateRoomWatchIncome > 0 {
		addIncomeAmountLocked(tb, LiveRoomIncomeTotalPrivateRoomWatchIncome, id, &dst.TotalPrivateRoomWatchIncome, src.TotalPrivateRoomWatchIncome)
	}
	if src.TotalVideoCallIncome != 0 {
		addIncomeAmountLocked(tb, LiveRoomIncomeTotalVideoCallIncome, id, &dst.TotalVideoCallIncome, src.TotalVideoCallIncome)
	}
	if src.TotalVideoCallTicketIncome != 0 {
		addIncomeAmountLocked(tb, LiveRoomIncomeTotalVideoCallTicketIncome, id, &dst.TotalVideoCallTicketIncome, src.TotalVideoCallTicketIncome)
	}
	if src.TotalVideoCallBillingIncome != 0 {
		addIncomeAmountLocked(tb, LiveRoomIncomeTotalVideoCallBillingIncome, id, &dst.TotalVideoCallBillingIncome, src.TotalVideoCallBillingIncome)
	}
	if src.TotalShortVideoIncome > 0 {
		addIncomeAmountLocked(tb, LiveRoomIncomeTotalShortVideoIncome, id, &dst.TotalShortVideoIncome, src.TotalShortVideoIncome)
	}
	if src.TotalLiveDuration > 0 {
		addIncomeAmountLocked(tb, LiveRoomIncomeTotalLiveDuration, id, &dst.TotalLiveDuration, src.TotalLiveDuration)
	}
}

// clearIncomeAmountsLocked 在已持锁前提下清零并写库
func clearIncomeAmountsLocked(tb db.TbName, id any, a *LiveRoomIncomeAmounts, updatedAt *time.Time) {
	if a == nil {
		return
	}
	a.clearAmounts()
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalIncome, id, 0)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalGiftIncome, id, 0)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalPaidDanmakuIncome, id, 0)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalPrivateRoomTicketIncome, id, 0)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalPrivateRoomWatchIncome, id, 0)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalVideoCallIncome, id, 0)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalVideoCallTicketIncome, id, 0)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalVideoCallBillingIncome, id, 0)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalShortVideoIncome, id, 0)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalLiveDuration, id, 0)
	touchIncomeUpdatedAt(tb, id, updatedAt)
}

func addIncomeAmount(tb db.TbName, col db.TbCol, id uint64, cur *float64, v float64, skipNonPositive bool, updatedAt *time.Time) {
	if skipNonPositive && v <= 0 {
		return
	}
	key := liveRoomIncomeLockKey(tb, id)
	gmlock.Lock(key)
	defer gmlock.Unlock(key)
	addIncomeAmountLocked(tb, col, id, cur, v)
	touchIncomeUpdatedAt(tb, id, updatedAt)
}

func addIncomeAmountLocked(tb db.TbName, col db.TbCol, id any, cur *float64, v float64) {
	*cur = math.AddFloat64(*cur, v)
	syndb.AddData(tb, col, &syndb.ColData{IdVal: id, ColVal: *cur})
}

func writeIncomeAmountLocked(tb db.TbName, col db.TbCol, id any, v float64) {
	syndb.AddData(tb, col, &syndb.ColData{IdVal: id, ColVal: v})
}

func touchIncomeUpdatedAt(tb db.TbName, id any, updatedAt *time.Time) {
	if updatedAt == nil {
		return
	}
	*updatedAt = time.Now()
	syndb.AddData(tb, db.UpdatedAtName, &syndb.ColData{IdVal: id, ColVal: *updatedAt})
}

// addIncomeEarn 一次加锁:总收益 + 单一细分项(细分项仅 amount>0 时累加)
func addIncomeEarn(tb db.TbName, id uint64, a *LiveRoomIncomeAmounts, updatedAt *time.Time, amount float64, extraCol db.TbCol, extra *float64) {
	if a == nil || id == 0 || extra == nil {
		return
	}
	addIncomeEarnWithLockKey(tb, liveRoomIncomeLockKey(tb, id), id, a, updatedAt, amount, extraCol, extra)
}

// addIncomeEarnWithLockKey 自定义锁键与主键类型的收益累加(用于日表等复合主键)
func addIncomeEarnWithLockKey(tb db.TbName, lockKey string, id any, a *LiveRoomIncomeAmounts, updatedAt *time.Time, amount float64, extraCol db.TbCol, extra *float64) {
	if a == nil || id == nil || extra == nil || lockKey == "" {
		return
	}
	gmlock.Lock(lockKey)
	defer gmlock.Unlock(lockKey)
	addIncomeAmountLocked(tb, LiveRoomIncomeTotalIncome, id, &a.TotalIncome, amount)
	if amount > 0 {
		addIncomeAmountLocked(tb, extraCol, id, extra, amount)
	}
	touchIncomeUpdatedAt(tb, id, updatedAt)
}

// ApplyVideoCallIncomeDelta 通话收益增减(支持负数退款),内部按表+房间加锁
func ApplyVideoCallIncomeDelta(tb db.TbName, id uint64, a *LiveRoomIncomeAmounts, updatedAt *time.Time, amount float64, ticket, billing bool) {
	if a == nil || id == 0 {
		return
	}
	ApplyVideoCallIncomeDeltaWithLockKey(tb, liveRoomIncomeLockKey(tb, id), id, a, updatedAt, amount, ticket, billing)
}

// ApplyVideoCallIncomeDeltaWithLockKey 自定义锁键与主键类型的通话收益增减
func ApplyVideoCallIncomeDeltaWithLockKey(tb db.TbName, lockKey string, id any, a *LiveRoomIncomeAmounts, updatedAt *time.Time, amount float64, ticket, billing bool) {
	if a == nil || id == nil || lockKey == "" {
		return
	}
	gmlock.Lock(lockKey)
	defer gmlock.Unlock(lockKey)
	addIncomeAmountLocked(tb, LiveRoomIncomeTotalIncome, id, &a.TotalIncome, amount)
	addIncomeAmountLocked(tb, LiveRoomIncomeTotalVideoCallIncome, id, &a.TotalVideoCallIncome, amount)
	if ticket {
		addIncomeAmountLocked(tb, LiveRoomIncomeTotalVideoCallTicketIncome, id, &a.TotalVideoCallTicketIncome, amount)
	}
	if billing {
		addIncomeAmountLocked(tb, LiveRoomIncomeTotalVideoCallBillingIncome, id, &a.TotalVideoCallBillingIncome, amount)
	}
	touchIncomeUpdatedAt(tb, id, updatedAt)
}

func regLiveRoomIncomeCols(tb db.TbName) {
	syndb.RegQuick(tb, db.CreatedAtName)
	syndb.RegQuick(tb, db.UpdatedAtName)
	syndb.RegQuick(tb, LiveRoomIncomeTotalIncome)
	syndb.RegQuick(tb, LiveRoomIncomeTotalGiftIncome)
	syndb.RegQuick(tb, LiveRoomIncomeTotalPaidDanmakuIncome)
	syndb.RegQuick(tb, LiveRoomIncomeTotalPrivateRoomTicketIncome)
	syndb.RegQuick(tb, LiveRoomIncomeTotalPrivateRoomWatchIncome)
	syndb.RegQuick(tb, LiveRoomIncomeTotalVideoCallIncome)
	syndb.RegQuick(tb, LiveRoomIncomeTotalVideoCallTicketIncome)
	syndb.RegQuick(tb, LiveRoomIncomeTotalVideoCallBillingIncome)
	syndb.RegQuick(tb, LiveRoomIncomeTotalShortVideoIncome)
	syndb.RegQuick(tb, LiveRoomIncomeTotalLiveDuration)
}

func initLiveRoomIncome() {
	initLiveRoomIncomeUnsettled()
	initLiveRoomIncomeUnsettledArchive()
	initLiveRoomIncomeSettled()
	initLiveRoomIncomeTotal()
	initGuildIncomeUnsettled()
	initGuildIncomeUnsettledArchive()
	initGuildIncomeSettled()
	initGuildIncomeTotal()
}
