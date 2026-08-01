package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbUserPersonalSystemMessage db.TbName = "user_personal_system_messages"
)

const (
	UserPersonalSystemMessageUserId        db.TbCol = "user_id"
	UserPersonalSystemMessageMessageTypeId db.TbCol = "message_type_id"
	UserPersonalSystemMessageParams        db.TbCol = "params"
)

const (
	PersonalSystemMessageTypeWelcome uint32 = 1 // 欢迎消息
)

// UserPersonalSystemMessage 用户个人系统消息
// 走 syndb 快速同步缓冲(quick chan),变更通过 Setter 推送到队列由 worker 周期性 Save 到 DB
type UserPersonalSystemMessage struct {
	migrate.OneModel
	UserId        uint64 `gorm:"index;default:0;comment:用户ID" json:"userId"`
	MessageTypeId uint32 `gorm:"default:0;comment:消息类型ID" json:"messageTypeId"`
	Params        string `gorm:"size:1024;default:'';comment:参数" json:"params"`
}

func NewUserPersonalSystemMessage(userId uint64, messageTypeId uint32, params string) *UserPersonalSystemMessage {
	row := &UserPersonalSystemMessage{}
	row.ID = snowflake.GetId()
	now := time.Now()
	row.SetCreatedAt(now)
	row.SetUpdatedAt(now)
	row.SetUserId(userId)
	row.SetMessageTypeId(messageTypeId)
	row.SetParams(params)
	return row
}

func (m *UserPersonalSystemMessage) SetUserId(v uint64) {
	m.UserId = v
	syndb.AddDataToQuickChan(TbUserPersonalSystemMessage, UserPersonalSystemMessageUserId, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserPersonalSystemMessage) SetMessageTypeId(v uint32) {
	m.MessageTypeId = v
	syndb.AddDataToQuickChan(TbUserPersonalSystemMessage, UserPersonalSystemMessageMessageTypeId, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserPersonalSystemMessage) SetParams(v string) {
	m.Params = v
	syndb.AddDataToQuickChan(TbUserPersonalSystemMessage, UserPersonalSystemMessageParams, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserPersonalSystemMessage) SetCreatedAt(v time.Time) {
	m.CreatedAt = v
	syndb.AddDataToQuickChan(TbUserPersonalSystemMessage, db.CreatedAtName, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserPersonalSystemMessage) SetUpdatedAt(v time.Time) {
	m.UpdatedAt = v
	syndb.AddDataToQuickChan(TbUserPersonalSystemMessage, db.UpdatedAtName, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func initUserPersonalSystemMessage() {
	syndb.RegQuick(TbUserPersonalSystemMessage, db.CreatedAtName)
	syndb.RegQuick(TbUserPersonalSystemMessage, db.UpdatedAtName)
	syndb.RegQuick(TbUserPersonalSystemMessage, UserPersonalSystemMessageUserId)
	syndb.RegQuick(TbUserPersonalSystemMessage, UserPersonalSystemMessageMessageTypeId)
	syndb.RegQuick(TbUserPersonalSystemMessage, UserPersonalSystemMessageParams)
	migrate.AutoMigrate(&UserPersonalSystemMessage{})
}
