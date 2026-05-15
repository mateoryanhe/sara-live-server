package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbCMSUser db.TbName = "cms_users"
)

// 列名
const (
	CMSUserName db.TbCol = "name"
)

type CMSUser struct {
	migrate.OneModel
	Name   string `gorm:"comment:名称" json:"name"`
	Pwd    string `gorm:"comment:密码" json:"pwd"`
	Status uint8  `gorm:"default:0;comment:状态" json:"status"`
	Admin  bool   `gorm:"default:0;comment:是否是管理员" json:"admin"`
	RoleId uint64 `gorm:"default:0;comment:角色ID" json:"roleId"`
}

func InitCMSUser() {
	migrate.AutoMigrate(&CMSUser{})
}
