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

// 消息类型
const (
	UserMessageTypeSystem  uint8 = 1 // 系统消息
	UserMessageTypePrivate uint8 = 2 // 私信消息
)

const (
	UserMessageType       db.TbCol = "type"
	UserMessageSenderId   db.TbCol = "sender_id"
	UserMessageReceiverId db.TbCol = "receiver_id"
	UserMessageTitle      db.TbCol = "title"
	UserMessageContent    db.TbCol = "content"
)

// UserMessage 用户消息(系统消息/私信)
// 走 syndb 快速同步缓冲(quick chan),变更通过 Setter 推送到队列由 worker 周期性 Save 到 DB
type UserMessage struct {
	migrate.OneModel
	Type       uint8  `gorm:"index:idx_receiver_type,priority:2;default:0;comment:消息类型(1-系统消息,2-私信)" json:"type"`
	SenderId   uint64 `gorm:"index;default:0;comment:发送者ID(系统消息可为0)" json:"senderId"`
	ReceiverId uint64 `gorm:"index:idx_receiver_type,priority:1;default:0;comment:接收者ID" json:"receiverId"`
	Title      string `gorm:"size:128;default:'';comment:标题(系统消息使用)" json:"title"`
	Content    string `gorm:"size:1024;default:'';comment:消息内容" json:"content"`
}

// NewUserMessage 构造一条新消息
func NewUserMessage(msgType uint8, senderId, receiverId uint64, title, content string) *UserMessage {
	ret := &UserMessage{}
	ret.ID = snowflake.GetId()
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	ret.SetType(msgType)
	ret.SetSenderId(senderId)
	ret.SetReceiverId(receiverId)
	ret.SetTitle(title)
	ret.SetContent(content)
	return ret
}

func (m *UserMessage) SetType(v uint8) {
	m.Type = v
	syndb.AddDataToQuickChan(TbUserMessage, UserMessageType, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *UserMessage) SetSenderId(v uint64) {
	m.SenderId = v
	syndb.AddDataToQuickChan(TbUserMessage, UserMessageSenderId, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *UserMessage) SetReceiverId(v uint64) {
	m.ReceiverId = v
	syndb.AddDataToQuickChan(TbUserMessage, UserMessageReceiverId, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *UserMessage) SetTitle(v string) {
	m.Title = v
	syndb.AddDataToQuickChan(TbUserMessage, UserMessageTitle, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *UserMessage) SetContent(v string) {
	m.Content = v
	syndb.AddDataToQuickChan(TbUserMessage, UserMessageContent, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *UserMessage) SetCreatedAt(v time.Time) {
	m.CreatedAt = v
	syndb.AddDataToQuickChan(TbUserMessage, db.CreatedAtName, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *UserMessage) SetUpdatedAt(v time.Time) {
	m.UpdatedAt = v
	syndb.AddDataToQuickChan(TbUserMessage, db.UpdatedAtName, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func initUserMessage() {
	syndb.RegQuickWithLarge(TbUserMessage, db.CreatedAtName)
	syndb.RegQuickWithLarge(TbUserMessage, db.UpdatedAtName)
	syndb.RegQuickWithLarge(TbUserMessage, UserMessageType)
	syndb.RegQuickWithLarge(TbUserMessage, UserMessageSenderId)
	syndb.RegQuickWithLarge(TbUserMessage, UserMessageReceiverId)
	syndb.RegQuickWithLarge(TbUserMessage, UserMessageTitle)
	syndb.RegQuickWithLarge(TbUserMessage, UserMessageContent)
	migrate.AutoMigrate(&UserMessage{})
}
