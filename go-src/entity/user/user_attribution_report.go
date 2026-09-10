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
	UserAttributionReportUserId               db.TbCol = "user_id"
	UserAttributionReportPackageName          db.TbCol = "package_name"
	UserAttributionReportAttributionEnabled   db.TbCol = "attribution_enabled"
	UserAttributionReportAttributionProvider  db.TbCol = "attribution_provider"
	UserAttributionReportAppsFlyerDevKey      db.TbCol = "apps_flyer_dev_key"
	UserAttributionReportAppsFlyerAppId       db.TbCol = "apps_flyer_app_id"
)

// UserAttributionReport App 归因上报记录(append-only, syndb 缓冲入库)
type UserAttributionReport struct {
	migrate.OneModel
	UserId               uint64 `gorm:"index;default:0;comment:用户ID" json:"userId"`
	PackageName          string `gorm:"size:128;default:'';index;comment:包名" json:"packageName"`
	AttributionEnabled   bool   `gorm:"default:0;comment:是否启用归因" json:"attributionEnabled"`
	AttributionProvider  string `gorm:"size:64;default:'';comment:归因渠道" json:"attributionProvider"`
	AppsFlyerDevKey      string `gorm:"size:128;default:'';comment:AppsFlyer Dev Key" json:"appsFlyerDevKey"`
	AppsFlyerAppId       string `gorm:"size:128;default:'';comment:AppsFlyer App ID" json:"appsFlyerAppId"`
}

func NewUserAttributionReport(
	userId uint64,
	packageName string,
	attributionEnabled bool,
	attributionProvider, appsFlyerDevKey, appsFlyerAppId string,
) *UserAttributionReport {
	now := time.Now()
	row := &UserAttributionReport{}
	row.ID = snowflake.GetId()
	row.SetCreatedAt(now)
	row.SetUpdatedAt(now)
	row.SetUserId(userId)
	row.SetPackageName(packageName)
	row.SetAttributionEnabled(attributionEnabled)
	row.SetAttributionProvider(attributionProvider)
	row.SetAppsFlyerDevKey(appsFlyerDevKey)
	row.SetAppsFlyerAppId(appsFlyerAppId)
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

func (r *UserAttributionReport) SetAttributionProvider(v string) {
	r.AttributionProvider = v
	syndb.AddData(TbUserAttributionReport, UserAttributionReportAttributionProvider, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *UserAttributionReport) SetAppsFlyerDevKey(v string) {
	r.AppsFlyerDevKey = v
	syndb.AddData(TbUserAttributionReport, UserAttributionReportAppsFlyerDevKey, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *UserAttributionReport) SetAppsFlyerAppId(v string) {
	r.AppsFlyerAppId = v
	syndb.AddData(TbUserAttributionReport, UserAttributionReportAppsFlyerAppId, &syndb.ColData{IdVal: r.ID, ColVal: v})
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
	syndb.RegQuick(TbUserAttributionReport, UserAttributionReportAttributionProvider)
	syndb.RegQuick(TbUserAttributionReport, UserAttributionReportAppsFlyerDevKey)
	syndb.RegQuick(TbUserAttributionReport, UserAttributionReportAppsFlyerAppId)
	migrate.AutoMigrate(&UserAttributionReport{})
}
