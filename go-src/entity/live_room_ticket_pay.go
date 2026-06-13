package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbLiveRoomTicketPay         db.TbName = "live_room_ticket_pays"
	LiveRoomTicketValidDuration           = 24 * time.Hour
)

const (
	LiveRoomTicketPayUserId     db.TbCol = "user_id"
	LiveRoomTicketPayRoomId     db.TbCol = "room_id"
	LiveRoomTicketPayLastPaidAt db.TbCol = "last_paid_at"
)

// LiveRoomTicketPay 私密直播间门票扣费记录(主键 ID = userId_roomId)
type LiveRoomTicketPay struct {
	ID         string    `gorm:"primaryKey;size:64;comment:复合ID(userId_roomId)" json:"id"`
	UserId     uint64    `gorm:"index;default:0;comment:用户ID" json:"userId"`
	RoomId     uint64    `gorm:"index;default:0;comment:直播间ID" json:"roomId"`
	LastPaidAt time.Time `gorm:"comment:最近一次扣门票时间" json:"lastPaidAt"`
}

func BuildLiveRoomTicketPayId(userId, roomId uint64) string {
	return fmt.Sprintf("%d_%d", userId, roomId)
}

func NewLiveRoomTicketPay(userId, roomId uint64) *LiveRoomTicketPay {
	r := &LiveRoomTicketPay{}
	r.ID = BuildLiveRoomTicketPayId(userId, roomId)
	r.SetUserId(userId)
	r.SetRoomId(roomId)
	return r
}

func (r *LiveRoomTicketPay) IsValidWithin24h(now time.Time) bool {
	if r.LastPaidAt.IsZero() {
		return false
	}
	return now.Sub(r.LastPaidAt) < LiveRoomTicketValidDuration
}

func (r *LiveRoomTicketPay) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddDataToLazyChan(TbLiveRoomTicketPay, LiveRoomTicketPayUserId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoomTicketPay) SetRoomId(v uint64) {
	r.RoomId = v
	syndb.AddDataToLazyChan(TbLiveRoomTicketPay, LiveRoomTicketPayRoomId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoomTicketPay) SetLastPaidAt(v time.Time) {
	r.LastPaidAt = v
	syndb.AddDataToQuickChan(TbLiveRoomTicketPay, LiveRoomTicketPayLastPaidAt, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func initLiveRoomTicketPay() {
	syndb.RegLazyWithMiddle(TbLiveRoomTicketPay, LiveRoomTicketPayUserId)
	syndb.RegLazyWithMiddle(TbLiveRoomTicketPay, LiveRoomTicketPayRoomId)
	syndb.RegQuickWithMiddle(TbLiveRoomTicketPay, LiveRoomTicketPayLastPaidAt)
	migrate.AutoMigrate(&LiveRoomTicketPay{})
}
