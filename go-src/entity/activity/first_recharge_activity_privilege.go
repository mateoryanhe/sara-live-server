package activity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbFirstRechargeActivityPrivilege db.TbName = "first_recharge_activity_privileges"
)

// FirstRechargeActivityPrivilege 首充活动特权项(CMS 管理,保存时全量覆盖)
type FirstRechargeActivityPrivilege struct {
	migrate.OneModel
	Icon   string `gorm:"size:255;default:'';comment:特权图标资源名" json:"icon"`
	DescEn string `gorm:"size:255;default:'';comment:特权描述(英文)" json:"descEn"`
	DescEs string `gorm:"size:255;default:'';comment:特权描述(西班牙语)" json:"descEs"`
	DescPt string `gorm:"size:255;default:'';comment:特权描述(葡萄牙语)" json:"descPt"`
	DescHi string `gorm:"size:255;default:'';comment:特权描述(印地语)" json:"descHi"`
	DescId string `gorm:"size:255;default:'';comment:特权描述(印尼语)" json:"descId"`
	Sort   int    `gorm:"default:0;comment:排序值(越小越靠前)" json:"sort"`
}

func initFirstRechargeActivityPrivilege() {
	migrate.AutoMigrate(&FirstRechargeActivityPrivilege{})
}
