package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const TbUserReportRecord db.TbName = "user_report_records"

const (
	UserReportRecordUserId  db.TbCol = "user_id"
	UserReportRecordContent db.TbCol = "content"
)

// UserReportRecord App 用户举报上报(append-only, syndb 缓冲入库)
type UserReportRecord struct {
	migrate.OneModel
	UserId  uint64 `gorm:"index;default:0;comment:用户ID" json:"userId"`
	Content string `gorm:"type:mediumtext;comment:App上报原文" json:"content"`
}

func NewUserReportRecord(userId uint64, content string) *UserReportRecord {
	now := time.Now()
	row := &UserReportRecord{}
	row.ID = snowflake.GetId()
	row.SetCreatedAt(now)
	row.SetUpdatedAt(now)
	row.SetUserId(userId)
	row.SetContent(content)
	return row
}

func (r *UserReportRecord) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddData(TbUserReportRecord, UserReportRecordUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *UserReportRecord) SetContent(v string) {
	r.Content = v
	syndb.AddData(TbUserReportRecord, UserReportRecordContent, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *UserReportRecord) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbUserReportRecord, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *UserReportRecord) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddData(TbUserReportRecord, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initUserReportRecord() {
	syndb.RegQuick(TbUserReportRecord, db.CreatedAtName)
	syndb.RegQuick(TbUserReportRecord, db.UpdatedAtName)
	syndb.RegQuick(TbUserReportRecord, UserReportRecordUserId)
	syndb.RegQuick(TbUserReportRecord, UserReportRecordContent)
	migrate.AutoMigrate(&UserReportRecord{})
}
