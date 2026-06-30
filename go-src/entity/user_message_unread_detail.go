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
	UserMessageUnreadDetailSenderId    db.TbCol = "sender_id"
	UserMessageUnreadDetailUserId      db.TbCol = "user_id"
	UserMessageUnreadDetailUnreadCount db.TbCol = "unread_count"
)

// UserMessageUnreadDetail 用户未读消息明细(按接收者+发送者维度)
// 走 syndb 快速同步缓冲(quick chan),变更通过 Setter 推送到队列由 worker 周期性 Save 到 DB
type UserMessageUnreadDetail struct {
	ID          string    `gorm:"primaryKey;size:64;comment:复合ID(userId_senderId)" json:"id"`
	UnreadCount uint64    `gorm:"default:0;comment:未读数量" json:"unreadCount"`
	UserId      uint64    `gorm:"default:0;comment:用户id" json:"userId"`
	SenderId    uint64    `gorm:"default:0;comment:用户id" json:"senderId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func BuildUserMessageUnreadDetailId(receiverId, senderId uint64) string {
	return fmt.Sprintf("%d_%d", receiverId, senderId)
}

func NewUserMessageUnreadDetail(userId, senderId uint64) *UserMessageUnreadDetail {
	row := &UserMessageUnreadDetail{}
	row.ID = BuildUserMessageUnreadDetailId(userId, senderId)
	row.SetSenderId(senderId)
	now := time.Now()
	row.SetCreatedAt(now)
	row.SetUpdatedAt(now)
	row.SetUserId(userId)
	return row
}

func (m *UserMessageUnreadDetail) SetUserId(v uint64) {
	m.UserId = v
	syndb.AddDataToQuickChan(TbUserMessageUnreadDetail, UserMessageUnreadDetailUserId, &syndb.ColData{
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

func (m *UserMessageUnreadDetail) ClearUnread(v uint64) {
	m.UnreadCount = math.Sub(m.UnreadCount, v)
	syndb.AddDataToQuickChan(TbUserMessageUnreadDetail, UserMessageUnreadDetailUnreadCount, &syndb.ColData{
		IdVal: m.ID, ColVal: m.UnreadCount,
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
	syndb.RegQuickWithLarge(TbUserMessageUnreadDetail, UserMessageUnreadDetailSenderId)
	syndb.RegQuickWithLarge(TbUserMessageUnreadDetail, UserMessageUnreadDetailUserId)
	syndb.RegQuickWithLarge(TbUserMessageUnreadDetail, UserMessageUnreadDetailUnreadCount)
	migrate.AutoMigrate(&UserMessageUnreadDetail{})
}
