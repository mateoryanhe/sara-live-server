package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const TbVoiceChatMicSeatTier db.TbName = "voice_chat_mic_seat_tiers"

const (
	VoiceChatMicSeatTierSeatCount db.TbCol = "seat_count"
	VoiceChatMicSeatTierSort      db.TbCol = "sort"
	VoiceChatMicSeatTierStatus    db.TbCol = "status"
)

const (
	VoiceChatMicSeatTierStatusOffShelf uint8 = 0
	VoiceChatMicSeatTierStatusOnShelf  uint8 = 1
)

// VoiceChatMicSeatTier 语聊麦位挡位(CMS 管理,上架项缓存)
type VoiceChatMicSeatTier struct {
	migrate.OneModel
	SeatCount int   `gorm:"uniqueIndex;comment:麦位总数(含主播位,2~25)" json:"seatCount"`
	Sort      int   `gorm:"default:0;comment:排序值(越大越靠前)" json:"sort"`
	Status    uint8 `gorm:"default:0;comment:状态(0下架,1上架)" json:"status"`
}

func (r *VoiceChatMicSeatTier) queue(column db.TbCol, value any) {
	syndb.AddData(TbVoiceChatMicSeatTier, column, &syndb.ColData{IdVal: r.ID, ColVal: value})
}

func (r *VoiceChatMicSeatTier) touchUpdatedAt() {
	r.UpdatedAt = time.Now()
	r.queue(db.UpdatedAtName, r.UpdatedAt)
}

func (r *VoiceChatMicSeatTier) SetSeatCount(v int) {
	r.SeatCount = v
	r.touchUpdatedAt()
	r.queue(VoiceChatMicSeatTierSeatCount, v)
}

func (r *VoiceChatMicSeatTier) SetSort(v int) {
	r.Sort = v
	r.touchUpdatedAt()
	r.queue(VoiceChatMicSeatTierSort, v)
}

func (r *VoiceChatMicSeatTier) SetStatus(v uint8) {
	r.Status = v
	r.touchUpdatedAt()
	r.queue(VoiceChatMicSeatTierStatus, v)
}

func initVoiceChatMicSeatTier() {
	syndb.RegQuick(TbVoiceChatMicSeatTier, db.CreatedAtName)
	syndb.RegQuick(TbVoiceChatMicSeatTier, db.UpdatedAtName)
	syndb.RegQuick(TbVoiceChatMicSeatTier, VoiceChatMicSeatTierSeatCount)
	syndb.RegQuick(TbVoiceChatMicSeatTier, VoiceChatMicSeatTierSort)
	syndb.RegQuick(TbVoiceChatMicSeatTier, VoiceChatMicSeatTierStatus)
	migrate.AutoMigrate(&VoiceChatMicSeatTier{})
}
