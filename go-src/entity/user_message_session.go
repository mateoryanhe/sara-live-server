package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbUserMessageSession db.TbName = "user_message_sessions"
)

const (
	UserMessageSessionMessageId db.TbCol = "message_id"
	UserMessageSessionSessionId db.TbCol = "session_id"
)

// UserMessageSession 私信消息会话索引(消息ID -> 会话ID)
// 走 syndb 快速同步缓冲(quick chan),变更通过 Setter 推送到队列由 worker 周期性 Save 到 DB
type UserMessageSession struct {
	migrate.OneModel
	MessageId uint64 `gorm:"index;default:0;comment:消息ID" json:"messageId"`
	SessionId string `gorm:"index;size:250;default:'';comment:会话ID(userId_targetId)" json:"sessionId"`
}

// BuildUserMessageSessionId 构建会话ID: 当前用户_目标用户
func BuildUserMessageSessionId(userId, targetId uint64) string {
	return fmt.Sprintf("%d_%d", userId, targetId)
}

// NewUserMessageSession 构造消息会话索引(ID 即消息ID)
func NewUserMessageSession(messageId uint64, sessionId string) *UserMessageSession {
	row := &UserMessageSession{}
	row.ID = snowflake.GetId()
	now := time.Now()
	row.SetCreatedAt(now)
	row.SetUpdatedAt(now)
	row.SetMessageId(messageId)
	row.SetSessionId(sessionId)
	return row
}

func (m *UserMessageSession) SetMessageId(v uint64) {
	m.MessageId = v
	syndb.AddDataToQuickChan(TbUserMessageSession, UserMessageSessionMessageId, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserMessageSession) SetSessionId(v string) {
	m.SessionId = v
	syndb.AddDataToQuickChan(TbUserMessageSession, UserMessageSessionSessionId, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserMessageSession) SetCreatedAt(v time.Time) {
	m.CreatedAt = v
	syndb.AddDataToQuickChan(TbUserMessageSession, db.CreatedAtName, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserMessageSession) SetUpdatedAt(v time.Time) {
	m.UpdatedAt = v
	syndb.AddDataToQuickChan(TbUserMessageSession, db.UpdatedAtName, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func initUserMessageSession() {
	syndb.RegQuickWithLarge(TbUserMessageSession, db.CreatedAtName)
	syndb.RegQuickWithLarge(TbUserMessageSession, db.UpdatedAtName)
	syndb.RegQuickWithLarge(TbUserMessageSession, UserMessageSessionMessageId)
	syndb.RegQuickWithLarge(TbUserMessageSession, UserMessageSessionSessionId)
	migrate.AutoMigrate(&UserMessageSession{})
}
