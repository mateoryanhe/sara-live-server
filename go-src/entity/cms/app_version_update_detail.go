package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbAppVersionUpdateDetail db.TbName = "app_version_update_details"
)

// AppVersionUpdateDetail App版本更新明细(CMS 管理,保存时全量覆盖)
type AppVersionUpdateDetail struct {
	migrate.OneModel
	Content string `gorm:"size:512;default:'';comment:更新明细" json:"content"`
	Sort    int    `gorm:"default:0;comment:排序值(越小越靠前)" json:"sort"`
}

func initAppVersionUpdateDetail() {
	migrate.AutoMigrate(&AppVersionUpdateDetail{})
}
