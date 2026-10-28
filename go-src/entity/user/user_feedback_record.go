package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const TbUserFeedbackRecord db.TbName = "user_feedback_records"

const (
	UserFeedbackRecordUserId  db.TbCol = "user_id"
	UserFeedbackRecordContent db.TbCol = "content"
)

// UserFeedbackRecord App 用户反馈上报(append-only, syndb 缓冲入库)
type UserFeedbackRecord struct {
	migrate.OneModel
	UserId  uint64 `gorm:"index;default:0;comment:用户ID" json:"userId"`
	Content string `gorm:"type:mediumtext;comment:App上报原文" json:"content"`
}

func NewUserFeedbackRecord(userId uint64, content string) *UserFeedbackRecord {
	now := time.Now()
	row := &UserFeedbackRecord{}
	row.ID = snowflake.GetId()
	row.SetCreatedAt(now)
	row.SetUpdatedAt(now)
	row.SetUserId(userId)
	row.SetContent(content)
	return row
}

func (r *UserFeedbackRecord) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddData(TbUserFeedbackRecord, UserFeedbackRecordUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *UserFeedbackRecord) SetContent(v string) {
	r.Content = v
	syndb.AddData(TbUserFeedbackRecord, UserFeedbackRecordContent, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *UserFeedbackRecord) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbUserFeedbackRecord, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *UserFeedbackRecord) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddData(TbUserFeedbackRecord, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initUserFeedbackRecord() {
	syndb.RegQuick(TbUserFeedbackRecord, db.CreatedAtName)
	syndb.RegQuick(TbUserFeedbackRecord, db.UpdatedAtName)
	syndb.RegQuick(TbUserFeedbackRecord, UserFeedbackRecordUserId)
	syndb.RegQuick(TbUserFeedbackRecord, UserFeedbackRecordContent)
	migrate.AutoMigrate(&UserFeedbackRecord{})
}
