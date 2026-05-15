package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbCMSRole db.TbName = "cms_roles"
)

type CMSRole struct {
	migrate.OneModel
	Name        string `gorm:"comment:角色名称" json:"name"`
	Description string `gorm:"comment:角色描述" json:"description"`
	Status      uint8  `gorm:"default:1;comment:状态(0-禁用,1-启用)" json:"status"`
}

func InitCMSRole() {
	migrate.AutoMigrate(&CMSRole{})
}
