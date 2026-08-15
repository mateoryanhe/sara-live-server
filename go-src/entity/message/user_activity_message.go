package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbUserActivityMessage db.TbName = "user_activity_messages"
)

const (
	UserActivityMessageActivityMessageId db.TbCol = "activity_message_id"
	UserActivityMessageUserId            db.TbCol = "user_id"
	UserActivityMessagePublishedAt       db.TbCol = "published_at"
)

// UserActivityMessage 用户活动消息
// 走 syndb 快速同步缓冲(quick chan),变更通过 Setter 推送到队列由 worker 周期性 Save 到 DB
type UserActivityMessage struct {
	migrate.OneModel
	ActivityMessageId uint64     `gorm:"index;default:0;comment:活动消息ID" json:"activityMessageId"`
	UserId            uint64     `gorm:"index;default:0;comment:用户ID" json:"userId"`
	PublishedAt       *time.Time `gorm:"comment:发布时间(同活动消息发布时间)" json:"publishedAt"`
}

func NewUserActivityMessage(userId, activityMessageId uint64, publishedAt *time.Time) *UserActivityMessage {
	row := &UserActivityMessage{}
	row.ID = snowflake.GetId()
	now := time.Now()
	row.SetCreatedAt(now)
	row.SetUpdatedAt(now)
	row.SetUserId(userId)
	row.SetActivityMessageId(activityMessageId)
	row.SetPublishedAt(publishedAt)
	return row
}

func (m *UserActivityMessage) SetActivityMessageId(v uint64) {
	m.ActivityMessageId = v
	syndb.AddData(TbUserActivityMessage, UserActivityMessageActivityMessageId, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserActivityMessage) SetUserId(v uint64) {
	m.UserId = v
	syndb.AddData(TbUserActivityMessage, UserActivityMessageUserId, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserActivityMessage) SetPublishedAt(v *time.Time) {
	m.PublishedAt = v
	syndb.AddData(TbUserActivityMessage, UserActivityMessagePublishedAt, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserActivityMessage) SetCreatedAt(v time.Time) {
	m.CreatedAt = v
	syndb.AddData(TbUserActivityMessage, db.CreatedAtName, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserActivityMessage) SetUpdatedAt(v time.Time) {
	m.UpdatedAt = v
	syndb.AddData(TbUserActivityMessage, db.UpdatedAtName, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func initUserActivityMessage() {
	syndb.RegQuick(TbUserActivityMessage, db.CreatedAtName)
	syndb.RegQuick(TbUserActivityMessage, db.UpdatedAtName)
	syndb.RegQuick(TbUserActivityMessage, UserActivityMessageActivityMessageId)
	syndb.RegQuick(TbUserActivityMessage, UserActivityMessageUserId)
	syndb.RegQuick(TbUserActivityMessage, UserActivityMessagePublishedAt)
	migrate.AutoMigrate(&UserActivityMessage{})
}
