package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbUserMessage db.TbName = "user_messages"
)

const (
	UserMessageSenderId   db.TbCol = "sender_id"
	UserMessageReceiverId db.TbCol = "receiver_id"
	UserMessageContent    db.TbCol = "content"
)

// UserMessage 用户消息(系统消息/私信)
// 系统消息: sender_id=0; 私信: sender_id>0
// 走 syndb 快速同步缓冲(quick chan),变更通过 Setter 推送到队列由 worker 周期性 Save 到 DB
type UserMessage struct {
	migrate.MoreModel
	SenderId   uint64 `gorm:"index;default:0;comment:发送者ID(系统消息为0)" json:"senderId"`
	ReceiverId uint64 `gorm:"index;default:0;comment:接收者ID" json:"receiverId"`
	Content    string `gorm:"size:1024;default:'';comment:消息内容" json:"content"`
}

// IsSystemMessage 是否为系统消息
func (m *UserMessage) IsSystemMessage() bool {
	return m != nil && m.SenderId == 0
}

// IsPrivateMessage 是否为私信
func (m *UserMessage) IsPrivateMessage() bool {
	return m != nil && m.SenderId > 0
}

// SessionIdForUser 返回指定用户视角的私信会话ID
func (m *UserMessage) SessionIdForUser(userId uint64) string {
	if m == nil || !m.IsPrivateMessage() {
		return ""
	}
	if m.SenderId == userId {
		return BuildUserMessageSessionId(userId, m.ReceiverId)
	}
	if m.ReceiverId == userId {
		return BuildUserMessageSessionId(userId, m.SenderId)
	}
	return ""
}

// NewUserMessage 构造一条新消息
func NewUserMessage(senderId, receiverId uint64, content string) *UserMessage {
	ret := &UserMessage{}
	ret.ID = snowflake.GetId()
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	ret.SetSenderId(senderId)
	ret.SetReceiverId(receiverId)
	ret.SetContent(content)
	return ret
}

func (m *UserMessage) SetSenderId(v uint64) {
	m.SenderId = v
	syndb.AddData(TbUserMessage, UserMessageSenderId, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *UserMessage) SetReceiverId(v uint64) {
	m.ReceiverId = v
	syndb.AddData(TbUserMessage, UserMessageReceiverId, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *UserMessage) SetContent(v string) {
	m.Content = v
	syndb.AddData(TbUserMessage, UserMessageContent, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *UserMessage) SetCreatedAt(v time.Time) {
	m.CreatedAt = v
	syndb.AddData(TbUserMessage, db.CreatedAtName, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *UserMessage) SetUpdatedAt(v time.Time) {
	m.UpdatedAt = v
	syndb.AddData(TbUserMessage, db.UpdatedAtName, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *UserMessage) SetIsDeleted(v bool) {
	m.IsDeleted = v
	syndb.AddData(TbUserMessage, db.IsDeletedName, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *UserMessage) SetDeletedAt(v time.Time) {
	m.DeletedAt = v
	syndb.AddData(TbUserMessage, db.DeletedAtName, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

// MarkDeleted 标记消息为已删除
func (m *UserMessage) MarkDeleted() {
	now := time.Now()
	m.SetIsDeleted(true)
	m.SetDeletedAt(now)
	m.SetUpdatedAt(now)
}

// NewMarkDeletedUserMessage 标记消息删除(走 syndb,物理删库由定时任务执行)
func NewMarkDeletedUserMessage(id uint64) {
	if id == 0 {
		return
	}
	row := &UserMessage{}
	row.ID = id
	row.MarkDeleted()
}

func initUserMessage() {
	syndb.RegQuick(TbUserMessage, db.CreatedAtName)
	syndb.RegQuick(TbUserMessage, db.UpdatedAtName)
	syndb.RegQuick(TbUserMessage, db.IsDeletedName)
	syndb.RegQuick(TbUserMessage, db.DeletedAtName)
	syndb.RegQuick(TbUserMessage, UserMessageSenderId)
	syndb.RegQuick(TbUserMessage, UserMessageReceiverId)
	syndb.RegQuick(TbUserMessage, UserMessageContent)
	migrate.AutoMigrate(&UserMessage{})
}
