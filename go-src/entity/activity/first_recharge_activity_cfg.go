package activity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbFirstRechargeActivityCfg db.TbName = "first_recharge_activity_cfgs"
)

// FirstRechargeActivityCfg 首充活动配置(CMS 管理,通常仅一条)
type FirstRechargeActivityCfg struct {
	migrate.OneModel
	Enabled           bool   `gorm:"default:0;comment:活动开关" json:"enabled"`
	Icon              string `gorm:"size:255;default:'';comment:小图标资源名" json:"icon"`
	TitleEn           string `gorm:"size:128;default:'';comment:标题(英文)" json:"titleEn"`
	TitleEs           string `gorm:"size:128;default:'';comment:标题(西班牙语)" json:"titleEs"`
	TitlePt           string `gorm:"size:128;default:'';comment:标题(葡萄牙语)" json:"titlePt"`
	TitleHi           string `gorm:"size:128;default:'';comment:标题(印地语)" json:"titleHi"`
	TitleId           string `gorm:"size:128;default:'';comment:标题(印尼语)" json:"titleId"`
	RechargeBtnTextEn string `gorm:"size:64;default:'';comment:充值按钮文案(英文)" json:"rechargeBtnTextEn"`
	RechargeBtnTextEs string `gorm:"size:64;default:'';comment:充值按钮文案(西班牙语)" json:"rechargeBtnTextEs"`
	RechargeBtnTextPt string `gorm:"size:64;default:'';comment:充值按钮文案(葡萄牙语)" json:"rechargeBtnTextPt"`
	RechargeBtnTextHi string  `gorm:"size:64;default:'';comment:充值按钮文案(印地语)" json:"rechargeBtnTextHi"`
	RechargeBtnTextId string  `gorm:"size:64;default:'';comment:充值按钮文案(印尼语)" json:"rechargeBtnTextId"`
	FirstRechargeRatio float64 `gorm:"type:decimal(6,2);default:20;comment:首充比例(%)" json:"firstRechargeRatio"`
}

func initFirstRechargeActivityCfg() {
	migrate.AutoMigrate(&FirstRechargeActivityCfg{})
}
