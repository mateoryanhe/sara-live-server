package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const TbUserAttributionReport db.TbName = "user_attribution_reports"

const (
	UserAttributionReportUserId             db.TbCol = "user_id"
	UserAttributionReportPackageName        db.TbCol = "package_name"
	UserAttributionReportAttributionEnabled db.TbCol = "attribution_enabled"
	UserAttributionReportAppsFlyerDevKey    db.TbCol = "apps_flyer_dev_key"
	UserAttributionReportSource             db.TbCol = "source"
	UserAttributionReportMediaSource        db.TbCol = "media_source"
	UserAttributionReportCampaign           db.TbCol = "campaign"
	UserAttributionReportCreativeId         db.TbCol = "creative_id"
)

// UserAttributionReport App 归因上报记录(append-only, syndb 缓冲入库)
type UserAttributionReport struct {
	migrate.OneModel
	UserId             uint64 `gorm:"index;default:0;comment:用户ID" json:"userId"`
	PackageName        string `gorm:"size:128;default:'';index;comment:包名" json:"packageName"`
	AttributionEnabled bool   `gorm:"default:0;comment:是否启用归因" json:"attributionEnabled"`
	AppsFlyerDevKey    string `gorm:"size:128;default:'';comment:AppsFlyer Dev Key" json:"appsFlyerDevKey"`
	Source             string `gorm:"size:32;default:'';comment:AF source(如 paid/organic)" json:"source"`
	MediaSource        string `gorm:"size:128;default:'';comment:AF media_source" json:"mediaSource"`
	Campaign           string `gorm:"size:256;default:'';comment:AF campaign" json:"campaign"`
	CreativeId         string `gorm:"size:128;default:'';comment:AF creative_id" json:"creativeId"`
}

func NewUserAttributionReport(
	userId uint64,
	packageName string,
	attributionEnabled bool,
	appsFlyerDevKey, source, mediaSource, campaign, creativeId string,
) *UserAttributionReport {
	now := time.Now()
	row := &UserAttributionReport{}
	row.ID = snowflake.GetId()
	row.SetCreatedAt(now)
	row.SetUpdatedAt(now)
	row.SetUserId(userId)
	row.SetPackageName(packageName)
	row.SetAttributionEnabled(attributionEnabled)
	row.SetAppsFlyerDevKey(appsFlyerDevKey)
	row.SetSource(source)
	row.SetMediaSource(mediaSource)
	row.SetCampaign(campaign)
	row.SetCreativeId(creativeId)
	return row
}

func (r *UserAttributionReport) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddData(TbUserAttributionReport, UserAttributionReportUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *UserAttributionReport) SetPackageName(v string) {
	r.PackageName = v
	syndb.AddData(TbUserAttributionReport, UserAttributionReportPackageName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *UserAttributionReport) SetAttributionEnabled(v bool) {
	r.AttributionEnabled = v
	syndb.AddData(TbUserAttributionReport, UserAttributionReportAttributionEnabled, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *UserAttributionReport) SetAppsFlyerDevKey(v string) {
	r.AppsFlyerDevKey = v
	syndb.AddData(TbUserAttributionReport, UserAttributionReportAppsFlyerDevKey, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *UserAttributionReport) SetSource(v string) {
	r.Source = v
	syndb.AddData(TbUserAttributionReport, UserAttributionReportSource, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *UserAttributionReport) SetMediaSource(v string) {
	r.MediaSource = v
	syndb.AddData(TbUserAttributionReport, UserAttributionReportMediaSource, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *UserAttributionReport) SetCampaign(v string) {
	r.Campaign = v
	syndb.AddData(TbUserAttributionReport, UserAttributionReportCampaign, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *UserAttributionReport) SetCreativeId(v string) {
	r.CreativeId = v
	syndb.AddData(TbUserAttributionReport, UserAttributionReportCreativeId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *UserAttributionReport) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbUserAttributionReport, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *UserAttributionReport) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddData(TbUserAttributionReport, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initUserAttributionReport() {
	syndb.RegQuick(TbUserAttributionReport, db.CreatedAtName)
	syndb.RegQuick(TbUserAttributionReport, db.UpdatedAtName)
	syndb.RegQuick(TbUserAttributionReport, UserAttributionReportUserId)
	syndb.RegQuick(TbUserAttributionReport, UserAttributionReportPackageName)
	syndb.RegQuick(TbUserAttributionReport, UserAttributionReportAttributionEnabled)
	syndb.RegQuick(TbUserAttributionReport, UserAttributionReportAppsFlyerDevKey)
	syndb.RegQuick(TbUserAttributionReport, UserAttributionReportSource)
	syndb.RegQuick(TbUserAttributionReport, UserAttributionReportMediaSource)
	syndb.RegQuick(TbUserAttributionReport, UserAttributionReportCampaign)
	syndb.RegQuick(TbUserAttributionReport, UserAttributionReportCreativeId)
	migrate.AutoMigrate(&UserAttributionReport{})
}