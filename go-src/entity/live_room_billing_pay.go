package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbLiveRoomBillingPay db.TbName = "live_room_billing_pays"
	// LiveRoomBillingFreeDuration 私密直播间观众免费观看时长
	LiveRoomBillingFreeDuration = 7 * time.Minute
)

const (
	LiveRoomBillingPayUserId       db.TbCol = "user_id"
	LiveRoomBillingPayRoomId       db.TbCol = "room_id"
	LiveRoomBillingPayLiveRecordId db.TbCol = "live_record_id"
	LiveRoomBillingPayLastBilledAt db.TbCol = "last_billed_at"
)

// LiveRoomBillingPay 私密直播间按分钟计费记录(主键 ID = userId_roomId_liveRecordId)
type LiveRoomBillingPay struct {
	ID           string    `gorm:"primaryKey;size:96;comment:复合ID(userId_roomId_liveRecordId)" json:"id"`
	UserId       uint64    `gorm:"index;default:0;comment:用户ID" json:"userId"`
	RoomId       uint64    `gorm:"index;default:0;comment:直播间ID" json:"roomId"`
	LiveRecordId uint64    `gorm:"index;default:0;comment:直播记录ID" json:"liveRecordId"`
	LastBilledAt time.Time `gorm:"comment:最近一次按分钟扣费时间" json:"lastBilledAt"`
}

func BuildLiveRoomBillingPayId(userId, roomId, liveRecordId uint64) string {
	return fmt.Sprintf("%d_%d_%d", userId, roomId, liveRecordId)
}

func NewLiveRoomBillingPay(userId, roomId, liveRecordId uint64) *LiveRoomBillingPay {
	r := &LiveRoomBillingPay{}
	r.ID = BuildLiveRoomBillingPayId(userId, roomId, liveRecordId)
	r.SetUserId(userId)
	r.SetRoomId(roomId)
	r.SetLiveRecordId(liveRecordId)
	return r
}

func (r *LiveRoomBillingPay) ShouldChargeMinute(joinTime *time.Time, now time.Time) bool {
	if joinTime == nil || joinTime.IsZero() {
		return false
	}
	freeUntil := joinTime.Add(LiveRoomBillingFreeDuration)
	if now.Before(freeUntil) {
		return false
	}
	anchor := r.LastBilledAt
	if anchor.IsZero() || anchor.Before(freeUntil) {
		anchor = freeUntil
	}
	return now.Sub(anchor) >= time.Minute
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

func (r *LiveRoomBillingPay) SetLiveRecordId(v uint64) {
	r.LiveRecordId = v
	syndb.AddDataToLazyChan(TbLiveRoomBillingPay, LiveRoomBillingPayLiveRecordId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoomBillingPay) SetLastBilledAt(v time.Time) {
	r.LastBilledAt = v
	syndb.AddDataToQuickChan(TbLiveRoomBillingPay, LiveRoomBillingPayLastBilledAt, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func initLiveRoomBillingPay() {
	syndb.RegLazyWithMiddle(TbLiveRoomBillingPay, LiveRoomBillingPayUserId)
	syndb.RegLazyWithMiddle(TbLiveRoomBillingPay, LiveRoomBillingPayRoomId)
	syndb.RegLazyWithMiddle(TbLiveRoomBillingPay, LiveRoomBillingPayLiveRecordId)
	syndb.RegQuickWithMiddle(TbLiveRoomBillingPay, LiveRoomBillingPayLastBilledAt)
	migrate.AutoMigrate(&LiveRoomBillingPay{})
}
