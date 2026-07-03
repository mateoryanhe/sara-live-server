package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
	"xr-game-server/core/xrtime"
)

const (
	TbLiveRoomBillingPay        db.TbName = "live_room_billing_pays"
	LiveRoomTicketValidDuration           = 24 * time.Hour
)

const (
	LiveRoomBillingPayUserId         db.TbCol = "user_id"
	LiveRoomBillingPayRoomId         db.TbCol = "room_id"
	LiveRoomBillingPayFreeTime       db.TbCol = "free_time"
	LiveRoomBillingPayLastLastFreeAt db.TbCol = "last_free_at"
	LiveRoomBillingPayLastPaidAt     db.TbCol = "last_paid_at"
	LiveRoomBillingPayBillTime       db.TbCol = "bill_time"
)

// LiveRoomBillingPay 私密直播间按分钟计费记录(主键 ID = userId_roomId_liveRecordId)
type LiveRoomBillingPay struct {
	ID         string     `gorm:"primaryKey;size:96;comment:复合ID(userId_roomId)" json:"id"`
	UserId     uint64     `gorm:"index;default:0;comment:用户ID" json:"userId"`
	RoomId     uint64     `gorm:"index;default:0;comment:直播间ID" json:"roomId"`
	FreeTime   uint64     `gorm:"index;default:0;comment:免费时长" json:"free_time"`
	BillTime   uint64     `gorm:"index;default:0;comment:计费时长" json:"bill_time"`
	LastFreeAt *time.Time `gorm:"comment:最近一次按分钟扣费时间" json:"lastBilledAt"`
	LastPaidAt *time.Time `gorm:"comment:最近一次扣门票时间" json:"lastPaidAt"`
}

func BuildLiveRoomBillingPayId(userId, roomId uint64) string {
	return fmt.Sprintf("%d_%d", userId, roomId)
}

func NewLiveRoomBillingPay(userId, roomId uint64) *LiveRoomBillingPay {
	r := &LiveRoomBillingPay{}
	r.ID = BuildLiveRoomBillingPayId(userId, roomId)
	r.SetUserId(userId)
	r.SetRoomId(roomId)
	return r
}

func (r *LiveRoomBillingPay) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddDataToLazyChan(TbLiveRoomBillingPay, LiveRoomBillingPayUserId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoomBillingPay) SetRoomId(v uint64) {
	r.RoomId = v
	syndb.AddDataToLazyChan(TbLiveRoomBillingPay, LiveRoomBillingPayRoomId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoomBillingPay) SubFreeTime(v uint64) {
	r.FreeTime -= v
	if r.FreeTime < 0 {
		r.FreeTime = 0
	}
	r.SetFreeTime(r.FreeTime)
}

func (r *LiveRoomBillingPay) GetFreeTime(freeTimeCfg uint64) uint64 {
	now := time.Now()
	if r.LastFreeAt == nil || (!xrtime.IsSameDay(now, *r.LastFreeAt)) {
		r.SetLastFreeAt(now)
		r.SetFreeTime(freeTimeCfg)
	}
	return r.FreeTime
}

func (r *LiveRoomBillingPay) IsValidWithin24h(now time.Time) bool {
	if r.LastPaidAt == nil {
		return false
	}
	return now.Sub(*r.LastPaidAt) < LiveRoomTicketValidDuration
}

func (r *LiveRoomBillingPay) SetFreeTime(v uint64) {
	r.FreeTime = v
	syndb.AddDataToLazyChan(TbLiveRoomBillingPay, LiveRoomBillingPayFreeTime, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoomBillingPay) SetBillTime(v uint64) {
	r.FreeTime = v
	syndb.AddDataToLazyChan(TbLiveRoomBillingPay, LiveRoomBillingPayBillTime, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoomBillingPay) SubBillTime(v uint64) {
	r.BillTime -= v
	r.SetBillTime(r.BillTime)
}

func (r *LiveRoomBillingPay) AddBillTime(v uint64) {
	r.BillTime += v
	r.SetBillTime(r.BillTime)
}

func (r *LiveRoomBillingPay) SetLastFreeAt(v time.Time) {
	r.LastFreeAt = &v
	syndb.AddDataToQuickChan(TbLiveRoomBillingPay, LiveRoomBillingPayLastLastFreeAt, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoomBillingPay) SetLastPaidAt(v time.Time) {
	r.LastPaidAt = &v
	syndb.AddDataToQuickChan(TbLiveRoomBillingPay, LiveRoomBillingPayLastPaidAt, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func initLiveRoomBillingPay() {
	syndb.RegLazyWithMiddle(TbLiveRoomBillingPay, LiveRoomBillingPayUserId)
	syndb.RegLazyWithMiddle(TbLiveRoomBillingPay, LiveRoomBillingPayRoomId)
	syndb.RegLazyWithMiddle(TbLiveRoomBillingPay, LiveRoomBillingPayFreeTime)
	syndb.RegQuickWithMiddle(TbLiveRoomBillingPay, LiveRoomBillingPayLastLastFreeAt)
	syndb.RegQuickWithMiddle(TbLiveRoomBillingPay, LiveRoomBillingPayLastPaidAt)

	syndb.RegQuickWithMiddle(TbLiveRoomBillingPay, LiveRoomBillingPayBillTime)

	migrate.AutoMigrate(&LiveRoomBillingPay{})
}
