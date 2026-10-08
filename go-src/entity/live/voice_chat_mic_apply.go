package entity

import (
	"fmt"
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const TbVoiceChatMicApply db.TbName = "voice_chat_mic_applies"

const (
	VoiceChatMicApplyRoomId       db.TbCol = "room_id"
	VoiceChatMicApplyUserId       db.TbCol = "user_id"
	VoiceChatMicApplyLiveRecordId db.TbCol = "live_record_id"
	VoiceChatMicApplyStatus       db.TbCol = "status"
	VoiceChatMicApplyAppliedAt    db.TbCol = "applied_at"
)

const (
	VoiceChatMicApplyStatusCancelled uint8 = 0
	VoiceChatMicApplyStatusActive    uint8 = 1
)

// VoiceChatMicApply 语聊房上麦申请(主键 roomId_userId)
type VoiceChatMicApply struct {
	ID           string     `gorm:"primaryKey;size:64;comment:复合ID(roomId_userId)" json:"id"`
	RoomId       uint64     `gorm:"index:idx_voice_chat_mic_apply_room_live,priority:1;default:0;comment:直播间ID" json:"roomId"`
	UserId       uint64     `gorm:"index;default:0;comment:申请用户ID" json:"userId"`
	LiveRecordId uint64     `gorm:"index:idx_voice_chat_mic_apply_room_live,priority:2;default:0;comment:本场直播记录ID" json:"liveRecordId"`
	Status       uint8      `gorm:"index:idx_voice_chat_mic_apply_room_live,priority:3;default:1;comment:状态(0取消,1有效)" json:"status"`
	AppliedAt    *time.Time `gorm:"comment:申请时间" json:"appliedAt"`
	CreatedAt    time.Time  `json:"-"`
	UpdatedAt    time.Time  `json:"-"`
}

func BuildVoiceChatMicApplyId(roomId, userId uint64) string {
	return fmt.Sprintf("%d_a%d", roomId, userId)
}

func NewVoiceChatMicApply(roomId, userId, liveRecordId uint64, appliedAt time.Time) *VoiceChatMicApply {
	row := &VoiceChatMicApply{}
	row.ID = BuildVoiceChatMicApplyId(roomId, userId)
	now := time.Now()
	row.CreatedAt = now
	row.UpdatedAt = now
	row.SetRoomId(roomId)
	row.SetUserId(userId)
	row.SetLiveRecordId(liveRecordId)
	row.SetStatus(VoiceChatMicApplyStatusActive)
	row.SetAppliedAt(&appliedAt)
	row.SetCreatedAt(now)
	row.SetUpdatedAt(now)
	return row
}

func (r *VoiceChatMicApply) SetRoomId(v uint64) {
	r.RoomId = v
	r.touchUpdatedAt()
	syndb.AddData(TbVoiceChatMicApply, VoiceChatMicApplyRoomId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *VoiceChatMicApply) SetUserId(v uint64) {
	r.UserId = v
	r.touchUpdatedAt()
	syndb.AddData(TbVoiceChatMicApply, VoiceChatMicApplyUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *VoiceChatMicApply) SetLiveRecordId(v uint64) {
	r.LiveRecordId = v
	r.touchUpdatedAt()
	syndb.AddData(TbVoiceChatMicApply, VoiceChatMicApplyLiveRecordId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *VoiceChatMicApply) SetStatus(v uint8) {
	r.Status = v
	r.touchUpdatedAt()
	syndb.AddData(TbVoiceChatMicApply, VoiceChatMicApplyStatus, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *VoiceChatMicApply) SetAppliedAt(v *time.Time) {
	r.AppliedAt = v
	r.touchUpdatedAt()
	syndb.AddData(TbVoiceChatMicApply, VoiceChatMicApplyAppliedAt, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *VoiceChatMicApply) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbVoiceChatMicApply, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *VoiceChatMicApply) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddData(TbVoiceChatMicApply, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *VoiceChatMicApply) touchUpdatedAt() {
	r.UpdatedAt = time.Now()
	syndb.AddData(TbVoiceChatMicApply, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: r.UpdatedAt})
}

func initVoiceChatMicApply() {
	migrate.AutoMigrate(&VoiceChatMicApply{})
	syndb.RegQuick(TbVoiceChatMicApply, db.CreatedAtName)
	syndb.RegQuick(TbVoiceChatMicApply, db.UpdatedAtName)
	syndb.RegQuick(TbVoiceChatMicApply, VoiceChatMicApplyRoomId)
	syndb.RegQuick(TbVoiceChatMicApply, VoiceChatMicApplyLiveRecordId)
	syndb.RegQuick(TbVoiceChatMicApply, VoiceChatMicApplyStatus)
	syndb.RegLazy(TbVoiceChatMicApply, VoiceChatMicApplyAppliedAt)
}
