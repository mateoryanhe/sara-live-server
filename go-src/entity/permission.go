package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbPermission db.TbName = "permissions"
)

type Permission struct {
	migrate.OneModel
	Module string `gorm:"comment:模块" json:"module"`
	RoleId uint64 `gorm:"comment:角色ID" json:"roleId"`
}

func InitPermission() {
	migrate.AutoMigrate(&Permission{})
}
