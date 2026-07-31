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
	UserPersonalSystemMessageUserId    db.TbCol = "user_id"
	UserPersonalSystemMessageIcon      db.TbCol = "icon"
	UserPersonalSystemMessageTitleEn   db.TbCol = "title_en"
	UserPersonalSystemMessageTitleEs   db.TbCol = "title_es"
	UserPersonalSystemMessageTitlePt   db.TbCol = "title_pt"
	UserPersonalSystemMessageTitleHi   db.TbCol = "title_hi"
	UserPersonalSystemMessageContentEn db.TbCol = "content_en"
	UserPersonalSystemMessageContentEs db.TbCol = "content_es"
	UserPersonalSystemMessageContentPt db.TbCol = "content_pt"
	UserPersonalSystemMessageContentHi db.TbCol = "content_hi"
	UserPersonalSystemMessageParams    db.TbCol = "params"
)

// UserPersonalSystemMessage 用户个人系统消息
// 走 syndb 快速同步缓冲(quick chan),变更通过 Setter 推送到队列由 worker 周期性 Save 到 DB
type UserPersonalSystemMessage struct {
	migrate.OneModel
	UserId    uint64 `gorm:"index;default:0;comment:用户ID" json:"userId"`
	Icon      string `gorm:"size:255;default:'';comment:图标" json:"icon"`
	TitleEn   string `gorm:"size:128;default:'';comment:标题(英文)" json:"titleEn"`
	TitleEs   string `gorm:"size:128;default:'';comment:标题(西班牙语)" json:"titleEs"`
	TitlePt   string `gorm:"size:128;default:'';comment:标题(葡萄牙语)" json:"titlePt"`
	TitleHi   string `gorm:"size:128;default:'';comment:标题(印地语)" json:"titleHi"`
	ContentEn string `gorm:"type:text;comment:内容(英文)" json:"contentEn"`
	ContentEs string `gorm:"type:text;comment:内容(西班牙语)" json:"contentEs"`
	ContentPt string `gorm:"type:text;comment:内容(葡萄牙语)" json:"contentPt"`
	ContentHi string `gorm:"type:text;comment:内容(印地语)" json:"contentHi"`
	Params    string `gorm:"size:1024;default:'';comment:参数" json:"params"`
}

func NewUserPersonalSystemMessage(
	userId uint64,
	icon, params string,
	titleEn, titleEs, titlePt, titleHi string,
	contentEn, contentEs, contentPt, contentHi string,
) *UserPersonalSystemMessage {
	row := &UserPersonalSystemMessage{}
	row.ID = snowflake.GetId()
	now := time.Now()
	row.SetCreatedAt(now)
	row.SetUpdatedAt(now)
	row.SetUserId(userId)
	row.SetIcon(icon)
	row.SetParams(params)
	row.SetTitleEn(titleEn)
	row.SetTitleEs(titleEs)
	row.SetTitlePt(titlePt)
	row.SetTitleHi(titleHi)
	row.SetContentEn(contentEn)
	row.SetContentEs(contentEs)
	row.SetContentPt(contentPt)
	row.SetContentHi(contentHi)
	return row
}

func (m *UserPersonalSystemMessage) SetUserId(v uint64) {
	m.UserId = v
	syndb.AddDataToQuickChan(TbUserPersonalSystemMessage, UserPersonalSystemMessageUserId, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserPersonalSystemMessage) SetIcon(v string) {
	m.Icon = v
	syndb.AddDataToQuickChan(TbUserPersonalSystemMessage, UserPersonalSystemMessageIcon, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserPersonalSystemMessage) SetTitleEn(v string) {
	m.TitleEn = v
	syndb.AddDataToQuickChan(TbUserPersonalSystemMessage, UserPersonalSystemMessageTitleEn, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserPersonalSystemMessage) SetTitleEs(v string) {
	m.TitleEs = v
	syndb.AddDataToQuickChan(TbUserPersonalSystemMessage, UserPersonalSystemMessageTitleEs, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserPersonalSystemMessage) SetTitlePt(v string) {
	m.TitlePt = v
	syndb.AddDataToQuickChan(TbUserPersonalSystemMessage, UserPersonalSystemMessageTitlePt, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserPersonalSystemMessage) SetTitleHi(v string) {
	m.TitleHi = v
	syndb.AddDataToQuickChan(TbUserPersonalSystemMessage, UserPersonalSystemMessageTitleHi, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserPersonalSystemMessage) SetContentEn(v string) {
	m.ContentEn = v
	syndb.AddDataToQuickChan(TbUserPersonalSystemMessage, UserPersonalSystemMessageContentEn, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserPersonalSystemMessage) SetContentEs(v string) {
	m.ContentEs = v
	syndb.AddDataToQuickChan(TbUserPersonalSystemMessage, UserPersonalSystemMessageContentEs, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserPersonalSystemMessage) SetContentPt(v string) {
	m.ContentPt = v
	syndb.AddDataToQuickChan(TbUserPersonalSystemMessage, UserPersonalSystemMessageContentPt, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserPersonalSystemMessage) SetContentHi(v string) {
	m.ContentHi = v
	syndb.AddDataToQuickChan(TbUserPersonalSystemMessage, UserPersonalSystemMessageContentHi, &syndb.ColData{
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
	syndb.RegQuick(TbUserPersonalSystemMessage, UserPersonalSystemMessageIcon)
	syndb.RegQuick(TbUserPersonalSystemMessage, UserPersonalSystemMessageTitleEn)
	syndb.RegQuick(TbUserPersonalSystemMessage, UserPersonalSystemMessageTitleEs)
	syndb.RegQuick(TbUserPersonalSystemMessage, UserPersonalSystemMessageTitlePt)
	syndb.RegQuick(TbUserPersonalSystemMessage, UserPersonalSystemMessageTitleHi)
	syndb.RegQuick(TbUserPersonalSystemMessage, UserPersonalSystemMessageContentEn)
	syndb.RegQuick(TbUserPersonalSystemMessage, UserPersonalSystemMessageContentEs)
	syndb.RegQuick(TbUserPersonalSystemMessage, UserPersonalSystemMessageContentPt)
	syndb.RegQuick(TbUserPersonalSystemMessage, UserPersonalSystemMessageContentHi)
	syndb.RegQuick(TbUserPersonalSystemMessage, UserPersonalSystemMessageParams)
	migrate.AutoMigrate(&UserPersonalSystemMessage{})
}
