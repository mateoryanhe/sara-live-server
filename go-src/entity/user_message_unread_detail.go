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
	TbUserMessageUnreadDetail db.TbName = "user_message_unread_details"
)

const (
	UserMessageUnreadDetailReceiverId  db.TbCol = "receiver_id"
	UserMessageUnreadDetailSenderId    db.TbCol = "sender_id"
	UserMessageUnreadDetailUnreadCount db.TbCol = "unread_count"
)

// UserMessageUnreadDetail 用户未读消息明细(按接收者+发送者维度)
// 走 syndb 快速同步缓冲(quick chan),变更通过 Setter 推送到队列由 worker 周期性 Save 到 DB
type UserMessageUnreadDetail struct {
	ID          string    `gorm:"primaryKey;size:64;comment:复合ID(receiverId_senderId)" json:"id"`
	ReceiverId  uint64    `gorm:"index:idx_receiver_sender,priority:1;default:0;comment:接收者ID" json:"receiverId"`
	SenderId    uint64    `gorm:"index:idx_receiver_sender,priority:2;default:0;comment:发送者ID" json:"senderId"`
	UnreadCount uint64    `gorm:"default:0;comment:未读数量" json:"unreadCount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func BuildUserMessageUnreadDetailId(receiverId, senderId uint64) string {
	return fmt.Sprintf("%d_%d", receiverId, senderId)
}

func NewUserMessageUnreadDetail(receiverId, senderId uint64) *UserMessageUnreadDetail {
	row := &UserMessageUnreadDetail{}
	row.ID = BuildUserMessageUnreadDetailId(receiverId, senderId)
	now := time.Now()
	row.SetCreatedAt(now)
	row.SetUpdatedAt(now)
	row.SetReceiverId(receiverId)
	row.SetSenderId(senderId)
	return row
}

func (m *UserMessageUnreadDetail) SetReceiverId(v uint64) {
	m.ReceiverId = v
	syndb.AddDataToQuickChan(TbUserMessageUnreadDetail, UserMessageUnreadDetailReceiverId, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserMessageUnreadDetail) SetSenderId(v uint64) {
	m.SenderId = v
	syndb.AddDataToQuickChan(TbUserMessageUnreadDetail, UserMessageUnreadDetailSenderId, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserMessageUnreadDetail) AddUnread(v uint64) {
	m.UnreadCount = math.Add(m.UnreadCount, v)
	m.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserMessageUnreadDetail, UserMessageUnreadDetailUnreadCount, &syndb.ColData{
		IdVal: m.ID, ColVal: m.UnreadCount,
	})
}

func (m *UserMessageUnreadDetail) ClearUnread() {
	m.UnreadCount = 0
	m.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserMessageUnreadDetail, UserMessageUnreadDetailUnreadCount, &syndb.ColData{
		IdVal: m.ID, ColVal: 0,
	})
}

func (m *UserMessageUnreadDetail) SetCreatedAt(v time.Time) {
	m.CreatedAt = v
	syndb.AddDataToQuickChan(TbUserMessageUnreadDetail, db.CreatedAtName, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserMessageUnreadDetail) SetUpdatedAt(v time.Time) {
	m.UpdatedAt = v
	syndb.AddDataToQuickChan(TbUserMessageUnreadDetail, db.UpdatedAtName, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func initUserMessageUnreadDetail() {
	syndb.RegQuickWithLarge(TbUserMessageUnreadDetail, db.CreatedAtName)
	syndb.RegQuickWithLarge(TbUserMessageUnreadDetail, db.UpdatedAtName)
	syndb.RegQuickWithLarge(TbUserMessageUnreadDetail, UserMessageUnreadDetailReceiverId)
	syndb.RegQuickWithLarge(TbUserMessageUnreadDetail, UserMessageUnreadDetailSenderId)
	syndb.RegQuickWithLarge(TbUserMessageUnreadDetail, UserMessageUnreadDetailUnreadCount)
	migrate.AutoMigrate(&UserMessageUnreadDetail{})
}
