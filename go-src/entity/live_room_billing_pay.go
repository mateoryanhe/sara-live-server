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
	TbLiveRoomBillingPay        db.TbName = "live_room_billing_pays"
	LiveRoomTicketValidDuration           = 24 * time.Hour
)

const (
	LiveRoomBillingPayUserId       db.TbCol = "user_id"
	LiveRoomBillingPayRoomId       db.TbCol = "room_id"
	LiveRoomBillingPayFreeTime     db.TbCol = "free_time"
	LiveRoomBillingPayLastPaidAt   db.TbCol = "last_paid_at"
	LiveRoomBillingPayLastTicketAt db.TbCol = "last_ticket_at"
)

// LiveRoomBillingPay 私密直播间计费记录(主键 ID = userId_roomId)
type LiveRoomBillingPay struct {
	ID           string     `gorm:"primaryKey;size:96;comment:复合ID(userId_roomId)" json:"id"`
	UserId       uint64     `gorm:"index;default:0;comment:用户ID" json:"userId"`
	RoomId       uint64     `gorm:"index;default:0;comment:直播间ID" json:"roomId"`
	FreeTime     uint64     `gorm:"index;default:0;comment:免费时长" json:"freeTime"`
	LastPaidAt   *time.Time `gorm:"comment:最近一次观看时间" json:"lastPaidAt"`
	LastTicketAt *time.Time `gorm:"comment:最近一次扣门票时间" json:"lastTicketAt"`
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
	syndb.AddDataToQuickChan(TbLiveRoomBillingPay, LiveRoomBillingPayUserId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoomBillingPay) SetRoomId(v uint64) {
	r.RoomId = v
	syndb.AddDataToQuickChan(TbLiveRoomBillingPay, LiveRoomBillingPayRoomId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoomBillingPay) SubFreeTime(v uint64) {
	r.FreeTime = math.Sub(r.FreeTime, v)
	r.SetFreeTime(r.FreeTime)
}

func (r *LiveRoomBillingPay) IsValidWithin24h() bool {
	if r.LastTicketAt == nil {
		return false
	}
	now := time.Now()
	return now.Sub(*r.LastTicketAt) < LiveRoomTicketValidDuration
}

func (r *LiveRoomBillingPay) GetTicketTime() int64 {
	if r.LastTicketAt == nil {
		return 0
	}
	expireTime := r.LastTicketAt.Add(LiveRoomTicketValidDuration)
	sec := time.Since(expireTime).Seconds()
	if sec > 0 {
		return 0
	}
	return int64(-sec)
}

func (r *LiveRoomBillingPay) SetFreeTime(v uint64) {
	r.FreeTime = v
	syndb.AddDataToQuickChan(TbLiveRoomBillingPay, LiveRoomBillingPayFreeTime, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoomBillingPay) SetLastTicketAt(v time.Time) {
	r.LastTicketAt = &v
	syndb.AddDataToQuickChan(TbLiveRoomBillingPay, LiveRoomBillingPayLastTicketAt, &syndb.ColData{
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
	syndb.RegQuickWithLarge(TbLiveRoomBillingPay, LiveRoomBillingPayUserId)
	syndb.RegQuickWithLarge(TbLiveRoomBillingPay, LiveRoomBillingPayRoomId)
	syndb.RegQuickWithLarge(TbLiveRoomBillingPay, LiveRoomBillingPayFreeTime)
	syndb.RegQuickWithLarge(TbLiveRoomBillingPay, LiveRoomBillingPayLastPaidAt)
	syndb.RegQuickWithLarge(TbLiveRoomBillingPay, LiveRoomBillingPayLastTicketAt)
	migrate.AutoMigrate(&LiveRoomBillingPay{})
}
