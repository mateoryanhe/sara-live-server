package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbFirebaseAnalyticsCfg db.TbName = "firebase_analytics_cfgs"
)

// FirebaseAnalyticsCfg Firebase Analytics(GA4) 埋点配置(CMS 管理,通常仅一条).
type FirebaseAnalyticsCfg struct {
	migrate.OneModel
	Enabled            uint8  `gorm:"default:0;comment:是否启用埋点(0否1是)" json:"enabled"`
	ProjectId          string `gorm:"size:128;default:'';comment:Firebase Project ID" json:"projectId"`
	ClientConfigJson      string `gorm:"type:text;comment:Firebase客户端配置JSON(含measurementId等)" json:"clientConfigJson"`
	ServiceAccountJson    string `gorm:"type:text;comment:Firebase服务账号JSON(Admin SDK)" json:"serviceAccountJson"`
	MeasurementApiSecret  string `gorm:"size:128;default:'';comment:GA4 Measurement Protocol API Secret" json:"measurementApiSecret"`
}

func (FirebaseAnalyticsCfg) TableName() string {
	return string(TbFirebaseAnalyticsCfg)
}

func (c *FirebaseAnalyticsCfg) IsActive() bool {
	return c != nil && c.ID > 0 && c.Enabled == 1 &&
		c.ProjectId != "" && c.ClientConfigJson != "" && c.ServiceAccountJson != "" &&
		c.MeasurementApiSecret != ""
}

func initFirebaseAnalyticsCfg() {
	migrate.AutoMigrate(&FirebaseAnalyticsCfg{})
}
