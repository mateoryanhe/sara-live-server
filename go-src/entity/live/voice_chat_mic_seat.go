package entity

import (
	"fmt"
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const TbVoiceChatMicSeat db.TbName = "voice_chat_mic_seats"

const (
	VoiceChatMicSeatRoomId       db.TbCol = "room_id"
	VoiceChatMicSeatIndexCol     db.TbCol = "seat_index"
	VoiceChatMicSeatUserId       db.TbCol = "user_id"
	VoiceChatMicSeatLiveRecordId db.TbCol = "live_record_id"
	VoiceChatMicSeatStatusCol    db.TbCol = "status"
)

// VoiceChatMicSeat 语聊房麦位(主键 roomId_seatIndex,按场次 live_record_id 区分)
type VoiceChatMicSeat struct {
	ID           string    `gorm:"primaryKey;size:64;comment:复合ID(roomId_seatIndex)" json:"id"`
	RoomId       uint64    `gorm:"index:idx_voice_chat_mic_room_live,priority:1;default:0;comment:直播间ID" json:"roomId"`
	SeatIndex    uint8     `gorm:"index:idx_voice_chat_mic_room_live,priority:3;default:0;comment:麦位序号(0主播,1~8观众)" json:"seatIndex"`
	UserId       uint64    `gorm:"default:0;comment:麦上用户ID(0表示空位)" json:"userId"`
	LiveRecordId uint64    `gorm:"index:idx_voice_chat_mic_room_live,priority:2;default:0;comment:本场直播记录ID" json:"liveRecordId"`
	Status       uint8     `gorm:"default:1;comment:麦位状态(1空闲,2锁麦,3在麦,4在麦禁音)" json:"status"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

func BuildVoiceChatMicSeatId(roomId uint64, seatIndex int) string {
	return fmt.Sprintf("%d_s%d", roomId, seatIndex)
}

func NewVoiceChatMicSeat(roomId uint64, seatIndex int, liveRecordId, userId uint64) *VoiceChatMicSeat {
	row := &VoiceChatMicSeat{}
	row.ID = BuildVoiceChatMicSeatId(roomId, seatIndex)
	now := time.Now()
	row.CreatedAt = now
	row.UpdatedAt = now
	row.SetRoomId(roomId)
	row.SetSeatIndex(uint8(seatIndex))
	row.SetLiveRecordId(liveRecordId)
	row.SetUserId(userId)
	if userId > 0 {
		row.SetStatus(VoiceChatMicSeatStatusOnMic)
	} else {
		row.SetStatus(VoiceChatMicSeatStatusIdle)
	}
	row.SetCreatedAt(now)
	row.SetUpdatedAt(now)
	return row
}

func (r *VoiceChatMicSeat) SetRoomId(v uint64) {
	r.RoomId = v
	r.touchUpdatedAt()
	syndb.AddData(TbVoiceChatMicSeat, VoiceChatMicSeatRoomId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *VoiceChatMicSeat) SetSeatIndex(v uint8) {
	r.SeatIndex = v
	r.touchUpdatedAt()
	syndb.AddData(TbVoiceChatMicSeat, VoiceChatMicSeatIndexCol, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *VoiceChatMicSeat) SetUserId(v uint64) {
	r.UserId = v
	r.touchUpdatedAt()
	syndb.AddData(TbVoiceChatMicSeat, VoiceChatMicSeatUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *VoiceChatMicSeat) SetLiveRecordId(v uint64) {
	r.LiveRecordId = v
	r.touchUpdatedAt()
	syndb.AddData(TbVoiceChatMicSeat, VoiceChatMicSeatLiveRecordId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *VoiceChatMicSeat) SetStatus(v uint8) {
	if !IsValidVoiceChatMicSeatStatus(v) {
		v = VoiceChatMicSeatStatusIdle
	}
	r.Status = v
	r.touchUpdatedAt()
	syndb.AddData(TbVoiceChatMicSeat, VoiceChatMicSeatStatusCol, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *VoiceChatMicSeat) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbVoiceChatMicSeat, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *VoiceChatMicSeat) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddData(TbVoiceChatMicSeat, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *VoiceChatMicSeat) touchUpdatedAt() {
	r.UpdatedAt = time.Now()
	syndb.AddData(TbVoiceChatMicSeat, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: r.UpdatedAt})
}

func initVoiceChatMicSeat() {
	migrate.AutoMigrate(&VoiceChatMicSeat{})
	syndb.RegQuick(TbVoiceChatMicSeat, db.CreatedAtName)
	syndb.RegQuick(TbVoiceChatMicSeat, db.UpdatedAtName)
	syndb.RegQuick(TbVoiceChatMicSeat, VoiceChatMicSeatRoomId)
	syndb.RegQuick(TbVoiceChatMicSeat, VoiceChatMicSeatLiveRecordId)
	syndb.RegQuick(TbVoiceChatMicSeat, VoiceChatMicSeatIndexCol)
	syndb.RegLazy(TbVoiceChatMicSeat, VoiceChatMicSeatUserId)
	syndb.RegLazy(TbVoiceChatMicSeat, VoiceChatMicSeatStatusCol)
}
